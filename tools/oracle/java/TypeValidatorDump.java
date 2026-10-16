// API-sweep oracle for com.google.javascript.jscomp.TypeValidator
// (internal/jscomp/type_validator.go).
//
// The compiler is new Compiler() + initOptions(new CompilerOptions()) whose ErrorHandler
// (CompilerOptions.setErrorHandler) records every diagnostic Compiler.report lets through. On the
// compiler's own JSTypeRegistry it builds the type universe of ReverseInterpreterDump (every
// JSTypeNative, a few user classes, interfaces, a templated class and enums, then the type
// expressions of <types.txt>; the code below is a copy of ReverseInterpreterDump.buildTypes) plus a
// few classes for the interface/abstract-method checks, and then, for a TypeValidator in
// SubtypingMode NORMAL and one in IGNORE_NULL_UNDEFINED:
//   U <method> <i> -> <result>        every unary expect* method on every type i
//   P <method> <i> <j> [<k>] -> <r>   the binary expect* methods on <pairs> seeded pairs/triples
//                                     (with a fresh validator, "reset", every 500 of them)
// and finally scenarios on fresh types/scopes (expectSuperType, which mutates its types, and
// expectUndeclaredVariable with a TypedScope, CompilerInput and @suppress {duplicate}).
// After each call it prints the diagnostics that were reported
//   D <level> <key> <description> @<node>
// and the TypeMismatches the validator registered
//   M <found> | <required> @<node>
// An exception is "X <class>: <message>" (a NullPointerException without its JVM message).
// Type indices refer to the "type" lines. Strings are escaped like ReverseInterpreterDump.
//
// internal/jscomp/type_validator_oracle_test.go mirrors this file step by step (it reuses the
// universe builder of type_reverse_interpreter_oracle_test.go).
//
// Usage: java -cp closure-compiler.jar:<build> com.google.javascript.jscomp.TypeValidatorDump \
//          <types.txt> <out.gz> <pairs>
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.parsing.JsDocInfoParser;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.Token;
import com.google.javascript.rhino.jstype.EnumType;
import com.google.javascript.rhino.jstype.FunctionType;
import com.google.javascript.rhino.jstype.JSType;
import com.google.javascript.rhino.jstype.JSType.SubtypingMode;
import com.google.javascript.rhino.jstype.JSTypeNative;
import com.google.javascript.rhino.jstype.JSTypeRegistry;
import com.google.javascript.rhino.jstype.JSTypeResolver;
import com.google.javascript.rhino.jstype.ObjectType;
import com.google.javascript.rhino.jstype.Property;
import com.google.javascript.rhino.jstype.StaticTypedScope;
import com.google.javascript.rhino.jstype.StaticTypedSlot;
import com.google.javascript.rhino.jstype.TemplateType;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.zip.GZIPOutputStream;

public final class TypeValidatorDump {
  private final StringBuilder out = new StringBuilder();
  private final List<String> labels = new ArrayList<>();
  private final List<JSType> types = new ArrayList<>();
  private final List<JSType> ctors = new ArrayList<>();
  private final Compiler compiler;
  private final JSTypeRegistry registry;
  private final TestScope scope = new TestScope();
  private final List<String> reported = new ArrayList<>();
  private TypeValidator validator;
  private int mismatchCount;

  // ---------------------------------------------------------------------
  // Output helpers

  /** ReverseInterpreterDump.esc (backslash, \\n, \\t, and \\uXXXX outside printable ASCII). */
  static String esc(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c == '\n') {
        sb.append("\\n");
      } else if (c == '\t') {
        sb.append("\\t");
      } else if (c < 0x20 || c >= 0x7f) {
        sb.append(String.format("\\u%04x", (int) c));
      } else {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  static String str(JSType t) {
    return t == null ? "null" : esc(t.toString());
  }

  void line(String s) {
    out.append(s).append('\n');
  }

  // The same LCG as TypeAlgebraDump / ReverseInterpreterDump.
  private long seed = 0x5DEECE66DL;

  int nextInt(int n) {
    seed = seed * 6364136223846793005L + 1442695040888963407L;
    return (int) ((seed >>> 33) % n);
  }

  static final class TestScope implements StaticTypedScope {
    final Node root = new Node(Token.ROOT);

    @Override
    public Node getRootNode() {
      return root;
    }

    @Override
    public StaticTypedScope getParentScope() {
      return null;
    }

    @Override
    public StaticTypedSlot getSlot(String name) {
      return null;
    }

    @Override
    public StaticTypedSlot getOwnSlot(String name) {
      return null;
    }

    @Override
    public JSType getTypeOfThis() {
      return null;
    }
  }

  TypeValidatorDump() {
    CompilerOptions options = new CompilerOptions();
    options.setErrorHandler(
        (level, error) ->
            reported.add(
                "  D "
                    + level
                    + " "
                    + error.getType().key
                    + " "
                    + esc(error.getDescription())
                    + " @"
                    + nodeLabel(error.getNode())));
    compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    compiler.initOptions(options);
    registry = compiler.getTypeRegistry();
  }

  static String nodeLabel(Node n) {
    if (n == null) {
      return "null";
    }
    if (n.isName() || n.isStringLit()) {
      return n.getToken() + "(" + esc(n.getString()) + ")";
    }
    return n.getToken().toString();
  }

  // ---------------------------------------------------------------------
  // The type universe (a copy of ReverseInterpreterDump.buildTypes)

  JSType nat(JSTypeNative n) {
    return registry.getNativeType(n);
  }

  void add(String label, JSType t) {
    labels.add(label);
    types.add(t);
    line("type\t" + types.size() + "\t" + esc(label) + "\t" + str(t));
  }

  static Node fn(String name) {
    return IR.function(IR.name(name), IR.paramList(), IR.block());
  }

  void build(List<String> typeStrings) {
    try (JSTypeResolver.Closer closer = registry.getResolver().openForDefinition()) {
      buildTypes(typeStrings);
    }
    try (JSTypeResolver.Closer closer = registry.getResolver().openForDefinition()) {
      buildExtraTypes();
    }
  }

  void buildTypes(List<String> typeStrings) {
    JSType number = nat(JSTypeNative.NUMBER_TYPE);
    JSType string = nat(JSTypeNative.STRING_TYPE);
    JSType bool = nat(JSTypeNative.BOOLEAN_TYPE);

    FunctionType fooCtor =
        registry.createConstructorType(
            "Foo", fn("Foo"), registry.createParameters(number), null, null, false);
    ObjectType foo = fooCtor.getInstanceType();
    fooCtor.getPrototype().defineDeclaredProperty("bar", number, null);
    foo.defineDeclaredProperty("p", string, null);
    registry.declareType(scope, "Foo", foo);

    FunctionType subFooCtor =
        registry.createConstructorType(
            "SubFoo", fn("SubFoo"), registry.createParameters(), null, null, false);
    subFooCtor.setPrototypeBasedOn(foo);
    ObjectType subFoo = subFooCtor.getInstanceType();
    subFoo.defineDeclaredProperty("qux", bool, null);
    registry.declareType(scope, "SubFoo", subFoo);

    FunctionType iCtor =
        FunctionType.builder(registry).forInterface().withName("I").withSourceNode(fn("I")).build();
    iCtor.getPrototype().defineDeclaredProperty("bar", number, null);
    ObjectType iface = iCtor.getInstanceType();
    registry.declareType(scope, "I", iface);

    FunctionType implCtor =
        registry.createConstructorType(
            "Impl", fn("Impl"), registry.createParameters(), null, null, false);
    implCtor.setImplementedInterfaces(List.of(iface));
    implCtor.getPrototype().defineDeclaredProperty("bar", number, null);
    ObjectType impl = implCtor.getInstanceType();
    registry.declareType(scope, "Impl", impl);

    FunctionType recCtor =
        FunctionType.builder(registry).forInterface().withName("Rec").withSourceNode(fn("Rec")).build();
    recCtor.setImplicitMatch(true);
    recCtor.getPrototype().defineDeclaredProperty("bar", number, null);
    ObjectType rec = recCtor.getInstanceType();
    registry.declareType(scope, "Rec", rec);

    TemplateType tKey = registry.createTemplateType("T");
    FunctionType boxCtor =
        FunctionType.builder(registry)
            .forConstructor()
            .withName("Box")
            .withSourceNode(fn("Box"))
            .withParameters(registry.createParameters(tKey))
            .withTemplateKeys(tKey)
            .build();
    boxCtor.getPrototype().defineDeclaredProperty("value", tKey, null);
    ObjectType box = boxCtor.getInstanceType();
    registry.declareType(scope, "Box", box);

    EnumType numEnum = EnumType.builder(registry).setName("E").setElementType(number).build();
    numEnum.defineElement("A", null);
    registry.declareType(scope, "E", numEnum.getElementsType());
    EnumType strEnum = EnumType.builder(registry).setName("S").setElementType(string).build();
    strEnum.defineElement("X", null);
    registry.declareType(scope, "S", strEnum.getElementsType());
    EnumType unionEnum =
        EnumType.builder(registry)
            .setName("U")
            .setElementType(registry.createUnionType(string, number))
            .build();
    unionEnum.defineElement("Y", null);
    registry.declareType(scope, "U", unionEnum.getElementsType());
    EnumType objEnum = EnumType.builder(registry).setName("O").setElementType(foo).build();
    objEnum.defineElement("Z", null);
    registry.declareType(scope, "O", objEnum.getElementsType());

    for (JSTypeNative n : JSTypeNative.values()) {
      add("native:" + n.name(), nat(n));
    }
    add("Foo", foo);
    add("ctor:Foo", fooCtor);
    add("Foo.prototype", fooCtor.getPrototype());
    add("SubFoo", subFoo);
    add("ctor:SubFoo", subFooCtor);
    add("I", iface);
    add("ctor:I", iCtor);
    add("Impl", impl);
    add("Rec", rec);
    add("Box", box);
    add("ctor:Box", boxCtor);
    add("T", tKey);
    add("enum:E", numEnum);
    add("E", numEnum.getElementsType());
    add("enum:S", strEnum);
    add("S", strEnum.getElementsType());
    add("U", unionEnum.getElementsType());
    add("O", objEnum.getElementsType());

    // The instanceof targets.
    ctors.add(fooCtor);
    ctors.add(subFooCtor);
    ctors.add(iCtor);
    ctors.add(implCtor);
    ctors.add(boxCtor);
    ctors.add(nat(JSTypeNative.ARRAY_FUNCTION_TYPE));
    ctors.add(nat(JSTypeNative.OBJECT_FUNCTION_TYPE));
    ctors.add(nat(JSTypeNative.STRING_OBJECT_FUNCTION_TYPE));
    ctors.add(nat(JSTypeNative.FUNCTION_FUNCTION_TYPE));
    ctors.add(nat(JSTypeNative.DATE_FUNCTION_TYPE));
    ctors.add(nat(JSTypeNative.REGEXP_FUNCTION_TYPE));
    ctors.add(nat(JSTypeNative.PROMISE_FUNCTION_TYPE));
    ctors.add(nat(JSTypeNative.I_THENABLE_FUNCTION_TYPE));
    ctors.add(nat(JSTypeNative.LEAST_FUNCTION_TYPE));
    ctors.add(nat(JSTypeNative.FUNCTION_TYPE));
    ctors.add(null);
    ctors.add(number);
    ctors.add(nat(JSTypeNative.UNKNOWN_TYPE));

    StaticTypedScope templateScope = registry.createScopeWithTemplates(scope, List.of(tKey));
    for (String s : typeStrings) {
      Node n = JsDocInfoParser.parseTypeString(s);
      if (n == null) {
        line("unparsed\t" + esc(s));
        continue;
      }
      JSType t;
      try {
        t = registry.createTypeFromCommentNode(n, "src.js", templateScope);
      } catch (RuntimeException e) {
        line("typeexception\t" + esc(s) + "\t" + e.getClass().getSimpleName());
        continue;
      }
      add("expr:" + s, t);
    }
  }

  /** Classes for the interface-property, abstract-method and struct checks. */
  void buildExtraTypes() {
    JSType number = nat(JSTypeNative.NUMBER_TYPE);
    JSType string = nat(JSTypeNative.STRING_TYPE);
    JSType numberOrUndefined = registry.createUnionType(number, nat(JSTypeNative.VOID_TYPE));

    FunctionType i2Ctor =
        FunctionType.builder(registry).forInterface().withName("I2").withSourceNode(fn("I2")).build();
    i2Ctor.getPrototype().defineDeclaredProperty("m", number, null);
    i2Ctor.getPrototype().defineDeclaredProperty("opt", numberOrUndefined, null);
    ObjectType i2 = i2Ctor.getInstanceType();
    i2.defineDeclaredProperty("own", string, null);

    FunctionType badImplCtor =
        registry.createConstructorType(
            "BadImpl", fn("BadImpl"), registry.createParameters(), null, null, false);
    badImplCtor.setImplementedInterfaces(List.of(i2));

    FunctionType wrongImplCtor =
        registry.createConstructorType(
            "WrongImpl", fn("WrongImpl"), registry.createParameters(), null, null, false);
    wrongImplCtor.setImplementedInterfaces(List.of(i2));
    wrongImplCtor.getPrototype().defineDeclaredProperty("m", string, null);
    wrongImplCtor.getInstanceType().defineDeclaredProperty("own", number, null);

    FunctionType goodImplCtor =
        registry.createConstructorType(
            "GoodImpl", fn("GoodImpl"), registry.createParameters(), null, null, false);
    goodImplCtor.setImplementedInterfaces(List.of(i2));
    goodImplCtor.getPrototype().defineDeclaredProperty("m", number, null);
    goodImplCtor.getInstanceType().defineDeclaredProperty("own", string, null);

    FunctionType i3Ctor =
        FunctionType.builder(registry).forInterface().withName("I3").withSourceNode(fn("I3")).build();
    i3Ctor.setExtendedInterfaces(List.of(i2));
    i3Ctor.getPrototype().defineDeclaredProperty("n", string, null);

    FunctionType abstractMethod =
        FunctionType.builder(registry).withReturnType(number).withIsAbstract(true).build();
    FunctionType concreteMethod = FunctionType.builder(registry).withReturnType(number).build();
    FunctionType absCtor =
        FunctionType.builder(registry)
            .forConstructor()
            .withName("Abs")
            .withSourceNode(fn("Abs"))
            .withIsAbstract(true)
            .build();
    absCtor.getPrototype().defineDeclaredProperty("am", abstractMethod, null);
    absCtor.getPrototype().defineDeclaredProperty("cm", concreteMethod, null);
    FunctionType abs2Ctor =
        FunctionType.builder(registry)
            .forConstructor()
            .withName("Abs2")
            .withSourceNode(fn("Abs2"))
            .withIsAbstract(true)
            .build();
    abs2Ctor.setPrototypeBasedOn(absCtor.getInstanceType());
    abs2Ctor.getPrototype().defineDeclaredProperty("am2", abstractMethod, null);
    FunctionType concreteCtor =
        registry.createConstructorType(
            "Concrete", fn("Concrete"), registry.createParameters(), null, null, false);
    concreteCtor.setPrototypeBasedOn(abs2Ctor.getInstanceType());
    concreteCtor.getPrototype().defineDeclaredProperty("am2", concreteMethod, null);
    FunctionType concrete2Ctor =
        registry.createConstructorType(
            "Concrete2", fn("Concrete2"), registry.createParameters(), null, null, false);
    concrete2Ctor.setPrototypeBasedOn(absCtor.getInstanceType());
    concrete2Ctor.getPrototype().defineDeclaredProperty("am", concreteMethod, null);

    FunctionType structCtor =
        registry.createConstructorType(
            "Struct", fn("Struct"), registry.createParameters(), null, null, false);
    structCtor.setStruct();
    structCtor.getPrototype().defineDeclaredProperty("s", number, null);

    add("I2", i2);
    add("ctor:I2", i2Ctor);
    add("ctor:BadImpl", badImplCtor);
    add("BadImpl", badImplCtor.getInstanceType());
    add("ctor:WrongImpl", wrongImplCtor);
    add("WrongImpl", wrongImplCtor.getInstanceType());
    add("ctor:GoodImpl", goodImplCtor);
    add("GoodImpl", goodImplCtor.getInstanceType());
    add("ctor:I3", i3Ctor);
    add("I3", i3Ctor.getInstanceType());
    add("abstractMethod", abstractMethod);
    add("ctor:Abs", absCtor);
    add("Abs", absCtor.getInstanceType());
    add("ctor:Abs2", abs2Ctor);
    add("ctor:Concrete", concreteCtor);
    add("Concrete", concreteCtor.getInstanceType());
    add("ctor:Concrete2", concrete2Ctor);
    add("ctor:Struct", structCtor);
    add("Struct", structCtor.getInstanceType());
  }

  // ---------------------------------------------------------------------
  // Calls

  interface Call {
    String run();
  }

  void call(String desc, Call c) {
    String result;
    try {
      result = c.run();
    } catch (NullPointerException e) {
      result = "X java.lang.NullPointerException";
    } catch (RuntimeException e) {
      result = "X " + e.getClass().getName() + ": " + esc(String.valueOf(e.getMessage()));
    }
    line(desc + " -> " + result);
    for (String r : reported) {
      line(r);
    }
    reported.clear();
    List<TypeMismatch> all = new ArrayList<>();
    for (TypeMismatch m : validator.getMismatches()) {
      all.add(m);
    }
    for (int i = mismatchCount; i < all.size(); i++) {
      TypeMismatch m = all.get(i);
      line(
          "  M " + str(m.found()) + " | " + str(m.required()) + " @" + nodeLabel(m.location()));
    }
    mismatchCount = all.size();
  }

  static String b(boolean v) {
    return v ? "true" : "false";
  }

  static Node typed(Node n, JSType t) {
    if (t != null) {
      n.setJSType(t);
    }
    return n;
  }

  static final String[] PROPS = {"p", "bar", "m", "own", "length", "value"};

  void unary(int i) {
    JSType t = types.get(i);
    String idx = " " + (i + 1);
    call("U expectObject" + idx, () -> b(validator.expectObject(IR.name("n"), t, "msg")));
    call(
        "U expectActualObject" + idx,
        () -> {
          validator.expectActualObject(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectAnyObject" + idx,
        () -> {
          validator.expectAnyObject(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectAutoboxesToIterable" + idx,
        () -> b(validator.expectAutoboxesToIterable(IR.name("n"), t, "msg")));
    call(
        "U expectAutoboxesToIterableOrAsyncIterable" + idx,
        () -> {
          Optional<JSType> r =
              validator.expectAutoboxesToIterableOrAsyncIterable(IR.name("n"), t, "msg");
          return r.isPresent() ? str(r.get()) : "empty";
        });
    call(
        "U expectGeneratorSupertype" + idx,
        () -> {
          validator.expectGeneratorSupertype(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectAsyncGeneratorSupertype" + idx,
        () -> {
          validator.expectAsyncGeneratorSupertype(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectValidAsyncReturnType" + idx,
        () -> {
          validator.expectValidAsyncReturnType(IR.name("n"), t);
          return "";
        });
    call(
        "U expectITemplateArraySupertype" + idx,
        () -> {
          validator.expectITemplateArraySupertype(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectString" + idx,
        () -> {
          validator.expectString(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectNumber" + idx,
        () -> {
          validator.expectNumber(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectNumberStrict" + idx,
        () -> {
          validator.expectNumberStrict(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectBitwiseable" + idx,
        () -> {
          validator.expectBitwiseable(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectNumberOrSymbol" + idx,
        () -> {
          validator.expectNumberOrSymbol(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectStringOrSymbol" + idx,
        () -> {
          validator.expectStringOrSymbol(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectUnknownOrComparable" + idx,
        () -> {
          validator.expectUnknownOrComparable(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectStringOrNumberOrSymbol" + idx,
        () -> {
          validator.expectStringOrNumberOrSymbol(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectStringOrNumberOrSymbolStrict" + idx,
        () -> {
          validator.expectStringOrNumberOrSymbolStrict(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectBigIntOrNumber" + idx,
        () -> {
          validator.expectBigIntOrNumber(IR.name("n"), t, "msg");
          return "";
        });
    call(
        "U expectNotNullOrUndefined" + idx,
        () ->
            b(
                validator.expectNotNullOrUndefined(
                    null, IR.name("n"), t, "msg", nat(JSTypeNative.OBJECT_TYPE))));
    call(
        "U expectWellFormedTemplatizedType" + idx,
        () -> {
          validator.expectWellFormedTemplatizedType(typed(IR.name("n"), t));
          return "";
        });
    call("U getInstanceOfCtor" + idx, () -> str(TypeValidator.getInstanceOfCtor(t)));
    FunctionType fnType = t.toMaybeFunctionType();
    if (fnType != null && (fnType.isConstructor() || fnType.isInterface())) {
      call(
          "U expectAllInterfaceProperties" + idx,
          () -> {
            validator.expectAllInterfaceProperties(IR.name("n"), fnType);
            return "";
          });
    }
    if (fnType != null && fnType.isConstructor()) {
      call(
          "U expectAbstractMethodsImplemented" + idx,
          () -> {
            validator.expectAbstractMethodsImplemented(IR.name("n"), fnType);
            return "";
          });
    }
  }

  void binary() {
    int n = types.size();
    int i = nextInt(n);
    int j = nextInt(n);
    int k = nextInt(n);
    JSType a = types.get(i);
    JSType bt = types.get(j);
    JSType c = types.get(k);
    String ij = " " + (i + 1) + " " + (j + 1);
    String ijk = ij + " " + (k + 1);
    switch (nextInt(11)) {
      case 0 ->
          call(
              "P expectMatchingTypesStrict" + ij,
              () -> {
                validator.expectMatchingTypesStrict(IR.name("n"), a, bt, "msg");
                return "";
              });
      case 1 ->
          call(
              "P expectSwitchMatchesCase" + ij,
              () -> {
                validator.expectSwitchMatchesCase(
                    new Node(Token.CASE, IR.name("c"), IR.block()), a, bt);
                return "";
              });
      case 2 ->
          call(
              "P expectIndexMatch" + ij,
              () -> {
                validator.expectIndexMatch(IR.getelem(IR.name("o"), IR.name("i")), a, bt);
                return "";
              });
      case 3 ->
          call(
              "P expectCanAssignTo" + ij,
              () -> b(validator.expectCanAssignTo(IR.name("n"), a, bt, "msg")));
      case 4 -> {
        String prop = PROPS[nextInt(PROPS.length)];
        call(
            "P expectCanAssignToPropertyOf" + ijk + " " + prop,
            () ->
                b(
                    validator.expectCanAssignToPropertyOf(
                        IR.name("n"), a, bt, typed(IR.name("owner"), c), prop)));
      }
      case 5 -> {
        int ordinal = 1 + nextInt(3);
        call(
            "P expectArgumentMatchesParameter" + ijk + " " + ordinal,
            () -> {
              validator.expectArgumentMatchesParameter(
                  IR.name("n"), a, bt, IR.call(typed(IR.name("f"), c)), ordinal);
              return "";
            });
      }
      case 6 ->
          call(
              "P expectCanCast" + ij,
              () -> {
                validator.expectCanCast(IR.name("n"), a, bt);
                return "";
              });
      case 7 ->
          call(
              "P expectCanAssignToPrototype" + ij,
              () -> {
                validator.expectCanAssignToPrototype(a, IR.name("n"), bt);
                return "";
              });
      case 8 -> {
        FunctionType sub = a.toMaybeFunctionType();
        FunctionType sup = bt.toMaybeFunctionType();
        if (sub != null) {
          call(
              "P expectExtends" + ij,
              () -> {
                validator.expectExtends(IR.name("n"), sub, sup);
                return "";
              });
        }
      }
      case 9 -> {
        ObjectType parent = bt.toMaybeObjectType();
        String prop = PROPS[nextInt(PROPS.length)];
        if (parent != null) {
          call(
              "P checkPropertyType" + ijk + " " + prop,
              () -> {
                validator.checkPropertyType(
                    IR.name("n"), a, parent, new Property.StringKey(prop), c);
                return "";
              });
        }
      }
      default ->
          call(
              "P expectNotNullOrUndefined" + ij,
              () -> b(validator.expectNotNullOrUndefined(null, IR.name("n"), a, "msg", bt)));
    }
  }

  void sweep(SubtypingMode mode, int pairs) {
    line("mode " + mode);
    validator = new TypeValidator(compiler);
    validator.setSubtypingMode(mode);
    mismatchCount = 0;
    call(
        "U expectValidTypeofName",
        () -> {
          validator.expectValidTypeofName(IR.name("n"), "foo");
          return "";
        });
    for (int i = 0; i < types.size(); i++) {
      unary(i);
    }
    for (int p = 0; p < pairs; p++) {
      if (p % 500 == 0) {
        // A fresh validator every 500 pairs keeps getMismatches() (a copy) short.
        line("reset");
        validator = new TypeValidator(compiler);
        validator.setSubtypingMode(mode);
        mismatchCount = 0;
      }
      binary();
    }
  }

  // ---------------------------------------------------------------------
  // Scenarios on fresh types

  void scenarios() {
    line("scenarios");
    validator = new TypeValidator(compiler);
    mismatchCount = 0;
    JSType number = nat(JSTypeNative.NUMBER_TYPE);
    JSType string = nat(JSTypeNative.STRING_TYPE);
    try (JSTypeResolver.Closer closer = registry.getResolver().openForDefinition()) {
      FunctionType aCtor =
          registry.createConstructorType(
              "SA", fn("SA"), registry.createParameters(), null, null, false);
      FunctionType bCtor =
          registry.createConstructorType(
              "SB", fn("SB"), registry.createParameters(), null, null, false);
      FunctionType cCtor =
          registry.createConstructorType(
              "SC", fn("SC"), registry.createParameters(), null, null, false);
      cCtor.setPrototypeBasedOn(aCtor.getInstanceType());
      FunctionType dCtor =
          registry.createConstructorType(
              "SD", fn("SD"), registry.createParameters(), null, null, false);
      FunctionType boxCtor = types.get(labels.indexOf("ctor:Box")).toMaybeFunctionType();
      FunctionType eCtor =
          registry.createConstructorType(
              "SE", fn("SE"), registry.createParameters(), null, null, false);
      eCtor.setPrototypeBasedOn(
          registry.createTemplatizedType(boxCtor.getInstanceType(), number));
      FunctionType fCtor =
          registry.createConstructorType(
              "SF", fn("SF"), registry.createParameters(), null, null, false);
      fCtor.setPrototypeBasedOn(aCtor.getInstanceType());
      call(
          "S expectSuperType SB SC",
          () -> {
            validator.expectSuperType(IR.name("n"), bCtor.getInstanceType(), cCtor.getInstanceType());
            return str(cCtor.getInstanceType().getImplicitPrototype().getImplicitPrototype());
          });
      call(
          "S expectSuperType SA SC",
          () -> {
            validator.expectSuperType(IR.name("n"), aCtor.getInstanceType(), cCtor.getInstanceType());
            return str(cCtor.getInstanceType().getImplicitPrototype().getImplicitPrototype());
          });
      call(
          "S expectSuperType SA SD",
          () -> {
            validator.expectSuperType(IR.name("n"), aCtor.getInstanceType(), dCtor.getInstanceType());
            return str(dCtor.getInstanceType().getImplicitPrototype().getImplicitPrototype());
          });
      call(
          "S expectSuperType Box SE",
          () -> {
            validator.expectSuperType(
                IR.name("n"), boxCtor.getInstanceType(), eCtor.getInstanceType());
            return str(eCtor.getInstanceType().getImplicitPrototype().getImplicitPrototype());
          });
      call(
          "S expectSuperType ? SF",
          () -> {
            validator.expectSuperType(
                IR.name("n"),
                (ObjectType) nat(JSTypeNative.UNKNOWN_TYPE),
                fCtor.getInstanceType());
            return str(fCtor.getInstanceType().getImplicitPrototype().getImplicitPrototype());
          });
      fCtor.getInstanceType().getPropertyType("toString"); // caches values
      call(
          "S expectSuperType SB SF (cached)",
          () -> {
            validator.expectSuperType(IR.name("n"), bCtor.getInstanceType(), fCtor.getInstanceType());
            return str(fCtor.getInstanceType().getImplicitPrototype().getImplicitPrototype());
          });
    }

    // expectUndeclaredVariable
    Node script = IR.script();
    Node root = IR.root(script);
    TypedScope global = TypedScope.createGlobalScope(root);
    CompilerInput input1 = new CompilerInput(SourceFile.fromCode("first.js", ""));
    CompilerInput input2 = new CompilerInput(SourceFile.fromCode("second.js", ""));
    String[] names = {"a", "b", "c", "d", "e", "f", "g"};
    for (int s = 0; s < names.length; s++) {
      String name = names[s];
      Node oldName = IR.name(name);
      Node oldDecl;
      if (s == 3) {
        oldDecl = IR.exprResult(oldName);
      } else {
        oldDecl = IR.var(oldName);
      }
      if (s == 5) {
        JSDocInfo.Builder info = JSDocInfo.builder();
        info.recordSuppressions(Set.of("duplicate"));
        oldDecl.setJSDocInfo(info.build());
      }
      script.addChildToBack(oldDecl);
      oldName.setLinenoCharno(10 + s, 4);
      JSType oldType = s == 4 ? nat(JSTypeNative.UNKNOWN_TYPE) : number;
      TypedVar var = global.declare(name, oldName, oldType, s == 0 ? null : input1, false);
      Node newName = IR.name(name);
      Node newDecl;
      if (s == 6) {
        newName = IR.function(IR.name(name), IR.paramList(), IR.block());
        newDecl = IR.exprResult(newName);
      } else {
        newDecl = IR.var(newName);
        newName.addChildToBack(IR.number(1));
      }
      script.addChildToBack(newDecl);
      Node nn = newName;
      JSType newType = (s == 2 || s == 5) ? string : number;
      final int fs = s;
      call(
          "S expectUndeclaredVariable " + name,
          () -> {
            Node parent = fs == 6 ? nn : nn.getParent();
            Node target = fs == 6 ? nn.getFirstChild() : nn;
            TypedVar result =
                validator.expectUndeclaredVariable(
                    "second.js", input2, target, parent, var, name, newType);
            return (result == var ? "same" : "new")
                + " "
                + esc(result.getInputName())
                + " "
                + str(result.getType())
                + " "
                + str(target.getJSType())
                + " "
                + (target.hasChildren() ? str(target.getFirstChild().getJSType()) : "-")
                + " "
                + str(parent.getJSType());
          });
    }
  }

  public static void main(String[] args) throws Exception {
    List<String> typeStrings = new ArrayList<>();
    for (String l : Files.readAllLines(Paths.get(args[0]), StandardCharsets.UTF_8)) {
      if (!l.isEmpty() && !l.startsWith("#")) {
        typeStrings.add(l);
      }
    }
    int pairs = Integer.parseInt(args[2]);
    TypeValidatorDump d = new TypeValidatorDump();
    d.build(typeStrings);
    d.line("built");
    for (String r : d.reported) {
      d.line(r);
    }
    d.reported.clear();
    d.sweep(SubtypingMode.NORMAL, pairs);
    d.sweep(SubtypingMode.IGNORE_NULL_UNDEFINED, pairs);
    d.scenarios();
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(Paths.get(args[1])))) {
      os.write(d.out.toString().getBytes(StandardCharsets.UTF_8));
    }
  }
}

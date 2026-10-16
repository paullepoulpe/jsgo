// Type-algebra oracle for internal/rhino (rhino.jstype).
//
// Builds a large, deterministic set of types through the JSTypeRegistry API
// and prints, for each type, its string forms and the results of the unary
// type operations, then the binary operations (subtyping, equivalence,
// equality tests, least supertype, greatest subtype, ...) for every pair of
// types. internal/rhino/jstype_oracle_test.go builds the same types with the
// Go port, in the same order, and must reproduce this output byte for byte.
// Keep the two files in step: every API call made here is made there too, in
// the same order, because many operations mutate the types (lazy
// prototypes, call/apply/bind properties, caches).
//
// The class lives in the jstype package so that it can reach the
// package-private constructors (ProxyObjectType, NamedType.builder, ...).
//
// Usage: java -cp closure-compiler.jar:<build> \
//          com.google.javascript.rhino.jstype.TypeAlgebraDump [--pairs=all|<n>] > out.txt
package com.google.javascript.rhino.jstype;

import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.base.Tri;
import com.google.javascript.jscomp.parsing.JsDocInfoParser;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.Outcome;
import com.google.javascript.rhino.QualifiedName;
import com.google.javascript.rhino.StaticScope;
import com.google.javascript.rhino.StaticSlot;
import com.google.javascript.rhino.Token;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Supplier;

public final class TypeAlgebraDump {
  private final StringBuilder out = new StringBuilder();
  private final List<String> labels = new ArrayList<>();
  private final List<JSType> types = new ArrayList<>();
  private final JSTypeRegistry registry;
  private final TestScope scope;
  // Set by build(); used by the facts.
  private TemplateType tKey;
  private TemplateTypeMap partialMap;
  private TemplateTypeMap totalMap;
  private final Map<TemplateType, JSType> inferMap = new LinkedHashMap<>();

  // ---------------------------------------------------------------------
  // Output helpers

  static String esc(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c == '\n') {
        sb.append("\\n");
      } else if (c < 0x20 || c >= 0x7f) {
        sb.append(String.format("\\u%04x", (int) c));
      } else {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  static String str(Object o) {
    if (o == null) {
      return "null";
    }
    if (o instanceof Iterable<?> it) {
      StringBuilder sb = new StringBuilder("[");
      boolean first = true;
      for (Object e : it) {
        if (!first) {
          sb.append(", ");
        }
        first = false;
        sb.append(str(e));
      }
      return sb.append("]").toString();
    }
    if (o instanceof FunctionType.Parameter p) {
      return "P(" + str(p.getJSType()) + "," + p.isOptional() + "," + p.isVariadic() + ")";
    }
    if (o instanceof JSType.TypePair p) {
      return "<" + str(p.typeA) + " ; " + str(p.typeB) + ">";
    }
    return String.valueOf(o);
  }

  void line(String s) {
    out.append(s).append('\n');
  }

  void rec(String label, String value) {
    line(label + "\t" + esc(value));
  }

  /** Runs s, printing its result or the class of the exception it threw. */
  void fact(String name, Supplier<Object> s) {
    String v;
    try {
      v = str(s.get());
    } catch (RuntimeException e) {
      v = "!" + e.getClass().getSimpleName();
    }
    rec(name, v);
  }

  // ---------------------------------------------------------------------
  // Deterministic sequences

  private long seed = 0x5DEECE66DL;

  int nextInt(int n) {
    seed = seed * 6364136223846793005L + 1442695040888963407L;
    return (int) ((seed >>> 33) % n);
  }

  // ---------------------------------------------------------------------
  // A minimal scope and reporter

  static final class TestReporter implements ErrorReporter {
    final StringBuilder sb;

    TestReporter(StringBuilder sb) {
      this.sb = sb;
    }

    @Override
    public void warning(String message, String sourceName, int line, int lineOffset) {
      sb.append("WARNING\t").append(esc(message + " @" + sourceName + ":" + line + ":" + lineOffset)).append('\n');
    }

    @Override
    public void error(String message, String sourceName, int line, int lineOffset) {
      sb.append("ERROR\t").append(esc(message + " @" + sourceName + ":" + line + ":" + lineOffset)).append('\n');
    }
  }

  static final class TestScope implements StaticTypedScope {
    final Node root = new Node(Token.ROOT);
    final Map<String, SimpleSlot> slots = new LinkedHashMap<>();

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
      return slots.get(name);
    }

    @Override
    public StaticTypedSlot getOwnSlot(String name) {
      return slots.get(name);
    }

    @Override
    public JSType getTypeOfThis() {
      return null;
    }
  }

  TypeAlgebraDump() {
    this.registry = new JSTypeRegistry(new TestReporter(out), java.util.Set.of("Forwarded"));
    this.scope = new TestScope();
  }

  JSType nat(JSTypeNative n) {
    return registry.getNativeType(n);
  }

  ObjectType natObj(JSTypeNative n) {
    return registry.getNativeObjectType(n);
  }

  <T extends JSType> T add(String label, T t) {
    labels.add(label);
    types.add(t);
    line("add\t" + types.size() + "\t" + label);
    return t;
  }

  static Node fn(String name) {
    return IR.function(IR.name(name), IR.paramList(), IR.block());
  }

  /** Serializes a type AST: (TOKEN[ "string"] line:char child...). */
  static String sexp(Node n) {
    StringBuilder sb = new StringBuilder();
    sexp(n, sb);
    return sb.toString();
  }

  static void sexp(Node n, StringBuilder sb) {
    sb.append('(').append(n.getToken().name());
    String s = null;
    try {
      s = n.getString();
    } catch (ClassCastException e) {
      // not a string node
    }
    if (s != null) {
      sb.append(" \"").append(esc(s).replace("\"", "\\\"")).append('"');
    }
    sb.append(' ').append(n.getLineno()).append(':').append(n.getCharno());
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      sb.append(' ');
      sexp(c, sb);
    }
    sb.append(')');
  }

  JSType parsed(String typeString, StaticTypedScope s) {
    Node n = JsDocInfoParser.parseTypeString(typeString);
    line("ast\t" + esc(typeString) + "\t" + esc(sexp(n)));
    return registry.createTypeFromCommentNode(n, "src.js", s);
  }

  // ---------------------------------------------------------------------
  // The type universe

  void build() {
    // Like BaseJSTypeTestCase: define everything with the resolver open, so
    // that types are resolved together when it closes.
    try (JSTypeResolver.Closer closer = registry.getResolver().openForDefinition()) {
      buildTypes();
    }
  }

  void buildTypes() {
    for (JSTypeNative n : JSTypeNative.values()) {
      add("native:" + n.name(), nat(n));
    }
    JSType number = nat(JSTypeNative.NUMBER_TYPE);
    JSType string = nat(JSTypeNative.STRING_TYPE);
    JSType bool = nat(JSTypeNative.BOOLEAN_TYPE);
    JSType nul = nat(JSTypeNative.NULL_TYPE);
    JSType undef = nat(JSTypeNative.VOID_TYPE);
    JSType unknown = nat(JSTypeNative.UNKNOWN_TYPE);
    JSType all = nat(JSTypeNative.ALL_TYPE);
    JSType symbol = nat(JSTypeNative.SYMBOL_TYPE);
    ObjectType array = natObj(JSTypeNative.ARRAY_TYPE);
    ObjectType object = natObj(JSTypeNative.OBJECT_TYPE);
    ObjectType readonlyArray = natObj(JSTypeNative.READONLY_ARRAY_TYPE);
    ObjectType iterable = natObj(JSTypeNative.ITERABLE_TYPE);
    ObjectType promise = natObj(JSTypeNative.PROMISE_TYPE);
    ObjectType ithenable = natObj(JSTypeNative.I_THENABLE_TYPE);
    ObjectType iobject = natObj(JSTypeNative.I_OBJECT_TYPE);
    ObjectType map = natObj(JSTypeNative.MAP_TYPE);

    // Classes, interfaces and inheritance.
    FunctionType fooCtor =
        registry.createConstructorType(
            "Foo", fn("Foo"), registry.createParameters(number), null, null, false);
    ObjectType foo = fooCtor.getInstanceType();
    fooCtor.getPrototype().defineDeclaredProperty("bar", number, null);
    foo.defineDeclaredProperty("baz", string, null);
    fooCtor.defineDeclaredProperty("staticProp", bool, null);
    registry.declareType(scope, "Foo", foo);
    add("Foo", foo);
    add("ctor:Foo", fooCtor);
    add("Foo.prototype", fooCtor.getPrototype());

    FunctionType subFooCtor =
        registry.createConstructorType(
            "SubFoo", fn("SubFoo"), registry.createParameters(), null, null, false);
    subFooCtor.setPrototypeBasedOn(foo);
    ObjectType subFoo = subFooCtor.getInstanceType();
    subFoo.defineDeclaredProperty("qux", bool, null);
    registry.declareType(scope, "SubFoo", subFoo);
    add("SubFoo", subFoo);
    add("ctor:SubFoo", subFooCtor);
    add("SubFoo.prototype", subFooCtor.getPrototype());

    FunctionType iCtor =
        FunctionType.builder(registry).forInterface().withName("I").withSourceNode(fn("I")).build();
    iCtor.getPrototype().defineDeclaredProperty("bar", number, null);
    ObjectType iface = iCtor.getInstanceType();
    registry.declareType(scope, "I", iface);
    add("I", iface);
    add("ctor:I", iCtor);

    FunctionType jCtor =
        FunctionType.builder(registry).forInterface().withName("J").withSourceNode(fn("J")).build();
    jCtor.setExtendedInterfaces(ImmutableList.of(iface));
    jCtor.getPrototype().defineDeclaredProperty("jay", string, null);
    ObjectType jface = jCtor.getInstanceType();
    add("J", jface);
    add("ctor:J", jCtor);

    FunctionType implCtor =
        registry.createConstructorType(
            "Impl", fn("Impl"), registry.createParameters(), null, null, false);
    implCtor.setImplementedInterfaces(ImmutableList.of(jface));
    implCtor.getPrototype().defineDeclaredProperty("bar", number, null);
    ObjectType impl = implCtor.getInstanceType();
    registry.declareType(scope, "Impl", impl);
    add("Impl", impl);

    FunctionType recCtor =
        FunctionType.builder(registry).forInterface().withName("Rec").withSourceNode(fn("Rec")).build();
    recCtor.setImplicitMatch(true);
    recCtor.getPrototype().defineDeclaredProperty("bar", number, null);
    recCtor.getPrototype().defineDeclaredProperty("opt", registry.createOptionalType(string), null);
    ObjectType rec = recCtor.getInstanceType();
    registry.declareType(scope, "Rec", rec);
    add("Rec", rec);
    add("ctor:Rec", recCtor);

    TemplateType tKey = registry.createTemplateType("T");
    this.tKey = tKey;
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
    add("Box", box);
    add("ctor:Box", boxCtor);
    add("T", tKey);
    TemplatizedType boxString = add("Box<string>", registry.createTemplatizedType(box, string));
    TemplatizedType boxNumber = add("Box<number>", registry.createTemplatizedType(box, number));
    add("Box<?>", registry.createTemplatizedType(box, unknown));
    add("Box<string>.value", boxString.getPropertyType("value"));

    FunctionType structCtor =
        registry.createConstructorType(
            "Struct", fn("Struct"), registry.createParameters(), null, null, false);
    structCtor.setStruct();
    structCtor.getPrototype().defineDeclaredProperty("s", number, null);
    add("Struct", structCtor.getInstanceType());
    FunctionType dictCtor =
        registry.createConstructorType(
            "Dict", fn("Dict"), registry.createParameters(), null, null, false);
    dictCtor.setDict();
    add("Dict", dictCtor.getInstanceType());
    FunctionType abstractCtor =
        registry.createConstructorType(
            "Abs", fn("Abs"), registry.createParameters(), null, null, true);
    abstractCtor.setImplementedInterfaces(ImmutableList.of(iface));
    add("Abs", abstractCtor.getInstanceType());
    FunctionType subStructCtor =
        registry.createConstructorType(
            "SubStruct", fn("SubStruct"), registry.createParameters(), null, null, false);
    subStructCtor.setPrototypeBasedOn(structCtor.getInstanceType());
    add("SubStruct", subStructCtor.getInstanceType());
    FunctionType anonCtor =
        registry.createConstructorType("", null, registry.createParameters(), null, null, false);
    add("anonCtor", anonCtor);
    add("anonInstance", anonCtor.getInstanceType());
    FunctionType structuralCtor =
        FunctionType.builder(registry).forConstructor().withTypeOfThis(foo).build();
    add("structuralCtor", structuralCtor);

    ObjectType lit = registry.createAnonymousObjectType(null);
    lit.defineInferredProperty("a", number, null);
    lit.defineInferredProperty("b", string, null);
    add("objlit", lit);
    ObjectType emptyLit = registry.createAnonymousObjectType(null);
    add("emptyObjlit", emptyLit);
    ObjectType named = registry.createObjectType("Named", object);
    named.defineDeclaredProperty("n", number, null);
    add("Named", named);

    // Records.
    Map<String, JSType> rA = new LinkedHashMap<>();
    rA.put("a", number);
    JSType recA = add("{a:number}", registry.createRecordType(rA));
    Map<String, JSType> rAB = new LinkedHashMap<>();
    rAB.put("b", string);
    rAB.put("a", number);
    JSType recAB = add("{a:number,b:string}", registry.createRecordType(rAB));
    Map<String, JSType> rAopt = new LinkedHashMap<>();
    rAopt.put("a", registry.createOptionalType(number));
    add("{a:number=}", registry.createRecordType(rAopt));
    Map<String, JSType> rAS = new LinkedHashMap<>();
    rAS.put("a", string);
    JSType recAS = add("{a:string}", registry.createRecordType(rAS));
    Map<String, JSType> rBar = new LinkedHashMap<>();
    rBar.put("bar", number);
    add("{bar:number}", registry.createRecordType(rBar));
    Map<String, JSType> rNested = new LinkedHashMap<>();
    rNested.put("r", recA);
    rNested.put("f", registry.createFunctionType(number, string));
    add("{r:{a:number},f:...}", registry.createRecordType(rNested));
    Map<String, JSType> rMany = new LinkedHashMap<>();
    for (int i = 0; i < 12; i++) {
      rMany.put("p" + (char) ('a' + i), i % 2 == 0 ? number : string);
    }
    add("{many}", registry.createRecordType(rMany));
    RecordTypeBuilder synth = new RecordTypeBuilder(registry);
    synth.setSynthesized(true);
    synth.addProperty("x", bool, null);
    add("{x:boolean}(synthetic)", synth.build());
    add("emptyRecord", new RecordTypeBuilder(registry).build());
    add("buildRecordTypeFromObject(objlit)", registry.buildRecordTypeFromObject(lit));

    // Unions.
    add("number|string", registry.createUnionType(number, string));
    add("number|null", registry.createNullableType(number));
    add("Foo|null", registry.createNullableType(foo));
    add("Foo|SubFoo", registry.createUnionType(foo, subFoo));
    add("SubFoo|Foo", registry.createUnionType(subFoo, foo));
    add(
        "(number|string)|boolean",
        registry.createUnionType(registry.createUnionType(number, string), bool));
    add("?|number", registry.createUnionType(unknown, number));
    add("*|number", registry.createUnionType(all, number));
    add("?|undefined", registry.createUnionType(unknown, undef));
    add("number|undefined", registry.createOptionalType(number));
    add("Foo|null|undefined", registry.createOptionalNullableType(foo));
    add(
        "Array<string>|Array<number>",
        registry.createUnionType(
            registry.createTemplatizedType(array, string),
            registry.createTemplatizedType(array, number)));
    add(
        "ReadonlyArray<string>|ReadonlyArray<number>",
        registry.createUnionType(
            registry.createTemplatizedType(readonlyArray, string),
            registry.createTemplatizedType(readonlyArray, number)));
    add("Array|Array<string>", registry.createUnionType(array, registry.createTemplatizedType(array, string)));
    add("Object|Array<string>", registry.createUnionType(object, registry.createTemplatizedType(array, string)));
    add("{a:number}|{a:number,b:string}", registry.createUnionType(recA, recAB));
    add("{a:number}|{a:string}", registry.createUnionType(recA, recAS));
    add(
        "fn|fn",
        registry.createUnionType(
            registry.createFunctionType(number, string), registry.createFunctionType(string, string)));
    add(
        "fn(number)|fn(string)",
        registry.createUnionType(
            registry.createFunctionType(number, number), registry.createFunctionType(number, string)));
    add("null|undefined|number|string|boolean", registry.createUnionType(nul, undef, number, string, bool));
    add("symbol|number", registry.createUnionType(symbol, number));
    add("Foo|Impl|I", registry.createUnionType(foo, impl, iface));
    add("T|number", registry.createUnionType(tKey, number));
    add("Box<string>|Box<number>", registry.createUnionType(boxString, boxNumber));
    add("emptyUnion", UnionType.builder(registry).build());
    add("checked?|?", registry.createUnionType(nat(JSTypeNative.CHECKED_UNKNOWN_TYPE), unknown));
    add("checked?|number", registry.createUnionType(nat(JSTypeNative.CHECKED_UNKNOWN_TYPE), number));
    UnionType.Builder big = UnionType.builder(registry);
    for (int i = 0; i < 35; i++) {
      Map<String, JSType> m = new LinkedHashMap<>();
      m.put("k" + i, number);
      big.addAlternate(registry.createRecordType(m));
    }
    add("bigUnion", big.build());

    // Functions.
    add("function(number):string", registry.createFunctionType(string, number));
    add(
        "function(number=):string",
        registry.createFunctionType(string, registry.createOptionalParameters(number)));
    add("function(...number):undefined", registry.createFunctionTypeWithVarArgs(undef, number));
    add("function():?", registry.createFunctionType(unknown));
    FunctionParamBuilder pb = new FunctionParamBuilder(registry);
    pb.addRequiredParams(number);
    pb.addOptionalParams(string);
    pb.addVarArgs(bool);
    add(
        "function(number,string=,...boolean):Foo",
        FunctionType.builder(registry).withParameters(pb.build()).withReturnType(foo).build());
    add(
        "function(this:Foo):undefined",
        FunctionType.builder(registry).withParameters().withReturnType(undef).withTypeOfThis(foo).build());
    add(
        "function(new:Foo):Foo",
        FunctionType.builder(registry)
            .forConstructor()
            .withParameters()
            .withReturnType(foo)
            .withTypeOfThis(foo)
            .build());
    add("function(?):number", registry.createFunctionType(number, unknown));
    add("function(T):T", FunctionType.builder(registry)
        .withParameters(registry.createParameters(tKey))
        .withReturnType(tKey)
        .withTemplateKeys(tKey)
        .build());
    add(
        "function(number):string(inferred)",
        FunctionType.builder(registry)
            .withParameters(registry.createParameters(number))
            .withInferredReturnType(string)
            .build());
    add(
        "function(this:Foo, number):number",
        registry.createFunctionTypeWithInstanceType(foo, number, ImmutableList.of(number)));
    FunctionType withProps = registry.createFunctionType(number, string);
    withProps.defineDeclaredProperty("extra", bool, null);
    add("fnWithProps", withProps);
    add("bind(-1)", fooCtor.getBindReturnType(-1));
    add("bind(1)", registry.createFunctionType(number, string, bool).getBindReturnType(1));
    add("bind(2)", registry.createFunctionType(number, string, bool).getBindReturnType(2));
    add("forget(Foo)", fooCtor.forgetParameterAndReturnTypes());
    add(
        "newReturnType",
        registry.createFunctionTypeWithNewReturnType(registry.createFunctionType(number, string), bool));
    add("toBuilder(SubFoo)", subFooCtor.toBuilder().build());

    // Enums.
    EnumType eNum =
        EnumType.builder(registry).setName("E").setElementType(number).setSource(null).build();
    eNum.defineElement("A", null);
    eNum.defineElement("B", null);
    registry.declareType(scope, "E", eNum.getElementsType());
    registry.identifyNonNullableName(scope, "E");
    add("enum E", eNum);
    add("E", eNum.getElementsType());
    EnumType eStr = EnumType.builder(registry).setName("S").setElementType(string).build();
    eStr.defineElement("X", null);
    add("enum S", eStr);
    add("S", eStr.getElementsType());
    EnumType eFoo = EnumType.builder(registry).setName("EF").setElementType(foo).build();
    eFoo.defineElement("F", null);
    add("enum EF", eFoo);
    add("EF", eFoo.getElementsType());
    EnumType eUnion =
        EnumType.builder(registry)
            .setName("EU")
            .setElementType(registry.createUnionType(number, string))
            .build();
    add("EU", eUnion.getElementsType());
    add("E|S", registry.createUnionType(eNum.getElementsType(), eStr.getElementsType()));

    // Templatized and template types.
    add("Array<string>", registry.createTemplatizedType(array, string));
    add("Array<number>", registry.createTemplatizedType(array, number));
    add("Array<?>", registry.createTemplatizedType(array, unknown));
    add("Array<Foo>", registry.createTemplatizedType(array, foo));
    add("Array<SubFoo>", registry.createTemplatizedType(array, subFoo));
    add("Object<string,number>", registry.createTemplatizedType(object, string, number));
    add("IObject<number,string>", registry.createTemplatizedType(iobject, number, string));
    add("Iterable<string>", registry.createTemplatizedType(iterable, string));
    add("ReadonlyArray<string>", registry.createTemplatizedType(readonlyArray, string));
    add("ReadonlyArray<Foo>", registry.createTemplatizedType(readonlyArray, foo));
    add("Promise<number>", registry.createTemplatizedType(promise, number));
    add("IThenable<string>", registry.createTemplatizedType(ithenable, string));
    add("Map<string,number>", registry.createTemplatizedType(map, string, number));
    add("Array<Array<string>>", registry.createTemplatizedType(array, registry.createTemplatizedType(array, string)));
    TemplateType uKey = add("U extends number", registry.createTemplateType("U", number));
    add("Box<T>", registry.createTemplatizedType(box, tKey));
    Map<TemplateType, JSType> tmap = new LinkedHashMap<>();
    tmap.put(tKey, bool);
    add("Box<boolean>(map)", registry.createTemplatizedType(box, tmap));
    TemplateTypeReplacer replacer = TemplateTypeReplacer.forInference(registry, tmap);
    add("replace(function(T):T)", types.get(labels.indexOf("function(T):T")).visit(replacer));
    add("replace(Box<T>)", types.get(labels.indexOf("Box<T>")).visit(replacer));
    add("replace(T|number)", types.get(labels.indexOf("T|number")).visit(replacer));
    add("U|string", registry.createUnionType(uKey, string));

    // Proxies and named types.
    add("proxy(number)", new ProxyObjectType(registry, number));
    add("proxy(Foo)", new ProxyObjectType(registry, foo));
    add("proxy(number|string)", new ProxyObjectType(registry, registry.createUnionType(number, string)));
    add("named:Foo", registry.createNamedType(scope, "Foo", "src.js", 1, 2));
    add("named:Missing", registry.createNamedType(scope, "Missing", "src.js", 3, 4));
    add("named:Forwarded", registry.createNamedType(scope, "Forwarded", "src.js", 5, 6));
    add("named:Box", registry.createNamedType(scope, "Box", "src.js", 7, 8));
    add("noResolved", new NoResolvedType(registry, "NR", ImmutableList.of()));
    KnownSymbolType iter = add("knownSymbol", new KnownSymbolType(registry, "Symbol.iterator"));
    ObjectType withSym = registry.createAnonymousObjectType(null);
    withSym.defineDeclaredProperty(new Property.SymbolKey(iter), number, null);
    add("objlitWithSymbol", withSym);

    // More generics, structural interfaces and forward references.
    TemplateType kKey = registry.createTemplateType("K");
    TemplateType vKey = registry.createTemplateType("V");
    FunctionType pairCtor =
        FunctionType.builder(registry)
            .forInterface()
            .withName("Pair")
            .withSourceNode(fn("Pair"))
            .withTemplateKeys(kKey, vKey)
            .build();
    pairCtor.setImplicitMatch(true);
    pairCtor.getPrototype().defineDeclaredProperty("k", kKey, null);
    pairCtor.getPrototype().defineDeclaredProperty("v", vKey, null);
    ObjectType pairType = pairCtor.getInstanceType();
    registry.declareType(scope, "Pair", pairType);
    add("Pair", pairType);
    TemplatizedType pairSN = add("Pair<string,number>", registry.createTemplatizedType(pairType, string, number));
    add("Pair<number,number>", registry.createTemplatizedType(pairType, number, number));
    add("Pair<K,V>", registry.createTemplatizedType(pairType, kKey, vKey));
    Map<String, JSType> rKV = new LinkedHashMap<>();
    rKV.put("k", string);
    rKV.put("v", number);
    add("{k:string,v:number}", registry.createRecordType(rKV));
    Map<String, JSType> rK = new LinkedHashMap<>();
    rK.put("k", string);
    add("{k:string}", registry.createRecordType(rK));

    FunctionType strBoxCtor =
        registry.createConstructorType(
            "StrBox", fn("StrBox"), registry.createParameters(), null, null, false);
    strBoxCtor.setPrototypeBasedOn(boxString);
    registry.declareType(scope, "StrBox", strBoxCtor.getInstanceType());
    add("StrBox", strBoxCtor.getInstanceType());
    add("ctor:StrBox", strBoxCtor);

    TemplateType elKey = registry.createTemplateType("EL");
    FunctionType listCtor =
        FunctionType.builder(registry)
            .forInterface()
            .withName("List")
            .withSourceNode(fn("List"))
            .withTemplateKeys(elKey)
            .build();
    listCtor.setExtendedInterfaces(ImmutableList.of(registry.createTemplatizedType(iterable, elKey)));
    listCtor.getPrototype().defineDeclaredProperty("get", registry.createFunctionType(elKey, number), null);
    registry.declareType(scope, "List", listCtor.getInstanceType());
    add("List", listCtor.getInstanceType());
    add("ctor:List", listCtor);
    add("List<string>", registry.createTemplatizedType(listCtor.getInstanceType(), string));
    add("Box<Box<string>>", registry.createTemplatizedType(box, boxString));
    add("Box<T>|null", registry.createNullableType(registry.createTemplatizedType(box, tKey)));
    add(
        "function(K,V):Pair<K,V>",
        FunctionType.builder(registry)
            .withParameters(registry.createParameters(kKey, vKey))
            .withReturnType(registry.createTemplatizedType(pairType, kKey, vKey))
            .withTemplateKeys(kKey, vKey)
            .build());
    add(
        "function(this:T):T",
        FunctionType.builder(registry)
            .withParameters()
            .withReturnType(tKey)
            .withTypeOfThis(tKey)
            .withTemplateKeys(tKey)
            .build());
    add(
        "function(function(T):T, ...T):Array<T>",
        FunctionType.builder(registry)
            .withParameters(
                registry.createParametersWithVarArgs(
                    types.get(labels.indexOf("function(T):T")), tKey))
            .withReturnType(registry.createTemplatizedType(array, tKey))
            .withTemplateKeys(tKey)
            .build());
    add(
        "function(new:Box<T>, T)",
        FunctionType.builder(registry)
            .forConstructor()
            .withParameters(registry.createParameters(tKey))
            .withTypeOfThis(registry.createTemplatizedType(box, tKey))
            .withTemplateKeys(tKey)
            .build());
    add("Pair<string,number>|{k:string}", registry.createUnionType(pairSN, types.get(labels.indexOf("{k:string}"))));
    add("T|K|V", registry.createUnionType(tKey, kKey, vKey));

    // Forward references, resolved when the resolver closes.
    NamedType late = registry.createNamedType(scope, "Late", "src.js", 9, 10);
    late.defineDeclaredProperty("lateProp", number, null);
    add("named:Late", late);
    NamedType lateMissing = registry.createNamedType(scope, "LateMissing", "src.js", 11, 12);
    lateMissing.defineDeclaredProperty("p", number, null);
    add("named:LateMissing", lateMissing);
    add("Box<Late>", registry.createTemplatizedType(box, late));
    add("Late|null", registry.createNullableType(late));
    FunctionType lateCtor =
        registry.createConstructorType(
            "Late", fn("Late"), registry.createParameters(), null, null, false);
    registry.declareType(scope, "Late", lateCtor.getInstanceType());
    add("Late", lateCtor.getInstanceType());

    // The property index behind canPropertyBeDefined.
    registry.registerPropertyOnType("bar", foo);
    registry.registerPropertyOnType("bar", rec);
    registry.registerPropertyOnType("bar", registry.createUnionType(impl, nul));
    registry.registerPropertyOnType("a", recA);
    registry.registerPropertyOnType("a", lit);
    registry.registerPropertyOnType("k", pairType);
    registry.registerPropertyOnType("zzz", types.get(labels.indexOf("{x:boolean}(synthetic)")));
    registry.registerPropertyOnType("zzz", lateCtor.getInstanceType());

    // Bindings for the replacer facts.
    partialMap = boxString.getTemplateTypeMap();
    totalMap = pairSN.getTemplateTypeMap();
    inferMap.put(tKey, number);
    inferMap.put(kKey, string);

    // Types from JSDoc type expressions.
    scope.slots.put("ns", new SimpleSlot("ns", lit, false));
    String[] typeStrings = {
      "number",
      "!Array<string>",
      "Array<string>",
      "?Foo",
      "Foo",
      "!Foo",
      "{a: number, b: (string|undefined)}",
      "{a: number, 'quoted': string}",
      "function(this:Foo, number=, ...string): boolean",
      "function(new:Foo, number)",
      "function(new:number)",
      "function(...?)",
      "function(number=, string)",
      "Object<string, number>",
      "Object<number>",
      "(number|string|null)",
      "*",
      "?",
      "undefined",
      "void",
      "Box<string>",
      "Box<string, number>",
      "!Box",
      "Missing",
      "Forwarded",
      "E",
      "!E",
      "?E",
      "Rec",
      "SubFoo",
      "typeof ns",
      "!Map<string, !Array<?>>",
      "function(): undefined",
      "string=",
      "Pair<string, number>",
      "StrBox",
      "List<number>",
      "Late",
      "Box<Box<string>>",
      "{k: string, v: number}",
      "function(this:StrBox): string",
      "LateMissing",
    };
    for (String ts : typeStrings) {
      add("jsdoc:" + ts, parsed(ts, scope));
    }
    StaticTypedScope tScope = registry.createScopeWithTemplates(scope, ImmutableList.of(tKey));
    String[] templateTypeStrings = {
      "T", "Box<T>", "function(T): T", "(Array<T>|T)", "Pair<T, string>", "K",
    };
    for (String ts : templateTypeStrings) {
      add("jsdocT:" + ts, parsed(ts, tScope));
    }
  }

  // ---------------------------------------------------------------------
  // Facts

  static final String[] PROPS = {
    "toString", "valueOf", "length", "prototype", "call", "apply", "bind", "bar", "baz", "value", "a",
    "constructor", "then", "missing"
  };

  void unaryFacts(int i) {
    JSType t = types.get(i);
    line("=== " + (i + 1) + " " + labels.get(i));
    fact("toString", () -> t.toString());
    fact("annot", () -> t.toAnnotationString(JSType.Nullability.EXPLICIT));
    fact("annotImplicit", () -> t.toAnnotationString(JSType.Nullability.IMPLICIT));
    fact("typeClass", () -> t.getTypeClass());
    fact(
        "preds",
        () ->
            bits(
                t.isNoType(), t.isNoResolvedType(), t.isNoObjectType(), t.isEmptyType(),
                t.isNumberObjectType(), t.isNumberValueType(), t.isBigIntObjectType(),
                t.isBigIntValueType(), t.isFunctionPrototypeType(), t.isStringObjectType(),
                t.isSymbolObjectType(), t.isTheObjectType(), t.isStringValueType(),
                t.isSymbolValueType(), t.isKnownSymbolValueType(), t.isString(), t.isNumber(),
                t.isSymbol(), t.isBigIntOrNumber(), t.isOnlyBigInt(), t.isArrayType(),
                t.isReadonlyArrayType(), t.isBooleanObjectType(), t.isBooleanValueType(),
                t.isRegexpType(), t.isDateType(), t.isNullType(), t.isVoidType(), t.isAllType(),
                t.isUnknownType(), t.isSomeUnknownType(), t.isCheckedUnknownType(),
                t.isUnionType(), t.isRawTypeOfTemplatizedType(), t.isLiteralObject(),
                t.isGlobalThisType(), t.isFunctionType(), t.isEnumElementType(), t.isEnumType(),
                t.isNamedType(), t.isRecordType(), t.isStructuralInterface(),
                t.isStructuralType(), t.isTemplatizedType(), t.isTemplateType(),
                t.hasAnyTemplateTypes(), t.isObject(), t.isObjectType(), t.isConstructor(),
                t.isNominalType(), t.isNominalConstructorOrInterface(), t.isNativeObjectType(),
                t.loosenTypecheckingDueToForwardReferencedSupertype(), t.isInstanceType(),
                t.isInterface(), t.isOrdinaryFunction()));
    fact(
        "preds2",
        () ->
            bits(
                t.matchesNumberContext(), t.matchesStringContext(), t.matchesSymbolContext(),
                t.matchesObjectContext(), t.canBeCalled(), t.isBoxableScalar(), t.isNullable(),
                t.isVoidable(), t.isExplicitlyVoidable(), t.isResolved(),
                t.isSuccessfullyResolved(), t.isUnsuccessfullyResolved(), t.hasDisplayName()));
    fact("isStruct", () -> t.isStruct());
    fact("isDict", () -> t.isDict());
    fact("displayName", () -> t.getDisplayName());
    fact("autoboxesTo", () -> t.autoboxesTo());
    fact("autobox", () -> t.autobox());
    fact("dereference", () -> t.dereference());
    fact("restrictByNotNullOrUndefined", () -> t.restrictByNotNullOrUndefined());
    fact("restrictByNotUndefined", () -> t.restrictByNotUndefined());
    fact("restrictByNotNull", () -> t.restrictByNotNull());
    fact("collapseUnion", () -> t.collapseUnion());
    fact("toBoolean", () -> t.getPossibleToBooleanOutcomes());
    for (Outcome o : Outcome.values()) {
      fact("givenOutcome:" + o, () -> t.getRestrictedTypeGivenOutcome(o));
    }
    fact("templateParamCount", () -> t.getTemplateParamCount());
    fact("typeParameters", () -> t.getTypeParameters());
    fact("templateTypeMap", () -> t.getTemplateTypeMap().toString());
    fact("enumeratedTypeOfEnumElement", () -> t.getEnumeratedTypeOfEnumElement());
    fact("unionMembers", () -> t.getUnionMembers());
    for (String p : PROPS) {
      fact("findPropertyType:" + p, () -> t.findPropertyType(p));
      fact("propertyKind:" + p, () -> t.getPropertyKind(p) + "/" + t.getPropertyKind(p, false));
    }
    for (String p : new String[] {"bar", "a", "k", "zzz"}) {
      fact("canPropertyBeDefined:" + p, () -> registry.canPropertyBeDefined(t, p));
    }
    fact(
        "replacePartial",
        () -> {
          TemplateTypeReplacer r = TemplateTypeReplacer.forPartialReplacement(registry, partialMap);
          return t.visit(r) + " " + r.hasMadeReplacement();
        });
    fact(
        "replaceTotal",
        () -> {
          TemplateTypeReplacer r = TemplateTypeReplacer.forTotalReplacement(registry, totalMap);
          return t.visit(r) + " " + r.hasMadeReplacement();
        });
    fact(
        "replaceInfer",
        () -> {
          TemplateTypeReplacer r = TemplateTypeReplacer.forInference(registry, inferMap);
          return t.visit(r) + " " + r.hasMadeReplacement();
        });
    fact("ttm.keys", () -> t.getTemplateTypeMap().getTemplateKeys());
    fact("ttm.values", () -> t.getTemplateTypeMap().getTemplateValues());
    fact(
        "ttm.resolved",
        () -> {
          TemplateTypeMap m = t.getTemplateTypeMap();
          List<JSType> l = new ArrayList<>();
          for (TemplateType k : m.getTemplateKeys()) {
            l.add(m.getResolvedTemplateType(k));
          }
          return l;
        });
    fact(
        "ttm.hasT",
        () -> {
          TemplateTypeMap m = t.getTemplateTypeMap();
          return m.hasTemplateKey(tKey) + "/" + m.hasTemplateType(tKey) + "/" + m.size() + "/" + m.isEmpty();
        });
    fact("ttm.lastByName(T)", () -> t.getTemplateTypeMap().getLastTemplateTypeKeyByName("T"));
    fact(
        "readableTypeName",
        () -> {
          Node n = IR.name("x");
          n.setJSType(t);
          return registry.getReadableTypeName(n) + "/" + registry.getReadableTypeNameNoDeref(n);
        });
    ObjectType o = t.toMaybeObjectType();
    if (o != null) {
      fact("o.referenceName", () -> o.getReferenceName());
      fact("o.normalizedReferenceName", () -> o.getNormalizedReferenceName());
      fact("o.hasReferenceName", () -> o.hasReferenceName());
      fact("o.constructor", () -> o.getConstructor());
      fact("o.implicitPrototype", () -> o.getImplicitPrototype());
      fact("o.ownPropertyNames", () -> o.getOwnPropertyNames());
      fact("o.propertyNames", () -> o.getPropertyNames());
      fact("o.propertiesCount", () -> o.getPropertiesCount());
      fact("o.templateTypes", () -> o.getTemplateTypes());
      fact("o.ownerFunction", () -> o.getOwnerFunction());
      fact("o.ctorImplementedInterfaces", () -> o.getCtorImplementedInterfaces());
      fact("o.ctorExtendedInterfaces", () -> o.getCtorExtendedInterfaces());
      fact("o.typeOfThis", () -> o.getTypeOfThis());
      fact("o.superClassConstructor", () -> o.getSuperClassConstructor());
      fact("o.rawType", () -> o.getRawType());
      fact("o.propertyTypeMap", () -> o.getPropertyTypeMap().toString());
      fact("o.ownPropertyKeys", () -> o.getOwnPropertyKeys());
      fact("o.hasCachedValues", () -> o.hasCachedValues());
      fact("o.enumeratedTypeOfEnumObject", () -> o.getEnumeratedTypeOfEnumObject());
      fact("o.detectInheritanceCycle", () -> o.detectInheritanceCycle());
      fact("o.slot:bar", () -> o.getSlot("bar"));
      fact("o.ownSlot:bar", () -> o.getOwnSlot("bar"));
      fact("o.slot:k", () -> o.getSlot("k"));
      for (String p : PROPS) {
        fact("o.propertyType:" + p, () -> o.getPropertyType(p));
        fact("o.hasOwnProperty:" + p, () -> o.hasOwnProperty(p) + "/" + o.isPropertyTypeDeclared(p) + "/" + o.isPropertyTypeInferred(p));
        fact("o.closestDefiningType:" + p, () -> o.getClosestDefiningType(p));
      }
    }
    FunctionType f = t.toMaybeFunctionType();
    if (f != null) {
      fact("f.kind", () -> f.getKind());
      fact("f.parameters", () -> f.getParameters());
      fact("f.minArity", () -> f.getMinArity());
      fact("f.maxArity", () -> f.getMaxArity());
      fact("f.returnType", () -> f.getReturnType());
      fact("f.isReturnTypeInferred", () -> f.isReturnTypeInferred());
      fact("f.typeOfThis", () -> f.getTypeOfThis());
      fact("f.hasInstanceType", () -> f.hasInstanceType());
      fact("f.instanceType", () -> f.getInstanceType());
      fact("f.prototype", () -> f.getPrototype());
      fact("f.implementedInterfaces", () -> f.getImplementedInterfaces());
      fact("f.extendedInterfaces", () -> f.getExtendedInterfaces());
      fact("f.allImplementedInterfaces", () -> f.getAllImplementedInterfaces());
      fact("f.ancestorInterfaces", () -> f.getAncestorInterfaces());
      fact("f.isAbstract", () -> f.isAbstract());
      fact("f.makesStructs", () -> f.makesStructs());
      fact("f.makesDicts", () -> f.makesDicts());
      fact("f.hasProperties", () -> f.hasProperties());
      fact("f.isAmbiguousConstructor", () -> f.isAmbiguousConstructor());
      fact("f.checkExtendsLoop", () -> f.checkExtendsLoop());
      fact("f.constructorOnlyTemplateParameters", () -> f.getConstructorOnlyTemplateParameters());
      fact("f.bindReturnType(0)", () -> f.getBindReturnType(0));
      fact("f.bindReturnType(1)", () -> f.getBindReturnType(1));
      fact("f.acceptsArguments()", () -> f.acceptsArguments(ImmutableList.of()));
      fact("f.acceptsArguments(number)", () -> f.acceptsArguments(ImmutableList.of(nat(JSTypeNative.NUMBER_TYPE))));
      fact(
          "f.acceptsArguments(number,string)",
          () -> f.acceptsArguments(ImmutableList.of(nat(JSTypeNative.NUMBER_TYPE), nat(JSTypeNative.STRING_TYPE))));
      fact("f.propertyTypeMap", () -> f.getPropertyTypeMap().toString());
    }
    UnionType u = t.toMaybeUnionType();
    if (u != null) {
      fact("u.alternates", () -> u.getAlternates());
      fact("u.restrictedUnion(null)", () -> u.getRestrictedUnion(nat(JSTypeNative.NULL_TYPE)));
      fact("u.restrictedUnion(Object)", () -> u.getRestrictedUnion(nat(JSTypeNative.OBJECT_TYPE)));
      fact("u.contains(number)", () -> u.contains(nat(JSTypeNative.NUMBER_TYPE)));
    }
    EnumType e = t.toMaybeEnumType();
    if (e != null) {
      fact("e.elements", () -> e.getElements());
      fact("e.elementsType", () -> e.getElementsType());
    }
    TemplatizedType tt = t.toMaybeTemplatizedType();
    if (tt != null) {
      fact("tt.referencedType", () -> tt.getReferencedType());
      fact("tt.templateTypes", () -> tt.getTemplateTypes());
    }
    TemplateType tk = t.toMaybeTemplateType();
    if (tk != null) {
      fact("tk.bound", () -> tk.getBound());
      fact("tk.containsCycle", () -> tk.containsCycle());
    }
    NamedType nt = t.toMaybeNamedType();
    if (nt != null) {
      fact("nt.referencedType", () -> nt.getReferencedType());
    }
  }

  static String bits(boolean... bs) {
    StringBuilder sb = new StringBuilder();
    for (boolean b : bs) {
      sb.append(b ? '1' : '0');
    }
    return sb.toString();
  }

  static String ident(JSType r, JSType a, JSType b) {
    if (r == a && r == b) {
      return "=AB";
    } else if (r == a) {
      return "=A";
    } else if (r == b) {
      return "=B";
    }
    return str(r);
  }

  void pair(int i, int j) {
    JSType a = types.get(i);
    JSType b = types.get(j);
    StringBuilder sb = new StringBuilder();
    sb.append(i + 1).append(' ').append(j + 1).append('\t');
    sb.append(pf(() -> bits(
        a.isSubtypeOf(b),
        a.isSubtypeOf(b, JSType.SubtypingMode.IGNORE_NULL_UNDEFINED),
        a.isSubtypeWithoutStructuralTyping(b),
        a.equals(b),
        a.differsFrom(b),
        a.canTestForEqualityWith(b),
        a.canTestForShallowEqualityWith(b),
        a.canCastTo(b),
        JSType.areSimilar(a, b))));
    sb.append('\t').append(pf(() -> a.testForEquality(b)));
    sb.append('\t').append(pf(() -> ident(a.getLeastSupertype(b), a, b)));
    sb.append('\t').append(pf(() -> ident(a.getGreatestSubtype(b), a, b)));
    sb.append('\t').append(pf(() -> a.getTypesUnderEquality(b)));
    sb.append('\t').append(pf(() -> a.getTypesUnderInequality(b)));
    sb.append('\t').append(pf(() -> a.getTypesUnderShallowInequality(b)));
    sb.append('\t').append(pf(() -> a.getTypesUnderShallowEquality(b)));
    sb.append('\t').append(pf(() -> {
      if (a instanceof FunctionType fa && b instanceof FunctionType fb) {
        return "fn:" + fa.hasEqualCallType(fb);
      }
      if (a instanceof ObjectType oa && b instanceof ObjectType ob) {
        return "obj:" + registry.findCommonSuperObject(oa, ob);
      }
      return "-";
    }));
    line(esc(sb.toString()));
  }

  static String pf(Supplier<Object> s) {
    try {
      return str(s.get());
    } catch (RuntimeException e) {
      return "!" + e.getClass().getSimpleName();
    }
  }

  void registryFacts() {
    line("=== registry");
    for (String p : new String[] {"bar", "a", "k", "zzz"}) {
      fact("eachReferenceTypeWithProperty:" + p, () -> registry.getEachReferenceTypeWithProperty(p));
    }
    for (String n : new String[] {"Foo", "Box", "Late", "Missing", "E", "Pair", "Forwarded", "List", "ns"}) {
      fact("getType:" + n, () -> registry.getType(scope, n));
      fact("getGlobalType:" + n, () -> registry.getGlobalType(n));
      fact("isNonNullableName:" + n, () -> registry.isNonNullableName(scope, n));
    }
  }

  void run(String pairsArg) {
    build();
    registryFacts();
    for (int i = 0; i < types.size(); i++) {
      unaryFacts(i);
    }
    line("=== pairs");
    int n = types.size();
    if (pairsArg.equals("all")) {
      for (int i = 0; i < n; i++) {
        for (int j = 0; j < n; j++) {
          pair(i, j);
        }
      }
    } else {
      int count = Integer.parseInt(pairsArg);
      for (int k = 0; k < count; k++) {
        pair(nextInt(n), nextInt(n));
      }
    }
    line("=== end");
  }

  public static void main(String[] args) {
    String pairs = "all";
    for (String a : args) {
      if (a.startsWith("--pairs=")) {
        pairs = a.substring("--pairs=".length());
      }
    }
    TypeAlgebraDump d = new TypeAlgebraDump();
    d.run(pairs);
    PrintStream ps = new PrintStream(System.out, false, StandardCharsets.UTF_8);
    ps.print(d.out);
    ps.flush();
  }
}

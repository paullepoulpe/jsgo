// Differential oracle for the jscomp classes that sit on the type system
// (internal/jscomp): the reverse abstract interpreters
// (jscomp.type.SemanticReverseAbstractInterpreter, ClosureReverseAbstractInterpreter
// and their ChainableReverseAbstractInterpreter base: the type refinement that
// type inference applies on conditions), and JsIterables, Promises,
// TypeMatchingStrategy, TypeMismatch.Accumulator, DestructuredTarget and
// serialization.JSTypeColorIdHasher.
//
// Builds a type universe (every JSTypeNative, a few user-defined classes,
// interfaces, a templated class and enums, then every type expression of
// <types.txt> parsed with JsDocInfoParser.parseTypeString and evaluated in a
// scope that declares the user types).
//
// The "rai" output is a deterministic, seeded set of conditions
// (typeof/instanceof/==/!=/===/!==/!/&&/||/in/</switch-case, comparisons with
// null/undefined, Array.isArray, goog.isObject, qualified names, ...) over
// variables of those types. For every condition it asks four interpreter chains
// (Semantic, Closure, Closure->Semantic, Semantic->Closure) for
// getPreciserScopeKnowingConditionOutcome under each Outcome, with a simple
// copy-on-write FlowScope, and prints the refined slots.
//
// The "deps" output (on a fresh universe) prints, for every type, the results
// of JsIterables, Promises, DestructuredTarget and JSTypeColorIdHasher, and for
// seeded pairs of types TypeMatchingStrategy.match and the mismatches a
// TypeMismatch.Accumulator registers.
//
// internal/jscomp/type_reverse_interpreter_oracle_test.go mirrors this file
// step by step (same API calls, same order, same seeded choices) and must
// reproduce both outputs byte for byte. Keep the two in step. The class lives
// in package com.google.javascript.jscomp to reach the package-private API.
//
// Usage: java -cp closure-compiler.jar:<build> com.google.javascript.jscomp.ReverseInterpreterDump \
//          <types.txt> <rai-out.tsv> <deps-out.tsv> <serialization-out.tsv>
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.parsing.JsDocInfoParser;
import com.google.javascript.jscomp.type.ChainableReverseAbstractInterpreter;
import com.google.javascript.jscomp.type.ClosureReverseAbstractInterpreter;
import com.google.javascript.jscomp.type.FlowScope;
import com.google.javascript.jscomp.type.ReverseAbstractInterpreter;
import com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.Outcome;
import com.google.javascript.rhino.Token;
import com.google.javascript.rhino.jstype.EnumType;
import com.google.javascript.rhino.jstype.FunctionType;
import com.google.javascript.rhino.jstype.JSType;
import com.google.javascript.rhino.jstype.JSTypeNative;
import com.google.javascript.rhino.jstype.JSTypeRegistry;
import com.google.javascript.rhino.jstype.JSTypeResolver;
import com.google.javascript.rhino.jstype.ObjectType;
import com.google.javascript.rhino.jstype.SimpleSlot;
import com.google.javascript.rhino.jstype.StaticTypedScope;
import com.google.javascript.rhino.jstype.StaticTypedSlot;
import com.google.javascript.rhino.jstype.TemplateType;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public final class ReverseInterpreterDump {
  private final StringBuilder out = new StringBuilder();
  private final List<String> labels = new ArrayList<>();
  private final List<JSType> types = new ArrayList<>();
  private final List<JSType> ctors = new ArrayList<>();
  private final JSTypeRegistry registry;
  private final TestScope scope = new TestScope();
  private final List<String> chainNames = new ArrayList<>();
  private final List<ReverseAbstractInterpreter> chains = new ArrayList<>();

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

  // ---------------------------------------------------------------------
  // Deterministic sequences (the same LCG as TypeAlgebraDump)

  private long seed = 0x5DEECE66DL;

  int nextInt(int n) {
    seed = seed * 6364136223846793005L + 1442695040888963407L;
    return (int) ((seed >>> 33) % n);
  }

  // ---------------------------------------------------------------------
  // A reporter, a declaration scope and a copy-on-write flow scope

  final class TestReporter implements ErrorReporter {
    @Override
    public void warning(String message, String sourceName, int line, int lineOffset) {
      line("WARNING\t" + esc(message));
    }

    @Override
    public void error(String message, String sourceName, int line, int lineOffset) {
      line("ERROR\t" + esc(message));
    }
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

  /** A flow scope whose updates copy the slot map (insertion-ordered). */
  static final class SimpleFlowScope implements FlowScope {
    final LinkedHashMap<String, StaticTypedSlot> slots;

    SimpleFlowScope(LinkedHashMap<String, StaticTypedSlot> slots) {
      this.slots = slots;
    }

    @Override
    public FlowScope withSyntacticScope(StaticTypedScope scope) {
      throw new UnsupportedOperationException();
    }

    @Override
    public FlowScope inferSlotType(String symbol, JSType type) {
      LinkedHashMap<String, StaticTypedSlot> copy = new LinkedHashMap<>(slots);
      copy.put(symbol, new SimpleSlot(symbol, type, true));
      return new SimpleFlowScope(copy);
    }

    @Override
    public FlowScope inferQualifiedSlot(
        Node node, String symbol, JSType bottomType, JSType inferredType, boolean declare) {
      LinkedHashMap<String, StaticTypedSlot> copy = new LinkedHashMap<>(slots);
      copy.put(symbol, new SimpleSlot(symbol, inferredType, true));
      return new SimpleFlowScope(copy);
    }

    @Override
    public StaticTypedScope getDeclarationScope() {
      return null;
    }

    @Override
    public Node getRootNode() {
      return null;
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

  ReverseInterpreterDump() {
    this.registry = new JSTypeRegistry(new TestReporter());
  }

  // ---------------------------------------------------------------------
  // The type universe

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

  // ---------------------------------------------------------------------
  // Conditions

  static final String[] TYPEOFS = {
    "string", "number", "boolean", "object", "function", "undefined", "symbol", "bigint", "unknown"
  };
  static final Token[] EQ_OPS = {Token.EQ, Token.NE, Token.SHEQ, Token.SHNE};
  static final Token[] CMP_OPS = {Token.LT, Token.LE, Token.GT, Token.GE};
  static final String[] PROPS = {"p", "bar", "length", "then", "prototype", "qux"};

  /** The state of one case: the blind scope's slots and their declaration order. */
  final class Case {
    final LinkedHashMap<String, StaticTypedSlot> slots = new LinkedHashMap<>();
    final StringBuilder decls = new StringBuilder();

    void declare(String name, JSType t) {
      slots.put(name, new SimpleSlot(name, t, true));
      if (decls.length() > 0) {
        decls.append("; ");
      }
      decls.append(name).append(':').append(str(t));
    }

    /** A declared variable whose NAME node is typed like its slot. */
    Node var(String name, JSType t) {
      declare(name, t);
      return typed(Node.newString(Token.NAME, name), t);
    }
  }

  static Node typed(Node n, JSType t) {
    if (t != null) {
      n.setJSType(t);
    }
    return n;
  }

  Node num(double d) {
    return typed(Node.newNumber(d), nat(JSTypeNative.NUMBER_TYPE));
  }

  Node strLit(String s) {
    return typed(Node.newString(s), nat(JSTypeNative.STRING_TYPE));
  }

  Node nul() {
    return typed(new Node(Token.NULL), nat(JSTypeNative.NULL_TYPE));
  }

  Node undef() {
    return typed(Node.newString(Token.NAME, "undefined"), nat(JSTypeNative.VOID_TYPE));
  }

  Node voidZero() {
    return typed(new Node(Token.VOID, num(0)), nat(JSTypeNative.VOID_TYPE));
  }

  Node tru() {
    return typed(new Node(Token.TRUE), nat(JSTypeNative.BOOLEAN_TYPE));
  }

  Node literal(int k) {
    return switch (k) {
      case 0 -> strLit("x");
      case 1 -> num(0);
      case 2 -> nul();
      case 3 -> undef();
      case 4 -> voidZero();
      default -> tru();
    };
  }

  static Node binop(Token op, Node a, Node b, boolean swap) {
    return swap ? new Node(op, b, a) : new Node(op, a, b);
  }

  Node typeofStr(Node operand, Token op, String s, boolean swap) {
    return binop(op, new Node(Token.TYPEOF, operand), strLit(s), swap);
  }

  static Node call(Node callee, Node arg) {
    return new Node(Token.CALL, callee, arg);
  }

  static Node getprop(Node obj, String prop) {
    return IR.getprop(obj, prop);
  }

  /** goog.<fn>(arg). */
  static Node goog(String fn, Node arg) {
    return call(getprop(Node.newString(Token.NAME, "goog"), fn), arg);
  }

  /** switch (cond) { case expr: } — returns the CASE node. */
  static Node switchCase(Node cond, Node expr) {
    Node c = new Node(Token.CASE, expr, new Node(Token.BLOCK));
    new Node(Token.SWITCH, cond, new Node(Token.SWITCH_BODY, c));
    return c;
  }

  /** o.p with its slot and node type chosen by variant. */
  Node qualified(Case c, String obj, String prop, JSType t) {
    int variant = nextInt(3);
    Node n = getprop(Node.newString(Token.NAME, obj), prop);
    if (variant != 1) {
      c.declare(obj + "." + prop, t);
    }
    if (variant != 2) {
      typed(n, t);
    }
    return n;
  }

  JSType randomType() {
    return types.get(nextInt(types.size()));
  }

  static final int UNARY = 37;

  Node unary(Case c, int tmpl, JSType t) {
    switch (tmpl) {
      case 0:
        return c.var("a", t);
      case 1:
        return new Node(Token.NOT, c.var("a", t));
      case 2:
        return new Node(Token.NOT, new Node(Token.NOT, c.var("a", t)));
      case 3:
        {
          Token op = EQ_OPS[nextInt(4)];
          String s = TYPEOFS[nextInt(TYPEOFS.length)];
          return typeofStr(c.var("a", t), op, s, nextInt(2) == 1);
        }
      case 4:
      case 5:
      case 6:
      case 7:
      case 8:
      case 9:
        {
          Token op = EQ_OPS[nextInt(4)];
          boolean swap = nextInt(2) == 1;
          Node lit =
              switch (tmpl) {
                case 4 -> nul();
                case 5 -> undef();
                case 6 -> voidZero();
                case 7 -> num(0);
                case 8 -> strLit("");
                default -> tru();
              };
          return binop(op, c.var("a", t), lit, swap);
        }
      case 10:
        {
          Token op = CMP_OPS[nextInt(4)];
          return binop(op, c.var("a", t), num(3), nextInt(2) == 1);
        }
      case 11:
        {
          JSType ctor = ctors.get(nextInt(ctors.size()));
          return new Node(Token.INSTANCEOF, c.var("a", t), typed(Node.newString(Token.NAME, "C"), ctor));
        }
      case 12:
        return new Node(Token.IN, strLit(PROPS[nextInt(PROPS.length)]), c.var("a", t));
      case 13:
        return call(getprop(Node.newString(Token.NAME, "Array"), "isArray"), c.var("a", t));
      case 14:
        return goog("isObject", c.var("a", t));
      case 15:
        return goog("isString", c.var("a", t));
      case 16:
        return qualified(c, "o", "p", t);
      case 17:
        {
          Token op = EQ_OPS[nextInt(4)];
          String s = TYPEOFS[nextInt(TYPEOFS.length)];
          return typeofStr(qualified(c, "o", "p", t), op, s, nextInt(2) == 1);
        }
      case 18:
        {
          Token op = EQ_OPS[nextInt(4)];
          return binop(op, qualified(c, "o", "p", t), nul(), nextInt(2) == 1);
        }
      case 19:
        return goog("isObject", qualified(c, "o", "p", t));
      case 20:
        return new Node(Token.NOT, qualified(c, "o", "p", t));
      case 21:
        return switchCase(c.var("a", t), literal(nextInt(6)));
      case 22:
        return switchCase(
            new Node(Token.TYPEOF, c.var("a", t)), strLit(TYPEOFS[nextInt(TYPEOFS.length)]));
      case 23:
        {
          Node a = c.var("a", t);
          return new Node(Token.ASSIGN, a, c.var("b", randomType()));
        }
      case 24:
        {
          Node object = getprop(Node.newString(Token.NAME, "o"), "q");
          typed(object, t);
          String prop = PROPS[nextInt(PROPS.length)];
          if (nextInt(2) == 1) {
            c.declare("o.q." + prop, nat(JSTypeNative.NUMBER_TYPE));
          }
          return new Node(Token.IN, strLit(prop), object);
        }
      case 25:
        {
          Node a = c.var("a", t);
          return getprop(a, "b");
        }
      case 26:
        return typed(getprop(call(Node.newString(Token.NAME, "f"), num(1)), "p"), t);
      case 27:
        return nextInt(2) == 1 ? new Node(Token.THIS) : goog("isObject", new Node(Token.THIS));
      case 28:
        {
          Token op = EQ_OPS[nextInt(4)];
          Node a = c.var("a", t);
          return new Node(op, a, typed(Node.newString(Token.NAME, "a"), t));
        }
      case 29:
        {
          Token op = EQ_OPS[nextInt(4)];
          String s = TYPEOFS[nextInt(TYPEOFS.length)];
          return typeofStr(typed(Node.newString(Token.NAME, "x"), t), op, s, false);
        }
      case 30:
        return typed(Node.newString(Token.NAME, "x"), t);
      case 31:
        {
          c.declare("a", null);
          Node a = typed(Node.newString(Token.NAME, "a"), t);
          return nextInt(2) == 1 ? a : new Node(Token.NOT, a);
        }
      case 32:
        {
          Token op = EQ_OPS[nextInt(4)];
          Node a = c.var("a", t);
          return binop(op, a, typed(Node.newString(Token.NAME, "y"), randomType()), false);
        }
      case 33:
        return call(getprop(Node.newString(Token.NAME, "x"), "isObject"), c.var("a", t));
      case 34:
        {
          Node call = goog("isObject", c.var("a", t));
          call.addChildToBack(num(1));
          return call;
        }
      case 35:
        {
          Node abc = getprop(getprop(Node.newString(Token.NAME, "a"), "b"), "c");
          if (nextInt(2) == 1) {
            c.declare("a.b.c", t);
          }
          return goog("isObject", typed(abc, t));
        }
      default:
        return goog("isObject", typed(call(Node.newString(Token.NAME, "f"), num(2)), t));
    }
  }

  static final int BINARY = 20;

  Node binary(Case c, int tmpl, JSType t1, JSType t2) {
    switch (tmpl) {
      case 0:
        {
          Token op = EQ_OPS[nextInt(4)];
          return new Node(op, c.var("a", t1), c.var("b", t2));
        }
      case 1:
        {
          Token op = CMP_OPS[nextInt(4)];
          return new Node(op, c.var("a", t1), c.var("b", t2));
        }
      case 2:
        return new Node(Token.AND, c.var("a", t1), c.var("b", t2));
      case 3:
        return new Node(Token.OR, c.var("a", t1), c.var("b", t2));
      case 4:
        return new Node(Token.NOT, new Node(Token.AND, c.var("a", t1), c.var("b", t2)));
      case 5:
        return new Node(Token.NOT, new Node(Token.OR, c.var("a", t1), c.var("b", t2)));
      case 6:
        {
          Node a = c.var("a", t1);
          Node b = c.var("b", t2);
          return new Node(Token.AND, a, new Node(Token.OR, b, c.var("c", randomType())));
        }
      case 7:
        {
          Node a = c.var("a", t1);
          Node b = c.var("b", t2);
          return new Node(Token.AND, new Node(Token.OR, a, b), c.var("c", randomType()));
        }
      case 8:
        return new Node(Token.INSTANCEOF, c.var("a", t1), c.var("b", t2));
      case 9:
        {
          Node a = c.var("a", t1);
          return new Node(Token.AND, a, goog("isObject", c.var("b", t2)));
        }
      case 10:
        {
          Node a = goog("isObject", c.var("a", t1));
          return new Node(Token.OR, a, c.var("b", t2));
        }
      case 11:
        {
          Node a = c.var("a", t1);
          Node b = c.var("b", t2);
          return new Node(
              Token.AND, new Node(Token.ASSIGN, a, b), typed(Node.newString(Token.NAME, "b"), t2));
        }
      case 12:
        {
          Node a = typeofStr(c.var("a", t1), Token.EQ, TYPEOFS[nextInt(TYPEOFS.length)], false);
          return new Node(Token.AND, a, c.var("b", t2));
        }
      case 13:
        {
          Node p = qualified(c, "a", "p", t2);
          c.declare("a", t1);
          return new Node(Token.AND, p, typed(Node.newString(Token.NAME, "a"), t1));
        }
      case 14:
        {
          Token op = EQ_OPS[nextInt(4)];
          return new Node(Token.NOT, new Node(op, c.var("a", t1), c.var("b", t2)));
        }
      case 15:
        {
          Token op = EQ_OPS[nextInt(4)];
          return new Node(op, c.var("a", t1), typed(Node.newNumber(1), t2));
        }
      case 16:
        {
          Node a = c.var("a", t1);
          return switchCase(a, c.var("b", t2));
        }
      case 17:
        {
          Node a = c.var("a", t1);
          return new Node(Token.AND, a, new Node(Token.NOT, c.var("b", t2)));
        }
      case 18:
        {
          Node a = call(getprop(Node.newString(Token.NAME, "Array"), "isArray"), c.var("a", t1));
          return new Node(Token.OR, a, c.var("b", t2));
        }
      default:
        {
          Node a = c.var("a", t1);
          Node b = c.var("b", t2);
          Node c3 = c.var("c", randomType());
          return new Node(
              Token.OR,
              new Node(Token.AND, a, b),
              new Node(Token.AND, c3, typed(Node.newString(Token.NAME, "a"), t1)));
        }
    }
  }

  /** Serializes a condition: (TOKEN[ "string"][:type] child...). */
  static void sexp(Node n, StringBuilder sb) {
    sb.append('(').append(n.getToken().name());
    if (n.isName() || n.isStringLit() || n.isGetProp()) {
      sb.append(" \"").append(esc(n.getString())).append('"');
    } else if (n.isNumber()) {
      sb.append(' ').append(n.getDouble());
    }
    if (n.getJSType() != null) {
      sb.append(" :").append(str(n.getJSType()));
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      sb.append(' ');
      sexp(c, sb);
    }
    sb.append(')');
  }

  /** Prints how result differs from blind. */
  static String diff(SimpleFlowScope blind, FlowScope result) {
    if (result == blind) {
      return "=";
    }
    if (!(result instanceof SimpleFlowScope)) {
      return "W";
    }
    StringBuilder sb = new StringBuilder("~");
    for (Map.Entry<String, StaticTypedSlot> e : ((SimpleFlowScope) result).slots.entrySet()) {
      StaticTypedSlot old = blind.slots.get(e.getKey());
      JSType t = e.getValue().getType();
      if (old == null) {
        sb.append(" +").append(esc(e.getKey())).append(':').append(str(t));
      } else if (old.getType() != t) {
        sb.append(' ').append(esc(e.getKey())).append(':').append(str(t));
      }
    }
    return sb.toString();
  }

  void runCase(String kind, Case c, Node cond) {
    StringBuilder sb = new StringBuilder();
    // A CASE is shown with its switch.
    sexp(cond.isCase() ? cond.getGrandparent() : cond, sb);
    line("case\t" + kind + "\t" + sb + "\t" + c.decls);
    SimpleFlowScope blind = new SimpleFlowScope(c.slots);
    for (int i = 0; i < chains.size(); i++) {
      StringBuilder r = new StringBuilder(chainNames.get(i));
      for (Outcome o : Outcome.values()) {
        String v;
        try {
          v = diff(blind, chains.get(i).getPreciserScopeKnowingConditionOutcome(cond, blind, o));
        } catch (RuntimeException e) {
          v = "!" + e.getClass().getSimpleName();
        }
        r.append('\t').append(v);
      }
      line(r.toString());
    }
  }

  void run(List<String> typeStrings) {
    build(typeStrings);

    chainNames.add("S");
    chains.add(new SemanticReverseAbstractInterpreter(registry));
    chainNames.add("C");
    chains.add(new ClosureReverseAbstractInterpreter(registry));
    chainNames.add("CS");
    chains.add(
        new ClosureReverseAbstractInterpreter(registry)
            .append(new SemanticReverseAbstractInterpreter(registry))
            .getFirst());
    chainNames.add("SC");
    chains.add(
        new SemanticReverseAbstractInterpreter(registry)
            .append(new ClosureReverseAbstractInterpreter(registry))
            .getFirst());

    int n = types.size();
    // Every unary template on every type, then every typeof string on every type.
    for (int i = 0; i < n; i++) {
      for (int tmpl = 0; tmpl < UNARY; tmpl++) {
        Case c = new Case();
        Node cond = unary(c, tmpl, types.get(i));
        runCase("u" + tmpl + ":" + (i + 1), c, cond);
      }
      for (String s : TYPEOFS) {
        Case c = new Case();
        Token op = EQ_OPS[nextInt(4)];
        Node cond = typeofStr(c.var("a", types.get(i)), op, s, nextInt(2) == 1);
        runCase("t:" + (i + 1), c, cond);
      }
    }
    for (int k = 0; k < 6000; k++) {
      int i = nextInt(n);
      int j = nextInt(n);
      int tmpl = nextInt(BINARY);
      Case c = new Case();
      Node cond = binary(c, tmpl, types.get(i), types.get(j));
      runCase("b" + tmpl + ":" + (i + 1) + "," + (j + 1), c, cond);
    }
  }

  // ---------------------------------------------------------------------
  // The "deps" section

  /** Runs s, printing its result or the class of the exception it threw. */
  void fact(String name, java.util.function.Supplier<String> s) {
    String v;
    try {
      v = s.get();
    } catch (RuntimeException e) {
      v = "!" + e.getClass().getSimpleName();
    }
    line(name + "\t" + v);
  }

  static String boxed(JsIterables.MaybeBoxedType m) {
    String templated;
    try {
      templated = str(m.getTemplatedType());
    } catch (RuntimeException e) {
      templated = "!" + e.getClass().getSimpleName();
    }
    return m.isMatch() + " " + templated + " " + str(m.getMismatchType()) + " " + str(m.orElse(null));
  }

  /** const {p: x, ['k']: y = 1, ...r} or [a, b = 1, , ...c], typed by t. */
  String destructure(JSType t, boolean object) {
    Node pattern;
    if (object) {
      Node k = new Node(Token.COMPUTED_PROP, strLit("k"), new Node(Token.DEFAULT_VALUE, Node.newString(Token.NAME, "y"), num(1)));
      pattern =
          new Node(
              Token.OBJECT_PATTERN,
              IR.stringKey("p", Node.newString(Token.NAME, "x")),
              k);
      pattern.addChildToBack(IR.stringKey("q", new Node(Token.DEFAULT_VALUE, Node.newString(Token.NAME, "z"), strLit(""))));
      pattern.addChildToBack(new Node(Token.OBJECT_REST, Node.newString(Token.NAME, "r")));
    } else {
      pattern =
          new Node(
              Token.ARRAY_PATTERN,
              Node.newString(Token.NAME, "a"),
              new Node(Token.DEFAULT_VALUE, Node.newString(Token.NAME, "b"), num(1)),
              new Node(Token.EMPTY));
      pattern.addChildToBack(new Node(Token.ITER_REST, Node.newString(Token.NAME, "c")));
    }
    StringBuilder sb = new StringBuilder();
    for (DestructuredTarget target :
        DestructuredTarget.createAllNonEmptyTargetsInPattern(registry, t, pattern)) {
      sb.append(" [")
          .append(target.getNode().getString())
          .append(' ')
          .append(target.hasComputedProperty())
          .append(target.hasStringKey())
          .append(target.hasDefaultValue());
      String a;
      try {
        a = str(target.inferTypeWithoutUsingDefaultValue());
      } catch (RuntimeException e) {
        a = "!" + e.getClass().getSimpleName();
      }
      String b;
      try {
        b = str(target.inferType());
      } catch (RuntimeException e) {
        b = "!" + e.getClass().getSimpleName();
      }
      sb.append(' ').append(a).append(" | ").append(b).append(']');
    }
    return sb.toString();
  }

  void runDeps(List<String> typeStrings) throws Exception {
    build(typeStrings);
    java.lang.reflect.Constructor<?> hasherCtor =
        Class.forName("com.google.javascript.jscomp.serialization.JSTypeColorIdHasher")
            .getDeclaredConstructor(JSTypeRegistry.class);
    hasherCtor.setAccessible(true);
    Object hasher = hasherCtor.newInstance(registry);
    java.lang.reflect.Method hash =
        hasher.getClass().getDeclaredMethod("hashObjectType", ObjectType.class);
    hash.setAccessible(true);

    int n = types.size();
    for (int i = 0; i < n; i++) {
      JSType t = types.get(i);
      String p = "t" + (i + 1) + "\t";
      fact(p + "elem", () -> str(JsIterables.getElementType(t, registry)));
      fact(p + "ret", () -> str(JsIterables.getReturnElementType(t, registry)));
      fact(p + "box", () -> boxed(JsIterables.maybeBoxIterable(t, registry)));
      fact(p + "boxAsync", () -> boxed(JsIterables.maybeBoxIterableOrAsyncIterable(t, registry)));
      fact(p + "iterOf", () -> str(JsIterables.createIterableTypeOf(t, registry)));
      fact(p + "thenable", () -> str(Promises.getTemplateTypeOfThenable(registry, t)));
      fact(p + "resolved", () -> str(Promises.getResolvedType(registry, t)));
      fact(p + "wrap", () -> str(Promises.wrapInIThenable(registry, t)));
      fact(p + "asyncRet", () -> str(Promises.createAsyncReturnableType(registry, t)));
      fact(p + "objPattern", () -> destructure(t, true));
      fact(p + "arrPattern", () -> destructure(t, false));
      fact(
          p + "colorId",
          () -> {
            ObjectType o = t.toMaybeObjectType();
            if (o == null) {
              return "-";
            }
            try {
              return hash.invoke(hasher, o).toString();
            } catch (java.lang.reflect.InvocationTargetException e) {
              throw (RuntimeException) e.getCause();
            } catch (IllegalAccessException e) {
              throw new AssertionError(e);
            }
          });
    }
    Node location = IR.name("loc");
    for (int k = 0; k < 8000; k++) {
      int i = nextInt(n);
      int j = nextInt(n);
      JSType t1 = types.get(i);
      JSType t2 = types.get(j);
      String p = "p" + (i + 1) + "," + (j + 1) + "\t";
      fact(
          p + "match",
          () -> {
            StringBuilder sb = new StringBuilder();
            for (TypeMatchingStrategy s : TypeMatchingStrategy.values()) {
              sb.append(' ').append(s.match(t1, t2));
            }
            return sb.toString();
          });
      fact(
          p + "mismatches",
          () -> {
            TypeMismatch.Accumulator acc = new TypeMismatch.Accumulator();
            acc.registerMismatch(location, t1, t2);
            StringBuilder sb = new StringBuilder();
            for (TypeMismatch m : acc.getMismatches()) {
              sb.append(" <").append(str(m.found())).append(" ; ").append(str(m.required())).append('>');
            }
            return sb.toString();
          });
    }
  }

  // ---------------------------------------------------------------------
  // The "serialization" section: InvalidatingTypes and JSTypeReconserializer

  void runSerialization(List<String> typeStrings) throws Exception {
    build(typeStrings);
    int n = types.size();

    // Seeded mismatches, each with its own location.
    List<TypeMismatch> mismatches = new ArrayList<>();
    for (int k = 0; k < 300; k++) {
      JSType found = types.get(nextInt(n));
      JSType required = types.get(nextInt(n));
      mismatches.add(new TypeMismatch(found, required, IR.name("loc" + k)));
    }
    InvalidatingTypes invalidating =
        new InvalidatingTypes.Builder(registry).addAllTypeMismatches(mismatches).build();
    for (int i = 0; i < n; i++) {
      JSType t = types.get(i);
      fact("t" + (i + 1) + "\tinvalidating", () -> String.valueOf(invalidating.isInvalidating(t)));
    }
    for (JSType key : invalidating.getMismatchLocations().keySet()) {
      StringBuilder sb = new StringBuilder();
      for (Node loc : invalidating.getMismatchLocations().get(key)) {
        sb.append(' ').append(loc.getString());
      }
      line("location\t" + str(key) + "\t" + sb);
    }

    Class<?> cls = Class.forName("com.google.javascript.jscomp.serialization.JSTypeReconserializer");
    java.lang.reflect.Method create =
        cls.getDeclaredMethod(
            "create",
            JSTypeRegistry.class,
            InvalidatingTypes.class,
            com.google.javascript.jscomp.serialization.StringPool.Builder.class,
            java.util.function.Predicate.class,
            com.google.javascript.jscomp.serialization.SerializationOptions.class);
    create.setAccessible(true);
    java.lang.reflect.Method serializeType = cls.getDeclaredMethod("serializeType", JSType.class);
    serializeType.setAccessible(true);
    java.lang.reflect.Method generateTypePool = cls.getDeclaredMethod("generateTypePool");
    generateTypePool.setAccessible(true);
    java.lang.reflect.Method debugMap = cls.getDeclaredMethod("getColorIdToJSTypeMapForDebugging");
    debugMap.setAccessible(true);

    com.google.javascript.jscomp.serialization.StringPool.Builder stringPool =
        com.google.javascript.jscomp.serialization.StringPool.builder();
    java.util.function.Predicate<String> propagate = (name) -> !name.startsWith("q");
    Object serializer =
        create.invoke(
            null,
            registry,
            invalidating,
            stringPool,
            propagate,
            com.google.javascript.jscomp.serialization.SerializationOptions.builder()
                .setIncludeDebugInfo(true)
                .setRunValidation(true)
                .build());
    for (int i = 0; i < n; i++) {
      JSType t = types.get(i);
      fact(
          "t" + (i + 1) + "\tpointer",
          () -> {
            try {
              return String.valueOf(serializeType.invoke(serializer, t));
            } catch (java.lang.reflect.InvocationTargetException e) {
              throw (RuntimeException) e.getCause();
            } catch (IllegalAccessException e) {
              throw new AssertionError(e);
            }
          });
    }
    fact(
        "pool",
        () -> {
          try {
            return esc(generateTypePool.invoke(serializer).toString());
          } catch (java.lang.reflect.InvocationTargetException e) {
            throw (RuntimeException) e.getCause();
          } catch (IllegalAccessException e) {
            throw new AssertionError(e);
          }
        });
    line("strings\t" + esc(stringPool.build().toProto().toString()));
    @SuppressWarnings("unchecked")
    Map<String, java.util.Collection<JSType>> debug =
        ((com.google.javascript.jscomp.jarjar.com.google.common.collect.Multimap<String, JSType>)
                debugMap.invoke(serializer))
            .asMap();
    for (Map.Entry<String, java.util.Collection<JSType>> e : debug.entrySet()) {
      StringBuilder sb = new StringBuilder();
      for (JSType t : e.getValue()) {
        sb.append(" ").append(str(t));
      }
      line("debug\t" + e.getKey() + "\t" + sb);
    }
  }

  public static void main(String[] args) throws Exception {
    List<String> typeStrings = new ArrayList<>();
    for (String l : Files.readAllLines(Paths.get(args[0]), StandardCharsets.UTF_8)) {
      if (!l.isEmpty() && !l.startsWith("#")) {
        typeStrings.add(l);
      }
    }
    ReverseInterpreterDump d = new ReverseInterpreterDump();
    d.run(typeStrings);
    Files.writeString(Paths.get(args[1]), d.out, StandardCharsets.UTF_8);
    ReverseInterpreterDump deps = new ReverseInterpreterDump();
    deps.runDeps(typeStrings);
    Files.writeString(Paths.get(args[2]), deps.out, StandardCharsets.UTF_8);
    ReverseInterpreterDump ser = new ReverseInterpreterDump();
    ser.runSerialization(typeStrings);
    Files.writeString(Paths.get(args[3]), ser.out, StandardCharsets.UTF_8);
  }
}

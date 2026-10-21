package com.google.javascript.rhino;

import com.google.javascript.jscomp.serialization.NodeProperty;
import java.lang.reflect.Constructor;
import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;
import java.util.function.BooleanSupplier;

/** Checks dedicated getters against unchanged generic property reads and exercises list ownership. */
public final class NodeHotPropSummaryProbe {
  private static final Class<?> ITEM = nestedClass("PropListItem");
  private static final Method GET_HEAD = method("getPropListHeadForTesting");
  private static final Method SET_HEAD = method("setPropListHead", ITEM);
  private static final Constructor<?> INT_ITEM = constructor("IntPropListItem", byte.class, int.class, ITEM);
  private static final Constructor<?> OBJECT_ITEM = constructor("ObjectPropListItem", byte.class, Object.class, ITEM);
  private static Field summaryField;
  private static int observations;

  public static void main(String[] args) throws Exception {
    boolean requireSummary = args.length == 1 && args[0].equals("--require-summary");
    require(args.length == 0 || requireSummary, "unexpected arguments");
    if (requireSummary) {
      summaryField = ITEM.getDeclaredField("hotBooleanSummary");
      require(summaryField.getType() == byte.class && Modifier.isFinal(summaryField.getModifiers()),
          "summary must be an immutable byte on the property item");
      summaryField.setAccessible(true);
    }
    mutations();
    tokenGates();
    duplicateHeads();
    sharingAndCloning();
    sourceTransfers();
    deserialization();
    requireThrows(IllegalStateException.class,
        () -> create(INT_ITEM, (byte) Node.Prop.ARROW_FN.ordinal(), 0, null));
    requireThrows(NullPointerException.class,
        () -> create(OBJECT_ITEM, (byte) Node.Prop.SYNTHETIC.ordinal(), null, null));
    System.out.println("observations=" + observations);
    System.out.println("NodeHotPropSummaryProbe OK");
  }

  private static Node function() {
    return IR.function(IR.name("f"), IR.paramList(IR.name("x")), IR.block());
  }

  private static void mutations() {
    Node n = function();
    observe("initial", n);
    for (int value : new int[] {1, 2, -1, Integer.MIN_VALUE, Integer.MAX_VALUE, 0}) {
      n.putIntProp(Node.Prop.ARROW_FN, value);
      n.putIntProp(Node.Prop.SYNTHETIC, value);
      observe("integer/" + value, n);
      n.putIntProp(Node.Prop.CHANGE_TIME, 7);
      n.putProp(Node.Prop.INPUT_ID, new InputId("input"));
      observe("unrelated-prefix/" + value, n);
      n.removeProp(Node.Prop.CHANGE_TIME);
      n.removeProp(Node.Prop.INPUT_ID);
      n.removeProp(Node.Prop.ASYNC_FN); // missing-property rebuild retains the original list
      observe("remove-prefix/" + value, n);
    }
    n.setIsArrowFunction(true);
    n.putBooleanProp(Node.Prop.SYNTHETIC, true);
    Object marker = new Object();
    n.putProp(Node.Prop.ARROW_FN, marker);
    require(n.getProp(Node.Prop.ARROW_FN) == marker, "object identity");
    observe("object-arrow", n);
    n.putProp(Node.Prop.SYNTHETIC, marker);
    observe("object-both", n);
    n.putProp(Node.Prop.ARROW_FN, null);
    observe("remove-object-arrow", n);
    n.putIntProp(Node.Prop.SYNTHETIC, -9);
    observe("replace-object-synthetic", n);
    n.removeProp(Node.Prop.SYNTHETIC);
    n.putBooleanProp(Node.Prop.ARROW_FN, false);
    observe("cleared", n);
    Node block = IR.block();
    block.setIsSyntheticBlock(true);
    observe("synthetic-setter-true", block);
    block.setIsSyntheticBlock(false);
    observe("synthetic-setter-false", block);
    requireThrows(IllegalStateException.class, () -> block.setIsArrowFunction(true));
    requireThrows(IllegalStateException.class, () -> n.setIsSyntheticBlock(true));
    observe("setter-errors-preserve-properties", n);
  }

  private static void tokenGates() {
    Node n = function();
    n.putIntProp(Node.Prop.ARROW_FN, -1);
    n.putIntProp(Node.Prop.SYNTHETIC, 2);
    for (Token token : Token.values()) {
      n.setToken(token);
      observe("token/" + token.name(), n);
    }
    n.setToken(null);
    observe("token/null", n);
    n.setToken(Token.FUNCTION);
    observe("token/function-again", n);
    n.putProp(Node.Prop.ARROW_FN, new Object());
    n.setToken(Token.NAME);
    require(!n.isArrowFunction(), "nonfunction must skip malformed arrow property");
    n.setToken(Token.FUNCTION);
    requireThrows(UnsupportedOperationException.class, n::isArrowFunction);
    observe("token/object-arrow-restored", n);
  }

  private static void duplicateHeads() {
    Object marker = new Object();
    Object tail = intItem(Node.Prop.SYNTHETIC, -7, null);
    Object invalidArrow = objectItem(Node.Prop.ARROW_FN, marker, tail);
    Node n = function();
    setHead(n, intItem(Node.Prop.CHANGE_TIME, 5, intItem(Node.Prop.ARROW_FN, 1, invalidArrow)));
    observe("duplicate/int-before-object", n);
    n.removeProp(Node.Prop.ARROW_FN);
    observe("duplicate/reveal-object", n);
    n.removeProp(Node.Prop.ARROW_FN);
    observe("duplicate/remove-second-arrow", n);
    setHead(n, objectItem(Node.Prop.ARROW_FN, marker, intItem(Node.Prop.ARROW_FN, 4, tail)));
    observe("duplicate/object-before-int", n);
    n.removeProp(Node.Prop.ARROW_FN);
    observe("duplicate/reveal-int", n);
    setHead(n, objectItem(Node.Prop.SYNTHETIC, marker,
        intItem(Node.Prop.SYNTHETIC, 1, intItem(Node.Prop.ARROW_FN, -1, null))));
    observe("duplicate/synthetic-object-first", n);
    n.removeProp(Node.Prop.SYNTHETIC);
    observe("duplicate/synthetic-reveal-int", n);
    n.putIntProp(Node.Prop.CHANGE_TIME, 8);
    n.putProp(Node.Prop.INPUT_ID, new InputId("duplicate"));
    n.removeProp(Node.Prop.SYNTHETIC); // reconstruct both unrelated prefix kinds
    observe("duplicate/rebuild-mixed-prefix", n);
    setHead(n, null);
    observe("raw-head-clear", n);
  }

  private static void sharingAndCloning() {
    Node source = function();
    source.setIsArrowFunction(true);
    source.setSourceFileForTesting("clone.js");
    JSDocInfo.Builder builder = JSDocInfo.builder();
    require(builder.recordType(new JSTypeExpression(Node.newString("T"), "clone.js")), "type fixture");
    source.setJSDocInfo(builder.build());
    source.getLastChild().addChildToBack(IR.block());
    source.getLastChild().getFirstChild().setIsSyntheticBlock(true);
    observe("clone/source", source);
    Node props = function();
    observe("clone/empty-destination", props);
    props.clonePropsFrom(source);
    require(head(props) == head(source), "property list sharing");
    observe("clone/props", props);
    requireThrows(IllegalStateException.class, () -> props.clonePropsFrom(source));
    Node shallow = source.cloneNode();
    require(head(shallow) == head(source), "shallow clone list sharing");
    observe("clone/shallow", shallow);
    Node tree = source.cloneTree();
    Node typedTree = source.cloneTree(true);
    require(source.isEquivalentTo(tree) && source.isEquivalentTo(typedTree), "tree equivalence");
    require(source.toStringTree().equals(tree.toStringTree()), "tree debug text");
    require(typedTree.getJSDocInfo().getType().getRoot() != source.getJSDocInfo().getType().getRoot(),
        "type expression cloned");
    observeTree("clone/tree", tree);
    observeTree("clone/type-tree", typedTree);
    source.setIsArrowFunction(false);
    props.putProp(Node.Prop.ARROW_FN, new Object());
    observe("clone/source-mutated", source);
    observe("clone/shared-copy-mutated", props);
    observe("clone/shallow-independent", shallow);
    observeTree("clone/tree-independent", tree);
    require(shallow.isArrowFunction() && tree.isArrowFunction(), "shared immutable list isolation");

    Node[] specialized = {Node.newNumber(1.5), Node.newBigInt(BigInteger.TEN),
        Node.newString(Token.NAME, "name"), Node.newTemplateLitString("cooked", "raw")};
    for (int i = 0; i < specialized.length; i++) {
      Node n = specialized[i];
      n.putIntProp(Node.Prop.ARROW_FN, 2);
      n.putBooleanProp(Node.Prop.SYNTHETIC, true);
      Node clone = n.cloneNode();
      require(clone.getClass() == n.getClass() && head(clone) == head(n), "specialized clone");
      observe("specialized/original/" + i, n);
      observe("specialized/clone/" + i, clone);
      n.removeProp(Node.Prop.SYNTHETIC);
      observe("specialized/original-mutated/" + i, n);
      observe("specialized/clone-independent/" + i, clone);
    }
    Node holder = IR.name("shadow");
    Node shadowFunction = function();
    shadowFunction.setIsArrowFunction(true);
    holder.setClosureUnawareShadow(IR.root(IR.script(shadowFunction)));
    Node clonedHolder = holder.cloneTree();
    require(holder.getClosureUnawareShadow() != clonedHolder.getClosureUnawareShadow(), "shadow cloned");
    observeTree("clone/shadow", clonedHolder.getClosureUnawareShadow());
  }

  private static void sourceTransfers() {
    Node source = function();
    source.setSourceFileForTesting("source.js");
    source.setIsArrowFunction(true);
    Node empty = function();
    empty.setStaticSourceFileFrom(source); // direct shared-tail assignment
    require(Node.validateMemorySensitivePropertyGuarantees(source, false, empty, false), "shared tail");
    require(!empty.isArrowFunction(), "source transfer must not copy arrow flag");
    observe("source/direct-tail", empty);
    Node sourceOnly = function();
    sourceOnly.setSourceFileForTesting("old.js");
    sourceOnly.setStaticSourceFileFrom(source);
    observe("source/replace-only-tail", sourceOnly);
    Node decorated = function();
    decorated.setIsArrowFunction(true);
    decorated.setStaticSourceFileFrom(source); // generic property replacement path
    require(decorated.getStaticSourceFile() == source.getStaticSourceFile(), "source file identity");
    observe("source/decorated", decorated);
    decorated.setStaticSourceFile(null);
    observe("source/removed", decorated);
    decorated.srcref(source);
    observe("source/srcref", decorated);
    Node emptySource = function();
    decorated.setStaticSourceFileFrom(emptySource);
    observe("source/absent-source", decorated);
  }

  private static void deserialization() {
    Node arrow = function();
    arrow.setSourceFileForTesting("deserialize.js");
    observe("deserialize/before", arrow);
    long arrowBits = bit(NodeProperty.ARROW_FN) | bit(NodeProperty.MUTATES_GLOBAL_STATE)
        | bit(NodeProperty.THROWS);
    arrow.deserializeProperties(arrowBits, false);
    require(arrow.serializeProperties() == arrowBits, "arrow bits roundtrip");
    observe("deserialize/arrow-side-effects", arrow);
    Node block = IR.block();
    block.setSourceFileForTesting("deserialize.js");
    block.deserializeProperties(bit(NodeProperty.SYNTHETIC), false);
    observe("deserialize/synthetic", block);
    Node name = IR.name("constant");
    name.setSourceFileForTesting("deserialize.js");
    long constants = bit(NodeProperty.IS_DECLARED_CONSTANT) | bit(NodeProperty.IS_INFERRED_CONSTANT);
    name.deserializeProperties(constants, false);
    require(name.serializeProperties() == constants, "constant flags roundtrip");
    observe("deserialize/constant-flags", name);
    Node closureUnaware = function();
    closureUnaware.setSourceFileForTesting("closure-unaware.js");
    closureUnaware.setIsInClosureUnawareSubtree(true);
    closureUnaware.deserializeProperties(bit(NodeProperty.ARROW_FN), true);
    observe("deserialize/closure-unaware", closureUnaware);
    Node invalidAfterWrite = function();
    invalidAfterWrite.setSourceFileForTesting("invalid.js");
    requireThrows(IllegalStateException.class,
        () -> invalidAfterWrite.deserializeProperties(bit(NodeProperty.SYNTHETIC), false));
    require(invalidAfterWrite.isSyntheticBlock(), "post-validation failure retains mutation");
    observe("deserialize/validation-failure", invalidAfterWrite);
    Node partial = function();
    partial.setSourceFileForTesting("partial.js");
    int unknown = -1;
    for (int i = NodeProperty.ARROW_FN.getNumber() + 1; i < 63; i++) {
      if (NodeProperty.forNumber(i) == null) {
        unknown = i;
        break;
      }
    }
    require(unknown >= 0, "unknown higher bit fixture");
    long partialBits = bit(NodeProperty.ARROW_FN) | (1L << unknown);
    requireThrows(NullPointerException.class, () -> partial.deserializeProperties(partialBits, false));
    require(partial.isArrowFunction(), "translation failure retains earlier arrow mutation");
    observe("deserialize/partial-failure", partial);
    Node missingSource = function();
    requireThrows(IllegalStateException.class,
        () -> missingSource.deserializeProperties(bit(NodeProperty.ARROW_FN), false));
    observe("deserialize/precondition-failure", missingSource);
  }

  private static long bit(NodeProperty property) {
    return 1L << property.getNumber();
  }

  private static void observeTree(String label, Node n) {
    observe(label, n);
    int index = 0;
    for (Node child = n.getFirstChild(); child != null; child = child.getNext()) {
      observeTree(label + "/" + index++, child);
    }
  }

  private static void observe(String label, Node n) {
    String properties = n.getPropListDebugString();
    long serialized = n.serializeProperties();
    List<String> violations = new ArrayList<>();
    n.validateProperties(violations::add);
    String arrow = outcome(() -> n.isFunction() && n.getBooleanProp(Node.Prop.ARROW_FN));
    String synthetic = outcome(() -> n.getBooleanProp(Node.Prop.SYNTHETIC));
    for (int i = 0; i < 3; i++) {
      require(outcome(n::isArrowFunction).equals(arrow), label + " arrow result/error");
      require(outcome(n::isSyntheticBlock).equals(synthetic), label + " synthetic result/error");
      require(properties.equals(n.getPropListDebugString()), label + " property-list mutation");
      require(serialized == n.serializeProperties(), label + " serialized-property mutation");
    }
    List<String> afterViolations = new ArrayList<>();
    n.validateProperties(afterViolations::add);
    require(violations.equals(afterViolations), label + " property-validation mutation");
    if (summaryField != null && head(n) != null) {
      try {
        int expected = rawState(n, Node.Prop.ARROW_FN) | (rawState(n, Node.Prop.SYNTHETIC) << 2);
        require(summaryField.getByte(head(n)) == expected, label + " head suffix summary");
      } catch (IllegalAccessException e) {
        throw new AssertionError(e);
      }
    }
    observations++;
    System.out.println(label + " token=" + n.getToken() + " arrow=" + arrow
        + " synthetic=" + synthetic + " bits=" + serialized + " props=" + properties
        + " violations=" + violations);
  }

  private static int rawState(Node n, Node.Prop prop) {
    try {
      return n.getBooleanProp(prop) ? 1 : 0;
    } catch (UnsupportedOperationException expected) {
      return 2;
    }
  }

  private static String outcome(BooleanSupplier operation) {
    try {
      return Boolean.toString(operation.getAsBoolean());
    } catch (RuntimeException e) {
      return "throws:" + e.getClass().getName();
    }
  }

  private static Object head(Node node) {
    return invoke(GET_HEAD, node);
  }

  private static void setHead(Node node, Object head) {
    invoke(SET_HEAD, node, head);
  }

  private static Object intItem(Node.Prop prop, int value, Object next) {
    return create(INT_ITEM, (byte) prop.ordinal(), value, next);
  }

  private static Object objectItem(Node.Prop prop, Object value, Object next) {
    return create(OBJECT_ITEM, (byte) prop.ordinal(), value, next);
  }

  private static Class<?> nestedClass(String name) {
    try {
      return Class.forName(Node.class.getName() + "$" + name);
    } catch (ClassNotFoundException e) {
      throw new AssertionError(e);
    }
  }

  private static Method method(String name, Class<?>... parameters) {
    try {
      Method result = Node.class.getDeclaredMethod(name, parameters);
      result.setAccessible(true);
      return result;
    } catch (ReflectiveOperationException e) {
      throw new AssertionError(e);
    }
  }

  private static Constructor<?> constructor(String name, Class<?>... parameters) {
    try {
      Constructor<?> result = nestedClass(name).getDeclaredConstructor(parameters);
      result.setAccessible(true);
      return result;
    } catch (ReflectiveOperationException e) {
      throw new AssertionError(e);
    }
  }

  private static Object invoke(Method method, Object receiver, Object... args) {
    try {
      return method.invoke(receiver, args);
    } catch (InvocationTargetException e) {
      throw cause(e);
    } catch (ReflectiveOperationException e) {
      throw new AssertionError(e);
    }
  }

  private static Object create(Constructor<?> constructor, Object... args) {
    try {
      return constructor.newInstance(args);
    } catch (InvocationTargetException e) {
      throw cause(e);
    } catch (ReflectiveOperationException e) {
      throw new AssertionError(e);
    }
  }

  private static RuntimeException cause(InvocationTargetException e) {
    if (e.getCause() instanceof RuntimeException runtime) {
      return runtime;
    }
    throw new AssertionError(e.getCause());
  }

  private static void requireThrows(Class<? extends RuntimeException> type, Runnable operation) {
    try {
      operation.run();
      throw new AssertionError("Expected " + type.getName());
    } catch (RuntimeException e) {
      require(type.isInstance(e), "Expected " + type.getName() + ", got " + e);
    }
  }

  private static void require(boolean condition, String message) {
    if (!condition) {
      throw new AssertionError(message);
    }
  }
}

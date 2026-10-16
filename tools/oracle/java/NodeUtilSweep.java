// Oracle tool: an API sweep of com.google.javascript.jscomp.NodeUtil over real ASTs. Every file is
// parsed with ParserRunner.parse and the parsing.Config the compiler uses
// (Compiler.getParserConfig(ConfigContext.DEFAULT), language ECMASCRIPT_NEXT); then, for every node
// of the AST in pre-order, every NodeUtil predicate/getter whose arguments are that node (plus
// fixed constants, or the node's parent, or a `new Compiler()` initialized with default
// CompilerOptions) is called, and the results are printed on one line per node in a fixed order.
// The Go side is internal/jscomp/node_util_sweep_test.go; both must print byte-identical text.
//
// It lives in package com.google.javascript.jscomp to reach NodeUtil's package-private API
// (private methods are called reflectively).
//
// Build/run:
//   javac -cp closure-compiler.jar -d build NodeUtilSweep.java
//   java -Xss1g -cp closure-compiler.jar:build com.google.javascript.jscomp.NodeUtilSweep
//        [--gzip] <out-file> <root> <rel-file>... | @<listfile>
//   java ... NodeUtilSweep [--gzip] --snippets=<file> <out-file>
// A snippets file holds many small sources: a line "#### <name>" starts the source <name>, made of
// every following line (each with its '\n') up to the next such line.
//
// Format (all files into one output, in argument order):
//   # file <rel-file>
//   # exception <class>                         if ParserRunner.parse threw
//   # tree null                                  if there is no AST
//   <index> <TOKEN> <bools> <field>...           one line per node, pre-order index from 0,
//                                                fields separated by TAB
//   # file-summary <field>...                    per-file sweeps over the whole script
// <bools> is one character per boolean probe: '1' true, '0' false, '.' not applicable (the probe
// is guarded), or an exception code: S IllegalStateException, A IllegalArgumentException,
// N NullPointerException, U UnsupportedOperationException, C ClassCastException,
// I IndexOutOfBoundsException (and subclasses), E Error, X anything else.
// A field is '-' for null, '.' for not applicable, '!<code>' for an exception, a node index
// ('?' for a node outside the main tree), "[i,j,...]" for node lists, a quoted string, or a
// number/name. Every field is escaped like ParsingDump (backslash, and every char outside
// printable ASCII as \\uXXXX; a '"' inside a quoted string is preceded by a backslash first).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.base.Tri;
import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.jscomp.parsing.parser.FeatureSet;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.InputId;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.JSTypeExpression;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.QualifiedName;
import com.google.javascript.rhino.StaticSourceFile;
import com.google.javascript.rhino.Token;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.IdentityHashMap;
import java.util.List;
import java.util.zip.GZIPOutputStream;

public final class NodeUtilSweep {
  private static final QualifiedName GOOG_MODULE_GET = QualifiedName.of("goog.module.get");
  private static final Node GOOG_PROVIDE_NODE = IR.getprop(IR.name("goog"), "provide");

  private static Method isDeclarationParent;
  private static Method isKnownGlobalObjectReference;
  private static Method isRhsClass;
  private static Method isRhsFunction;
  private static Method isToStringMethodCall;
  private static Method getEnclosingTarget;
  private static Method getNonOptChainToken;
  private static Method isPureIterable;

  /** Files with more nodes skip the mutation probes, which clone the whole script per probe. */
  private static final int MAX_NODES_FOR_MUTATIONS = 1000;

  private IdentityHashMap<Node, Integer> index;
  private AbstractCompiler compiler;
  private Node script;
  private boolean mutate;

  public static void main(String[] args) throws Exception {
    boolean gzip = false;
    String snippets = null;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--gzip")) {
        gzip = true;
      } else if (a.startsWith("--snippets=")) {
        snippets = a.substring("--snippets=".length());
      } else {
        rest.add(a);
      }
    }
    Path out = Paths.get(rest.get(0));
    List<String> files = new ArrayList<>();
    List<String> sources = new ArrayList<>();
    if (snippets != null) {
      splitSnippets(Files.readString(Paths.get(snippets), StandardCharsets.UTF_8), files, sources);
    } else {
      Path root = Paths.get(rest.get(1));
      for (String f : rest.subList(2, rest.size())) {
        if (f.startsWith("@")) {
          for (String l : Files.readAllLines(Paths.get(f.substring(1)), StandardCharsets.UTF_8)) {
            if (!l.isEmpty()) {
              files.add(l);
            }
          }
        } else {
          files.add(f);
        }
      }
      for (String rel : files) {
        sources.add(null);
      }
      for (int i = 0; i < files.size(); i++) {
        sources.set(i, Files.readString(root.resolve(files.get(i)), StandardCharsets.UTF_8));
      }
    }
    isDeclarationParent = privateMethod("isDeclarationParent", Node.class);
    isKnownGlobalObjectReference = privateMethod("isKnownGlobalObjectReference", Node.class);
    isRhsClass = privateMethod("isRhsClass", Node.class);
    isRhsFunction = privateMethod("isRhsFunction", Node.class);
    isToStringMethodCall = privateMethod("isToStringMethodCall", Node.class);
    getEnclosingTarget = privateMethod("getEnclosingTarget", Node.class);
    getNonOptChainToken = privateMethod("getNonOptChainToken", Token.class);
    isPureIterable = privateMethod("isPureIterable", Node.class);

    Config config = compilerParserConfig();
    try (OutputStream os =
        gzip ? new GZIPOutputStream(Files.newOutputStream(out)) : Files.newOutputStream(out)) {
      for (int i = 0; i < files.size(); i++) {
        StringBuilder sb = new StringBuilder();
        new NodeUtilSweep().sweepFile(sb, config, files.get(i), sources.get(i));
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  /**
   * Splits a snippets file: a line "#### <name>" starts the snippet <name>, whose source is every
   * following line (each with its '\\n') up to the next such line.
   */
  static void splitSnippets(String text, List<String> names, List<String> sources) {
    StringBuilder cur = null;
    String[] lines = text.split("\n", -1);
    for (int i = 0; i < lines.length; i++) {
      String l = lines[i];
      if (i == lines.length - 1 && l.isEmpty()) {
        break;
      }
      if (l.startsWith("#### ")) {
        if (cur != null) {
          sources.add(cur.toString());
        }
        names.add(l.substring("#### ".length()));
        cur = new StringBuilder();
      } else if (cur != null) {
        cur.append(l).append('\n');
      }
    }
    if (cur != null) {
      sources.add(cur.toString());
    }
  }

  private static Method privateMethod(String name, Class<?>... types) throws Exception {
    Method m = NodeUtil.class.getDeclaredMethod(name, types);
    m.setAccessible(true);
    return m;
  }

  /** The parsing.Config that CompilerInput passes to ParserRunner.parse for non-extern inputs. */
  private static Config compilerParserConfig() {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    compiler.initOptions(options);
    return compiler.getParserConfig(AbstractCompiler.ConfigContext.DEFAULT);
  }

  private static Compiler shared;

  /** `new Compiler()` + `initOptions(new CompilerOptions())`, as the Go test compiler. */
  private static Compiler sharedCompiler() {
    if (shared == null) {
      shared = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
      shared.initOptions(new CompilerOptions());
    }
    return shared;
  }

  void sweepFile(StringBuilder sb, Config config, String name, String src) {
    sb.append("# file ").append(escape(name)).append('\n');
    ErrorReporter reporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int line, int lineOffset) {}

          @Override
          public void error(String message, String sourceName, int line, int lineOffset) {}
        };
    ParserRunner.ParseResult result;
    try {
      result = ParserRunner.parse(SourceFile.fromCode(name, src), src, config, reporter);
    } catch (RuntimeException e) {
      sb.append("# exception ").append(e.getClass().getName()).append('\n');
      return;
    }
    if (result.ast == null) {
      sb.append("# tree null\n");
      return;
    }
    compiler = sharedCompiler();
    index = new IdentityHashMap<>();
    List<Node> nodes = new ArrayList<>();
    collect(result.ast, nodes);
    this.script = result.ast;
    this.mutate = nodes.size() <= MAX_NODES_FOR_MUTATIONS;
    for (int i = 0; i < nodes.size(); i++) {
      sweepNode(sb, i, nodes.get(i));
    }
    sweepScript(sb, result.ast);
  }

  private void collect(Node n, List<Node> nodes) {
    index.put(n, nodes.size());
    nodes.add(n);
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      collect(c, nodes);
    }
  }

  // ---------------------------------------------------------------------------------------------
  // Encoding

  private interface BoolProbe {
    boolean get() throws Exception;
  }

  private interface Probe {
    String get() throws Exception;
  }

  private static char code(Throwable e) {
    if (e instanceof InvocationTargetException) {
      e = e.getCause();
    }
    if (e instanceof IllegalStateException) {
      return 'S';
    } else if (e instanceof IllegalArgumentException) {
      return 'A';
    } else if (e instanceof NullPointerException) {
      return 'N';
    } else if (e instanceof UnsupportedOperationException) {
      return 'U';
    } else if (e instanceof ClassCastException) {
      return 'C';
    } else if (e instanceof IndexOutOfBoundsException) {
      return 'I';
    } else if (e instanceof Error) {
      return 'E';
    }
    return 'X';
  }

  private static void bool(StringBuilder sb, boolean applicable, BoolProbe p) {
    if (!applicable) {
      sb.append('.');
      return;
    }
    try {
      sb.append(p.get() ? '1' : '0');
    } catch (Exception | Error e) {
      sb.append(code(e));
    }
  }

  private static void field(StringBuilder sb, boolean applicable, Probe p) {
    sb.append('\t');
    if (!applicable) {
      sb.append('.');
      return;
    }
    String s;
    try {
      s = p.get();
    } catch (Exception | Error e) {
      s = "!" + code(e);
    }
    sb.append(escape(s));
  }

  private String idx(Node n) {
    if (n == null) {
      return "-";
    }
    Integer i = index.get(n);
    return i == null ? "?" : Integer.toString(i);
  }

  private String list(Iterable<Node> nodes) {
    StringBuilder b = new StringBuilder("[");
    boolean first = true;
    for (Node n : nodes) {
      if (!first) {
        b.append(',');
      }
      first = false;
      b.append(idx(n));
    }
    return b.append(']').toString();
  }

  private static String str(String s) {
    return s == null ? "-" : quote(s);
  }

  private static String b(boolean v) {
    return v ? "1" : "0";
  }

  private static String dbl(Double d) {
    return d == null ? "-" : Double.toString(d);
  }

  private static String big(BigInteger x) {
    return x == null ? "-" : x.toString();
  }

  private static String tri(Tri t) {
    return t.name();
  }

  private static String tok(Token t) {
    return t == null ? "-" : t.name();
  }

  private static String quote(String s) {
    StringBuilder b = new StringBuilder(s.length() + 2).append('"');
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '"') {
        b.append("\\\"");
      } else {
        b.append(c);
      }
    }
    return b.append('"').toString();
  }

  /** Backslash and every char outside printable ASCII become escapes (as in ParsingDump). */
  static String escape(String s) {
    StringBuilder b = new StringBuilder(s.length());
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        b.append("\\\\");
      } else if (c >= 0x20 && c < 0x7f) {
        b.append(c);
      } else {
        b.append(String.format("\\u%04x", (int) c));
      }
    }
    return b.toString();
  }

  /**
   * A one-line rendering of a tree: TOKEN[string]{length@line:charno}c(if IS_CONSTANT_NAME)
   * o=<original name>(children).
   */
  private static String compact(Node n) {
    StringBuilder b = new StringBuilder(n.getToken().name());
    if (n.isName() || n.isGetProp() || n.isStringLit() || n.isStringKey()) {
      b.append('[').append(quote(n.getString())).append(']');
    }
    b.append('{')
        .append(n.getLength())
        .append('@')
        .append(n.getLineno())
        .append(':')
        .append(n.getCharno())
        .append('}');
    if (n.getBooleanProp(Node.IS_CONSTANT_NAME)) {
      b.append('c');
    }
    if (n.getOriginalName() != null) {
      b.append("o=").append(quote(n.getOriginalName()));
    }
    if (n.hasChildren()) {
      b.append('(');
      for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
        if (c != n.getFirstChild()) {
          b.append(',');
        }
        b.append(compact(c));
      }
      b.append(')');
    }
    return b.toString();
  }

  private static boolean hasString(Node n) {
    return n.isName() || n.isStringKey() || n.isGetProp() || n.isStringLit();
  }

  private static boolean isObjectOrClassLitKeyToken(Node n) {
    switch (n.getToken()) {
      case STRING_KEY:
      case GETTER_DEF:
      case SETTER_DEF:
      case MEMBER_FUNCTION_DEF:
      case MEMBER_FIELD_DEF:
      case COMPUTED_PROP:
      case COMPUTED_FIELD_DEF:
        return true;
      default:
        return false;
    }
  }

  private static String firstKey(Node n) {
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      if (c.isStringKey() || c.isMemberFunctionDef() || c.isGetterDef() || c.isSetterDef()) {
        return c.getString();
      }
    }
    return "x";
  }

  private static boolean isAssigningParent(Node n) {
    return NodeUtil.isNameDeclaration(n)
        || n.isParamList()
        || NodeUtil.isAssignmentOp(n)
        || n.isCatch()
        || n.isDestructuringLhs()
        || n.isDefaultValue()
        || n.isImport()
        || NodeUtil.isEnhancedFor(n);
  }

  private interface Mutation {
    String apply(Node clone, Node cloneOfN) throws Exception;
  }

  /**
   * Applies a mutation to the node with the same pre-order index in a fresh clone of the script,
   * then returns the mutation's own result, ':' and the hash of the compact rendering of the clone.
   */
  private String mutated(int i, Mutation m) throws Exception {
    Node clone = script.cloneTree();
    List<Node> cloneNodes = new ArrayList<>();
    collectInto(clone, cloneNodes);
    String r = m.apply(clone, cloneNodes.get(i));
    return r + ":" + compact(clone).hashCode();
  }

  private static void collectInto(Node n, List<Node> nodes) {
    nodes.add(n);
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      collectInto(c, nodes);
    }
  }

  // ---------------------------------------------------------------------------------------------
  // The sweep

  void sweepNode(StringBuilder sb, int i, Node n) {
    Node parent = n.getParent();
    Token t = n.getToken();
    StringBuilder line = new StringBuilder();
    line.append(i).append(' ').append(t.name()).append(' ');

    boolean isFn = n.isFunction();
    boolean isFnOrClass = n.isFunction() || n.isClass();
    boolean hasParent = parent != null;
    boolean isName = n.isName();
    boolean isGet = NodeUtil.isNormalOrOptChainGet(n);
    boolean isQName = n.isQualifiedName();

    bool(line, true, () -> NodeUtil.isImmutableValue(n));
    bool(line, true, () -> NodeUtil.isSymmetricOperation(n));
    bool(line, true, () -> NodeUtil.isRelationalOperation(n));
    bool(line, true, () -> NodeUtil.isLiteralValue(n, false));
    bool(line, true, () -> NodeUtil.isLiteralValue(n, true));
    bool(line, true, () -> NodeUtil.isSomeCompileTimeConstStringValue(n));
    bool(line, true, () -> NodeUtil.isEmptyBlock(n));
    bool(line, true, () -> NodeUtil.isBinaryOperator(n));
    bool(line, true, () -> NodeUtil.isUnaryOperator(n));
    bool(line, true, () -> NodeUtil.isUpdateOperator(n));
    bool(line, true, () -> NodeUtil.isSimpleOperator(n));
    bool(line, hasParent, () -> NodeUtil.isNamespaceDecl(n));
    bool(line, true, () -> (boolean) isPureIterable.invoke(null, n));
    bool(line, true, () -> NodeUtil.canBeSideEffected(n));
    bool(line, true, () -> NodeUtil.isUndefined(n));
    bool(line, true, () -> NodeUtil.isNullOrUndefined(n));
    bool(line, true, () -> NodeUtil.isNumericResult(n));
    bool(line, true, () -> NodeUtil.isBigIntResult(n));
    bool(line, true, () -> NodeUtil.isBooleanResult(n));
    bool(line, true, () -> NodeUtil.isStringResult(n));
    bool(line, true, () -> NodeUtil.isObjectResult(n));
    bool(line, true, () -> NodeUtil.mayBeString(n));
    bool(line, true, () -> NodeUtil.mayBeString(n, true));
    bool(line, true, () -> NodeUtil.mayBeObject(n));
    bool(line, true, () -> NodeUtil.isAssociative(t));
    bool(line, true, () -> NodeUtil.isCommutative(t));
    bool(line, true, () -> NodeUtil.isAssignmentOp(n));
    bool(line, true, () -> NodeUtil.isLogicalAssignmentOp(n));
    bool(line, true, () -> NodeUtil.isCompoundAssignmentOp(n));
    bool(line, true, () -> NodeUtil.isHoistScopeRoot(n));
    bool(line, true, () -> NodeUtil.isInFunction(n));
    bool(line, true, () -> NodeUtil.referencesEnclosingReceiver(n));
    bool(line, true, () -> NodeUtil.referencesSuper(n));
    bool(line, true, () -> NodeUtil.isNormalOrOptChainGet(n));
    bool(line, true, () -> NodeUtil.isNormalOrOptChainGetProp(n));
    bool(line, true, () -> NodeUtil.isNormalOrOptChainCall(n));
    bool(line, true, () -> NodeUtil.isNormalGet(n));
    bool(line, true, () -> NodeUtil.isOptChainGet(n));
    bool(line, true, () -> NodeUtil.isOptChainNode(n));
    bool(line, true, () -> NodeUtil.isEndOfFullOptChain(n));
    bool(line, true, () -> NodeUtil.isEndOfOptChainSegment(n));
    bool(line, true, () -> NodeUtil.isBlockScopedDeclaration(n));
    bool(line, true, () -> NodeUtil.isNameDeclaration(n));
    bool(line, true, () -> NodeUtil.isDestructuringDeclaration(n));
    bool(line, true, () -> NodeUtil.isExprAssign(n));
    bool(line, true, () -> NodeUtil.isExprCall(n));
    bool(line, true, () -> NodeUtil.isNonArrowFunction(n));
    bool(line, true, () -> NodeUtil.isEnhancedFor(n));
    bool(line, true, () -> NodeUtil.isAnyFor(n));
    bool(line, true, () -> NodeUtil.isLoopStructure(n));
    bool(line, true, () -> NodeUtil.isWithinLoop(n));
    bool(line, true, () -> NodeUtil.isControlStructure(n));
    bool(line, true, () -> NodeUtil.isStatementBlock(n));
    bool(line, true, () -> NodeUtil.createsBlockScope(n));
    bool(line, true, () -> NodeUtil.createsScope(n));
    bool(line, true, () -> NodeUtil.isValidCfgRoot(n));
    bool(line, true, () -> NodeUtil.isStatement(n));
    bool(line, true, () -> NodeUtil.isStatementParent(n));
    bool(line, true, () -> (boolean) isDeclarationParent.invoke(null, n));
    bool(line, true, () -> NodeUtil.isSwitchCase(n));
    bool(line, true, () -> NodeUtil.isReferenceName(n));
    bool(line, true, () -> NodeUtil.isInSyntheticScript(n));
    bool(line, true, () -> NodeUtil.isCallOrNew(n));
    bool(line, true, () -> NodeUtil.isDeclaration(n));
    bool(line, true, () -> NodeUtil.isNamedExportsLiteral(n));
    bool(line, true, () -> NodeUtil.isFunctionDeclaration(n));
    bool(line, true, () -> NodeUtil.isMethodDeclaration(n));
    bool(line, true, () -> NodeUtil.isClassDeclaration(n));
    bool(line, true, () -> NodeUtil.isHoistedFunctionDeclaration(n));
    bool(line, true, () -> NodeUtil.isBlockScopedFunctionDeclaration(n));
    bool(line, true, () -> NodeUtil.isFunctionBlock(n));
    bool(line, true, () -> NodeUtil.isClassStaticBlock(n));
    bool(line, true, () -> NodeUtil.isFunctionExpression(n));
    bool(line, true, () -> NodeUtil.isNamedFunctionExpression(n));
    bool(line, true, () -> NodeUtil.isClassExpression(n));
    bool(line, true, () -> NodeUtil.isNamedClassExpression(n));
    bool(line, true, () -> NodeUtil.isNamedFunction(n));
    bool(line, true, () -> NodeUtil.isNamedClass(n));
    bool(line, true, () -> NodeUtil.isBleedingFunctionName(n));
    bool(line, true, () -> NodeUtil.isEmptyFunctionExpression(n));
    bool(line, true, () -> NodeUtil.isFunctionObjectCall(n));
    bool(line, true, () -> NodeUtil.isFunctionObjectApply(n));
    bool(line, true, () -> NodeUtil.isLValue(n));
    bool(line, true, () -> NodeUtil.isDeclarationLValue(n));
    bool(line, true, () -> NodeUtil.isLhsOfAssign(n));
    bool(line, hasParent, () -> NodeUtil.isImportedName(n));
    bool(line, true, () -> NodeUtil.isLhsByDestructuring(n));
    bool(line, true, () -> NodeUtil.mayBeObjectLitKey(n));
    bool(line, hasParent, () -> NodeUtil.isObjectLitKey(n));
    bool(line, true, () -> NodeUtil.isGetOrSetKey(n));
    bool(line, true, () -> NodeUtil.canContainHoistedVarsDecls(n));
    bool(line, true, () -> NodeUtil.isObjectDefinePropertiesDefinition(n));
    bool(line, true, () -> NodeUtil.isObjectDefinePropertyDefinition(n));
    bool(line, true, () -> NodeUtil.isPrototypePropertyDeclaration(n));
    bool(line, true, () -> NodeUtil.isPrototypeProperty(n));
    bool(line, true, () -> NodeUtil.isPrototypeMethod(n));
    bool(line, true, () -> NodeUtil.isPrototypeAssignment(n));
    bool(line, hasParent, () -> NodeUtil.isPropertyAbsenceTest(n));
    bool(line, true, () -> NodeUtil.isShallowStatementTree(n));
    bool(line, true, () -> new NodeUtil.MatchShallowStatement().apply(n));
    bool(line, true, () -> new NodeUtil.MatchDeclaration().apply(n));
    bool(line, true, () -> NodeUtil.isConstantVar(n, null));
    bool(line, true, () -> NodeUtil.isConstantName(n));
    bool(line, true, () -> NodeUtil.functionHasInlineJsdocs(n));
    bool(line, true, () -> NodeUtil.evaluatesToLocalValue(n));
    bool(line, true, () -> NodeUtil.mayBeUndefined(n));
    bool(line, true, () -> NodeUtil.isDefinedValue(n));
    bool(line, true, () -> NodeUtil.isInvocationTarget(n));
    bool(line, true, () -> NodeUtil.isInvocation(n));
    bool(line, true, () -> NodeUtil.isCallOrNewArgument(n));
    bool(line, hasParent, () -> NodeUtil.isExpressionResultUsed(n));
    bool(line, hasParent, () -> NodeUtil.isExecutedExactlyOnce(n));
    bool(line, true, () -> NodeUtil.isNaN(n));
    bool(line, true, () -> NodeUtil.isGoogProvideCall(n));
    bool(line, true, () -> NodeUtil.isGoogModuleCall(n));
    bool(line, true, () -> NodeUtil.isGoogModuleGetCall(n));
    bool(line, true, () -> NodeUtil.isGoogRequireCall(n));
    bool(line, true, () -> NodeUtil.isGoogRequireTypeCall(n));
    bool(line, true, () -> NodeUtil.isGoogForwardDeclareCall(n));
    bool(line, true, () -> NodeUtil.isGoogRequireDynamicCall(n));
    bool(line, true, () -> NodeUtil.isGoogWeakUsageCall(n));
    bool(line, true, () -> NodeUtil.isModuleScopeRoot(n));
    bool(line, true, () -> NodeUtil.isBundledGoogModuleCall(n));
    bool(line, true, () -> NodeUtil.isBundledGoogModuleScopeRoot(n));
    bool(line, true, () -> NodeUtil.isGoogModuleDeclareLegacyNamespaceCall(n));
    bool(line, true, () -> NodeUtil.isGoogSetTestOnlyCall(n));
    bool(line, true, () -> NodeUtil.isTopLevel(n));
    bool(line, true, () -> NodeUtil.isGoogModuleFile(n));
    bool(line, true, () -> NodeUtil.isLegacyGoogModuleFile(n));
    bool(line, true, () -> NodeUtil.isConstructor(n));
    bool(line, true, () -> NodeUtil.isEs6ConstructorMemberFunctionDef(n));
    bool(line, true, () -> NodeUtil.isEs6Constructor(n));
    bool(line, true, () -> NodeUtil.isGetterOrSetter(n));
    bool(line, true, () -> NodeUtil.isCallTo(n, "goog.require"));
    bool(line, true, () -> NodeUtil.isCallTo(n, GOOG_MODULE_GET));
    bool(line, true, () -> NodeUtil.isCallTo(n, GOOG_PROVIDE_NODE));
    bool(line, true, () -> NodeUtil.isObjLitProperty(n));
    bool(line, true, () -> NodeUtil.isBlocklessArrowFunctionResult(n));
    bool(line, hasParent, () -> NodeUtil.isTryCatchNodeContainer(n));
    bool(line, hasParent, () -> NodeUtil.isTryFinallyNode(parent, n));
    bool(line, hasParent, () -> NodeUtil.isNameDeclOrSimpleAssignLhs(n, parent));
    bool(line, isName, () -> NodeUtil.isNonlocalModuleExportName(n));
    bool(line, isGet, () -> NodeUtil.isPropertyTest(compiler, n));
    bool(line, true, () -> (boolean) isKnownGlobalObjectReference.invoke(null, n));
    bool(line, true, () -> (boolean) isRhsClass.invoke(null, n));
    bool(line, true, () -> (boolean) isRhsFunction.invoke(null, n));
    bool(
        line,
        n.isCall() || n.isOptChainCall(),
        () -> (boolean) isToStringMethodCall.invoke(null, n));
    bool(line, isFn, () -> NodeUtil.doesFunctionReferenceOwnArgumentsObject(n));
    bool(line, isFn, () -> NodeUtil.referencesOwnReceiver(n));
    bool(line, n.isScript(), () -> NodeUtil.isFromTypeSummary(n));
    bool(line, true, () -> NodeUtil.canMergeBlock(n));
    bool(line, n.isExport(), () -> NodeUtil.isExportFrom(n));
    bool(line, n.isTry(), () -> NodeUtil.hasFinally(n));
    bool(line, n.isBlock(), () -> NodeUtil.hasCatchHandler(n));
    bool(
        line,
        n.isSpread() || n.isYield() || n.isForOf() || n.isForAwaitOf() || n.isRest(),
        () -> NodeUtil.iteratesImpureIterable(n));
    bool(line, n.isNew(), () -> NodeUtil.newHasLocalResult(n));
    bool(line, NodeUtil.isCallOrNew(n), () -> NodeUtil.allArgsUnescapedLocal(n));
    bool(line, isFn, () -> NodeUtil.isNameReferenced(n, "arguments"));
    bool(
        line,
        hasParent
            && (NodeUtil.isObjectLitKey(n)
                || NodeUtil.isLhsOfAssign(n)
                || (isName && NodeUtil.isDeclarationLValue(n))
                || n.isMemberFieldDef()
                || n.isComputedFieldDef()),
        () -> NodeUtil.isConstantDeclaration(NodeUtil.getBestJSDocInfo(n), n));
    bool(
        line,
        hasParent && NodeUtil.isControlStructure(parent),
        () -> NodeUtil.isControlStructureCodeBlock(parent, n));

    field(line, true, () -> tri(NodeUtil.getBooleanValue(n)));
    field(line, true, () -> str(NodeUtil.getStringValue(n)));
    field(line, true, () -> str(NodeUtil.getArrayElementStringValue(n)));
    field(line, n.isArrayLit(), () -> str(NodeUtil.arrayToString(n)));
    field(line, true, () -> dbl(NodeUtil.getNumberValue(n)));
    field(line, true, () -> dbl(NodeUtil.getNumberValueNoConversions(n)));
    field(line, true, () -> big(NodeUtil.getBigIntValue(n)));
    field(line, isFnOrClass, () -> str(NodeUtil.getName(n)));
    field(line, isFnOrClass, () -> idx(NodeUtil.getNameNode(n)));
    field(line, true, () -> str(NodeUtil.getNearestFunctionName(n)));
    field(line, n.isClass(), () -> idx(NodeUtil.getClassMembers(n)));
    field(line, n.isClass(), () -> idx(NodeUtil.getEs6ClassConstructorMemberFunctionDef(n)));
    field(
        line,
        NodeUtil.isRelationalOperation(n),
        () -> tok(NodeUtil.getInverseOperator(t)));
    field(
        line,
        n.isBlock(),
        () -> idx(NodeUtil.getInsertionPointAfterAllInnerFunctionDeclarations(n)));
    field(line, true, () -> Integer.toString(NodeUtil.precedence(t)));
    field(line, true, () -> NodeUtil.getKnownValueType(n).name());
    field(
        line, NodeUtil.isCompoundAssignmentOp(n), () -> tok(NodeUtil.getOpFromAssignmentOp(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingNonArrowFunction(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingClass(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingModuleIfPresent(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingFunction(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingScript(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingBlock(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingBlockScopeRoot(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingScopeRoot(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingHoistScopeRoot(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingStatement(n)));
    field(line, true, () -> idx(NodeUtil.getEnclosingType(n, Token.CALL)));
    field(
        line,
        n.isObjectLit() || n.isClassMembers(),
        () -> idx(NodeUtil.getFirstPropMatchingKey(n, firstKey(n))));
    field(
        line,
        n.isObjectLit() || n.isClassMembers(),
        () -> idx(NodeUtil.getFirstGetterMatchingKey(n, firstKey(n))));
    field(line, NodeUtil.isOptChainNode(n), () -> idx(NodeUtil.getStartOfOptChainSegment(n)));
    field(line, NodeUtil.isOptChainNode(n), () -> idx(NodeUtil.getEndOfOptChainSegment(n)));
    field(
        line,
        NodeUtil.isOptChainNode(n),
        () -> tok((Token) getNonOptChainToken.invoke(null, t)));
    field(line, isName || n.isGetProp(), () -> idx(NodeUtil.getAssignedValue(n)));
    field(line, true, () -> idx(NodeUtil.getLoopCodeBlock(n)));
    field(
        line, NodeUtil.isControlStructure(n), () -> idx(NodeUtil.getConditionExpression(n)));
    field(line, isFn, () -> idx(NodeUtil.getFunctionBody(n)));
    field(
        line,
        n.isCall() || n.isNew() || n.isTaggedTemplateLit(),
        () -> idx(NodeUtil.getCallTargetResolvingIndirectCalls(n)));
    boolean isTarget = hasParent && n.isValidAssignmentTarget();
    field(line, isTarget, () -> idx((Node) getEnclosingTarget.invoke(null, n)));
    field(line, isTarget, () -> idx(NodeUtil.getRootTarget(n)));
    field(line, isTarget, () -> idx(NodeUtil.getDeclaringParent(n)));
    field(
        line,
        isObjectOrClassLitKeyToken(n),
        () -> str(NodeUtil.getObjectOrClassLitKeyName(n)));
    field(
        line, isObjectOrClassLitKeyToken(n), () -> idx(NodeUtil.getObjectOrClassLitKeyNode(n)));
    field(line, true, () -> str(NodeUtil.opToStr(t)));
    field(line, true, () -> NodeUtil.opToStrNoFail(t));
    field(line, isQName, () -> idx(NodeUtil.getRootOfQualifiedName(n)));
    field(
        line, isQName, () -> str(NodeUtil.getRootOfQualifiedName(n.getQualifiedName())));
    field(
        line,
        isQName,
        () -> b(NodeUtil.isValidQualifiedName(FeatureSet.ES3, n.getQualifiedName())));
    field(
        line,
        isQName,
        () -> b(NodeUtil.isValidQualifiedName(FeatureSet.ES_NEXT, n.getQualifiedName())));
    field(
        line,
        hasString(n),
        () ->
            b(NodeUtil.isValidPropertyName(FeatureSet.ES3, n.getString()))
                + b(NodeUtil.isValidPropertyName(FeatureSet.ES_NEXT, n.getString()))
                + b(NodeUtil.isValidSimpleName(n.getString()))
                + b(NodeUtil.isLatin(n.getString())));
    field(line, n.isStringLit(), () -> dbl(NodeUtil.getStringNumberValue(n.getString())));
    field(line, n.isStringLit(), () -> big(NodeUtil.getStringBigIntValue(n.getString())));
    field(line, n.isStringLit(), () -> str(NodeUtil.trimJsWhiteSpace(n.getString())));
    field(
        line,
        isAssigningParent(n),
        () -> {
          List<Node> lhs = new ArrayList<>();
          NodeUtil.visitLhsNodesInNode(n, lhs::add);
          return list(lhs);
        });
    field(
        line,
        n.isDestructuringPattern(),
        () -> {
          List<Node> lhs = new ArrayList<>();
          NodeUtil.visitLhsNodesInDestructuringPattern(n, lhs::add);
          return list(lhs);
        });
    field(
        line,
        NodeUtil.isObjectDefinePropertiesDefinition(n),
        () -> list(NodeUtil.getObjectDefinedPropertiesKeys(n)));
    field(line, true, () -> idx(NodeUtil.getPrototypeClassName(n)));
    field(
        line,
        isQName && n.getQualifiedName().contains(".prototype."),
        () -> str(NodeUtil.getPrototypePropertyName(n)));
    field(line, true, () -> Integer.toString(NodeUtil.countAstSizeUpToLimit(n, 10)));
    field(line, isFn || n.isScript(), () -> Integer.toString(NodeUtil.countAstSize(n)));
    field(line, n.isTry(), () -> idx(NodeUtil.getCatchBlock(n)));
    field(line, isFn, () -> idx(NodeUtil.getFunctionParameters(n)));
    field(line, true, () -> str(NodeUtil.getSourceName(n)));
    field(
        line,
        true,
        () -> {
          StaticSourceFile f = NodeUtil.getSourceFile(n);
          return f == null ? "-" : str(f.getName());
        });
    field(
        line,
        true,
        () -> {
          InputId id = NodeUtil.getInputId(n);
          return id == null ? "-" : str(id.getIdName());
        });
    field(line, isFn, () -> idx(NodeUtil.getArgumentForFunction(n, 0)));
    field(line, isFn, () -> idx(NodeUtil.getArgumentForFunction(n, 1)));
    field(line, NodeUtil.isCallOrNew(n), () -> idx(NodeUtil.getArgumentForCallOrNew(n, 0)));
    field(
        line,
        isName || n.isStringKey(),
        () -> {
          JSTypeExpression e = NodeUtil.getDeclaredTypeExpression(n);
          return e == null ? "-" : str(e.getRoot().toString());
        });
    field(line, true, () -> idx(NodeUtil.getBestJSDocInfoNode(n)));
    field(
        line,
        true,
        () -> {
          JSDocInfo info = NodeUtil.getBestJSDocInfo(n);
          return info == null ? "-" : "J";
        });
    field(line, true, () -> idx(NodeUtil.getBestJsDocInfoNodeStrict(n)));
    field(line, true, () -> idx(NodeUtil.getBestLValue(n)));
    field(line, hasParent, () -> idx(NodeUtil.getRValueOfLValue(n)));
    field(line, true, () -> idx(NodeUtil.getBestLValueOwner(n)));
    field(line, true, () -> str(NodeUtil.getBestLValueName(n)));
    field(line, true, () -> idx(NodeUtil.getBestLValueRoot(n)));
    field(
        line,
        NodeUtil.isInvocation(n),
        () -> list(NodeUtil.getInvocationArgsAsIterable(n)));
    field(
        line,
        NodeUtil.isInvocation(n),
        () -> Integer.toString(NodeUtil.getInvocationArgsCount(n)));
    field(
        line,
        n.isScript(),
        () -> {
          FeatureSet fs = NodeUtil.getFeatureSetOfScript(n);
          return fs == null ? "-" : str(fs.toString());
        });
    field(
        line,
        isName || n.isImportStar(),
        () -> {
          NodeUtil.GoogRequire r = NodeUtil.getGoogRequireInfo(n);
          return r == null ? "-" : str(r.toString());
        });
    field(line, n.isScript(), () -> Integer.toString(NodeUtil.estimateNumLines(n)));
    field(
        line,
        n.isParamList() || n.isDestructuringPattern(),
        () -> {
          List<Node> names = new ArrayList<>();
          NodeUtil.getParamOrPatternNames(n, names::add);
          return list(names);
        });
    field(
        line,
        n.isScript() || n.isBlock(),
        () -> list(NodeUtil.getVarsDeclaredInBranch(n)));
    field(
        line,
        isQName && !(hasParent && parent.isGetProp() && parent.getFirstChild() == n),
        () -> compact(NodeUtil.newQName(compiler, n.getQualifiedName())));
    field(
        line,
        isFn || n.isScript(),
        () ->
            NodeUtil.getNodeTypeReferenceCount(n, Token.NAME, NodeUtil.MATCH_NOT_FUNCTION)
                + ","
                + NodeUtil.getNameReferenceCount(n, "a"));
    field(
        line,
        isFn || n.isScript(),
        () -> {
          List<Node> pre = new ArrayList<>();
          for (Node x : NodeUtil.preOrderIterable(n, NodeUtil.MATCH_NOT_FUNCTION)) {
            pre.add(x);
          }
          List<Node> post = new ArrayList<>();
          NodeUtil.visitPostOrder(n, post::add, NodeUtil.MATCH_NOT_FUNCTION);
          List<Node> pre2 = new ArrayList<>();
          NodeUtil.visitPreOrder(n, pre2::add);
          long h = 0;
          for (Node x : pre) {
            h = h * 31 + index.getOrDefault(x, -1);
          }
          for (Node x : post) {
            h = h * 37 + index.getOrDefault(x, -1);
          }
          return pre.size() + "," + post.size() + "," + pre2.size() + "," + h;
        });
    field(
        line,
        isFn || n.isScript(),
        () -> idx(NodeUtil.findPreorder(n, Node::isName, NodeUtil.MATCH_NOT_FUNCTION)));

    // Mutations, each on a fresh clone of the script.
    field(
        line,
        mutate && hasParent,
        () ->
            mutated(
                i,
                (clone, c) -> {
                  NodeUtil.removeChild(c.getParent(), c);
                  return "";
                }));
    field(
        line,
        mutate && n.isBlock() && hasParent,
        () -> mutated(i, (clone, c) -> b(NodeUtil.tryMergeBlock(c, false))));
    field(
        line,
        mutate && n.isBlock() && hasParent,
        () -> mutated(i, (clone, c) -> b(NodeUtil.tryMergeBlock(c, true))));
    field(
        line,
        mutate && hasParent,
        () ->
            mutated(
                i,
                (clone, c) -> {
                  NodeUtil.redeclareVarsInsideBranch(c);
                  return "";
                }));
    field(
        line,
        mutate && NodeUtil.isEndOfOptChainSegment(n),
        () ->
            mutated(
                i,
                (clone, c) -> {
                  NodeUtil.convertToNonOptionalChainSegment(c);
                  return "";
                }));
    field(
        line,
        mutate && hasParent && NodeUtil.isNameDeclaration(parent),
        () ->
            mutated(
                i,
                (clone, c) -> {
                  NodeUtil.replaceDeclarationChild(c, IR.exprResult(IR.number(1)));
                  return "";
                }));
    field(
        line,
        mutate && isFnOrClass,
        () ->
            mutated(
                i,
                (clone, c) -> {
                  NodeUtil.removeName(c);
                  return "";
                }));
    field(
        line,
        mutate && n.isTry(),
        () ->
            mutated(
                i,
                (clone, c) -> {
                  NodeUtil.maybeAddFinally(c);
                  return "";
                }));

    // Constructions (the AST is not changed).
    field(
        line,
        n.isNumber() || n.isName() || n.isNeg() || n.isStringLit(),
        () -> {
          Double v = NodeUtil.getNumberValue(n);
          return v == null ? "-" : compact(NodeUtil.numberNode(v, n));
        });
    field(line, n.isNumber() || n.isName(), () -> compact(NodeUtil.newUndefinedNode(n)));
    boolean isOuterQName = isQName && !(hasParent && parent.isGetProp() && parent.getFirstChild() == n);
    field(
        line,
        isOuterQName,
        () -> compact(NodeUtil.newQNameDeclaration(compiler, n.getQualifiedName(), null, null)));
    field(
        line,
        isOuterQName,
        () ->
            compact(
                NodeUtil.newQNameDeclaration(
                    compiler, n.getQualifiedName(), IR.number(0), null, Token.LET)));
    field(
        line,
        isOuterQName,
        () -> compact(NodeUtil.newQName(compiler, n.getQualifiedName(), n, "orig")));
    field(
        line,
        isOuterQName,
        () -> compact(NodeUtil.newQName(compiler, n.getQualifiedName(), n, n.getQualifiedName())));
    field(line, isName, () -> compact(NodeUtil.newVarNode(n.getString(), IR.number(1))));
    field(line, isName, () -> compact(NodeUtil.newName(compiler, n.getString(), n, "orig")));
    field(
        line,
        n.isGetProp(),
        () -> compact(NodeUtil.newPropertyAccess(compiler, IR.name("x"), n.getString())));

    sb.append(line).append('\n');
  }

  void sweepScript(StringBuilder sb, Node script) {
    StringBuilder line = new StringBuilder("# file-summary");
    field(
        line,
        true,
        () -> {
          int count = 0;
          for (Node x : NodeUtil.preOrderIterable(script)) {
            count++;
          }
          return Integer.toString(count);
        });
    field(
        line,
        true,
        () -> {
          List<Node> all = new ArrayList<>();
          NodeUtil.visitPostOrder(script, all::add);
          return Integer.toString(all.size());
        });
    field(line, true, () -> Integer.toString(NodeUtil.countAstSize(script)));
    field(
        line,
        true,
        () ->
            Integer.toString(
                NodeUtil.getCount(script, Node::isName, NodeUtil.MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION)));
    field(
        line,
        true,
        () -> b(NodeUtil.has(script, Node::isThis, NodeUtil.MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION)));
    sb.append(line).append('\n');
  }
}

// Oracle tool for the level-11 infrastructure classes (cluster l11-infra): for each file, parses it
// the way CompilerInput.getAstRoot does (a Compiler with default options and
// languageIn=ECMASCRIPT_NEXT) and sweeps
//   V/X  AstValidator (validateScriptFeatures=true, TypeInfoValidation.NONE) over the SCRIPT: every
//        violation (message + node) and any exception;
//   C    the same over deliberately corrupted clones of the SCRIPT (six corruption kinds applied to
//        every k-th node in pre-order);
//   M    MinimizedCondition on every NOT/AND/OR/HOOK/COMMA node and every IF/WHILE/DO/FOR
//        condition, for both MinimizationStyles: the minimized MeasuredNode's length and change
//        bits and the tree that applyTo builds;
//   I    NodeIterators.LocalVarMotion.forVar for every declared NAME and forAssign for every
//        `x = ...;` statement: the tokens of every node the iterator yields;
//   T    TemplateAstMatcher with seven seeded templates (LOOSE strategy; one also EXACT) on every
//        node: each match, its loose bit and getTemplateNodeToMatchMap();
//   D    DataFlowAnalysis on the CFG of the SCRIPT and of every function (ControlFlowAnalysis with
//        edge annotations): a forward analysis (names written), a branched forward analysis and a
//        backward analysis (names read); per CFG node the IN/OUT sets and step counts;
//   E    DataFlowAnalysis.computeEscaped for every function scope.
// Exceptions are printed as their toString() for IllegalState/IllegalArgumentException and as
// their class name otherwise.
//
// Build/run (see tools/oracle/README.md):
//   java -Xss1g -cp $CP com.google.javascript.jscomp.InfraDump <out.gz> <root> <rel-file>... |
//        --list=<listfile> [--every=N]
//
// Format: "## file <name>" then the section lines; every line escaped like ParsingDump
// ('\' -> "\\", chars outside printable ASCII -> \\uXXXX).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.ControlFlowGraph.Branch;
import com.google.javascript.jscomp.DataFlowAnalysis.FlowBrancher;
import com.google.javascript.jscomp.DataFlowAnalysis.FlowJoiner;
import com.google.javascript.jscomp.DataFlowAnalysis.LinearFlowState;
import com.google.javascript.jscomp.MinimizedCondition.MeasuredNode;
import com.google.javascript.jscomp.MinimizedCondition.MinimizationStyle;
import com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge;
import com.google.javascript.jscomp.graph.DiGraph.DiGraphNode;
import com.google.javascript.jscomp.graph.LatticeElement;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.Token;
import java.io.OutputStream;
import java.lang.reflect.Field;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.IdentityHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class InfraDump {

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Path root = Paths.get(args[1]);
    List<String> files = new ArrayList<>();
    int every = 1;
    for (int i = 2; i < args.length; i++) {
      String a = args[i];
      if (a.startsWith("--list=")) {
        for (String l :
            Files.readAllLines(Paths.get(a.substring("--list=".length())), StandardCharsets.UTF_8)) {
          if (!l.isEmpty()) {
            files.add(l);
          }
        }
      } else if (a.startsWith("--every=")) {
        every = Integer.parseInt(a.substring("--every=".length()));
      } else {
        files.add(a);
      }
    }
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(out))) {
      for (int i = 0; i < files.size(); i += every) {
        String rel = files.get(i);
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        StringBuilder sb = new StringBuilder();
        new InfraDump(sb).dumpOne(rel, src);
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  private final StringBuilder sb;
  private Compiler compiler;
  private final IdentityHashMap<Node, Integer> ids = new IdentityHashMap<>();

  InfraDump(StringBuilder sb) {
    this.sb = sb;
  }

  private void line(String s) {
    sb.append(escape(s)).append('\n');
  }

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

  static String exc(Throwable e) {
    if (e instanceof IllegalStateException || e instanceof IllegalArgumentException) {
      return e.toString();
    }
    return e.getClass().getName();
  }

  static String show(Node n) {
    return n == null ? "null" : n.toString(true, false, false);
  }

  /** A compact tree: "(" + toString(false, false, false) + children + ")". */
  static String sexpr(Node n) {
    StringBuilder b = new StringBuilder();
    sexpr(b, n);
    return b.toString();
  }

  private static void sexpr(StringBuilder b, Node n) {
    if (n == null) {
      b.append("null");
      return;
    }
    b.append('(').append(n.toString(false, false, false));
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      b.append(' ');
      sexpr(b, c);
    }
    b.append(')');
  }

  private String id(Node n) {
    if (n == null) {
      return "R";
    }
    Integer i = ids.get(n);
    return i == null ? "?" : i.toString();
  }

  private static void preorder(Node n, List<Node> out) {
    out.add(n);
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      preorder(c, out);
    }
  }

  private static int size(Node n, int cap) {
    int s = 1;
    for (Node c = n.getFirstChild(); c != null && s <= cap; c = c.getNext()) {
      s += size(c, cap);
    }
    return s;
  }

  void dumpOne(String name, String src) {
    line("## file " + name);
    compiler = new Compiler();
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    compiler.initOptions(options);
    Node script;
    try {
      script = new CompilerInput(SourceFile.fromCode(name, src)).getAstRoot(compiler);
    } catch (RuntimeException e) {
      line("# parse-exception " + exc(e));
      return;
    }
    line("# errors " + compiler.getErrorCount());
    List<Node> nodes = new ArrayList<>();
    preorder(script, nodes);
    for (int i = 0; i < nodes.size(); i++) {
      ids.put(nodes.get(i), i);
    }
    line("# nodes " + nodes.size());

    validate(script, "");
    corrupt(script, nodes.size());
    minimize(nodes);
    iterators(nodes);
    templates(nodes);
    dataflow(script, nodes);
    escaped(script);
  }

  // ---------------------------------------------------------------------------------------------
  // AstValidator

  private void validate(Node script, String indent) {
    List<String> violations = new ArrayList<>();
    AstValidator validator =
        new AstValidator(
            compiler,
            (message, n) -> violations.add(indent + "V " + message + "\t" + show(n)),
            /* validateScriptFeatures= */ true,
            /* shouldValidateRequiredInlinings= */ false);
    validator.setTypeValidationMode(AstValidator.TypeInfoValidation.NONE);
    String exception = null;
    try {
      validator.validateScript(script);
    } catch (RuntimeException e) {
      exception = exc(e);
    }
    for (String v : violations) {
      line(v);
    }
    if (exception != null) {
      line(indent + "X " + exception);
    }
  }

  private static final int CORRUPTIONS = 6;

  private void corrupt(Node script, int count) {
    if (count < 2) {
      return;
    }
    int step = Math.max(1, (count - 1) / 12);
    int kind = 0;
    for (int k = 1; k < count; k += step, kind++) {
      Node clone = script.cloneTree();
      List<Node> cloneNodes = new ArrayList<>();
      preorder(clone, cloneNodes);
      Node n = cloneNodes.get(k);
      String what;
      switch (kind % CORRUPTIONS) {
        case 0 -> {
          if (n.hasChildren()) {
            n.getFirstChild().detach();
            what = "remove-first-child";
          } else {
            n.addChildToBack(IR.number(0));
            what = "add-child";
          }
        }
        case 1 -> {
          n.addChildToBack(IR.name("extra"));
          what = "add-name-child";
        }
        case 2 -> {
          n.setToken(n.isName() ? Token.STRINGLIT : Token.NAME);
          what = "retoken";
        }
        case 3 -> {
          Node next = n.getNext();
          if (next != null) {
            next.detach();
            next.insertBefore(n);
            what = "swap-with-next";
          } else {
            n.getParent().addChildToFront(IR.empty());
            what = "add-empty-sibling";
          }
        }
        case 4 -> {
          n.replaceWith(IR.empty());
          what = "replace-with-empty";
        }
        default -> {
          Node block = new Node(Token.BLOCK);
          n.replaceWith(block);
          block.addChildToBack(n);
          what = "wrap-in-block";
        }
      }
      line("C " + k + " " + what + " " + n.getToken());
      validate(clone, "  ");
    }
  }

  // ---------------------------------------------------------------------------------------------
  // MinimizedCondition

  private static Field lengthField;
  private static Field changedField;

  static {
    try {
      lengthField = MeasuredNode.class.getDeclaredField("length");
      lengthField.setAccessible(true);
      changedField = MeasuredNode.class.getDeclaredField("changed");
      changedField.setAccessible(true);
    } catch (ReflectiveOperationException e) {
      throw new AssertionError(e);
    }
  }

  private static boolean isCondition(Node n) {
    Node parent = n.getParent();
    if (parent == null) {
      return false;
    }
    switch (parent.getToken()) {
      case IF, WHILE, DO, FOR -> {
        return NodeUtil.getConditionExpression(parent) == n && !n.isEmpty();
      }
      default -> {
        return false;
      }
    }
  }

  private void minimize(List<Node> nodes) {
    for (Node n : nodes) {
      boolean candidate =
          switch (n.getToken()) {
            case NOT, AND, OR, HOOK, COMMA -> true;
            default -> isCondition(n);
          };
      if (!candidate) {
        continue;
      }
      for (MinimizationStyle style : MinimizationStyle.values()) {
        Node clone = n.cloneTree();
        Node holder = new Node(Token.EXPR_RESULT, clone);
        String result;
        try {
          MinimizedCondition mc = MinimizedCondition.fromConditionNode(clone);
          MeasuredNode m = mc.getMinimized(style);
          boolean will = m.willChange(clone);
          result =
              "len="
                  + lengthField.getInt(m)
                  + " changed="
                  + changedField.getBoolean(m)
                  + " will="
                  + will;
          if (will) {
            m.applyTo(clone);
            result += " " + sexpr(holder.getFirstChild());
          }
        } catch (RuntimeException e) {
          result = "X " + exc(e);
        } catch (IllegalAccessException e) {
          throw new AssertionError(e);
        }
        line("M " + id(n) + " " + n.getToken() + " " + style + ": " + result);
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // NodeIterators

  private static final int MAX_MOTION = 40;

  private void motion(String header, NodeIterators.LocalVarMotion it) {
    StringBuilder b = new StringBuilder(header).append(':');
    int count = 0;
    while (it.hasNext()) {
      Node next = it.next();
      if (count++ == MAX_MOTION) {
        b.append(" ...");
        break;
      }
      b.append(' ').append(next.getToken()).append('@').append(id(next));
    }
    line(b.toString());
  }

  private void iterators(List<Node> nodes) {
    for (Node n : nodes) {
      if (NodeUtil.isNameDeclaration(n) && NodeUtil.isStatement(n)) {
        for (Node name = n.getFirstChild(); name != null; name = name.getNext()) {
          if (!name.isName()) {
            continue;
          }
          String header = "I var " + id(name);
          try {
            motion(
                header, NodeIterators.LocalVarMotion.forVar(compiler, name, n, n.getParent()));
          } catch (RuntimeException e) {
            line(header + ": X " + exc(e));
          }
        }
      } else if (n.isExprResult()
          && n.getFirstChild().isAssign()
          && n.getFirstFirstChild().isName()) {
        Node assign = n.getFirstChild();
        String header = "I assign " + id(assign.getFirstChild());
        try {
          motion(
              header,
              NodeIterators.LocalVarMotion.forAssign(
                  compiler, assign.getFirstChild(), assign, n, n.getParent()));
        } catch (RuntimeException e) {
          line(header + ": X " + exc(e));
        }
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // TemplateAstMatcher

  static final String[] TEMPLATES = {
    "function template() { 3; }",
    "/** @param {?} a */ function template(a) { a.length; }",
    "/** @param {?} a\n * @param {?} b */ function template(a, b) { a.indexOf(b) != -1; }",
    "function template() { var x = 1; x++; }",
    "/** @param {string} string_literal_s */\n"
        + "function template(string_literal_s) { goog.require(string_literal_s); }",
    "/** @param {number} n */ function template(n) { Math.floor(n); }",
    "/** @param {?} a\n * @param {?} b */ function template(a, b) { a = b; }",
  };

  private void templates(List<Node> nodes) {
    for (int t = 0; t < TEMPLATES.length; t++) {
      TemplateAstMatcher matcher;
      TemplateAstMatcher exact = null;
      try {
        Node templateScript =
            new CompilerInput(SourceFile.fromCode("template" + t, TEMPLATES[t]))
                .getAstRoot(compiler);
        Node fn = templateScript.getFirstChild();
        matcher = new TemplateAstMatcher(compiler, fn, TypeMatchingStrategy.LOOSE);
        if (t == 1) {
          exact = new TemplateAstMatcher(compiler, fn, TypeMatchingStrategy.EXACT);
        }
      } catch (RuntimeException e) {
        line("T " + t + " X " + exc(e));
        continue;
      }
      for (Node n : nodes) {
        match("T " + t, matcher, n);
        if (exact != null) {
          match("T " + t + "e", exact, n);
        }
      }
    }
  }

  private void match(String tag, TemplateAstMatcher matcher, Node n) {
    boolean matches;
    try {
      matches = matcher.matches(n);
    } catch (RuntimeException e) {
      line(tag + " " + id(n) + ": X " + exc(e));
      return;
    }
    if (!matches) {
      return;
    }
    StringBuilder b = new StringBuilder(tag).append(' ').append(id(n)).append(": loose=");
    b.append(matcher.isLooseMatch()).append(" {");
    boolean first = true;
    for (Map.Entry<String, Node> e : matcher.getTemplateNodeToMatchMap().entrySet()) {
      if (!first) {
        b.append(", ");
      }
      first = false;
      b.append(e.getKey()).append('=').append(sexpr(e.getValue()));
    }
    line(b.append('}').toString());
  }

  // ---------------------------------------------------------------------------------------------
  // DataFlowAnalysis

  /** A set of names, compared by value. */
  static final class NameSet implements LatticeElement {
    final TreeSet<String> names;

    NameSet() {
      this.names = new TreeSet<>();
    }

    NameSet(NameSet other) {
      this.names = new TreeSet<>(other.names);
    }

    @Override
    public boolean equals(Object o) {
      return o instanceof NameSet other && other.names.equals(names);
    }

    @Override
    public int hashCode() {
      return names.hashCode();
    }

    @Override
    public String toString() {
      return names.toString();
    }
  }

  static final class UnionJoiner implements FlowJoiner<NameSet> {
    private final NameSet result = new NameSet();

    @Override
    public void joinFlow(NameSet input) {
      result.names.addAll(input.names);
    }

    @Override
    public NameSet finish() {
      return result;
    }
  }

  /**
   * The NAMEs of a CFG node: its subtree without descending into BLOCK, FUNCTION, CLASS or SCRIPT
   * children; writes are declared names and simple assignment targets, reads all others.
   */
  static void names(Node n, boolean writes, Set<String> out) {
    if (n.isName()) {
      Node parent = n.getParent();
      boolean isWrite =
          parent != null
              && (NodeUtil.isNameDeclaration(parent)
                  || (parent.isAssign() && parent.getFirstChild() == n));
      if (isWrite == writes) {
        out.add(n.getString());
      }
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      switch (c.getToken()) {
        case BLOCK, FUNCTION, CLASS, SCRIPT -> {}
        default -> names(c, writes, out);
      }
    }
  }

  static final class WrittenNames extends DataFlowAnalysis<Node, NameSet> {
    private final boolean branched;

    WrittenNames(ControlFlowGraph<Node> cfg, boolean branched) {
      super(cfg);
      this.branched = branched;
    }

    @Override
    boolean isForward() {
      return true;
    }

    @Override
    boolean isBranched() {
      return branched;
    }

    @Override
    FlowJoiner<NameSet> createFlowJoiner() {
      return new UnionJoiner();
    }

    @Override
    NameSet flowThrough(Node node, NameSet input) {
      NameSet out = new NameSet(input);
      names(node, true, out.names);
      return out;
    }

    @Override
    FlowBrancher<NameSet> createFlowBrancher(Node node, NameSet output) {
      return (Branch branch) -> {
        if (!branch.isConditional()) {
          return output;
        }
        NameSet result = new NameSet(output);
        result.names.add("@" + branch);
        return result;
      };
    }

    @Override
    NameSet createInitialEstimateLattice() {
      return new NameSet();
    }

    @Override
    NameSet createEntryLattice() {
      return new NameSet();
    }
  }

  static final class ReadNames extends DataFlowAnalysis<Node, NameSet> {
    ReadNames(ControlFlowGraph<Node> cfg) {
      super(cfg);
    }

    @Override
    boolean isForward() {
      return false;
    }

    @Override
    FlowJoiner<NameSet> createFlowJoiner() {
      return new UnionJoiner();
    }

    @Override
    NameSet flowThrough(Node node, NameSet input) {
      NameSet out = new NameSet(input);
      names(node, false, out.names);
      return out;
    }

    @Override
    NameSet createInitialEstimateLattice() {
      return new NameSet();
    }

    @Override
    NameSet createEntryLattice() {
      return new NameSet();
    }
  }

  private void dataflow(Node script, List<Node> nodes) {
    List<Node> roots = new ArrayList<>();
    roots.add(script);
    for (Node n : nodes) {
      if (n.isFunction()) {
        roots.add(n);
      }
    }
    for (Node root : roots) {
      ControlFlowGraph<Node> cfg;
      try {
        cfg =
            ControlFlowAnalysis.builder()
                .setCompiler(compiler)
                .setCfgRoot(root)
                .setIncludeEdgeAnnotations(true)
                .computeCfg();
      } catch (RuntimeException e) {
        line("D " + id(root) + " cfg X " + exc(e));
        continue;
      }
      runAnalysis("D " + id(root) + " fwd", cfg, new WrittenNames(cfg, false), false);
      runAnalysis("D " + id(root) + " br", cfg, new WrittenNames(cfg, true), true);
      runAnalysis("D " + id(root) + " bwd", cfg, new ReadNames(cfg), false);
    }
  }

  private void runAnalysis(
      String header, ControlFlowGraph<Node> cfg, DataFlowAnalysis<Node, NameSet> a, boolean edges) {
    try {
      a.analyze();
    } catch (RuntimeException e) {
      line(header + " X " + exc(e));
      return;
    }
    line(header);
    for (DiGraphNode<Node, Branch> node : cfg.getNodes()) {
      LinearFlowState<NameSet> state = node.getAnnotation();
      line(
          "  "
              + id(node.getValue())
              + " in="
              + state.getIn()
              + " out="
              + state.getOut()
              + " steps="
              + state.getStepCount());
      if (edges) {
        for (DiGraphEdge<Node, Branch> e : node.getOutEdges()) {
          line(
              "    -> "
                  + id(e.getDestination().getValue())
                  + " "
                  + e.getValue()
                  + " "
                  + e.getAnnotation());
        }
      }
    }
  }

  private void escaped(Node script) {
    SyntacticScopeCreator scopeCreator = new SyntacticScopeCreator(compiler);
    List<String> lines = new ArrayList<>();
    try {
      NodeTraversal.builder()
          .setCompiler(compiler)
          .setScopeCreator(scopeCreator)
          .setCallback(
              new NodeTraversal.AbstractScopedCallback() {
                @Override
                public void enterScope(NodeTraversal t) {
                  Node scopeRoot = t.getScopeRoot();
                  if (!scopeRoot.isFunction()) {
                    return;
                  }
                  Scope scope = t.getScope();
                  String result;
                  try {
                    Set<Var> escaped = new HashSet<>();
                    NodeUtil.AllVarsDeclaredInFunction all =
                        NodeUtil.getAllVarsDeclaredInFunction(compiler, scopeCreator, scope);
                    DataFlowAnalysis.computeEscaped(
                        scope, escaped, compiler, scopeCreator, all.getAllVariables());
                    TreeSet<String> names = new TreeSet<>();
                    for (Var v : escaped) {
                      names.add(v.getName() + "@" + id(v.getNameNode()));
                    }
                    result = names.toString();
                  } catch (RuntimeException e) {
                    result = "X " + exc(e);
                  }
                  lines.add("E " + id(scopeRoot) + ": " + result);
                }

                @Override
                public void visit(NodeTraversal t, Node n, Node parent) {}
              })
          .traverse(script);
    } catch (RuntimeException e) {
      lines.add("E X " + exc(e));
    }
    for (String l : lines) {
      line(l);
    }
  }
}

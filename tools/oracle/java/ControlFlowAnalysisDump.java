// Differential oracle for com.google.javascript.jscomp.ControlFlowAnalysis and ControlFlowGraph
// (internal/jscomp/control_flow_analysis.go, control_flow_graph.go).
//
// For every file: parses it with ParserRunner.parse and the compiler's own parsing.Config
// (Compiler.getParserConfig(ConfigContext.DEFAULT) for --language_in=ECMASCRIPT_NEXT, as
// ParsingDump does), files with parse errors are skipped, and the SCRIPT is placed under the
// compiler's root structure: ROOT(global) -> [ROOT(externs, empty), ROOT(js) -> SCRIPT]. Every node
// gets an id: its preorder index in that tree (global ROOT = 0, externs ROOT = 1, js ROOT = 2,
// SCRIPT = 3, ...); the implicit return is "R".
//
// Then it prints
//   S <id> <mayThrowException 0|1> <isEnteringNewCfgNode 0|1> <computeFollowNode(n)> <getExceptionHandler(n)>
// for the SCRIPT and every node below it (static API, cfa == null; an exception is "X<class>"), and,
// for every valid CFG root (NodeUtil.isValidCfgRoot) in preorder, the CFG that passes compute:
//   mode A: setIncludeEdgeAnnotations(true) (NodeTraversal.AbstractCfgCallback, TypeInferencePass,
//           FlowSensitiveInlineVariables), not traversing functions;
// and for the js ROOT also
//   mode F: setTraverseFunctions(true), no edge annotations (Compiler.computeCFG).
// A CFG is printed as
//   cfg <rootId> <mode>
//   N <id> <TOKEN> <line>:<charno> p=<priority>      (cfg.getNodes() order; the implicit return is
//                                                      "N R - - p=<priority>")
//    > <BRANCH> <dest>                               (the node's out-edges, in order)
//    < <BRANCH> <src>                                (the node's in-edges, in order)
//   multi                                            (when a finallyMap key holds several values)
// or, if computing it throws,
//   exception <class>: <message>   [then  cause <class>: <message>]
//
// Determinism: ControlFlowAnalysis.finallyMap is a HashMultimap keyed by Nodes, which hash by
// identity, so Java iterates the values of one key in an order that changes from run to run; it
// decides the order of the ON_EX edges created when leaving a FINALLY block (and of the graph nodes
// they create). The Go port iterates them in insertion order (PORTING.md §6, like IdentityHashMap ->
// LinkedHashMap). To compare deterministically, this oracle constructs ControlFlowAnalysis
// reflectively (exactly what Builder.computeCfg does) and replaces finallyMap by a
// LinkedHashMultimap; "multi" marks the CFGs where the order could differ in production Java.
//
// Every line is escaped like ParsingDump (backslash, and every char outside printable ASCII as
// \\uXXXX). The output is gzipped, one "## file <rel>" section per file.
//
// Usage: java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.ControlFlowAnalysisDump \
//          <out.gz> <root> <rel-file>... | --list=<listfile>
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.ControlFlowGraph.Branch;
import com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge;
import com.google.javascript.jscomp.graph.DiGraph.DiGraphNode;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.LinkedHashMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.SetMultimap;
import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Constructor;
import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.IdentityHashMap;
import java.util.List;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class ControlFlowAnalysisDump {
  private static final Constructor<ControlFlowAnalysis> CTOR;
  private static final Method COMPUTE;
  private static final Field FINALLY;

  static {
    try {
      CTOR =
          ControlFlowAnalysis.class.getDeclaredConstructor(
              AbstractCompiler.class, boolean.class, boolean.class);
      CTOR.setAccessible(true);
      COMPUTE = ControlFlowAnalysis.class.getDeclaredMethod("computeCfg", Node.class);
      COMPUTE.setAccessible(true);
      FINALLY = ControlFlowAnalysis.class.getDeclaredField("finallyMap");
      FINALLY.setAccessible(true);
    } catch (ReflectiveOperationException e) {
      throw new RuntimeException(e);
    }
  }

  public static void main(String[] args) throws Exception {
    Path outFile = Paths.get(args[0]);
    Path root = Paths.get(args[1]);
    List<String> files = new ArrayList<>();
    for (int i = 2; i < args.length; i++) {
      String f = args[i];
      if (f.startsWith("--list=")) {
        for (String l :
            Files.readAllLines(Paths.get(f.substring("--list=".length())), StandardCharsets.UTF_8)) {
          if (!l.isEmpty()) {
            files.add(l);
          }
        }
      } else {
        files.add(f);
      }
    }
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    compiler.initOptions(options);
    Config config = compiler.getParserConfig(AbstractCompiler.ConfigContext.DEFAULT);
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(outFile), 1 << 16)) {
      for (String rel : files) {
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        StringBuilder sb = new StringBuilder();
        dumpOne(sb, compiler, config, rel, src);
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  static String esc(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c >= 0x20 && c < 0x7f) {
        sb.append(c);
      } else {
        sb.append(String.format("\\u%04x", (int) c));
      }
    }
    return sb.toString();
  }

  static void line(StringBuilder sb, String s) {
    sb.append(esc(s)).append('\n');
  }

  static void dumpOne(StringBuilder sb, Compiler compiler, Config config, String name, String src) {
    line(sb, "## file " + name);
    boolean[] hasError = {false};
    ErrorReporter reporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int line, int lineOffset) {}

          @Override
          public void error(String message, String sourceName, int line, int lineOffset) {
            hasError[0] = true;
          }
        };
    ParserRunner.ParseResult result;
    try {
      result = ParserRunner.parse(SourceFile.fromCode(name, src), src, config, reporter);
    } catch (RuntimeException e) {
      line(sb, "# parse-exception");
      return;
    }
    if (hasError[0] || result.ast == null) {
      line(sb, "# parse-error");
      return;
    }
    Node script = result.ast;
    Node externsRoot = IR.root();
    Node jsRoot = IR.root(script);
    Node globalRoot = IR.root(externsRoot, jsRoot);

    IdentityHashMap<Node, Integer> ids = new IdentityHashMap<>();
    List<Node> preorder = new ArrayList<>();
    number(globalRoot, ids, preorder);
    line(sb, "# nodes " + preorder.size());

    // The static API.
    for (int i = 3; i < preorder.size(); i++) {
      Node n = preorder.get(i);
      String follow;
      try {
        follow = id(ids, ControlFlowAnalysis.computeFollowNode(n));
      } catch (RuntimeException e) {
        follow = "X" + e.getClass().getSimpleName();
      }
      String handler;
      try {
        handler = id(ids, ControlFlowAnalysis.getExceptionHandler(n));
      } catch (RuntimeException e) {
        handler = "X" + e.getClass().getSimpleName();
      }
      line(
          sb,
          "S "
              + i
              + " "
              + (ControlFlowAnalysis.mayThrowException(n) ? 1 : 0)
              + " "
              + (ControlFlowGraph.isEnteringNewCfgNode(n) ? 1 : 0)
              + " "
              + follow
              + " "
              + handler);
    }

    // The CFGs.
    for (Node n : preorder) {
      if (NodeUtil.isValidCfgRoot(n)) {
        dumpCfg(sb, compiler, ids, n, "A", false, true);
      }
    }
    dumpCfg(sb, compiler, ids, jsRoot, "F", true, false);
  }

  static void number(Node n, IdentityHashMap<Node, Integer> ids, List<Node> preorder) {
    ids.put(n, preorder.size());
    preorder.add(n);
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      number(c, ids, preorder);
    }
  }

  static String id(IdentityHashMap<Node, Integer> ids, Node n) {
    if (n == null) {
      return "R";
    }
    Integer i = ids.get(n);
    return i == null ? "?" : String.valueOf(i);
  }

  @SuppressWarnings("unchecked")
  static void dumpCfg(
      StringBuilder sb,
      Compiler compiler,
      IdentityHashMap<Node, Integer> ids,
      Node root,
      String mode,
      boolean traverseFunctions,
      boolean edgeAnnotations) {
    line(sb, "cfg " + ids.get(root) + " " + mode);
    ControlFlowGraph<Node> cfg;
    SetMultimap<Node, Node> finallyMap = LinkedHashMultimap.create();
    try {
      ControlFlowAnalysis cfa = CTOR.newInstance(compiler, traverseFunctions, edgeAnnotations);
      FINALLY.set(cfa, finallyMap);
      COMPUTE.invoke(cfa, root);
      cfg = cfa.getCfg();
    } catch (InvocationTargetException e) {
      Throwable t = e.getCause();
      line(sb, "exception " + t.getClass().getName() + ": " + t.getMessage());
      if (t.getCause() != null) {
        line(sb, "cause " + t.getCause().getClass().getName() + ": " + t.getCause().getMessage());
      }
      return;
    } catch (ReflectiveOperationException e) {
      throw new RuntimeException(e);
    }
    for (DiGraphNode<Node, Branch> gn : cfg.getNodes()) {
      Node v = gn.getValue();
      String desc =
          v == null
              ? "N R - -"
              : "N " + id(ids, v) + " " + v.getToken() + " " + v.getLineno() + ":" + v.getCharno();
      line(sb, desc + " p=" + gn.getPriority());
      for (DiGraphEdge<Node, Branch> e : gn.getOutEdges()) {
        line(sb, " > " + e.getValue() + " " + id(ids, e.getDestination().getValue()));
      }
      for (DiGraphEdge<Node, Branch> e : gn.getInEdges()) {
        line(sb, " < " + e.getValue() + " " + id(ids, e.getSource().getValue()));
      }
    }
    for (Node k : finallyMap.keySet()) {
      if (finallyMap.get(k).size() > 1) {
        line(sb, "multi");
        break;
      }
    }
  }
}

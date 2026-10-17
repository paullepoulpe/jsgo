// Oracle tool for the peephole passes (PeepholeFoldConstants, PeepholeMinimizeConditions,
// PeepholeSubstituteAlternateSyntax, PeepholeReplaceKnownMethods, PeepholeRemoveDeadCode,
// PeepholeCollectPropertyAssignments, MinimizeExitPoints, StatementFusion, ExploitAssigns,
// ConvertToDottedProperties and PeepholeOptimizationsPass).
//
// For every input file and every configuration below it parses the file the way
// CompilerTestCase does (a Compiler with default CompilerOptions and languageIn
// ECMASCRIPT_NEXT, empty externs, compiler.parseInputs()), sets the accessor summary as
// CompilerTestCase.updateAccessorSummary does, normalizes it with
// Normalize.createNormalizeForOptimizations (CompilerTestCase.enableNormalize()), and runs the
// pass the way each pass's unit test's getProcessor does: a PeepholeOptimizationsPass wrapping
// the single optimization, with the late/useTypes flags DefaultPassConfig uses.
//
// Configurations (one fresh Compiler each):
//   normalize           no pass (the normalized input)
//   fold                PeepholeFoldConstants(late=false, useTypes=false)
//   fold-late           PeepholeFoldConstants(late=true, useTypes=false)
//   minimize            PeepholeMinimizeConditions(late=false)
//   minimize-late       PeepholeMinimizeConditions(late=true)
//   substitute          PeepholeSubstituteAlternateSyntax(late=false)
//   substitute-late     PeepholeSubstituteAlternateSyntax(late=true)
//   known-methods       PeepholeReplaceKnownMethods(late=false, useTypes=false)
//   known-methods-late  PeepholeReplaceKnownMethods(late=true, useTypes=false)
//   dead-code           PeepholeRemoveDeadCode()
//   collect-props       PeepholeCollectPropertyAssignments()
//   exit-points         MinimizeExitPoints()
//   fusion              StatementFusion()
//   exploit-assigns     ExploitAssigns()           (DefaultPassConfig.exploitAssign)
//   dotted              new ConvertToDottedProperties(compiler)
//   peephole            DefaultPassConfig.peepholeOptimizations (createPeepholeOptimizationsPass
//                       with expectAstIsNormalized=true: MinimizeExitPoints,
//                       PeepholeMinimizeConditions, PeepholeSubstituteAlternateSyntax,
//                       PeepholeReplaceKnownMethods, PeepholeRemoveDeadCode,
//                       PeepholeFoldConstants, PeepholeCollectPropertyAssignments; late=false)
//   peephole-late       DefaultPassConfig.latePeepholeOptimizations (StatementFusion,
//                       PeepholeRemoveDeadCode, PeepholeMinimizeConditions,
//                       PeepholeSubstituteAlternateSyntax, PeepholeReplaceKnownMethods,
//                       PeepholeFoldConstants; late=true)
//
// Build/run:
//   javac -cp closure-compiler.jar -d build PeepholeDump.java
//   java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:build \
//        com.google.javascript.jscomp.PeepholeDump
//        [--full] [--configs=a,b,...] <out-file(.gz)> <root> <rel-file>... | @<listfile>
//
// Output, per file and configuration, a block
//   === <rel-file> <config>
//   # parse errors              then  <JSError line>          (and no further sections)
//   # diagnostics               then  <JSError line>          (errors and warnings reported
//                                                              by the pass, in order)
//   # exception <class>: <message>                            (if the pass threw; the tree
//                                                              is not printed)
//   # tree                      then  Node.toStringTree() of the main (non-externs) root
//   # bigints                   then  "<line>:<charno> <value>" for every BIGINT node, in
//                                     pre-order (Node.toString does not print them)
// where a JSError line is "<ERROR|WARNING> <type key> <line>:<charno> <description>".
// Without --full only "<rel-file>\t<config>\t<sha256 of the block>" is printed per block.
// Every line is escaped like ParsingDump: '\' -> "\\", and every char outside printable ASCII
// -> \\uXXXX, so the output is ASCII.
package com.google.javascript.jscomp;

import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.function.Function;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class PeepholeDump {

  /** A configuration: builds the pass to run on the normalized AST (null: no pass). */
  private static final LinkedHashMap<String, Function<AbstractCompiler, CompilerPass>> CONFIGS =
      new LinkedHashMap<>();

  private static CompilerPass single(
      AbstractCompiler compiler, String name, AbstractPeepholeOptimization optimization) {
    return new PeepholeOptimizationsPass(compiler, name, optimization);
  }

  static {
    CONFIGS.put("normalize", c -> null);
    CONFIGS.put("fold", c -> single(c, "fold", new PeepholeFoldConstants(false, false)));
    CONFIGS.put("fold-late", c -> single(c, "fold-late", new PeepholeFoldConstants(true, false)));
    CONFIGS.put("minimize", c -> single(c, "minimize", new PeepholeMinimizeConditions(false)));
    CONFIGS.put(
        "minimize-late", c -> single(c, "minimize-late", new PeepholeMinimizeConditions(true)));
    CONFIGS.put(
        "substitute", c -> single(c, "substitute", new PeepholeSubstituteAlternateSyntax(false)));
    CONFIGS.put(
        "substitute-late",
        c -> single(c, "substitute-late", new PeepholeSubstituteAlternateSyntax(true)));
    CONFIGS.put(
        "known-methods",
        c -> single(c, "known-methods", new PeepholeReplaceKnownMethods(false, false)));
    CONFIGS.put(
        "known-methods-late",
        c -> single(c, "known-methods-late", new PeepholeReplaceKnownMethods(true, false)));
    CONFIGS.put("dead-code", c -> single(c, "dead-code", new PeepholeRemoveDeadCode()));
    CONFIGS.put(
        "collect-props", c -> single(c, "collect-props", new PeepholeCollectPropertyAssignments()));
    CONFIGS.put("exit-points", c -> single(c, "exit-points", new MinimizeExitPoints()));
    CONFIGS.put("fusion", c -> single(c, "fusion", new StatementFusion()));
    CONFIGS.put(
        "exploit-assigns",
        c -> new PeepholeOptimizationsPass(c, PassNames.EXPLOIT_ASSIGN, new ExploitAssigns()));
    CONFIGS.put("dotted", c -> new ConvertToDottedProperties(c));
    CONFIGS.put(
        "peephole",
        c -> {
          // DefaultPassConfig.createPeepholeOptimizationsPass(compiler, PEEPHOLE_OPTIMIZATIONS,
          // true) with useTypesForOptimization = options.shouldUseTypesForLocalOptimization().
          boolean late = false;
          boolean useTypes = c.getOptions().shouldUseTypesForLocalOptimization();
          List<AbstractPeepholeOptimization> optimizations = new ArrayList<>();
          optimizations.add(new MinimizeExitPoints());
          optimizations.add(new PeepholeMinimizeConditions(late));
          optimizations.add(new PeepholeSubstituteAlternateSyntax(late));
          optimizations.add(new PeepholeReplaceKnownMethods(late, useTypes));
          optimizations.add(new PeepholeRemoveDeadCode());
          optimizations.add(new PeepholeFoldConstants(late, useTypes));
          optimizations.add(new PeepholeCollectPropertyAssignments());
          return new PeepholeOptimizationsPass(c, PassNames.PEEPHOLE_OPTIMIZATIONS, optimizations);
        });
    CONFIGS.put(
        "peephole-late",
        c -> {
          // DefaultPassConfig.latePeepholeOptimizations.
          boolean late = true;
          boolean useTypes = c.getOptions().shouldUseTypesForLocalOptimization();
          return new PeepholeOptimizationsPass(
              c,
              "latePeepholeOptimizations",
              new StatementFusion(),
              new PeepholeRemoveDeadCode(),
              new PeepholeMinimizeConditions(late),
              new PeepholeSubstituteAlternateSyntax(late),
              new PeepholeReplaceKnownMethods(late, useTypes),
              new PeepholeFoldConstants(late, useTypes));
        });
  }

  public static void main(String[] args) throws Exception {
    boolean full = false;
    List<String> configs = new ArrayList<>(CONFIGS.keySet());
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.startsWith("--configs=")) {
        configs = Arrays.asList(a.substring("--configs=".length()).split(","));
      } else {
        rest.add(a);
      }
    }
    Path outFile = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    List<String> files = new ArrayList<>();
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
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    OutputStream os = Files.newOutputStream(outFile);
    if (outFile.toString().endsWith(".gz")) {
      os = new GZIPOutputStream(os, 1 << 16);
    }
    try (PrintStream out = new PrintStream(os, false, StandardCharsets.UTF_8)) {
      for (String rel : files) {
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        for (String config : configs) {
          String block = dumpOne(rel, src, config);
          if (full) {
            out.print(block);
          } else {
            out.print(rel + "\t" + config + "\t" + sha256(block) + "\n");
          }
        }
      }
    }
  }

  static String sha256(String s) throws Exception {
    byte[] d = MessageDigest.getInstance("SHA-256").digest(s.getBytes(StandardCharsets.UTF_8));
    StringBuilder sb = new StringBuilder();
    for (byte b : d) {
      sb.append(String.format("%02x", b & 0xff));
    }
    return sb.toString();
  }

  static String dumpOne(String rel, String src, String config) {
    StringBuilder sb = new StringBuilder();
    line(sb, "=== " + rel + " " + config);
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    compiler.init(
        List.of(SourceFile.fromCode("externs", "")),
        List.of(SourceFile.fromCode(rel, src)),
        options);
    Node root = compiler.parseInputs();
    if (root == null || compiler.getErrorCount() > 0) {
      line(sb, "# parse errors");
      for (JSError e : compiler.getErrors()) {
        line(sb, error("ERROR", e));
      }
      return sb.toString();
    }
    int parseWarnings = compiler.getWarningCount();
    Node externsRoot = root.getFirstChild();
    Node mainRoot = root.getLastChild();
    try {
      // CompilerTestCase.updateAccessorSummary
      compiler.setAccessorSummary(
          AccessorSummary.create(GatherGetterAndSetterProperties.gather(compiler, root)));
      // CompilerTestCase.normalizeActualCode
      Normalize.createNormalizeForOptimizations(compiler).process(externsRoot, mainRoot);
      CompilerPass pass = CONFIGS.get(config).apply(compiler);
      if (pass != null) {
        pass.process(externsRoot, mainRoot);
      }
    } catch (RuntimeException e) {
      diagnostics(sb, compiler, parseWarnings);
      String msg = String.valueOf(e.getMessage());
      int nl = msg.indexOf('\n');
      if (nl >= 0) {
        msg = msg.substring(0, nl);
      }
      line(sb, "# exception " + e.getClass().getName() + ": " + msg);
      return sb.toString();
    }
    diagnostics(sb, compiler, parseWarnings);
    line(sb, "# tree");
    for (String l : mainRoot.toStringTree().split("\n", -1)) {
      line(sb, l);
    }
    // Node.toString does not print BIGINT values.
    line(sb, "# bigints");
    bigints(sb, mainRoot);
    return sb.toString();
  }

  private static void bigints(StringBuilder sb, Node n) {
    if (n.isBigInt()) {
      line(sb, n.getLineno() + ":" + n.getCharno() + " " + n.getBigInt());
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      bigints(sb, c);
    }
  }

  private static void diagnostics(StringBuilder sb, Compiler compiler, int skipWarnings) {
    line(sb, "# diagnostics");
    for (JSError e : compiler.getErrors()) {
      line(sb, error("ERROR", e));
    }
    List<JSError> warnings = compiler.getWarnings();
    for (int i = skipWarnings; i < warnings.size(); i++) {
      line(sb, error("WARNING", warnings.get(i)));
    }
  }

  private static String error(String level, JSError e) {
    return level
        + " "
        + e.getType().key
        + " "
        + e.getLineno()
        + ":"
        + e.getCharno()
        + " "
        + e.getDescription();
  }

  /** Backslash and every char outside printable ASCII become escapes. */
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

  private static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }
}

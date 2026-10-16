// Differential oracle for the w5-remove-unused passes (internal/jscomp): RemoveUnusedCode,
// PureFunctionIdentifier, DeadAssignmentsElimination, CoalesceVariableNames,
// LiveVariablesAnalysis, OptimizeParameters, OptimizeReturns, OptimizeConstructors,
// DevirtualizeMethods, FunctionRewriter, CrossChunkReferenceCollector, CrossChunkCodeMotion and
// CrossChunkMethodMotion. (RemoveUnusedCode's polyfill removal is off in every run: it needs
// PolyfillUsageFinder, which is owned by another port; RemoveUnusedCodeTest covers it.)
//
// For every file, each run starts from a fresh Compiler set up the way CompilerTestCase sets one
// up: default CompilerOptions with --language_in=ECMASCRIPT_NEXT and --language_out=NO_TRANSPILE,
// Compiler.init([SourceFile.fromCode("externs.js", EXTERNS)], [SourceFile.fromCode(<rel>, <src>)],
// options) and Compiler.parse(). Files with parse errors print their diagnostics and are skipped.
// Before each pass run the prerequisites DefaultPassConfig runs before these passes: VarCheck
// (which declares undeclared names in the synthetic externs; the undefinedVars group is off),
// Normalize.createNormalizeForOptimizations, GatherGetterAndSetterProperties and
// GatherExternProperties(OPTIMIZE). The runs:
//   ruc-simple        RemoveUnusedCode as CompilationLevel.SIMPLE configures it (local vars only)
//   ruc-advanced      RemoveUnusedCode as ADVANCED configures it (locals, globals, prototype and
//                     this properties, Object.defineProperties definitions), polyfills off
//   pfi-ruc-advanced  PureFunctionIdentifier.Driver, then ruc-advanced (as DefaultPassConfig orders
//                     markPureFunctions before removeUnusedCode)
//   pfi               PureFunctionIdentifier.Driver; extra: the side-effect flags of every CALL,
//                     NEW, TAGGED_TEMPLATELIT and OPTCHAIN_CALL node (pre-order, externs first)
//   dae               DeadAssignmentsElimination
//   cvn               CoalesceVariableNames(usePseudoNames=false)
//   cvn-pseudo        CoalesceVariableNames(usePseudoNames=true)
//   lva               LiveVariablesAnalysis for every function, computed the way
//                     DeadAssignmentsElimination does (on exiting each function block scope, with
//                     NodeTraversal.AbstractCfgCallback's CFG); extra: per function its variables,
//                     escaped locals and the live-in/live-out sets of every CFG node
//   optparams         OptimizeParameters
//   optreturns        OptimizeReturns
//   optconstructors   OptimizeConstructors
//   devirtualize      OptimizeCalls.builder().addPass(new DevirtualizeMethods(compiler))
//                     (considerExterns false), as DevirtualizeMethodsTest runs it
//   rewrite-functions FunctionRewriter
// Then, for scripts (not ES modules) with at least 3 top-level statements, the cross-chunk runs on
// a 3-chunk chain m1 <- m2 <- m3 (JSChunkGraphBuilder.forChain in the unit tests). The source is
// split after top-level statements: with n statements s_0..s_{n-1} (children of the SCRIPT of the
// single-file parse), chunk k (k = 1, 2) starts at the end (source offset + length) of statement
// s_{floor(k*n/3)-1}; the line
//   ## chunks <offset1> <offset2>       (UTF-16 offsets into the file), or
//   ## chunks skipped module|small, or
//   ## chunks skipped offsets <offset1> <offset2>   when they are not 0 <= offset1 <= offset2 <=
//                                                    the file length
// and then each run's inputs are SourceFile.fromCode("m<k>.js", <part k>) in JSChunk "m<k>",
// set up with Compiler.initChunks and the same prerequisites:
//   crosschunk-refs    CrossChunkReferenceCollector(compiler, new SyntacticScopeCreator).process(js)
//                      extra: every top-level statement and every global variable's references
//   crosschunk-code    CrossChunkCodeMotion(compiler, chunkGraph, false)
//   crosschunk-method  CrossChunkMethodMotion(compiler, new IdGenerator(), false, false)
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside printable
// ASCII as \\uXXXX):
//   ## file <rel>
//   # parse ok|errors                then the parse diagnostics:
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   ## pass <name>
//   diag ...                         the diagnostics added since parsing, in ErrorManager order
//   exception <class>: <message>     if the prerequisites or the pass threw
//   tree unchanged                   if the AST (root.toStringTree()) is the same as after the
//                                    prerequisites, else
//   tree changed <sha256>            the first 16 hex digits of the SHA-256 of the escaped lines
//                                    of root.toStringTree() (each followed by '\n'), and with
//                                    --full the escaped lines themselves
//   changes <n> <sha256>             the change-scope roots (SCRIPT and FUNCTION nodes) whose
//                                    change time is non-zero, as "<token>@<line>:<charno>" (the
//                                    times themselves, here and as "[change_time: N]" in the trees,
//                                    are left out: RemoveUnusedCode's removal order, and so the
//                                    stamps, depends on identity hash codes in Java)
//   extra <n> <sha256>               the pass-specific lines: their count and hash, and with
//                                    --full the lines themselves
//
// Usage: java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//          com.google.javascript.jscomp.RemoveUnusedDump \
//          [--full] [--sample=<N>] <out.gz> <root> <rel-file>... | --list=<listfile>
// --sample=N dumps every Nth file of the list only (in list order, starting with the first).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.NodeUtil.AllVarsDeclaredInFunction;
import com.google.javascript.jscomp.graph.DiGraph.DiGraphNode;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.zip.GZIPOutputStream;

public final class RemoveUnusedDump {

  static final String EXTERNS =
      String.join(
          "\n",
          "var window;",
          "var document;",
          "var console;",
          "var Arguments;",
          "/** @constructor @param {*=} opt_value @return {!Object} */",
          "function Object(opt_value) {}",
          "/** @type {Function} */ Object.defineProperties = function() {};",
          "/** @type {Function} */ Object.defineProperty = function() {};",
          "/** @constructor @param {...*} var_args @return {!Array} */",
          "function Array(var_args) {}",
          "/** @constructor @param {*=} message */ function Error(message) {}",
          "/** @nosideeffects */ function externPure() {}",
          "function alert(a) {}",
          "");

  /** A pass run: gets the compiler and (externsRoot, jsRoot); may append extra output lines. */
  private interface Run {
    void run(Compiler compiler, Node externs, Node js, List<String> extra);
  }

  private static final class Pass {
    final String name;
    final Run run;

    Pass(String name, Run run) {
      this.name = name;
      this.run = run;
    }
  }

  private static boolean full = false;
  private static int sample = 1;

  public static void main(String[] args) throws Exception {
    List<String> rest = new ArrayList<>();
    List<String> files = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.startsWith("--sample=")) {
        sample = Integer.parseInt(a.substring("--sample=".length()));
      } else if (a.startsWith("--list=")) {
        for (String l :
            Files.readAllLines(Paths.get(a.substring("--list=".length())), StandardCharsets.UTF_8)) {
          if (!l.isEmpty()) {
            files.add(l);
          }
        }
      } else {
        rest.add(a);
      }
    }
    Path out = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    files.addAll(rest.subList(2, rest.size()));
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(out))) {
      int index = 0;
      for (String rel : files) {
        if (index++ % sample != 0) {
          continue;
        }
        StringBuilder sb = new StringBuilder();
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        dumpFile(sb, rel, src);
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  private static RemoveUnusedCode.Builder rucSimple(Compiler c) {
    // DefaultPassConfig.removeUnusedCode with CompilationLevel.SIMPLE's options.
    return new RemoveUnusedCode.Builder(c)
        .removeUnusedPrototypeProperties(false)
        .removeLocalVars(true)
        .removeGlobals(false)
        .preserveFunctionExpressionNames(false)
        .removeUnusedThisProperties(false)
        .removeUnusedObjectDefinePropertiesDefinitions(false)
        .assumeGettersArePure(true)
        .removeUnusedPolyfills(false);
  }

  private static RemoveUnusedCode.Builder rucAdvanced(Compiler c) {
    // DefaultPassConfig.removeUnusedCode with CompilationLevel.ADVANCED's options (polyfills off).
    return new RemoveUnusedCode.Builder(c)
        .removeUnusedPrototypeProperties(true)
        .removeLocalVars(true)
        .removeGlobals(true)
        .preserveFunctionExpressionNames(false)
        .removeUnusedThisProperties(true)
        .removeUnusedObjectDefinePropertiesDefinitions(true)
        .assumeGettersArePure(true)
        .removeUnusedPolyfills(false);
  }

  private static List<Pass> passes() {
    List<Pass> passes = new ArrayList<>();
    passes.add(new Pass("ruc-simple", (c, ext, js, extra) -> rucSimple(c).build().process(ext, js)));
    passes.add(
        new Pass("ruc-advanced", (c, ext, js, extra) -> rucAdvanced(c).build().process(ext, js)));
    passes.add(
        new Pass(
            "pfi-ruc-advanced",
            (c, ext, js, extra) -> {
              new PureFunctionIdentifier.Driver(c).process(ext, js);
              rucAdvanced(c).build().process(ext, js);
            }));
    passes.add(
        new Pass(
            "pfi",
            (c, ext, js, extra) -> {
              new PureFunctionIdentifier.Driver(c).process(ext, js);
              collectFlags(ext, extra);
              collectFlags(js, extra);
            }));
    passes.add(
        new Pass("dae", (c, ext, js, extra) -> new DeadAssignmentsElimination(c).process(ext, js)));
    passes.add(
        new Pass(
            "cvn", (c, ext, js, extra) -> new CoalesceVariableNames(c, false).process(ext, js)));
    passes.add(
        new Pass(
            "cvn-pseudo",
            (c, ext, js, extra) -> new CoalesceVariableNames(c, true).process(ext, js)));
    passes.add(new Pass("lva", (c, ext, js, extra) -> NodeTraversal.traverse(c, js, new Lva(c, extra))));
    passes.add(
        new Pass("optparams", (c, ext, js, extra) -> new OptimizeParameters(c).process(ext, js)));
    passes.add(
        new Pass("optreturns", (c, ext, js, extra) -> new OptimizeReturns(c).process(ext, js)));
    passes.add(
        new Pass(
            "optconstructors", (c, ext, js, extra) -> new OptimizeConstructors(c).process(ext, js)));
    passes.add(
        new Pass(
            "devirtualize",
            (c, ext, js, extra) ->
                OptimizeCalls.builder()
                    .setCompiler(c)
                    .setConsiderExterns(false)
                    .addPass(new DevirtualizeMethods(c))
                    .build()
                    .process(ext, js)));
    passes.add(
        new Pass(
            "rewrite-functions", (c, ext, js, extra) -> new FunctionRewriter(c).process(ext, js)));
    return passes;
  }

  private static List<Pass> chunkPasses() {
    List<Pass> passes = new ArrayList<>();
    passes.add(
        new Pass(
            "crosschunk-refs",
            (c, ext, js, extra) -> {
              CrossChunkReferenceCollector collector =
                  new CrossChunkReferenceCollector(c, new SyntacticScopeCreator(c));
              collector.process(js);
              for (CrossChunkReferenceCollector.TopLevelStatement s :
                  collector.getTopLevelStatements()) {
                StringBuilder l = new StringBuilder();
                l.append("stmt ")
                    .append(s.getOriginalOrder())
                    .append(' ')
                    .append(s.getChunk().getName())
                    .append(' ')
                    .append(pos(s.getStatementNode()));
                if (s.isDeclarationStatement()) {
                  l.append(" decl ")
                      .append(pos(s.getDeclaredNameReference().getNode()))
                      .append(" value ")
                      .append(s.getDeclaredValueNode() == null ? "-" : pos(s.getDeclaredValueNode()))
                      .append(" movable ")
                      .append(s.isMovableDeclaration());
                }
                l.append(" refs");
                for (Reference r : s.getNonDeclarationReferences()) {
                  l.append(' ').append(pos(r.getNode()));
                }
                extra.add(l.toString());
              }
              for (Map.Entry<String, Var> e : collector.getGlobalVariableNamesMap().entrySet()) {
                ReferenceCollection refs = collector.getReferences(e.getValue());
                StringBuilder l = new StringBuilder();
                l.append("var ").append(e.getKey());
                if (refs == null) {
                  l.append(" none");
                } else {
                  l.append(" wellDefined ")
                      .append(refs.isWellDefined())
                      .append(" assignedOnce ")
                      .append(refs.isAssignedOnceInLifetime())
                      .append(" refs");
                  for (Reference r : refs) {
                    l.append(' ').append(pos(r.getNode()));
                  }
                }
                extra.add(l.toString());
              }
            }));
    passes.add(
        new Pass(
            "crosschunk-code",
            (c, ext, js, extra) ->
                new CrossChunkCodeMotion(c, c.getChunkGraph(), false).process(ext, js)));
    passes.add(
        new Pass(
            "crosschunk-method",
            (c, ext, js, extra) ->
                new CrossChunkMethodMotion(c, new IdGenerator(), false, false).process(ext, js)));
    return passes;
  }

  /** LiveVariablesAnalysis for every function, as DeadAssignmentsElimination computes it. */
  private static final class Lva extends NodeTraversal.AbstractCfgCallback {
    private final Compiler compiler;
    private final List<String> extra;

    Lva(Compiler compiler, List<String> extra) {
      this.compiler = compiler;
      this.extra = extra;
    }

    @Override
    public void exitScopeWithCfg(NodeTraversal t) {
      if (!t.inFunctionBlockScope()) {
        return;
      }
      Scope blockScope = t.getScope();
      Scope functionScope = blockScope.getParent();
      Node function = t.getScopeRoot().getParent();
      ControlFlowGraph<Node> cfg = getControlFlowGraph(compiler);
      SyntacticScopeCreator scopeCreator = new SyntacticScopeCreator(compiler);
      AllVarsDeclaredInFunction allVarsDeclaredInFunction =
          NodeUtil.getAllVarsDeclaredInFunction(compiler, scopeCreator, functionScope);
      LiveVariablesAnalysis liveness =
          new LiveVariablesAnalysis(
              cfg, functionScope, blockScope, compiler, scopeCreator, allVarsDeclaredInFunction);
      liveness.analyze();
      List<String> escaped = new ArrayList<>();
      for (Var v : liveness.getEscapedLocals()) {
        escaped.add(v.getName());
      }
      Collections.sort(escaped);
      extra.add(
          "function "
              + pos(function)
              + " "
              + function.getFirstChild().getString()
              + " vars "
              + liveness.getAllVariablesInOrder().size()
              + " escaped "
              + String.join(",", escaped));
      for (Var v : liveness.getAllVariablesInOrder()) {
        extra.add("var " + liveness.getVarIndex(v.getName()) + " " + v.getName());
      }
      for (DiGraphNode<Node, ControlFlowGraph.Branch> cfgNode : cfg.getNodes()) {
        DataFlowAnalysis.LinearFlowState<LiveVariablesAnalysis.LiveVariableLattice> state =
            cfgNode.getAnnotation();
        Node n = cfgNode.getValue();
        extra.add(
            "node "
                + (n == null ? "IMPLICIT_RETURN" : pos(n))
                + " in "
                + bits(state.getIn())
                + " out "
                + bits(state.getOut()));
      }
    }
  }

  private static String bits(LiveVariablesAnalysis.LiveVariableLattice l) {
    StringBuilder sb = new StringBuilder("{");
    for (int i = l.nextSetBit(0); i >= 0; i = l.nextSetBit(i + 1)) {
      if (sb.length() > 1) {
        sb.append(',');
      }
      sb.append(i);
    }
    return sb.append('}').toString();
  }

  private static void collectFlags(Node n, List<String> extra) {
    if (n.isCall() || n.isNew() || n.isTaggedTemplateLit() || n.isOptChainCall()) {
      extra.add("flags " + pos(n) + " " + n.getSideEffectFlags());
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      collectFlags(c, extra);
    }
  }

  private static String pos(Node n) {
    return n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
  }

  private static CompilerOptions newOptions() {
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    // VarCheck declares the undeclared names (as the checks phase does before these passes run);
    // its "undefined variable" errors would only repeat in every run.
    options.setWarningLevel(DiagnosticGroups.UNDEFINED_VARIABLES, CheckLevel.OFF);
    return options;
  }

  private static Compiler newCompiler(String rel, String src) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    compiler.init(
        ImmutableList.of(SourceFile.fromCode("externs.js", EXTERNS)),
        ImmutableList.of(SourceFile.fromCode(rel, src)),
        newOptions());
    compiler.parse();
    return compiler;
  }

  private static Compiler newChunkCompiler(String[] parts) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    List<JSChunk> chunks = new ArrayList<>();
    for (int k = 0; k < parts.length; k++) {
      JSChunk chunk = new JSChunk("m" + (k + 1));
      chunk.add(SourceFile.fromCode("m" + (k + 1) + ".js", parts[k]));
      if (k > 0) {
        chunk.addDependency(chunks.get(k - 1));
      }
      chunks.add(chunk);
    }
    compiler.initChunks(
        ImmutableList.of(SourceFile.fromCode("externs.js", EXTERNS)), chunks, newOptions());
    compiler.parse();
    return compiler;
  }

  private static void prepare(Compiler c) {
    Node externs = c.getExternsRoot();
    Node js = c.getJsRoot();
    new VarCheck(c).process(externs, js);
    Normalize.createNormalizeForOptimizations(c).process(externs, js);
    new GatherGetterAndSetterProperties(c).process(externs, js);
    new GatherExternProperties(c, GatherExternProperties.Mode.OPTIMIZE).process(externs, js);
  }

  private static List<String> diagnostics(Compiler compiler) {
    List<String> out = new ArrayList<>();
    for (JSError e : compiler.getErrors()) {
      out.add(diag("ERROR", e));
    }
    for (JSError e : compiler.getWarnings()) {
      out.add(diag("WARNING", e));
    }
    return out;
  }

  private static String diag(String level, JSError e) {
    return "diag "
        + level
        + " "
        + e.getType().key
        + " "
        + e.getLineno()
        + ":"
        + e.getCharno()
        + " \""
        + e.getDescription()
        + "\"";
  }

  /**
   * Replaces the "[change_time: N]" annotation of a toStringTree() line by "[change_time]": the
   * stamps depend on the order in which RemoveUnusedCode removes the values of its
   * HashMultimap&lt;String, Removable&gt;, which Java iterates in identity-hash order, so they differ
   * from run to run of the JVM. Whether a node has a change time is kept.
   */
  static String dropChangeTime(String l) {
    String key = " [change_time: ";
    int i = l.indexOf(key);
    if (i < 0) {
      return l;
    }
    int end = l.indexOf(']', i);
    return l.substring(0, i) + " [change_time]" + l.substring(end + 1);
  }

  private static void collectChanges(Node n, List<String> out) {
    if ((n.isScript() || n.isFunction()) && n.getChangeTime() != 0) {
      out.add(pos(n));
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      collectChanges(c, out);
    }
  }

  private static String sha16(List<String> lines) throws Exception {
    MessageDigest md = MessageDigest.getInstance("SHA-256");
    for (String l : lines) {
      md.update((l + "\n").getBytes(StandardCharsets.UTF_8));
    }
    StringBuilder hex = new StringBuilder();
    for (byte b : md.digest()) {
      hex.append(String.format("%02x", b));
    }
    return hex.substring(0, 16);
  }

  /** The split points of the 3-chunk runs, or null (with the reason in skip[0]). */
  private static int[] chunkOffsets(Compiler parsed, String[] skip) {
    Node script = parsed.getJsRoot().getFirstChild();
    if (script.hasChildren() && script.getFirstChild().isModuleBody()) {
      skip[0] = "module";
      return null;
    }
    List<Node> statements = new ArrayList<>();
    for (Node c = script.getFirstChild(); c != null; c = c.getNext()) {
      statements.add(c);
    }
    int n = statements.size();
    if (n < 3) {
      skip[0] = "small";
      return null;
    }
    int[] offsets = new int[2];
    for (int k = 1; k <= 2; k++) {
      Node last = statements.get(k * n / 3 - 1);
      offsets[k - 1] = last.getSourceOffset() + last.getLength();
    }
    return offsets;
  }

  private static void runPass(StringBuilder sb, Pass pass, Compiler c, List<String> parseDiags)
      throws Exception {
    line(sb, "## pass " + pass.name);
    List<String> extra = new ArrayList<>();
    String before = null;
    String exception = null;
    try {
      prepare(c);
      before = c.getRoot().toStringTree();
      pass.run.run(c, c.getExternsRoot(), c.getJsRoot(), extra);
    } catch (Throwable t) {
      exception = "exception " + t.getClass().getName() + ": " + t.getMessage();
    }
    List<String> diags = diagnostics(c);
    for (String d : parseDiags) {
      diags.remove(d);
    }
    for (String d : diags) {
      line(sb, d);
    }
    if (exception != null) {
      line(sb, exception);
      return;
    }
    String after = c.getRoot().toStringTree();
    if (after.equals(before)) {
      line(sb, "tree unchanged");
    } else {
      List<String> lines = new ArrayList<>();
      for (String l : after.split("\n", -1)) {
        lines.add(escape(dropChangeTime(l)));
      }
      if (lines.get(lines.size() - 1).isEmpty()) {
        lines.remove(lines.size() - 1);
      }
      line(sb, "tree changed " + sha16(lines));
      if (full) {
        for (String l : lines) {
          sb.append(l).append('\n');
        }
      }
    }
    List<String> changes = new ArrayList<>();
    collectChanges(c.getRoot(), changes);
    emitBlock(sb, "changes", changes);
    emitBlock(sb, "extra", extra);
  }

  private static void emitBlock(StringBuilder sb, String kind, List<String> raw) throws Exception {
    List<String> escaped = new ArrayList<>();
    for (String l : raw) {
      escaped.add(escape(l));
    }
    line(sb, kind + " " + raw.size() + " " + sha16(escaped));
    if (full) {
      for (String l : escaped) {
        sb.append(l).append('\n');
      }
    }
  }

  private static void dumpFile(StringBuilder sb, String rel, String src) throws Exception {
    line(sb, "## file " + rel);
    Compiler compiler = newCompiler(rel, src);
    List<String> parseDiags = diagnostics(compiler);
    line(sb, "# parse " + (compiler.hasErrors() ? "errors" : "ok"));
    for (String d : parseDiags) {
      line(sb, d);
    }
    if (compiler.hasErrors()) {
      return;
    }
    for (Pass pass : passes()) {
      runPass(sb, pass, newCompiler(rel, src), parseDiags);
    }
    String[] skip = new String[1];
    int[] offsets = chunkOffsets(compiler, skip);
    if (offsets == null) {
      line(sb, "## chunks skipped " + skip[0]);
      return;
    }
    if (!(0 <= offsets[0] && offsets[0] <= offsets[1] && offsets[1] <= src.length())) {
      // A statement's source range can end after the next one's (e.g. when the parser's
      // positions come from a source with an unusual layout); such files get no chunk runs.
      line(sb, "## chunks skipped offsets " + offsets[0] + " " + offsets[1]);
      return;
    }
    line(sb, "## chunks " + offsets[0] + " " + offsets[1]);
    String[] parts = {
      src.substring(0, offsets[0]),
      src.substring(offsets[0], offsets[1]),
      src.substring(offsets[1])
    };
    Compiler probe = newChunkCompiler(parts);
    List<String> chunkParseDiags = diagnostics(probe);
    line(sb, "# chunks parse " + (probe.hasErrors() ? "errors" : "ok"));
    for (String d : chunkParseDiags) {
      line(sb, d);
    }
    if (probe.hasErrors()) {
      return;
    }
    for (Pass pass : chunkPasses()) {
      runPass(sb, pass, newChunkCompiler(parts), chunkParseDiags);
    }
  }

  private static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }

  static String escape(String s) {
    StringBuilder sb = new StringBuilder(s.length());
    for (int i = 0; i < s.length(); i++) {
      char ch = s.charAt(i);
      if (ch == '\\') {
        sb.append("\\\\");
      } else if (ch >= 0x20 && ch < 0x7f) {
        sb.append(ch);
      } else {
        sb.append(String.format("\\u%04x", (int) ch));
      }
    }
    return sb.toString();
  }
}

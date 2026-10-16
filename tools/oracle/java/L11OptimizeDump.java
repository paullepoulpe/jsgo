// Differential oracle for the l11-optimize passes (internal/jscomp): CollapseAnonymousFunctions,
// CollapseVariableDeclarations, Denormalize, DeadPropertyAssignmentElimination,
// InlineObjectLiterals, OptimizeCalls (its ReferenceMap), FunctionArgumentInjector,
// ExtractPrototypeMemberDeclarations, AnalyzePrototypeProperties, AliasStrings, RenameLabels,
// RenameVars (its VariableMap), ConcretizeStaticInheritanceForInlining, InferConsts,
// RemoveCastNodes, RemovePropertyRenamingCalls and RemoveUnnecessarySyntheticExterns.
// (InlineProperties needs type checking and colors, and RemoveWeakSources weak inputs: both are
// covered by their unit tests only.)
//
// For every file, each run starts from a fresh Compiler set up the way CompilerTestCase sets one
// up: default CompilerOptions with --language_in=ECMASCRIPT_NEXT and --language_out=NO_TRANSPILE,
// Compiler.init(no externs, [SourceFile.fromCode(<rel>, <src>)], options) and Compiler.parse().
// Files with parse errors print their diagnostics and are skipped. The prerequisites are the ones
// DefaultPassConfig runs before each pass:
//   raw         nothing (the passes DefaultPassConfig runs before normalization, and the
//               unnormalized unit-test setting of RenameVars/RenameLabels)
//   normalized  Normalize.createNormalizeForOptimizations; then, for
//               DeadPropertyAssignmentElimination and OptimizeCalls,
//               GatherGetterAndSetterProperties and GatherExternProperties(OPTIMIZE)
//   late        normalized, then markUnnormalized (LifeCycleStage.RAW) — for Denormalize and
//               CollapseVariableDeclarations
//   pipeline    late, then Denormalize(outputFeatureSet), MakeDeclaredNamesUnique's contextual
//               rename inverter, RenameVars (configured like DefaultPassConfig.runVariableRenaming
//               with default options) and RenameLabels
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
//   <pass-specific lines>
//
// Usage: java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//          com.google.javascript.jscomp.L11OptimizeDump \
//          [--full] [--sample=<N>] <out.gz> <root> <rel-file>... | --list=<listfile>
// -XX:-ShowCodeDetailsInExceptionMessages turns off the JVM's "helpful" NullPointerException
// messages (they describe bytecode, and a Go nil dereference has none): NPEs then have a null
// message, as jutil.AsThrowable gives them.
// --sample=N dumps every Nth file of the list only (in list order, starting with the first).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage;
import com.google.javascript.jscomp.CompilerOptions.AliasStringsMode;
import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.FunctionArgumentInjector.ParamArgPair;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableListMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableMap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSet;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.function.Supplier;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class L11OptimizeDump {

  /** A pass run: gets the compiler and (externsRoot, jsRoot); may append extra output lines. */
  private interface Run {
    void run(Compiler compiler, Node externs, Node js, List<String> extra);
  }

  private enum Stage {
    RAW,
    NORMALIZED,
    NORMALIZED_WITH_PROPERTIES,
    LATE,
    PIPELINE
  }

  private static final class Pass {
    final String name;
    final Stage stage;
    final Run run;

    Pass(String name, Stage stage, Run run) {
      this.name = name;
      this.stage = stage;
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
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
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

  private static List<Pass> passes() {
    List<Pass> p = new ArrayList<>();
    // raw
    p.add(
        new Pass(
            "RemoveCastNodes",
            Stage.RAW,
            (c, externs, js, extra) -> new RemoveCastNodes(c).process(externs, js)));
    p.add(
        new Pass(
            "RemovePropertyRenamingCalls",
            Stage.RAW,
            (c, externs, js, extra) -> new RemovePropertyRenamingCalls(c).process(externs, js)));
    p.add(
        new Pass(
            "RemoveUnnecessarySyntheticExterns",
            Stage.RAW,
            (c, externs, js, extra) ->
                new RemoveUnnecessarySyntheticExterns(c).process(externs, js)));
    p.add(
        new Pass(
            "RenameLabels",
            Stage.RAW,
            (c, externs, js, extra) -> new RenameLabels(c).process(externs, js)));
    p.add(
        new Pass(
            "RenameVars",
            Stage.RAW,
            (c, externs, js, extra) -> renameVars(c, externs, js, extra)));
    // normalized
    p.add(
        new Pass(
            "CollapseAnonymousFunctions",
            Stage.NORMALIZED,
            (c, externs, js, extra) -> new CollapseAnonymousFunctions(c).process(externs, js)));
    for (ExtractPrototypeMemberDeclarations.Pattern pattern :
        ExtractPrototypeMemberDeclarations.Pattern.values()) {
      p.add(
          new Pass(
              "ExtractPrototypeMemberDeclarations/" + pattern,
              Stage.NORMALIZED,
              (c, externs, js, extra) ->
                  new ExtractPrototypeMemberDeclarations(c, pattern).process(externs, js)));
    }
    for (AliasStringsMode mode :
        new AliasStringsMode[] {
          AliasStringsMode.ALL, AliasStringsMode.LARGE, AliasStringsMode.ALL_AGGRESSIVE
        }) {
      p.add(
          new Pass(
              "AliasStrings/" + mode,
              Stage.NORMALIZED,
              (c, externs, js, extra) ->
                  new AliasStrings(c, c.getChunkGraph(), false, mode).process(externs, js)));
    }
    p.add(
        new Pass(
            "InlineObjectLiterals",
            Stage.NORMALIZED,
            (c, externs, js, extra) ->
                new InlineObjectLiterals(c, c.getUniqueNameIdSupplier()).process(externs, js)));
    p.add(
        new Pass(
            "DeadPropertyAssignmentElimination",
            Stage.NORMALIZED_WITH_PROPERTIES,
            (c, externs, js, extra) ->
                new DeadPropertyAssignmentElimination(c).process(externs, js)));
    p.add(
        new Pass(
            "OptimizeCalls",
            Stage.NORMALIZED_WITH_PROPERTIES,
            (c, externs, js, extra) -> optimizeCalls(c, externs, js, extra)));
    p.add(
        new Pass(
            "AnalyzePrototypeProperties",
            Stage.NORMALIZED,
            (c, externs, js, extra) -> analyzePrototypeProperties(c, externs, js, extra)));
    p.add(
        new Pass(
            "ConcretizeStaticInheritanceForInlining",
            Stage.NORMALIZED,
            (c, externs, js, extra) ->
                new ConcretizeStaticInheritanceForInlining(c).process(externs, js)));
    p.add(
        new Pass(
            "InferConsts",
            Stage.NORMALIZED,
            (c, externs, js, extra) -> {
              new InferConsts(c).process(externs, js);
              dumpConsts(js, extra);
            }));
    p.add(
        new Pass(
            "FunctionArgumentInjector",
            Stage.NORMALIZED,
            (c, externs, js, extra) -> functionArgumentInjector(c, js, extra)));
    // late
    p.add(
        new Pass(
            "Denormalize",
            Stage.LATE,
            (c, externs, js, extra) ->
                new Denormalize(c, c.getOptions().getOutputFeatureSet()).process(externs, js)));
    p.add(
        new Pass(
            "CollapseVariableDeclarations",
            Stage.LATE,
            (c, externs, js, extra) -> new CollapseVariableDeclarations(c).process(externs, js)));
    // pipeline
    p.add(
        new Pass(
            "RenameVarsAndLabels",
            Stage.PIPELINE,
            (c, externs, js, extra) -> {
              renameVars(c, externs, js, extra);
              new RenameLabels(c).process(externs, js);
            }));
    return p;
  }

  /** RenameVars configured like DefaultPassConfig.runVariableRenaming with default options. */
  private static void renameVars(Compiler c, Node externs, Node js, List<String> extra) {
    CompilerOptions options = c.getOptions();
    Set<Character> reservedChars = ImmutableSet.of();
    Set<String> reservedNames = new LinkedHashSet<>();
    if (options.getRenamePrefixNamespace() != null) {
      reservedNames.add(options.getRenamePrefixNamespace());
    }
    reservedNames.addAll(c.getExportedNames());
    reservedNames.addAll(ParserRunner.getReservedVars());
    RenameVars rn =
        new RenameVars(
            c,
            options.getRenamePrefix(),
            options.getVariableRenaming() == VariableRenamingPolicy.LOCAL,
            options.shouldGeneratePseudoNames(),
            options.shouldPreferStableNames(),
            options.getInputVariableMap(),
            reservedChars,
            reservedNames,
            options.getNameGenerator());
    rn.process(externs, js);
    for (Map.Entry<String, String> e :
        rn.getVariableMap().getOriginalNameToNewNameMap().entrySet()) {
      extra.add("map " + e.getKey() + " -> " + e.getValue());
    }
  }

  private static void optimizeCalls(Compiler c, Node externs, Node js, List<String> extra) {
    for (boolean considerExterns : new boolean[] {false, true}) {
      OptimizeCalls.builder()
          .setCompiler(c)
          .setConsiderExterns(considerExterns)
          .addPass(
              (e, r, references) -> {
                extra.add("references considerExterns=" + considerExterns);
                for (Map.Entry<String, ArrayList<Node>> entry : references.getNameReferences()) {
                  extra.add("name " + entry.getKey() + refs(entry.getValue()));
                  functionNodes(entry.getValue(), extra);
                }
                for (Map.Entry<String, ArrayList<Node>> entry : references.getPropReferences()) {
                  extra.add("prop " + entry.getKey() + refs(entry.getValue()));
                  functionNodes(entry.getValue(), extra);
                }
                for (Map.Entry<String, ArrayList<Node>> entry : references.getNameReferences()) {
                  for (Node n : entry.getValue()) {
                    if (n.getParent() != null) {
                      callTarget(n, extra);
                    }
                  }
                }
              })
          .build()
          .process(externs, js);
    }
  }

  private static String refs(List<Node> nodes) {
    StringBuilder sb = new StringBuilder();
    for (Node n : nodes) {
      sb.append(' ').append(pos(n));
    }
    return sb.toString();
  }

  private static void functionNodes(List<Node> nodes, List<String> extra) {
    ImmutableListMultimap<Node, Node> fns = OptimizeCalls.ReferenceMap.getFunctionNodes(nodes);
    for (Map.Entry<Node, Node> e : fns.entries()) {
      extra.add("  fn " + pos(e.getKey()) + " -> " + pos(e.getValue()));
    }
  }

  private static void callTarget(Node n, List<String> extra) {
    StringBuilder sb = new StringBuilder("  target ").append(pos(n));
    boolean any = false;
    if (OptimizeCalls.ReferenceMap.isNormalOrOptChainCallOrNewTarget(n)) {
      any = true;
      sb.append(" call=").append(pos(OptimizeCalls.ReferenceMap.getCallOrNewNodeForTarget(n)));
      Node arg = OptimizeCalls.ReferenceMap.getFirstArgumentForCallOrNewOrDotCall(n);
      sb.append(" arg0=").append(arg == null ? "null" : pos(arg));
    }
    if (OptimizeCalls.ReferenceMap.isSimpleAssignmentTarget(n)) {
      any = true;
      sb.append(" simpleAssign");
    }
    if (OptimizeCalls.isAllowedReference(n)) {
      any = true;
      sb.append(" allowed");
    }
    if (any) {
      extra.add(sb.toString());
    }
  }

  private static void analyzePrototypeProperties(
      Compiler c, Node externs, Node js, List<String> extra) {
    AnalyzePrototypeProperties app =
        new AnalyzePrototypeProperties(
            c,
            c.getChunkGraph(),
            /* canModifyExterns= */ false,
            /* anchorUnusedVars= */ false,
            /* rootScopeUsesAreGlobal= */ false);
    app.process(externs, js);
    for (AnalyzePrototypeProperties.NameInfo info : app.getAllNameInfo()) {
      JSChunk chunk = info.getDeepestCommonChunkRef();
      extra.add(
          "info "
              + info
              + " referenced="
              + info.isReferenced()
              + " readsClosureVariables="
              + info.readsClosureVariables()
              + " referencesSuper="
              + info.referencesSuper()
              + " chunk="
              + (chunk == null ? "null" : chunk.getName()));
      for (AnalyzePrototypeProperties.Symbol s : info.getDeclarations()) {
        String kind;
        Node node;
        if (s instanceof AnalyzePrototypeProperties.GlobalFunction) {
          kind = "GlobalFunction";
          node = null;
        } else if (s instanceof AnalyzePrototypeProperties.ClassMemberFunction cmf) {
          kind = "ClassMemberFunction";
          node = cmf.getDefinitionNode();
        } else {
          AnalyzePrototypeProperties.PrototypeProperty pp =
              (AnalyzePrototypeProperties.PrototypeProperty) s;
          kind = s.getClass().getSimpleName() + " prototype=" + pos(pp.getPrototype());
          node = pp.getValue();
        }
        Var v = s.getRootVar();
        extra.add(
            "  decl "
                + kind
                + (node == null ? "" : " node=" + pos(node))
                + " rootVar="
                + (v == null ? "null" : v.getName())
                + " chunk="
                + (s.getChunk() == null ? "null" : s.getChunk().getName()));
      }
    }
  }

  private static void dumpConsts(Node n, List<String> extra) {
    if ((n.isName() || n.isImportStar())
        && (n.isDeclaredConstantVar() || n.isInferredConstantVar())) {
      extra.add(
          "const "
              + pos(n)
              + " "
              + n.getString()
              + " declared="
              + n.isDeclaredConstantVar()
              + " inferred="
              + n.isInferredConstantVar());
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      dumpConsts(c, extra);
    }
  }

  private static final int MAX_INJECTOR_PAIRS = 200;

  /**
   * For every named function declaration with plain (NAME or ...NAME) parameters and every direct
   * call to it by name in the same script (in preorder), prints what FunctionInjector asks of
   * FunctionArgumentInjector.
   */
  private static void functionArgumentInjector(Compiler c, Node js, List<String> extra) {
    FunctionArgumentInjector injector = new FunctionArgumentInjector(c.getAstAnalyzer());
    List<Node> fns = new ArrayList<>();
    List<Node> calls = new ArrayList<>();
    collect(js, fns, calls);
    int pairs = 0;
    for (Node fn : fns) {
      String name = fn.getFirstChild().getString();
      for (Node call : calls) {
        if (!call.getFirstChild().getString().equals(name)) {
          continue;
        }
        if (pairs++ >= MAX_INJECTOR_PAIRS) {
          return;
        }
        int[] counter = {0};
        Supplier<String> ids = () -> String.valueOf(counter[0]++);
        ImmutableMap<String, ParamArgPair> map =
            injector.getFunctionCallParameterMap(fn, call, ids);
        StringBuilder sb = new StringBuilder("inject fn=").append(pos(fn));
        sb.append(" call=").append(pos(call)).append(" map");
        for (Map.Entry<String, ParamArgPair> e : map.entrySet()) {
          sb.append(' ')
              .append(e.getKey())
              .append('=')
              .append(e.getValue().arg().getToken())
              .append('@')
              .append(pos(e.getValue().arg()));
        }
        extra.add(sb.toString());
        ImmutableSet<String> modified = injector.findModifiedParameters(fn);
        extra.add("  modified " + modified);
        ImmutableSet<String> temps =
            injector.gatherCallArgumentsNeedingTemps(
                c, fn, map, modified, c.getCodingConvention());
        extra.add("  temps " + temps);
        extra.add(
            "  conditional "
                + injector.bodyMayHaveConditionalCode(fn.getLastChild())
                + " "
                + injector.mayHaveConditionalCode(fn.getLastChild()));
      }
    }
  }

  private static void collect(Node n, List<Node> fns, List<Node> calls) {
    if (NodeUtil.isFunctionDeclaration(n) && !n.getFirstChild().getString().isEmpty()) {
      boolean plain = true;
      for (Node p = NodeUtil.getFunctionParameters(n).getFirstChild(); p != null; p = p.getNext()) {
        if (!(p.isName() || (p.isRest() && p.getOnlyChild().isName()))) {
          plain = false;
        }
      }
      if (plain) {
        fns.add(n);
      }
    }
    if (n.isCall() && n.getFirstChild().isName()) {
      calls.add(n);
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      collect(c, fns, calls);
    }
  }

  private static String pos(Node n) {
    return n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
  }

  private static Compiler newCompiler(String rel, String src) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    compiler.init(
        ImmutableList.<SourceFile>of(), ImmutableList.of(SourceFile.fromCode(rel, src)), options);
    compiler.parse();
    return compiler;
  }

  private static void prepare(Compiler c, Stage stage) {
    Node externs = c.getExternsRoot();
    Node js = c.getJsRoot();
    if (stage == Stage.RAW) {
      return;
    }
    Normalize.createNormalizeForOptimizations(c).process(externs, js);
    if (stage == Stage.NORMALIZED_WITH_PROPERTIES) {
      new GatherGetterAndSetterProperties(c).process(externs, js);
      new GatherExternProperties(c, GatherExternProperties.Mode.OPTIMIZE).process(externs, js);
    }
    if (stage == Stage.LATE || stage == Stage.PIPELINE) {
      c.setLifeCycleStage(LifeCycleStage.RAW);
    }
    if (stage == Stage.PIPELINE) {
      new Denormalize(c, c.getOptions().getOutputFeatureSet()).process(externs, js);
      MakeDeclaredNamesUnique.getContextualRenameInverter(c).process(externs, js);
    }
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
      line(sb, "## pass " + pass.name);
      Compiler c = newCompiler(rel, src);
      List<String> extra = new ArrayList<>();
      String before = null;
      String exception = null;
      try {
        prepare(c, pass.stage);
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
        continue;
      }
      String after = c.getRoot().toStringTree();
      if (after.equals(before)) {
        line(sb, "tree unchanged");
      } else {
        List<String> lines = new ArrayList<>();
        for (String l : after.split("\n", -1)) {
          lines.add(escape(l));
        }
        if (lines.get(lines.size() - 1).isEmpty()) {
          lines.remove(lines.size() - 1);
        }
        MessageDigest md = MessageDigest.getInstance("SHA-256");
        for (String l : lines) {
          md.update((l + "\n").getBytes(StandardCharsets.UTF_8));
        }
        StringBuilder hex = new StringBuilder();
        for (byte b : md.digest()) {
          hex.append(String.format("%02x", b));
        }
        line(sb, "tree changed " + hex.substring(0, 16));
        if (full) {
          for (String l : lines) {
            sb.append(l).append('\n');
          }
        }
      }
      for (String e : extra) {
        line(sb, e);
      }
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

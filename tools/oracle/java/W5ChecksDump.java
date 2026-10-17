// Oracle tool for the w5-checks cluster (Go: internal/jscomp/w5_checks_oracle_test.go):
// CheckJSDoc, CheckJSDocStyle, VariableReferenceCheck, CheckSideEffects (+ StripProtection),
// CheckAccessControls and CheckConformance (with ConformanceRules), run the way DefaultPassConfig
// runs them in getChecks(), on type-checked code.
//
// For each input, one Compiler (Compiler.init + parseInputs) with these options: those of
// TypeCheckDump (CompilerTypeTestCase.defaultOptions() + TypeCheckTestCase.setUp's levels +
// setCheckTypes(true)), JSDoc descriptions parsed (INCLUDE_DESCRIPTIONS_NO_WHITESPACE, as
// CheckJSDocStyleTest), and these levels so that the checks' diagnostics are visible:
// MISPLACED_SUPPRESS (as CheckJsDocTest), CheckJSDocStyle.ALL_DIAGNOSTICS (as CheckJSDocStyleTest),
// ACCESS_CONTROLS, DEPRECATED and CONSTANT_PROPERTY (as CheckAccessControlsTest, but at WARNING),
// and CheckJSDoc.MISPLACED_MSG_ANNOTATION and VariableReferenceCheck.UNUSED_LOCAL_ASSIGNMENT (both
// disabled by default) at WARNING; ModuleLoader.LOAD_WARNING and INVALID_MODULE_PATH are off (as
// ChecksDump), so that ES modules importing files outside the compilation still parse. The
// externs are the --externs files, in order, each a SourceFile named by its file name.
//
// Then, in DefaultPassConfig.getChecks() order, each in its own "## <name>" section:
//   modules                 GatherModuleMetadata(processCommonJsModules=false, BROWSER) and
//                           ModuleMapCreator (gatherModuleMetadataPass, createModuleMapPass)
//   CheckJSDoc              NodeTraversal.traverse(compiler, mainRoot, new CheckJSDoc(compiler)),
//                           as the combined checkJsDocAndEs6Modules pass runs it
//   CheckJSDocStyle         the same, as the combined lintChecks pass runs it
//   VariableReferenceCheck  new VariableReferenceCheck(compiler).process(externs, main)
//   CheckSideEffects        new CheckSideEffects(compiler, true, true).process(externs, main) (as
//                           CheckSideEffectsTest: report and protect), then "= protector
//                           <source>:<line>:<charno>" for each JSCOMPILER_PRESERVE call in the main
//                           AST (pre-order) and "= extern <line>" for the synthesized externs tree
//   TypeCheck               InferConsts, then the bodies of DefaultPassConfig.inferTypes and
//                           checkTypes, as TypeCheckDump does
//   CheckAccessControls     new CheckAccessControls(compiler).process(externs, main)
//   CheckConformance <name> new CheckConformance(compiler, ImmutableList.of(<config>),
//                           options.getConformanceReportingMode()).process(externs, main), for each
//                           --conformance config (TextFormat.merge of the file), in order; the
//                           step's diagnostic lines are sorted (String order): CheckConformance
//                           groups its rules in a HashMultimap, so the order in which two rules
//                           report on the same node depends on identity hash codes
//   StripProtection         new CheckSideEffects.StripProtection(compiler).process(externs, main),
//                           then "= protectors <n>" (JSCOMPILER_PRESERVE calls left in the main AST)
//
// An input whose text starts with "//// file: <name>" is split at every such line into several
// inputs (the marker lines are dropped), so that hand-written inputs can exercise cross-file
// rules (visibility, allowlists); any other input is one SourceFile named <rel-file>.
//
// Build/run (-XX:-ShowCodeDetailsInExceptionMessages: a NullPointerException has no message):
//   javac -cp closure-compiler.jar -d build W5ChecksDump.java
//   java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:build \
//       com.google.javascript.jscomp.W5ChecksDump --externs=<file>... --conformance=<name>:<file>... \
//       <out[.gz]> <root> <rel-file>... | @<listfile>
//
// Format (every line escaped like ParsingDump: '\' -> "\\", chars outside printable ASCII ->
// \\uXXXX):
//   # file <rel>
//   ## <pass>                          before each step (## parse first)
//   <LEVEL> <key> <source>:<line>:<charno> @<offset>+<length> <description>
//                                      each diagnostic reaching the ErrorManager, in report order
//   # parse-errors                     if parseInputs() reported errors (nothing else runs)
//   ! <exception class>                if the step threw, then a line per cause:
//   !! cause <class>[: <message>]      (no message for NullPointerException)
//   = <text>                           extra step output (see above)
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.deps.ModuleLoader;
import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.lint.CheckJSDocStyle;
import com.google.javascript.jscomp.modules.ModuleMapCreator;
import com.google.javascript.jscomp.parsing.Config.JsDocParsing;
import com.google.javascript.rhino.Node;
import com.google.javascript.jscomp.jarjar.com.google.protobuf.TextFormat;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class W5ChecksDump {
  static final String MARKER = "//// file: ";

  record NamedConfig(String name, ConformanceConfig config) {}

  public static void main(String[] args) throws Exception {
    List<SourceFile> externs = new ArrayList<>();
    List<NamedConfig> configs = new ArrayList<>();
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.startsWith("--externs=")) {
        Path p = Paths.get(a.substring("--externs=".length()));
        externs.add(
            SourceFile.fromCode(
                p.getFileName().toString(), Files.readString(p, StandardCharsets.UTF_8)));
      } else if (a.startsWith("--conformance=")) {
        String spec = a.substring("--conformance=".length());
        int colon = spec.indexOf(':');
        ConformanceConfig.Builder builder = ConformanceConfig.newBuilder();
        TextFormat.merge(
            Files.readString(Paths.get(spec.substring(colon + 1)), StandardCharsets.UTF_8),
            builder);
        configs.add(new NamedConfig(spec.substring(0, colon), builder.build()));
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
    if (outFile.getParent() != null) {
      Files.createDirectories(outFile.getParent());
    }
    try (OutputStream os =
        outFile.toString().endsWith(".gz")
            ? new GZIPOutputStream(Files.newOutputStream(outFile))
            : Files.newOutputStream(outFile)) {
      for (String rel : files) {
        StringBuilder sb = new StringBuilder();
        dumpOne(
            sb, externs, configs, rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8));
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  static CompilerOptions options() {
    CompilerOptions options = new CompilerOptions();
    // TypeCheckDump.options()
    options.setCodingConvention(new GoogleCodingConvention());
    options.setLanguage(LanguageMode.UNSUPPORTED);
    options.setWarningLevel(DiagnosticGroups.MISSING_PROPERTIES, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.MISPLACED_TYPE_ANNOTATION, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.INVALID_CASTS, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.LINT_CHECKS, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.JSDOC_MISSING_TYPE, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.BOUNDED_GENERICS, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.MISSING_OVERRIDE, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.STRICT_MISSING_PROPERTIES, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.STRICT_PRIMITIVE_OPERATORS, CheckLevel.WARNING);
    options.setCheckTypes(true);
    // The checks of this tool.
    options.setParseJsDocDocumentation(JsDocParsing.INCLUDE_DESCRIPTIONS_NO_WHITESPACE);
    options.setWarningLevel(DiagnosticGroups.MISPLACED_SUPPRESS, CheckLevel.WARNING);
    options.setWarningLevel(CheckJSDocStyle.ALL_DIAGNOSTICS, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.ACCESS_CONTROLS, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.DEPRECATED, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.CONSTANT_PROPERTY, CheckLevel.WARNING);
    options.setWarningLevel(
        new DiagnosticGroup(
            CheckJSDoc.MISPLACED_MSG_ANNOTATION, VariableReferenceCheck.UNUSED_LOCAL_ASSIGNMENT),
        CheckLevel.WARNING);
    options.setWarningLevel(
        new DiagnosticGroup(ModuleLoader.LOAD_WARNING, ModuleLoader.INVALID_MODULE_PATH),
        CheckLevel.OFF);
    return options;
  }

  /** Records every diagnostic that reaches the error manager, in report order. */
  static final class RecordingErrorManager extends BasicErrorManager {
    StringBuilder sb;

    RecordingErrorManager(StringBuilder sb) {
      this.sb = sb;
    }

    @Override
    public void report(CheckLevel level, JSError error) {
      line(
          sb,
          level
              + " "
              + error.getType().key
              + " "
              + error.getSourceName()
              + ":"
              + error.getLineno()
              + ":"
              + error.getCharno()
              + " @"
              + error.getNodeSourceOffset()
              + "+"
              + error.getLength()
              + " "
              + error.getDescription());
      super.report(level, error);
    }

    @Override
    public void println(CheckLevel level, JSError error) {}

    @Override
    protected void printSummary() {}
  }

  interface Step {
    void run(Compiler compiler, Node externsRoot, Node mainRoot, StringBuilder out);
  }

  /** The inputs of one file: split at the "//// file: " lines if it starts with one. */
  static List<SourceFile> inputs(String rel, String src) {
    if (!src.startsWith(MARKER)) {
      return List.of(SourceFile.fromCode(rel, src));
    }
    List<SourceFile> inputs = new ArrayList<>();
    String name = null;
    StringBuilder cur = new StringBuilder();
    for (String l : src.split("\n", -1)) {
      if (l.startsWith(MARKER)) {
        if (name != null) {
          inputs.add(SourceFile.fromCode(name, cur.toString()));
        }
        name = l.substring(MARKER.length()).trim();
        cur.setLength(0);
      } else {
        cur.append(l).append('\n');
      }
    }
    inputs.add(SourceFile.fromCode(name, cur.toString()));
    return inputs;
  }

  static void dumpOne(
      StringBuilder sb,
      List<SourceFile> externs,
      List<NamedConfig> configs,
      String rel,
      String src) {
    line(sb, "# file " + rel);
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    RecordingErrorManager em = new RecordingErrorManager(sb);
    compiler.setErrorManager(em);
    CompilerOptions options = options();
    compiler.init(externs, inputs(rel, src), options);
    line(sb, "## parse");
    Node root;
    try {
      root = compiler.parseInputs();
    } catch (RuntimeException | StackOverflowError e) {
      exception(sb, e);
      return;
    }
    if (root == null || compiler.getErrorCount() > 0) {
      line(sb, "# parse-errors");
      return;
    }
    Node externsRoot = root.getFirstChild();
    Node mainRoot = root.getLastChild();

    run(sb, "modules", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> {
          new GatherModuleMetadata(c, /* processCommonJsModules= */ false, ResolutionMode.BROWSER)
              .process(e, m);
          new ModuleMapCreator(c, c.getModuleMetadataMap()).process(e, m);
        });
    run(sb, "CheckJSDoc", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> NodeTraversal.traverse(c, m, new CheckJSDoc(c)));
    run(sb, "CheckJSDocStyle", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> NodeTraversal.traverse(c, m, new CheckJSDocStyle(c)));
    run(sb, "VariableReferenceCheck", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new VariableReferenceCheck(c).process(e, m));
    run(sb, "CheckSideEffects", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> {
          new CheckSideEffects(c, true, true).process(e, m);
          protectors(out, m);
          Node synthesized = c.getSynthesizedExternsInput().getAstRoot(c);
          for (String l : synthesized.toStringTree().split("\n", -1)) {
            if (!l.isEmpty()) {
              line(out, "= extern " + l);
            }
          }
        });
    run(sb, "TypeCheck", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> {
          new InferConsts(c).process(e, m);
          // DefaultPassConfig.inferTypes
          TypeInferencePass inferencePass =
              new TypeInferencePass(
                  c, c.getReverseAbstractInterpreter(), (TypedScopeCreator) c.getTypedScopeCreator());
          c.setTypeCheckingHasRun(true);
          c.setTopScope(inferencePass.inferAllScopes(m.getParent()));
          // DefaultPassConfig.checkTypes
          TypeCheck check =
              new TypeCheck(
                      c,
                      c.getReverseAbstractInterpreter(),
                      c.getTypeRegistry(),
                      c.getTopScope(),
                      (TypedScopeCreator) c.getTypedScopeCreator())
                  .reportUnknownTypes(
                      options.enables(DiagnosticGroup.forType(TypeCheck.UNKNOWN_EXPR_TYPE)))
                  .reportMissingProperties(
                      !options.disables(DiagnosticGroup.forType(TypeCheck.INEXISTENT_PROPERTY)));
          check.process(e, m);
          c.getErrorManager().setTypedPercent(check.getTypedPercent());
        });
    run(sb, "CheckAccessControls", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new CheckAccessControls(c).process(e, m));
    for (NamedConfig config : configs) {
      run(sb, "CheckConformance " + config.name(), compiler, externsRoot, mainRoot,
          (c, e, m, out) -> {
            StringBuilder diagnostics = new StringBuilder();
            em.sb = diagnostics;
            try {
              new CheckConformance(
                      c, ImmutableList.of(config.config()), options.getConformanceReportingMode())
                  .process(e, m);
            } finally {
              em.sb = out;
              List<String> lines = new ArrayList<>(List.of(diagnostics.toString().split("\n")));
              lines.removeIf(String::isEmpty);
              Collections.sort(lines);
              for (String l : lines) {
                out.append(l).append('\n');
              }
            }
          });
    }
    run(sb, "StripProtection", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> {
          new CheckSideEffects.StripProtection(c).process(e, m);
          line(out, "= protectors " + countProtectors(m));
        });
  }

  static boolean isProtector(Node n) {
    return n.isCall()
        && n.getFirstChild().isName()
        && n.getFirstChild().getString().equals("JSCOMPILER_PRESERVE");
  }

  static void protectors(StringBuilder out, Node n) {
    if (isProtector(n)) {
      line(out, "= protector " + n.getSourceFileName() + ":" + n.getLineno() + ":" + n.getCharno());
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      protectors(out, c);
    }
  }

  static int countProtectors(Node n) {
    int count = isProtector(n) ? 1 : 0;
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      count += countProtectors(c);
    }
    return count;
  }

  static void run(
      StringBuilder sb, String name, Compiler compiler, Node externsRoot, Node mainRoot, Step step) {
    line(sb, "## " + name);
    try {
      step.run(compiler, externsRoot, mainRoot, sb);
    } catch (RuntimeException | StackOverflowError e) {
      exception(sb, e);
    }
  }

  /** As ChecksDump.exception: the class, then each cause's class and message. */
  static void exception(StringBuilder sb, Throwable e) {
    line(sb, "! " + e.getClass().getName());
    for (Throwable c = e.getCause(); c != null && c != c.getCause(); c = c.getCause()) {
      if (c instanceof NullPointerException) {
        line(sb, "!! cause " + c.getClass().getName());
      } else {
        line(sb, "!! cause " + c.getClass().getName() + ": " + c.getMessage());
      }
    }
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

  static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }
}

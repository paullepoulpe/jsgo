// Oracle tool for the Polymer passes (PolymerPass with PolymerClassRewriter,
// PolymerBehaviorExtractor, PolymerClassDefinition, PolymerPassStaticUtils and
// PolymerPassSuppressBehaviorsAndProtectKeys) and the J2CL passes (J2clPass,
// J2clUtilGetDefineRewriterPass, J2clConstantHoisterPass, J2clClinitPrunerPass,
// J2clPropertyInlinerPass, J2clAssertRemovalPass, J2clEqualitySameRewriterPass,
// J2clStringValueOfRewriterPass, J2clUndefinedChecksRewriterPass, J2clSourceFileChecker).
//
// The corpus has no Polymer or J2CL code, so the inputs are the JavaScript strings of the Java
// unit tests (PolymerPassTest, PolymerBehaviorExtractorTest, ..., J2cl*Test), packed in one file:
//   === <name> <byte length>\n<UTF-8 source>\n   (repeated)
// Every input runs on a fresh Compiler set up the way CompilerTestCase does: its getOptions()
// (languageIn ECMASCRIPT_NEXT, languageOut NO_TRANSPILE, GoogleCodingConvention, ...),
// compiler.init(externs, [input]), compiler.parseInputs().
//
// Modes:
//   polymer   externs = <externs file> (PolymerPassTest's EXTERNS), options as PolymerPassTest
//             (also setPreventLibraryInjection(true)). As CompilerTestCase with
//             enableCreateModuleMap/enableTypeCheck/enableRunTypeCheckAfterProcessing:
//             GatherModuleMetadata + ModuleMapCreator, then PolymerPass (checked with a
//             ChangeVerifier), the diagnostics and Node.toStringTree() of the whole root (externs
//             and sources), then TypeCheck.processForTesting and all the diagnostics again.
//   j2cl      every input of the J2cl<Pass>Test it came from runs the configuration of that test
//             (J2clSourceFileChecker.markToRunJ2clPasses first, except for "checker"):
//               J2clPassTest                    normalize, J2clPass
//               J2clUtilGetDefineRewriterPassTest  ProcessDefines(OPTIMIZE), J2clUtilGetDefineRewriterPass
//               J2clConstantHoisterPassTest     J2clConstantHoisterPass
//               J2clClinitPrunerPassTest        normalize, PureFunctionIdentifier.Driver,
//                                               J2clClinitPrunerPass(changed scopes for the pass)
//               J2clPropertyInlinerPassTest     normalize, J2clPropertyInlinerPass
//               J2clAssertRemovalPassTest       J2clAssertRemovalPass
//               J2clEqualitySameRewriterPassTest  PeepholeOptimizationsPass(J2clEqualitySameRewriterPass(false))
//               J2clStringValueOfRewriterPassTest PeepholeOptimizationsPass(J2clStringValueOfRewriterPass)
//               J2clUndefinedChecksRewriterPassTest PeepholeOptimizationsPass(J2clUndefinedChecksRewriterPass)
//               J2clSourceFileCheckerTest       J2clSourceFileChecker, then shouldRunJ2clPasses
//             with empty externs; the accessor summary is set as CompilerTestCase does; prints the
//             diagnostics and Node.toStringTree() of the sources root.
//
// Build/run:
//   javac -cp closure-compiler.jar -d build PolymerJ2clDump.java
//   java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:build \
//        com.google.javascript.jscomp.PolymerJ2clDump <polymer|j2cl> <out(.gz)|-> <inputs.txt> [<externs.js>]
//
// Output, per input, a block
//   === <name>
//   # parse errors              then <JSError line>        (and nothing else)
//   # diagnostics               then <JSError line>        (Compiler.getErrors(), then getWarnings())
//   # exception <class>: <first message line>[ cause <class>[: <first line>]]  (if the pass
//                                     threw; no cause message for ClassCastException and NPE)
//   # tree                      then the tree lines
//   # typecheck-diagnostics / # typecheck-exception ...   (polymer only)
//   # shouldRunJ2clPasses <bool>                            (checker only)
// where a JSError line is "<ERROR|WARNING> <type key> <source>:<line>:<charno> <description>".
// Every line is escaped like ParsingDump: '\' -> "\\", every char outside printable ASCII ->
// \\uXXXX. The Go side is internal/jscomp/polymer_j2cl_oracle_test.go.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.js.RuntimeJsLibManager;
import com.google.javascript.jscomp.modules.ModuleMapCreator;
import com.google.javascript.jscomp.parsing.Config.JsDocParsing;
import com.google.javascript.jscomp.type.SemanticReverseAbstractInterpreter;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.zip.GZIPOutputStream;

public final class PolymerJ2clDump {
  private PolymerJ2clDump() {}

  record Input(String name, String code) {}

  public static void main(String[] args) throws Exception {
    String mode = args[0];
    OutputStream os =
        args[1].equals("-")
            ? System.out
            : args[1].endsWith(".gz")
                ? new GZIPOutputStream(Files.newOutputStream(Paths.get(args[1])))
                : Files.newOutputStream(Paths.get(args[1]));
    List<Input> inputs = readInputs(Files.readAllBytes(Paths.get(args[2])));
    String externs =
        args.length > 3 ? Files.readString(Paths.get(args[3]), StandardCharsets.UTF_8) : "";
    try (PrintStream out = new PrintStream(os, false, StandardCharsets.UTF_8)) {
      for (Input in : inputs) {
        String block =
            switch (mode) {
              case "polymer" -> dumpPolymer(in, externs);
              case "j2cl" -> dumpJ2cl(in);
              default -> throw new IllegalArgumentException(mode);
            };
        out.print(block);
      }
    }
  }

  static List<Input> readInputs(byte[] data) {
    List<Input> inputs = new ArrayList<>();
    int i = 0;
    while (i < data.length) {
      int nl = i;
      while (data[nl] != '\n') {
        nl++;
      }
      String header = new String(data, i, nl - i, StandardCharsets.UTF_8);
      String[] parts = header.split(" ");
      if (parts.length != 3 || !parts[0].equals("===")) {
        throw new IllegalStateException("bad header: " + header);
      }
      int len = Integer.parseInt(parts[2]);
      String code = new String(data, nl + 1, len, StandardCharsets.UTF_8);
      inputs.add(new Input(parts[1], code));
      i = nl + 1 + len + 1;
    }
    return inputs;
  }

  /** CompilerTestCase.getOptions() with CompilerTestCase.setUp()'s defaults. */
  static CompilerOptions testCaseOptions() {
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    options.setEmitUseStrict(false);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    options.setModuleResolutionMode(ResolutionMode.BROWSER);
    options.setParseJsDocDocumentation(JsDocParsing.TYPES_ONLY);
    options.setAssumeStaticInheritanceIsNotUsed(true);
    options.setPreserveTypeAnnotations(true);
    options.setAssumeGettersArePure(false);
    options.setCheckSymbols(true);
    options.setWarningLevel(DiagnosticGroups.INVALID_CASTS, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.MISPLACED_MSG_ANNOTATION, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.MISSING_PROPERTIES, CheckLevel.WARNING);
    options.setCodingConvention(new GoogleCodingConvention());
    options.setPolymerVersion(1);
    options.setRuntimeLibraryMode(RuntimeJsLibManager.RuntimeLibraryMode.RECORD_AND_VALIDATE_FIELDS);
    return options;
  }

  static String dumpPolymer(Input in, String externs) {
    StringBuilder sb = new StringBuilder();
    line(sb, "=== " + in.name());
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = testCaseOptions();
    // PolymerPassTest.getOptions()
    options.setPreventLibraryInjection(true);
    compiler.init(
        List.of(SourceFile.fromCode("externs", externs)),
        List.of(SourceFile.fromCode(in.name(), in.code())),
        options);
    Node root = compiler.parseInputs();
    if (root == null || compiler.getErrorCount() > 0) {
      parseErrors(sb, compiler);
      return sb.toString();
    }
    Node externsRoot = root.getFirstChild();
    Node mainRoot = root.getLastChild();
    try {
      // CompilerTestCase.enableCreateModuleMap()
      new GatherModuleMetadata(compiler, /* processCommonJsModules= */ false, ResolutionMode.BROWSER)
          .process(externsRoot, mainRoot);
      new ModuleMapCreator(compiler, compiler.getModuleMetadataMap()).process(externsRoot, mainRoot);
      updateAccessorSummary(compiler, root);
      ChangeVerifier changeVerifier = new ChangeVerifier(compiler).snapshot(mainRoot);
      new PolymerPass(compiler).process(externsRoot, mainRoot);
      changeVerifier.checkRecordedChanges(mainRoot);
    } catch (RuntimeException e) {
      diagnostics(sb, "# diagnostics", compiler);
      exception(sb, "# exception ", e);
      return sb.toString();
    }
    diagnostics(sb, "# diagnostics", compiler);
    tree(sb, root);
    try {
      // CompilerTestCase.createTypeCheck(compiler).processForTesting(externsRoot, mainRoot)
      TypeCheck check =
          new TypeCheck(
              compiler,
              new SemanticReverseAbstractInterpreter(compiler.getTypeRegistry()),
              compiler.getTypeRegistry());
      compiler.setTypeCheckingHasRun(true);
      check.processForTesting(externsRoot, mainRoot);
    } catch (RuntimeException e) {
      exception(sb, "# typecheck-exception ", e);
      return sb.toString();
    }
    diagnostics(sb, "# typecheck-diagnostics", compiler);
    return sb.toString();
  }

  static String dumpJ2cl(Input in) {
    StringBuilder sb = new StringBuilder();
    line(sb, "=== " + in.name());
    String test = in.name().substring(0, in.name().indexOf('-'));
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    if (!test.equals("J2clSourceFileCheckerTest")) {
      // The tests' createCompiler()
      J2clSourceFileChecker.markToRunJ2clPasses(compiler);
    }
    compiler.init(
        List.of(SourceFile.fromCode("externs", "")),
        List.of(SourceFile.fromCode(in.name(), in.code())),
        testCaseOptions());
    Node root = compiler.parseInputs();
    if (root == null || compiler.getErrorCount() > 0) {
      parseErrors(sb, compiler);
      return sb.toString();
    }
    Node externsRoot = root.getFirstChild();
    Node mainRoot = root.getLastChild();
    try {
      boolean normalize =
          test.equals("J2clPassTest")
              || test.equals("J2clClinitPrunerPassTest")
              || test.equals("J2clPropertyInlinerPassTest");
      if (normalize) {
        Normalize.createNormalizeForOptimizations(compiler).process(externsRoot, mainRoot);
      }
      updateAccessorSummary(compiler, root);
      if (test.equals("J2clClinitPrunerPassTest")) {
        new PureFunctionIdentifier.Driver(compiler).process(externsRoot, mainRoot);
      }
      ChangeVerifier changeVerifier = new ChangeVerifier(compiler).snapshot(mainRoot);
      CompilerPass pass =
          switch (test) {
            case "J2clPassTest" -> new J2clPass(compiler);
            case "J2clUtilGetDefineRewriterPassTest" ->
                (externs, js) -> {
                  new ProcessDefines.Builder(compiler)
                      .setMode(ProcessDefines.Mode.OPTIMIZE)
                      .build()
                      .process(externs, js);
                  new J2clUtilGetDefineRewriterPass(compiler).process(externs, js);
                };
            case "J2clConstantHoisterPassTest" -> new J2clConstantHoisterPass(compiler);
            case "J2clClinitPrunerPassTest" ->
                new J2clClinitPrunerPass(
                    compiler,
                    compiler.getChangeTracker().getChangedScopeNodesForPass("J2clClinitPrunerPass"));
            case "J2clPropertyInlinerPassTest" -> new J2clPropertyInlinerPass(compiler);
            case "J2clAssertRemovalPassTest" -> new J2clAssertRemovalPass(compiler);
            case "J2clEqualitySameRewriterPassTest" ->
                new PeepholeOptimizationsPass(
                    compiler, "equality", new J2clEqualitySameRewriterPass(false));
            case "J2clStringValueOfRewriterPassTest" ->
                new PeepholeOptimizationsPass(
                    compiler, "valueOf", new J2clStringValueOfRewriterPass());
            case "J2clUndefinedChecksRewriterPassTest" ->
                new PeepholeOptimizationsPass(
                    compiler, "undefined", new J2clUndefinedChecksRewriterPass());
            case "J2clSourceFileCheckerTest" -> new J2clSourceFileChecker(compiler);
            default -> throw new IllegalArgumentException(test);
          };
      compiler.beforePass(test);
      pass.process(externsRoot, mainRoot);
      changeVerifier.checkRecordedChanges(mainRoot);
    } catch (RuntimeException e) {
      diagnostics(sb, "# diagnostics", compiler);
      exception(sb, "# exception ", e);
      return sb.toString();
    }
    diagnostics(sb, "# diagnostics", compiler);
    tree(sb, mainRoot);
    if (test.equals("J2clSourceFileCheckerTest")) {
      line(sb, "# shouldRunJ2clPasses " + J2clSourceFileChecker.shouldRunJ2clPasses(compiler));
    }
    return sb.toString();
  }

  /** CompilerTestCase.updateAccessorSummary */
  static void updateAccessorSummary(Compiler compiler, Node root) {
    compiler.setAccessorSummary(
        AccessorSummary.create(GatherGetterAndSetterProperties.gather(compiler, root)));
  }

  static void parseErrors(StringBuilder sb, Compiler compiler) {
    line(sb, "# parse errors");
    for (JSError e : compiler.getErrors()) {
      line(sb, error("ERROR", e));
    }
  }

  static void exception(StringBuilder sb, String prefix, RuntimeException e) {
    String text = prefix + e.getClass().getName() + ": " + firstLine(e.getMessage());
    Throwable cause = e.getCause();
    if (cause != null) {
      text += " cause " + cause.getClass().getName();
      // The JVM's own messages for these name its class loaders.
      if (!(cause instanceof ClassCastException) && !(cause instanceof NullPointerException)) {
        text += ": " + firstLine(cause.getMessage());
      }
    }
    line(sb, text);
  }

  static String firstLine(String message) {
    String msg = String.valueOf(message);
    int nl = msg.indexOf('\n');
    return nl >= 0 ? msg.substring(0, nl) : msg;
  }

  static void tree(StringBuilder sb, Node n) {
    line(sb, "# tree");
    for (String l : n.toStringTree().split("\n", -1)) {
      line(sb, l);
    }
  }

  static void diagnostics(StringBuilder sb, String title, Compiler compiler) {
    line(sb, title);
    for (JSError e : compiler.getErrors()) {
      line(sb, error("ERROR", e));
    }
    for (JSError e : compiler.getWarnings()) {
      line(sb, error("WARNING", e));
    }
  }

  static String error(String level, JSError e) {
    return level
        + " "
        + e.getType().key
        + " "
        + e.getSourceName()
        + ":"
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

  static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }
}

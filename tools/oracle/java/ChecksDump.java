// Oracle tool for the level-11 checks (Go: internal/jscomp/checks_oracle_test.go).
//
// For each input file, compiles it the way CompilerTestCase does by default (one SourceFile,
// CompilerTestCase.getOptions() defaults, GoogleCodingConvention), parses it with
// Compiler.parseInputs(), and then runs these passes in order on the same compiler, printing every
// diagnostic that reaches the ErrorManager, in report order:
//
//   CheckDebuggerStatement, CheckGlobalThis, CheckSuspiciousCode, CheckRegExp(true, true),
//   CheckSuper, CheckTypeImportCodeReferences, StrictModeCheck(WARNING), SourceInfoCheck,
//   CheckUnreachableCode, CheckMissingReturn, J2clChecksPass, VarCheck, ConstCheck,
//   VarCheck(validityCheck=true), ForbidDynamicImportUsage
//
// Callbacks that are not CompilerPasses (CheckGlobalThis, CheckSuspiciousCode,
// CheckUnreachableCode, CheckMissingReturn) run through NodeTraversal.traverse(compiler, mainRoot,
// cb), as CombinedCompilerPass does. VarCheck runs before ConstCheck, as in DefaultPassConfig
// (ConstCheck needs every name declared). After VarCheck the synthesized externs script is dumped.
// No type checking runs, so the type-based parts of the checks see no types. CheckMissingRequires
// runs only with --program (first, after GatherModuleMetadata, as CheckMissingRequiresTest does):
// GatherModuleMetadata is level 14, and Closure namespaces only resolve across files.
//
// Diagnostic types that are DISABLED by default (DEBUGGER_STATEMENT_PRESENT,
// STRICT_CHUNK_DEP_ERROR, CheckMissingRequires' MISSING_REQUIRE_TYPE & co.) are turned on as
// warnings so that they are visible, and ModuleLoader.LOAD_WARNING and INVALID_MODULE_PATH are
// turned off (VarCheckTest turns off the former) so that ES modules importing files that are not
// part of the compilation still parse.
//
// Modes: by default each file is the single code input and the externs are
// "var window; function alert() {}" (VarCheckTest's); with --as-externs each file is the single
// externs input and the code is empty, which exercises VarCheck's extern checks; with --program all
// the files are the code inputs of one compilation (in list order), printed as one section
// "# program <n> files", and CheckMissingRequires runs first.
//
// Build/run:
//   javac -cp closure-compiler.jar -d build ChecksDump.java
//   java -Xss1g -cp closure-compiler.jar:build com.google.javascript.jscomp.ChecksDump
//        [--as-externs | --program] <out.gz> <root> <rel-file>... | @<listfile>
//
// Format (gzipped; every line escaped like ParsingDump: '\' -> "\\", chars outside printable
// ASCII -> \\uXXXX):
//   # file <rel>
//   # parse-errors                    if parseInputs() reported errors (then the passes are skipped)
//   ## <Pass>                          before each pass
//   <LEVEL> <key> <source>:<line>:<charno> @<offset>+<length> <description>
//   ! <exception class>               if the pass threw, followed by a line per cause:
//   !! cause <class>[: <message>]     (no message for NullPointerException)
//   = <text>                          extra pass output (CheckRegExp: globalRegExpPropertiesUsed;
//                                     VarCheck: the synthesized externs tree, one node per line)
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.parsing.Config.JsDocParsing;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.zip.GZIPOutputStream;

public final class ChecksDump {
  static final String EXTERNS = "var window; function alert() {}";

  public static void main(String[] args) throws Exception {
    boolean asExterns = false;
    boolean program = false;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--as-externs")) {
        asExterns = true;
      } else if (a.equals("--program")) {
        program = true;
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
    if (outFile.getParent() != null) {
      Files.createDirectories(outFile.getParent());
    }
    try (OutputStream os =
        outFile.toString().endsWith(".gz")
            ? new GZIPOutputStream(Files.newOutputStream(outFile))
            : Files.newOutputStream(outFile)) {
      if (program) {
        List<SourceFile> inputs = new ArrayList<>();
        for (String rel : files) {
          inputs.add(
              SourceFile.fromCode(rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8)));
        }
        StringBuilder sb = new StringBuilder();
        line(sb, "# program " + files.size() + " files");
        dump(sb, List.of(SourceFile.fromCode("externs", EXTERNS)), inputs, true);
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      } else {
        for (String rel : files) {
          StringBuilder sb = new StringBuilder();
          dumpOne(sb, rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8), asExterns);
          os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
        }
      }
    }
  }

  /** CompilerTestCase.getOptions() with its setUp() defaults. */
  static CompilerOptions options() {
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.UNSUPPORTED);
    options.setEmitUseStrict(false);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    options.setModuleResolutionMode(ResolutionMode.BROWSER);
    options.setParseJsDocDocumentation(JsDocParsing.TYPES_ONLY);
    options.setAssumeStaticInheritanceIsNotUsed(true);
    options.setPreserveTypeAnnotations(true);
    options.setAssumeGettersArePure(false);
    options.setCheckSymbols(true);
    options.setCodingConvention(new GoogleCodingConvention());
    options.setCheckTypes(false);
    options.setWarningLevel(
        new DiagnosticGroup(
            CheckDebuggerStatement.DEBUGGER_STATEMENT_PRESENT,
            VarCheck.STRICT_CHUNK_DEP_ERROR,
            CheckMissingRequires.MISSING_REQUIRE_TYPE,
            CheckMissingRequires.INCORRECT_NAMESPACE_ALIAS_REQUIRE,
            CheckMissingRequires.INCORRECT_NAMESPACE_ALIAS_REQUIRE_TYPE,
            CheckMissingRequires.INDIRECT_NAMESPACE_REF_REQUIRE,
            CheckMissingRequires.INDIRECT_NAMESPACE_REF_REQUIRE_TYPE,
            CheckMissingRequires.MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE),
        CheckLevel.WARNING);
    // Like VarCheckTest's DiagnosticGroups.MODULE_LOAD: imports of files that are not part of the
    // compilation (or that are not paths, like 'react') are not errors.
    options.setWarningLevel(
        new DiagnosticGroup(
            com.google.javascript.jscomp.deps.ModuleLoader.LOAD_WARNING,
            com.google.javascript.jscomp.deps.ModuleLoader.INVALID_MODULE_PATH),
        CheckLevel.OFF);
    return options;
  }

  /** Records every diagnostic that reaches the error manager, in report order. */
  static final class RecordingErrorManager extends BasicErrorManager {
    final StringBuilder sb;

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

  interface PassRunner {
    void run(Compiler compiler, Node externsRoot, Node mainRoot, StringBuilder sb);
  }

  static void dumpOne(StringBuilder sb, String rel, String src, boolean asExterns) {
    line(sb, "# file " + rel);
    if (asExterns) {
      dump(
          sb,
          List.of(SourceFile.fromCode(rel, src)),
          List.of(SourceFile.fromCode("code.js", "")),
          false);
    } else {
      dump(
          sb,
          List.of(SourceFile.fromCode("externs", EXTERNS)),
          List.of(SourceFile.fromCode(rel, src)),
          false);
    }
  }

  static void dump(
      StringBuilder sb, List<SourceFile> externs, List<SourceFile> inputs, boolean program) {
    Compiler compiler = new Compiler();
    RecordingErrorManager em = new RecordingErrorManager(sb);
    compiler.setErrorManager(em);
    compiler.init(externs, inputs, options());
    line(sb, "## parse");
    Node root;
    try {
      root = compiler.parseInputs();
    } catch (RuntimeException e) {
      exception(sb, e);
      return;
    }
    if (root == null || compiler.getErrorCount() > 0) {
      line(sb, "# parse-errors");
      return;
    }
    Node externsRoot = root.getFirstChild();
    Node mainRoot = root.getLastChild();

    if (program) {
      // As CheckMissingRequiresTest does.
      run(sb, "CheckMissingRequires", compiler, externsRoot, mainRoot,
          (c, e, m, out) -> {
            new GatherModuleMetadata(c, /* processCommonJsModules= */ false, ResolutionMode.BROWSER)
                .process(e, m);
            new CheckMissingRequires(c, c.getModuleMetadataMap()).process(e, m);
          });
    }

    run(sb, "CheckDebuggerStatement", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new CheckDebuggerStatement(c).process(e, m));
    run(sb, "CheckGlobalThis", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> NodeTraversal.traverse(c, m, new CheckGlobalThis(c)));
    run(sb, "CheckSuspiciousCode", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> NodeTraversal.traverse(c, m, new CheckSuspiciousCode()));
    run(sb, "CheckRegExp", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> {
          CheckRegExp pass = new CheckRegExp(c, true, true);
          pass.process(e, m);
          line(out, "= globalRegExpPropertiesUsed=" + pass.isGlobalRegExpPropertiesUsed());
        });
    run(sb, "CheckSuper", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new CheckSuper(c).process(e, m));
    run(sb, "CheckTypeImportCodeReferences", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new CheckTypeImportCodeReferences(c).process(e, m));
    run(sb, "StrictModeCheck", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new StrictModeCheck(c, CheckLevel.WARNING).process(e, m));
    run(sb, "SourceInfoCheck", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new SourceInfoCheck(c).process(e, m));
    run(sb, "CheckUnreachableCode", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> NodeTraversal.traverse(c, m, new CheckUnreachableCode(c)));
    run(sb, "CheckMissingReturn", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> NodeTraversal.traverse(c, m, new CheckMissingReturn(c)));
    run(sb, "J2clChecksPass", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new J2clChecksPass(c).process(e, m));
    run(sb, "VarCheck", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> {
          new VarCheck(c).process(e, m);
          Node synthesized = c.getSynthesizedExternsInput().getAstRoot(c);
          for (String l : synthesized.toStringTree().split("\n", -1)) {
            if (!l.isEmpty()) {
              line(out, "= " + l);
            }
          }
        });
    run(sb, "ConstCheck", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new ConstCheck(c).process(e, m));
    run(sb, "VarCheckValidity", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new VarCheck(c, true).process(e, m));
    run(sb, "ForbidDynamicImportUsage", compiler, externsRoot, mainRoot,
        (c, e, m, out) -> new ForbidDynamicImportUsage(c).process(e, m));
  }

  static void run(
      StringBuilder sb,
      String name,
      Compiler compiler,
      Node externsRoot,
      Node mainRoot,
      PassRunner runner) {
    line(sb, "## " + name);
    try {
      runner.run(compiler, externsRoot, mainRoot, sb);
    } catch (RuntimeException | StackOverflowError e) {
      exception(sb, e);
    }
  }

  /**
   * The exception's class, then each cause's class and message. The outer message (which
   * NodeTraversal builds from the cause and the node context) and the messages of
   * NullPointerExceptions (generated by the JVM) are left out.
   */
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

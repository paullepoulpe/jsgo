// Oracle tool: type-checks each file the way DefaultPassConfig's inferTypes and checkTypes passes
// do, and dumps every diagnostic plus the inferred JSType of every node of the main AST. It is the
// differential oracle for TypeCheck, TypeInference and TypeInferencePass (and, through them, for
// TypedScopeCreator, FunctionTypeBuilder and TypeValidator).
//
// It lives in package com.google.javascript.jscomp to reach package-private API (TypeCheck's
// reportUnknownTypes/getTypedPercent, TypeInferencePass, InferConsts, Compiler.parse):
//   javac -cp closure-compiler.jar -d build TypeCheckDump.java
//   java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages \
//       -cp closure-compiler.jar:build com.google.javascript.jscomp.TypeCheckDump \
//       [--gzip] --externs=<file> <outdir> <root> <rel-file>...|@<listfile>
// (-XX:-ShowCodeDetailsInExceptionMessages turns off the JVM's "helpful" NullPointerException
// messages, which the Go side cannot reproduce: a NullPointerException then has a null message.)
//
// For each <rel-file> writes <outdir>/<rel-file>.typecheck (or .typecheck.gz). The externs are the
// contents of <file> as one SourceFile "externs.js" (tools/oracle/typecheck-externs.js is
// CompilerTypeTestCase.DEFAULT_EXTERNS, i.e. CompilerTestCase.DEFAULT_EXTERNS); the input is one
// SourceFile named <rel-file>.
//
// Configuration: CompilerTypeTestCase.defaultOptions() (GoogleCodingConvention, language
// UNSUPPORTED, MISSING_PROPERTIES / MISPLACED_TYPE_ANNOTATION / INVALID_CASTS / LINT_CHECKS /
// JSDOC_MISSING_TYPE / BOUNDED_GENERICS at WARNING), plus the levels TypeCheckTestCase.setUp adds
// (MISSING_OVERRIDE, STRICT_MISSING_PROPERTIES, STRICT_PRIMITIVE_OPERATORS at WARNING) and
// setCheckTypes(true). Steps, as TypeCheckTestCase.parseAndTypeCheckWithScope prepares them and
// DefaultPassConfig runs them: Compiler.init, Compiler.parse; if parsing reported no error,
// GatherModuleMetadata(compiler, false, BROWSER), ModuleMapCreator, InferConsts, then the body of
// DefaultPassConfig.inferTypes (TypeInferencePass.inferAllScopes on the externs-and-js root, with
// compiler.getReverseAbstractInterpreter() and compiler.getTypedScopeCreator()) and of
// DefaultPassConfig.checkTypes (TypeCheck with the top scope, reportUnknownTypes(options.enables(
// UNKNOWN_EXPR_TYPE)), reportMissingProperties(!options.disables(INEXISTENT_PROPERTY)), process,
// getTypedPercent).
//
// Format:
//   # diagnostics          then  <reported level> <key> <source>:<lineno>:<charno> "<description>"
//                                in the order the ErrorManager received them (report order)
//   # typedPercent <Double.toString(TypeCheck.getTypedPercent())>   (absent if not type-checked)
//   # tree                 then  the main js ROOT, one node per line: 4 spaces per depth +
//                                Node.toString(true, true, true) (source position, annotations
//                                and " : " + JSType.toString())
//   # exception <class>: <message>   if a step threw (the diagnostics so far come first)
// Every line is escaped as in ParsingDump.java: '\' -> "\\", every char outside printable ASCII
// -> \\uXXXX. The Go side is internal/jscomp/type_check_oracle_test.go.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.modules.ModuleMapCreator;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class TypeCheckDump {
  public static void main(String[] args) throws Exception {
    boolean gzip = false;
    String externsFile = null;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--gzip")) {
        gzip = true;
      } else if (a.startsWith("--externs=")) {
        externsFile = a.substring("--externs=".length());
      } else {
        rest.add(a);
      }
    }
    if (externsFile == null) {
      throw new IllegalArgumentException("--externs=<file> is required");
    }
    String externs = Files.readString(Paths.get(externsFile), StandardCharsets.UTF_8);
    Path outDir = Paths.get(rest.get(0));
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
    for (String rel : files) {
      StringBuilder sb = new StringBuilder();
      dumpOne(sb, externs, rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8));
      Path out = outDir.resolve(rel + (gzip ? ".typecheck.gz" : ".typecheck"));
      Files.createDirectories(out.getParent());
      try (OutputStream os =
          gzip ? new GZIPOutputStream(Files.newOutputStream(out)) : Files.newOutputStream(out)) {
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  /** CompilerTypeTestCase.defaultOptions() + TypeCheckTestCase.setUp's levels + checkTypes. */
  static CompilerOptions options() {
    CompilerOptions options = new CompilerOptions();
    options.setCodingConvention(new GoogleCodingConvention());
    options.setLanguage(CompilerOptions.LanguageMode.UNSUPPORTED);
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
    return options;
  }

  /** Records every report(level, error) in order; BasicErrorManager keeps the counts. */
  static final class RecordingErrorManager extends BasicErrorManager {
    final List<String> lines = new ArrayList<>();

    @Override
    public void report(CheckLevel level, JSError error) {
      lines.add(
          level
              + " "
              + error.getType().key
              + " "
              + error.getSourceName()
              + ":"
              + error.getLineno()
              + ":"
              + error.getCharno()
              + " \""
              + error.getDescription()
              + "\"");
      super.report(level, error);
    }

    @Override
    public void println(CheckLevel level, JSError error) {}

    @Override
    protected void printSummary() {}
  }

  static void dumpOne(StringBuilder sb, String externs, String name, String src) {
    RecordingErrorManager errors = new RecordingErrorManager();
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    compiler.setErrorManager(errors);
    CompilerOptions options = options();
    Double typedPercent = null;
    Throwable thrown = null;
    try {
      compiler.init(
          List.of(SourceFile.fromCode("externs.js", externs)),
          List.of(SourceFile.fromCode(name, src)),
          options);
      compiler.parse();
      if (errors.getErrorCount() == 0) {
        Node externsRoot = compiler.getExternsRoot();
        Node jsRoot = compiler.getJsRoot();
        new GatherModuleMetadata(compiler, false, ResolutionMode.BROWSER)
            .process(externsRoot, jsRoot);
        new ModuleMapCreator(compiler, compiler.getModuleMetadataMap())
            .process(externsRoot, jsRoot);
        new InferConsts(compiler).process(externsRoot, jsRoot);

        // DefaultPassConfig.inferTypes
        Node globalRoot = jsRoot.getParent();
        TypeInferencePass inferencePass =
            new TypeInferencePass(
                compiler,
                compiler.getReverseAbstractInterpreter(),
                (TypedScopeCreator) compiler.getTypedScopeCreator());
        compiler.setTypeCheckingHasRun(true);
        compiler.setTopScope(inferencePass.inferAllScopes(globalRoot));

        // DefaultPassConfig.checkTypes
        TypeCheck check =
            new TypeCheck(
                    compiler,
                    compiler.getReverseAbstractInterpreter(),
                    compiler.getTypeRegistry(),
                    compiler.getTopScope(),
                    (TypedScopeCreator) compiler.getTypedScopeCreator())
                .reportUnknownTypes(
                    options.enables(DiagnosticGroup.forType(TypeCheck.UNKNOWN_EXPR_TYPE)))
                .reportMissingProperties(
                    !options.disables(DiagnosticGroup.forType(TypeCheck.INEXISTENT_PROPERTY)));
        check.process(externsRoot, jsRoot);
        typedPercent = check.getTypedPercent();
        compiler.getErrorManager().setTypedPercent(typedPercent);
      }
    } catch (RuntimeException | StackOverflowError e) {
      thrown = e;
    }
    line(sb, "# diagnostics");
    for (String d : errors.lines) {
      line(sb, d);
    }
    if (thrown != null) {
      line(sb, "# exception " + thrown.getClass().getName() + ": " + thrown.getMessage());
      return;
    }
    if (typedPercent != null) {
      line(sb, "# typedPercent " + typedPercent);
    }
    line(sb, "# tree");
    Node jsRoot = compiler.getJsRoot();
    if (jsRoot == null) {
      line(sb, "null");
    } else {
      dumpTree(sb, jsRoot, 0);
    }
  }

  private static void dumpTree(StringBuilder sb, Node n, int level) {
    line(sb, "    ".repeat(level) + n.toString(true, true, true));
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      dumpTree(sb, c, level + 1);
    }
  }

  /** Backslash and every char outside printable ASCII become escapes (as ParsingDump). */
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

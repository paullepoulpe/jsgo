// Oracle tool: runs each of the 23 single-file lint checks of com.google.javascript.jscomp.lint
// (CheckArrayWithGoogObject ... CheckVar) on every input file and dumps every diagnostic each
// check reports, in report order.
//
// Each file is compiled on its own, the way CompilerTestCase runs a pass: a fresh Compiler is
// initialized with one empty externs file ("externs") and the input, with the parsing options of
// CompilerTestCase.getOptions() (language in UNSUPPORTED, out NO_TRANSPILE, BROWSER module
// resolution, JSDoc TYPES_ONLY) plus module-load diagnostics off (see options()), and
// Compiler.parseInputs() builds the AST. Then each check, as a fresh instance, runs over that AST
// in a fixed order (the order of the list below): CompilerPass checks through
// process(externsRoot, mainRoot), the two NodeTraversal.Callback-only checks (CheckProvidesSorted,
// CheckRequiresSorted, in COLLECT_AND_REPORT mode, as DefaultPassConfig creates them) through
// NodeTraversal.traverse(compiler, mainRoot, callback). (DefaultPassConfig runs them combined in
// one traversal; run one by one, each check's diagnostics stay together.)
// CheckExtraRequires gets requiresToRemove = null. No type checking runs, so the checks that read
// JSTypes (CheckArrayWithGoogObject, CheckNestedNames) see none.
//
// Diagnostics are captured raw, by overriding Compiler.report(JSError): every JSError a check
// reports is printed, whatever its level (the lint diagnostics are mostly DiagnosticType.disabled)
// and before any warnings guard.
//
// It lives in package com.google.javascript.jscomp to reach the package-private
// Compiler.parseInputs().
//
// Build/run (see tools/oracle/README.md):
//   javac -cp closure-compiler.jar -d build LintChecksDump.java
//   java -Xss1g -cp closure-compiler.jar:build com.google.javascript.jscomp.LintChecksDump
//        [--convention=google|closure|default] [--gzip] <outfile> <root> <rel-file>... | @<listfile>
//
// --convention picks the CodingConvention: google (the default) is CompilerTestCase's
// GoogleCodingConvention; closure is the Compiler's default (no convention set: a
// ClosureCodingConvention); default is CodingConventions.getDefault().
//
// Format (one block per file, in input order):
//   # file <rel-file>
//   # parse-failed <KEY>...                 if parseInputs() returned null or reported an error
//                                           (the keys of Compiler.getErrors(), in its order);
//                                           the file's checks are then skipped
//   <Check> <KEY> <defaultLevel> <lineno>:<charno> +<length> @<nodeSourceOffset> "<message>"
//                                           one line per reported JSError, in report order
//   <Check> # exception <class> [<- <cause class>]...   if the check threw (the diagnostics it
//                                           reported before that precede this line)
// The diagnostics of CheckNoMutatedEs6Exports are sorted (as strings): its Java report order
// depends on a HashMap of Node HashSets, i.e. on identity hash codes, and is not deterministic.
// Every line is escaped: '\' -> "\\", and every char outside printable ASCII -> \\uXXXX
// (newlines included), so the output is ASCII.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.deps.ModuleLoader;
import com.google.javascript.jscomp.lint.CheckArrayWithGoogObject;
import com.google.javascript.jscomp.lint.CheckConstPrivateProperties;
import com.google.javascript.jscomp.lint.CheckConstantCaseNames;
import com.google.javascript.jscomp.lint.CheckDuplicateCase;
import com.google.javascript.jscomp.lint.CheckEmptyStatements;
import com.google.javascript.jscomp.lint.CheckEnums;
import com.google.javascript.jscomp.lint.CheckEs6ModuleFileStructure;
import com.google.javascript.jscomp.lint.CheckEs6Modules;
import com.google.javascript.jscomp.lint.CheckExtraRequires;
import com.google.javascript.jscomp.lint.CheckGoogModuleTypeScriptName;
import com.google.javascript.jscomp.lint.CheckInterfaces;
import com.google.javascript.jscomp.lint.CheckMissingSemicolon;
import com.google.javascript.jscomp.lint.CheckNestedNames;
import com.google.javascript.jscomp.lint.CheckNoMutatedEs6Exports;
import com.google.javascript.jscomp.lint.CheckNullabilityModifiers;
import com.google.javascript.jscomp.lint.CheckPrimitiveAsObject;
import com.google.javascript.jscomp.lint.CheckPrototypeProperties;
import com.google.javascript.jscomp.lint.CheckProvidesSorted;
import com.google.javascript.jscomp.lint.CheckRequiresSorted;
import com.google.javascript.jscomp.lint.CheckUnusedLabels;
import com.google.javascript.jscomp.lint.CheckUnusedPrivateProperties;
import com.google.javascript.jscomp.lint.CheckUselessBlocks;
import com.google.javascript.jscomp.lint.CheckVar;
import com.google.javascript.jscomp.parsing.Config.JsDocParsing;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.function.BiConsumer;
import java.util.zip.GZIPOutputStream;

public final class LintChecksDump {
  private static String convention = "google";

  /** One lint check: its name and how to run it over (externsRoot, mainRoot). */
  private record Check(String name, BiConsumer<Compiler, Node[]> run) {}

  private static final List<Check> CHECKS =
      List.of(
          pass("CheckArrayWithGoogObject", CheckArrayWithGoogObject::new),
          pass("CheckConstPrivateProperties", CheckConstPrivateProperties::new),
          pass("CheckConstantCaseNames", CheckConstantCaseNames::new),
          pass("CheckDuplicateCase", CheckDuplicateCase::new),
          pass("CheckEmptyStatements", CheckEmptyStatements::new),
          pass("CheckEnums", CheckEnums::new),
          pass("CheckEs6ModuleFileStructure", CheckEs6ModuleFileStructure::new),
          pass("CheckEs6Modules", CheckEs6Modules::new),
          pass("CheckExtraRequires", c -> new CheckExtraRequires(c, null)),
          pass("CheckGoogModuleTypeScriptName", CheckGoogModuleTypeScriptName::new),
          pass("CheckInterfaces", CheckInterfaces::new),
          pass("CheckMissingSemicolon", CheckMissingSemicolon::new),
          pass("CheckNestedNames", CheckNestedNames::new),
          pass("CheckNoMutatedEs6Exports", CheckNoMutatedEs6Exports::new),
          pass("CheckNullabilityModifiers", CheckNullabilityModifiers::new),
          pass("CheckPrimitiveAsObject", CheckPrimitiveAsObject::new),
          pass("CheckPrototypeProperties", CheckPrototypeProperties::new),
          new Check(
              "CheckProvidesSorted",
              (c, roots) ->
                  NodeTraversal.traverse(
                      c,
                      roots[1],
                      new CheckProvidesSorted(CheckProvidesSorted.Mode.COLLECT_AND_REPORT))),
          new Check(
              "CheckRequiresSorted",
              (c, roots) ->
                  NodeTraversal.traverse(
                      c,
                      roots[1],
                      new CheckRequiresSorted(CheckRequiresSorted.Mode.COLLECT_AND_REPORT))),
          pass("CheckUnusedLabels", CheckUnusedLabels::new),
          pass("CheckUnusedPrivateProperties", CheckUnusedPrivateProperties::new),
          pass("CheckUselessBlocks", CheckUselessBlocks::new),
          pass("CheckVar", CheckVar::new));

  private interface PassFactory {
    CompilerPass create(AbstractCompiler compiler);
  }

  private static Check pass(String name, PassFactory factory) {
    return new Check(name, (c, roots) -> factory.create(c).process(roots[0], roots[1]));
  }

  public static void main(String[] args) throws Exception {
    boolean gzip = false;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--gzip")) {
        gzip = true;
      } else if (a.startsWith("--convention=")) {
        convention = a.substring("--convention=".length());
      } else {
        rest.add(a);
      }
    }
    Path out = Paths.get(rest.get(0));
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
    if (out.getParent() != null) {
      Files.createDirectories(out.getParent());
    }
    try (OutputStream os =
        gzip ? new GZIPOutputStream(Files.newOutputStream(out)) : Files.newOutputStream(out)) {
      for (String rel : files) {
        StringBuilder sb = new StringBuilder();
        dumpOne(sb, rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8));
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  /**
   * The parsing-relevant part of CompilerTestCase.getOptions() with its setUp() defaults, plus
   * module-load diagnostics off. (The rest of getOptions() only configures type checking and
   * other passes, and would make the Go mirror depend on classes above the lint checks, such as
   * DiagnosticGroups.)
   */
  private static CompilerOptions options() {
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.UNSUPPORTED);
    options.setEmitUseStrict(false);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    options.setModuleResolutionMode(ModuleLoader.ResolutionMode.BROWSER);
    options.setParseJsDocDocumentation(JsDocParsing.TYPES_ONLY);
    options.setPreserveTypeAnnotations(true);
    // As CheckExtraRequiresTest does (DiagnosticGroups.MODULE_LOAD off), so that an ES module
    // importing a file that is not an input, or a bare specifier ("three", invalid under BROWSER
    // resolution), still compiles.
    options.setWarningLevel(DiagnosticGroup.forType(ModuleLoader.LOAD_WARNING), CheckLevel.OFF);
    options.setWarningLevel(
        DiagnosticGroup.forType(ModuleLoader.INVALID_MODULE_PATH), CheckLevel.OFF);
    switch (convention) {
      case "google" -> options.setCodingConvention(new GoogleCodingConvention());
      case "closure" -> {}
      case "default" -> options.setCodingConvention(CodingConventions.getDefault());
      default -> throw new IllegalArgumentException("unknown convention " + convention);
    }
    return options;
  }

  private static void dumpOne(StringBuilder sb, String name, String src) {
    line(sb, "# file " + name);
    List<String> diags = new ArrayList<>();
    String[] phase = {null};
    Compiler compiler =
        new Compiler(new PrintStream(OutputStream.nullOutputStream())) {
          @Override
          public void report(JSError error) {
            if (phase[0] != null) {
              diags.add(phase[0] + " " + format(error));
            }
            super.report(error);
          }
        };
    compiler.init(
        List.of(SourceFile.fromCode("externs", "")),
        List.of(SourceFile.fromCode(name, src)),
        options());
    Node root;
    try {
      root = compiler.parseInputs();
    } catch (RuntimeException e) {
      root = null;
    }
    if (root == null || compiler.hasErrors()) {
      StringBuilder failed = new StringBuilder("# parse-failed");
      for (JSError e : compiler.getErrors()) {
        failed.append(' ').append(e.getType().key);
      }
      line(sb, failed.toString());
      return;
    }
    Node[] roots = {root.getFirstChild(), root.getSecondChild()};
    for (Check check : CHECKS) {
      phase[0] = check.name();
      int start = diags.size();
      try {
        check.run().accept(compiler, roots);
      } catch (RuntimeException e) {
        StringBuilder ex = new StringBuilder(check.name() + " # exception " + e.getClass().getName());
        for (Throwable c = e.getCause(); c != null; c = c.getCause()) {
          ex.append(" <- ").append(c.getClass().getName());
        }
        diags.add(ex.toString());
      }
      if (check.name().equals("CheckNoMutatedEs6Exports")) {
        Collections.sort(diags.subList(start, diags.size()));
      }
    }
    phase[0] = null;
    for (String d : diags) {
      line(sb, d);
    }
  }

  private static String format(JSError e) {
    return e.getType().key
        + " "
        + e.getDefaultLevel()
        + " "
        + e.getLineno()
        + ":"
        + e.getCharno()
        + " +"
        + e.getLength()
        + " @"
        + e.getNodeSourceOffset()
        + " \""
        + e.getDescription()
        + "\"";
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

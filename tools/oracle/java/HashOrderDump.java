// Oracle tool for Java hash-iteration order that reaches diagnostics (Go:
// internal/jscomp/hash_order_oracle_test.go).
//
// Where a pass iterates a java.util.HashMap (or a Guava hash multimap) whose keys hash by value,
// the order of its reports follows the JDK's hash order. Most diagnostic dumps sort such reports
// away; this one keeps them in report order, for hand-written cases built to exercise that
// order (keys that share a bucket, a table that grows, a table that is cleared and reused).
//
// Each case names one pass and lists its input files. The files are compiled together, the way
// LintChecksDump compiles one (a fresh Compiler, one empty externs file "externs",
// CompilerTestCase's parsing options, module-load diagnostics off, GoogleCodingConvention), and
// Compiler.parseInputs() builds the AST. Then the pass runs once over (externsRoot, mainRoot),
// and every JSError it reports through Compiler.report(JSError) is printed, in report order,
// whatever its level.
//
// Passes: CheckNoMutatedEs6Exports (its mutatedNames is MultimapBuilder.hashKeys()
// .hashSetValues(); the cases mutate each exported name once, since the Node values of one key
// iterate in identity-hash order).
//
// It lives in package com.google.javascript.jscomp to reach the package-private
// Compiler.parseInputs().
//
// Build/run (see tools/oracle/README.md):
//   javac -cp closure-compiler.jar -d build HashOrderDump.java
//   java -cp closure-compiler.jar:build com.google.javascript.jscomp.HashOrderDump
//        <cases-file> <outfile>
//
// Cases file (internal/jscomp/testdata/hash_order/cases.txt):
//   # case <name> <Pass>
//   ## file <file-name>
//   <source lines, up to the next "## file" or "# case" line>
//
// Output (one block per case, in input order):
//   # case <name> <Pass>
//   <KEY> <sourceName>:<lineno>:<charno> "<description>"   one line per reported JSError
//   # parse-failed <KEY>...                                  instead, if parsing failed
//   # exception <class>                                      if the pass threw
// Every line is escaped: '\' -> "\\", and every char outside printable ASCII -> \\uXXXX.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.deps.ModuleLoader;
import com.google.javascript.jscomp.lint.CheckNoMutatedEs6Exports;
import com.google.javascript.jscomp.parsing.Config.JsDocParsing;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

public final class HashOrderDump {
  private record Case(String name, String pass, List<SourceFile> files) {}

  public static void main(String[] args) throws Exception {
    List<Case> cases = readCases(Files.readAllLines(Paths.get(args[0]), StandardCharsets.UTF_8));
    StringBuilder sb = new StringBuilder();
    for (Case c : cases) {
      dumpOne(sb, c);
    }
    Files.writeString(Paths.get(args[1]), sb.toString(), StandardCharsets.UTF_8);
  }

  private static List<Case> readCases(List<String> lines) {
    List<Case> cases = new ArrayList<>();
    String caseName = null;
    String pass = null;
    List<SourceFile> files = null;
    String fileName = null;
    StringBuilder src = null;
    for (String l : lines) {
      if (l.startsWith("# case ") || l.startsWith("## file ")) {
        if (fileName != null) {
          files.add(SourceFile.fromCode(fileName, src.toString()));
          fileName = null;
        }
        if (l.startsWith("# case ")) {
          if (caseName != null) {
            cases.add(new Case(caseName, pass, files));
          }
          String[] parts = l.substring("# case ".length()).split(" ");
          caseName = parts[0];
          pass = parts[1];
          files = new ArrayList<>();
        } else {
          fileName = l.substring("## file ".length());
          src = new StringBuilder();
        }
      } else if (fileName != null) {
        src.append(l).append('\n');
      }
    }
    if (fileName != null) {
      files.add(SourceFile.fromCode(fileName, src.toString()));
    }
    if (caseName != null) {
      cases.add(new Case(caseName, pass, files));
    }
    return cases;
  }

  /** LintChecksDump.options() with the google convention. */
  private static CompilerOptions options() {
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.UNSUPPORTED);
    options.setEmitUseStrict(false);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    options.setModuleResolutionMode(ModuleLoader.ResolutionMode.BROWSER);
    options.setParseJsDocDocumentation(JsDocParsing.TYPES_ONLY);
    options.setPreserveTypeAnnotations(true);
    options.setWarningLevel(DiagnosticGroup.forType(ModuleLoader.LOAD_WARNING), CheckLevel.OFF);
    options.setWarningLevel(
        DiagnosticGroup.forType(ModuleLoader.INVALID_MODULE_PATH), CheckLevel.OFF);
    options.setCodingConvention(new GoogleCodingConvention());
    return options;
  }

  private static CompilerPass createPass(String pass, AbstractCompiler compiler) {
    return switch (pass) {
      case "CheckNoMutatedEs6Exports" -> new CheckNoMutatedEs6Exports(compiler);
      default -> throw new IllegalArgumentException("unknown pass " + pass);
    };
  }

  private static void dumpOne(StringBuilder sb, Case c) {
    line(sb, "# case " + c.name() + " " + c.pass());
    List<String> diags = new ArrayList<>();
    boolean[] recording = {false};
    Compiler compiler =
        new Compiler(new PrintStream(OutputStream.nullOutputStream())) {
          @Override
          public void report(JSError error) {
            if (recording[0]) {
              diags.add(format(error));
            }
            super.report(error);
          }
        };
    compiler.init(List.of(SourceFile.fromCode("externs", "")), c.files(), options());
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
    recording[0] = true;
    try {
      createPass(c.pass(), compiler).process(root.getFirstChild(), root.getSecondChild());
    } catch (RuntimeException e) {
      diags.add("# exception " + e.getClass().getName());
    }
    recording[0] = false;
    for (String d : diags) {
      line(sb, d);
    }
  }

  private static String format(JSError e) {
    return e.getType().key
        + " "
        + e.getSourceName()
        + ":"
        + e.getLineno()
        + ":"
        + e.getCharno()
        + " \""
        + e.getDescription()
        + "\"";
  }

  /** Backslash and every char outside printable ASCII become escapes. */
  private static String escape(String s) {
    StringBuilder b = new StringBuilder(s.length());
    for (int i = 0; i < s.length(); i++) {
      char ch = s.charAt(i);
      if (ch == '\\') {
        b.append("\\\\");
      } else if (ch >= 0x20 && ch < 0x7f) {
        b.append(ch);
      } else {
        b.append(String.format("\\u%04x", (int) ch));
      }
    }
    return b.toString();
  }

  private static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }
}

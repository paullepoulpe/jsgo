// Differential oracle for the ES transpilation passes (the level-16 transpilation SCC:
// internal/jscomp/transpilation_passes.go and the passes it adds — PeepholeTranspilationsPass,
// Es6NormalizeClasses, RewriteLogicalAssignmentOperatorsPass, RewriteOptionalChainingOperator,
// RewriteNullishCoalesceOperator, RewriteAsyncIteration, RewriteObjectSpread,
// Es6RewriteDestructuring, RewriteAsyncFunctions, Es7RewriteExponentialOperator,
// Es6ConvertSuper, Es6RewriteArrowFunction, Es6RewriteClass, Es6RewriteRestAndSpread,
// LateEs6ToEs3Converter, Es6ForOfConverter, Es6RewriteBlockScoped*, Es6RewriteGenerators).
// InstrumentAsyncContext is not run: the pinned jar's runtime field table has no
// $jscomp.asyncContextStart, so its constructor always throws there.
//
// For every file, each variant starts from a fresh Compiler set up the way CompilerTestCase sets
// one up: default CompilerOptions with --language_in=ECMASCRIPT_NEXT, the variant's
// --language_out, the moduleLoad group off, Compiler.init([], [SourceFile.fromCode(<rel>, <src>)],
// options) and Compiler.parse(). Files with parse errors print their diagnostics and are skipped.
// Then, in DefaultPassConfig's transpile-only order (getTranspileOnlyPasses, without the checks,
// module and polyfill passes), each pass factory is created and run on (externs, js) as
// CompilerTestCase.transpileToEs5 does:
//   TranspilationPasses.addTranspilationRuntimeLibraries   (es6InjectRuntimeLibraries)
//   Normalize.createNormalizeForOptimizations              (normalize)
//   TranspilationPasses.addTranspilationPasses(passes, options)
// The variants:
//   es5            --language_out=ECMASCRIPT5
//   es2015         --language_out=ECMASCRIPT_2015
//   es2017         --language_out=ECMASCRIPT_2017 (object rest only: REWRITE_OBJECT_REST)
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside printable
// ASCII as \\uXXXX):
//   ## file <rel>
//   # parse ok|errors                then the parse diagnostics:
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   ## variant <name>
//   diag ...                         the diagnostics added since parsing, in ErrorManager order
//   exception <class>: <message>     if a pass threw (then the variant ends after "passes")
//   passes <n>                       the number of passes run (completed or not)
//   tree unchanged                   if root.toStringTree() is the same as after parsing, else
//   tree changed <sha256>            the first 16 hex digits of the SHA-256 of the escaped lines
//                                    of root.toStringTree() (each followed by '\n'), and with
//                                    --full the escaped lines themselves
//   changes <n> <sha256>             the change-scope roots (SCRIPT and FUNCTION nodes) whose
//                                    change time is non-zero, as "<token>@<line>:<charno> <time>"
//                                    lines: their count and the first 16 hex digits of the
//                                    SHA-256 of those lines, and with --full the lines themselves
//   libs <name>,<name>,...           RuntimeJsLibManager.getInjectedLibraries(), in order
//   features <name> <FeatureSet>     NodeUtil.getFeatureSetOfScript of each SCRIPT of the js root
//   allowable <FeatureSet>           Compiler.getAllowableFeatures()
//   stage <LifeCycleStage>
//
// Usage: java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//          com.google.javascript.jscomp.TranspileDump \
//          [--full] [--sample=<N>] [--variants=<a,b>] <out.gz> <root> <rel-file>... | --list=<f>
// --sample=N dumps every Nth file of the list only (in list order, starting with the first).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
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
import java.util.Arrays;
import java.util.List;
import java.util.zip.GZIPOutputStream;

public final class TranspileDump {

  private static boolean full = false;
  private static int sample = 1;
  private static List<String> variants = Arrays.asList("es5", "es2015", "es2017");

  public static void main(String[] args) throws Exception {
    List<String> rest = new ArrayList<>();
    List<String> files = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.startsWith("--sample=")) {
        sample = Integer.parseInt(a.substring("--sample=".length()));
      } else if (a.startsWith("--variants=")) {
        variants = Arrays.asList(a.substring("--variants=".length()).split(","));
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

  private static Compiler newCompiler(String rel, String src, String variant) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    switch (variant) {
      case "es5" -> options.setLanguageOut(LanguageMode.ECMASCRIPT5);
      case "es2015" -> options.setLanguageOut(LanguageMode.ECMASCRIPT_2015);
      case "es2017" -> options.setLanguageOut(LanguageMode.ECMASCRIPT_2017);
      case "parse" -> options.setLanguageOut(LanguageMode.NO_TRANSPILE);
      default -> throw new IllegalArgumentException("unknown variant " + variant);
    }
    options.setWarningLevel(DiagnosticGroups.MODULE_LOAD, CheckLevel.OFF);
    compiler.init(
        ImmutableList.<SourceFile>of(), ImmutableList.of(SourceFile.fromCode(rel, src)), options);
    compiler.parse();
    return compiler;
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

  private static void collectChanges(Node n, List<String> out) {
    if ((n.isScript() || n.isFunction()) && n.getChangeTime() != 0) {
      out.add(n.getToken() + "@" + n.getLineno() + ":" + n.getCharno() + " " + n.getChangeTime());
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

  private static String exception(Throwable t) {
    return t.getClass().getName() + ": " + t.getMessage();
  }

  /** The passes DefaultPassConfig.getTranspileOnlyPasses runs from the runtime libraries on. */
  private static List<PassFactory> passes(Compiler compiler) {
    CompilerOptions options = compiler.getOptions();
    PassListBuilder passes = new PassListBuilder(options);
    TranspilationPasses.addTranspilationRuntimeLibraries(passes);
    passes.maybeAdd(
        PassFactory.builder()
            .setName(PassNames.NORMALIZE)
            .setInternalFactory(Normalize::createNormalizeForOptimizations)
            .build());
    TranspilationPasses.addTranspilationPasses(passes, options);
    return passes.build();
  }

  private static void dumpFile(StringBuilder sb, String rel, String src) throws Exception {
    line(sb, "## file " + rel);
    Compiler compiler = newCompiler(rel, src, "parse");
    List<String> parseDiags = diagnostics(compiler);
    line(sb, "# parse " + (compiler.hasErrors() ? "errors" : "ok"));
    for (String d : parseDiags) {
      line(sb, d);
    }
    if (compiler.hasErrors()) {
      return;
    }
    for (String variant : variants) {
      line(sb, "## variant " + variant);
      Compiler c = newCompiler(rel, src, variant);
      List<String> baseDiags = diagnostics(c);
      String before = c.getRoot().toStringTree();
      String exception = null;
      int count = 0;
      try {
        for (PassFactory factory : passes(c)) {
          count++;
          factory.create(c).process(c.getExternsRoot(), c.getJsRoot());
        }
      } catch (Throwable t) {
        exception = "exception " + exception(t);
      }
      List<String> diags = diagnostics(c);
      for (String d : baseDiags) {
        diags.remove(d);
      }
      for (String d : diags) {
        line(sb, d);
      }
      if (exception != null) {
        line(sb, exception);
        line(sb, "passes " + count);
        continue;
      }
      line(sb, "passes " + count);
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
        line(sb, "tree changed " + sha16(lines));
        if (full) {
          for (String l : lines) {
            sb.append(l).append('\n');
          }
        }
      }
      List<String> changes = new ArrayList<>();
      collectChanges(c.getRoot(), changes);
      List<String> escapedChanges = new ArrayList<>();
      for (String l : changes) {
        escapedChanges.add(escape(l));
      }
      line(sb, "changes " + changes.size() + " " + sha16(escapedChanges));
      if (full) {
        for (String l : escapedChanges) {
          sb.append(l).append('\n');
        }
      }
      line(sb, "libs " + String.join(",", c.getRuntimeJsLibManager().getInjectedLibraries()));
      for (Node script = c.getJsRoot().getFirstChild();
          script != null;
          script = script.getNext()) {
        line(
            sb,
            "features " + script.getSourceFileName() + " " + NodeUtil.getFeatureSetOfScript(script));
      }
      line(sb, "allowable " + c.getAllowableFeatures());
      line(sb, "stage " + c.getLifeCycleStage());
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

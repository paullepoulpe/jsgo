// Differential oracle for the w5-transpile-b passes (internal/jscomp): Es6NormalizeClasses,
// Es6ConvertSuper, Es6RewriteClass (with Es6ConvertSuperConstructorCalls),
// Es6RewriteBlockScopedFunctionDeclaration, Es6RewriteBlockScopedDeclaration, the peephole
// transpilations RewriteNewDotTarget / ReportUntranspilableFeatures, RewritePolyfills (with
// PolyfillUsageFinder) and IsolatePolyfills.
//
// For every file and variant, a fresh Compiler is set up with default CompilerOptions,
// --language_in=ECMASCRIPT_NEXT, the variant's --language_out, --rewrite_polyfills (and
// --isolate_polyfills for es5-isolate), the moduleLoad group off and no externs;
// Compiler.init(externs, [SourceFile.fromCode("transpileb:first.js", ""),
// SourceFile.fromCode(<rel>, <src>)], options) and Compiler.parse() (the empty first input takes the
// injected runtime libraries: Compiler.getNodeForCodeInsertion cannot insert into a module). Files
// with parse errors print their diagnostics and are skipped. Then the passes run in the order
// DefaultPassConfig.getOptimizations() runs them around transpilation, each the way PhaseOptimizer
// runs a pass (validatePreconditions, create, process; stop when the compiler has halting errors):
//   es6InjectRuntimeLibraries  new InjectTranspilationRuntimeLibraries(compiler)
//   RewritePolyfills           new RewritePolyfills(compiler, rewritePolyfills, isolatePolyfills,
//                              injectPolyfillsNewerThan), as TranspilationPasses.rewritePolyfills
//   InjectRuntimeLibraries     new InjectRuntimeLibraries(compiler, forceLibraryInjectionList)
//   normalize                  Normalize.createNormalizeForOptimizations
//   gatherGettersAndSetters    new GatherGetterAndSetterProperties(compiler)
//   ...                        TranspilationPasses.addTranspilationPasses(passes, options)
//   IsolatePolyfills           (es5-isolate only; DefaultPassConfig runs it in finalizations)
// Variants: es5 (ECMASCRIPT5), es2015 (ECMASCRIPT_2015), es5-isolate (ECMASCRIPT5 +
// isolatePolyfills).
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside
// printable ASCII as \\uXXXX):
//   ## file <rel>
//   # parse ok|errors                then the parse diagnostics:
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   ## variant <name>
//   pass <name> same|<sha256>        after each pass: "same" if root.toStringTree() did not
//                                    change, else the first 16 hex digits of the SHA-256 of its
//                                    escaped lines (each followed by '\n')
//   diag ...                         the diagnostics the pass added, in ErrorManager order
//   exception <class>: <message>     if the pass threw (then the variant ends); the JVM's helpful
//                                    NullPointerException messages are replaced by "null"
//   halt                             if the compiler has halting errors after the pass (ends)
//   libs <n> <lib>,<lib>,...         RuntimeJsLibManager.getInjectedLibraries()
//   stage <LifeCycleStage>
//   tree <sha256>                    as for pass, of the final tree
//   script <name> <sha256>           for each SCRIPT of the JS root: the hash of its escaped
//                                    toStringTree() lines; with --full followed by those lines
//                                    (except for transpileb:first.js)
//
// Usage: java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.TranspileBDump \
//          [--full] [--sample=<N>] [--offset=<K>] <out.gz> <root> <rel-file>... | --list=<listfile>
// --sample=N dumps every Nth file of the list only (in list order, starting with the K-th, K
// default 0).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSet;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Function;
import java.util.regex.Pattern;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class TranspileBDump {

  private static boolean full = false;
  private static int sample = 1;
  private static int offset = 0;

  static final String[] VARIANTS = {"es5", "es2015", "es5-isolate"};

  /**
   * An empty first input: the runtime libraries are injected into the first input
   * (Compiler.getNodeForCodeInsertion), which must not be an ES module.
   */
  static final String FIRST = "transpileb:first.js";

  public static void main(String[] args) throws Exception {
    List<String> rest = new ArrayList<>();
    List<String> files = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.startsWith("--sample=")) {
        sample = Integer.parseInt(a.substring("--sample=".length()));
      } else if (a.startsWith("--offset=")) {
        offset = Integer.parseInt(a.substring("--offset=".length()));
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
        if (index++ % sample != offset % sample) {
          continue;
        }
        StringBuilder sb = new StringBuilder();
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        dumpFile(sb, rel, src);
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
        os.flush();
      }
    }
  }

  private static Compiler newCompiler(String rel, String src, String variant) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    if (variant == null) {
      options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    } else if (variant.equals("es2015")) {
      options.setLanguageOut(LanguageMode.ECMASCRIPT_2015);
    } else {
      options.setLanguageOut(LanguageMode.ECMASCRIPT5);
    }
    if (variant != null) {
      options.setRewritePolyfills(true);
      options.setIsolatePolyfills(variant.equals("es5-isolate"));
    }
    options.setWarningLevel(DiagnosticGroups.MODULE_LOAD, CheckLevel.OFF);
    compiler.init(
        ImmutableList.<SourceFile>of(),
        ImmutableList.of(SourceFile.fromCode(FIRST, ""), SourceFile.fromCode(rel, src)),
        options);
    compiler.parse();
    return compiler;
  }

  private static final class NamedPass {
    final String name;
    final Function<Compiler, CompilerPass> factory;
    final PassFactory passFactory;

    NamedPass(String name, Function<Compiler, CompilerPass> factory) {
      this.name = name;
      this.factory = factory;
      this.passFactory = null;
    }

    NamedPass(PassFactory passFactory) {
      this.name = passFactory.getName();
      this.factory = null;
      this.passFactory = passFactory;
    }

    CompilerPass create(Compiler compiler) {
      if (passFactory != null) {
        passFactory.validatePreconditions(compiler.getOptions());
        return passFactory.create(compiler);
      }
      return factory.apply(compiler);
    }
  }

  private static List<NamedPass> passes(Compiler compiler, String variant) {
    CompilerOptions options = compiler.getOptions();
    List<NamedPass> passes = new ArrayList<>();
    // As TranspilationPasses.injectTranspilationRuntimeLibraries and .rewritePolyfills create them.
    passes.add(new NamedPass("es6InjectRuntimeLibraries", InjectTranspilationRuntimeLibraries::new));
    if (options.getRewritePolyfills()
        || options.getIsolatePolyfills()
        || options.getInjectPolyfillsNewerThan() != null) {
      passes.add(
          new NamedPass(
              "RewritePolyfills",
              (c) ->
                  new RewritePolyfills(
                      c,
                      c.getOptions().getRewritePolyfills(),
                      c.getOptions().getIsolatePolyfills(),
                      c.getOptions().getInjectPolyfillsNewerThan())));
    }
    passes.add(
        new NamedPass(
            "InjectRuntimeLibraries",
            (c) ->
                new InjectRuntimeLibraries(
                    c, ImmutableSet.copyOf(c.getOptions().getForceLibraryInjectionList()))));
    passes.add(new NamedPass(PassNames.NORMALIZE, Normalize::createNormalizeForOptimizations));
    passes.add(
        new NamedPass(PassNames.GATHER_GETTERS_AND_SETTERS, GatherGetterAndSetterProperties::new));
    PassListBuilder transpilation = new PassListBuilder(options);
    TranspilationPasses.addTranspilationPasses(transpilation, options);
    for (PassFactory f : transpilation.build()) {
      passes.add(new NamedPass(f));
    }
    if (variant.equals("es5-isolate")) {
      passes.add(new NamedPass("IsolatePolyfills", IsolatePolyfills::new));
    }
    return passes;
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

  private static List<String> treeLines(String tree) {
    List<String> lines = new ArrayList<>();
    for (String l : tree.split("\n", -1)) {
      lines.add(escape(l));
    }
    if (lines.get(lines.size() - 1).isEmpty()) {
      lines.remove(lines.size() - 1);
    }
    return lines;
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

  /**
   * getClass().getName() + ": " + getMessage(), where the JVM's "helpful" NullPointerException
   * messages (of the exception or any cause), which describe Java bytecode, are replaced by "null"
   * (the Go port's nil dereference carries no message).
   */
  /**
   * The JVM's helpful NullPointerException messages (JEP 358), e.g. {@code Cannot invoke "..."
   * because "x" is null}; explicit messages (Preconditions.checkNotNull(x, "msg")) are kept.
   */
  private static final Pattern HELPFUL_NPE_MESSAGE =
      Pattern.compile("Cannot (invoke|read|load|store|assign|throw|enter|exit) .* because .* is null");

  private static String exception(Throwable t) {
    String message = String.valueOf(t.getMessage());
    for (Throwable c = t; c != null; c = c.getCause()) {
      if (c instanceof NullPointerException
          && c.getMessage() != null
          && HELPFUL_NPE_MESSAGE.matcher(c.getMessage()).matches()) {
        message = message.replace(c.getMessage(), "null");
      }
    }
    return t.getClass().getName() + ": " + message;
  }

  private static void dumpFile(StringBuilder sb, String rel, String src) throws Exception {
    line(sb, "## file " + rel);
    Compiler compiler = newCompiler(rel, src, null);
    List<String> parseDiags = diagnostics(compiler);
    line(sb, "# parse " + (compiler.hasErrors() ? "errors" : "ok"));
    for (String d : parseDiags) {
      line(sb, d);
    }
    if (compiler.hasErrors()) {
      return;
    }
    for (String variant : VARIANTS) {
      line(sb, "## variant " + variant);
      Compiler c = newCompiler(rel, src, variant);
      List<String> seen = diagnostics(c);
      String tree = c.getRoot().toStringTree();
      boolean ended = false;
      for (NamedPass pass : passes(c, variant)) {
        String exception = null;
        try {
          pass.create(c).process(c.getExternsRoot(), c.getJsRoot());
        } catch (Throwable t) {
          exception = "exception " + exception(t);
        }
        String after = c.getRoot().toStringTree();
        line(sb, "pass " + pass.name + " " + (after.equals(tree) ? "same" : sha16(treeLines(after))));
        tree = after;
        List<String> diags = diagnostics(c);
        List<String> added = new ArrayList<>(diags);
        for (String d : seen) {
          added.remove(d);
        }
        for (String d : added) {
          line(sb, d);
        }
        seen = diags;
        if (exception != null) {
          line(sb, exception);
          ended = true;
          break;
        }
        if (c.hasHaltingErrors()) {
          line(sb, "halt");
          ended = true;
          break;
        }
      }
      if (ended) {
        continue;
      }
      List<String> libs = c.getRuntimeJsLibManager().getInjectedLibraries();
      line(sb, "libs " + libs.size() + " " + String.join(",", libs));
      line(sb, "stage " + c.getLifeCycleStage());
      line(sb, "tree " + sha16(treeLines(tree)));
      for (Node script = c.getJsRoot().getFirstChild();
          script != null;
          script = script.getNext()) {
        List<String> lines = treeLines(script.toStringTree());
        line(sb, "script " + script.getSourceFileName() + " " + sha16(lines));
        if (full && !FIRST.equals(script.getSourceFileName())) {
          for (String l : lines) {
            sb.append(l).append('\n');
          }
        }
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

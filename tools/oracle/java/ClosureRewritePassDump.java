// Differential oracle for the w5-closure-modules passes (internal/jscomp): ClosureRewriteModule,
// ScopedAliases, CheckClosureImports, RewriteGoogJsImports, modules/ModuleMapCreator (with
// EsModuleProcessor, ClosureModuleProcessor and NonEsModuleProcessor), ProcessDefines,
// ClosureCodeRemoval, GenerateExports and ExternExportsPass.
//
// Each pass runs on a fresh Compiler set up the way CompilerTestCase sets one up: default
// CompilerOptions with --language_in=ECMASCRIPT_NEXT and --language_out=NO_TRANSPILE,
// Compiler.init(no externs, inputs, options) and Compiler.parse(). Files with parse errors print
// their diagnostics and are skipped. Every pass first gets what DefaultPassConfig runs before it:
// the synthesized externs input is created (as gatherModuleMetadataPass does), then
// GatherModuleMetadata (processCommonJsModules=false, ResolutionMode.BROWSER) and, where
// DefaultPassConfig has a module map, ModuleMapCreator. The passes are constructed as
// DefaultPassConfig constructs them, without a type-checked top scope (closureRewriteModule gets
// compiler.getTopScope(), which is null before type checking) and without a preprocessor symbol
// table (the default options).
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside
// printable ASCII as \\uXXXX):
//   ## file <rel>                    (with --program: ## program <first rel> (+N), one section)
//   # parse ok|errors                then the parse diagnostics (or # parse exception <e>):
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   ## pass <name>
//   diag ...                         the diagnostics the pass added (the parse ones are removed,
//                                    one occurrence each), in ErrorManager order
//   exception <Throwable.toString()> if the pass threw (a NullPointerException: its class only)
//   tree unchanged                   if jsRoot.toStringTree() did not change, else
//   tree changed <sha256>            the first 16 hex digits of the SHA-256 of the escaped lines
//                                    of jsRoot.toStringTree() (each followed by '\n'), and with
//                                    --full the escaped lines themselves
//   externs unchanged|changed <sha>  the same for externsRoot.toStringTree()
//   <pass-specific lines>
//
// Usage: java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.ClosureRewritePassDump \
//          [--full] [--program] <out.gz> <root> <rel-file>... | --list=<listfile>
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.ScopedAliases.InvalidModuleGetHandling;
import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.modules.Binding;
import com.google.javascript.jscomp.modules.Export;
import com.google.javascript.jscomp.modules.Module;
import com.google.javascript.jscomp.modules.ModuleMap;
import com.google.javascript.jscomp.modules.ModuleMapCreator;
import com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.Token;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class ClosureRewritePassDump {

  /** A pass run: gets the compiler and (externsRoot, jsRoot); may append extra output lines. */
  private interface Run {
    void run(Compiler compiler, Node externs, Node js, List<String> extra);
  }

  private static boolean full = false;

  public static void main(String[] args) throws Exception {
    List<String> rest = new ArrayList<>();
    List<String> files = new ArrayList<>();
    boolean program = false;
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.equals("--program")) {
        program = true;
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
      if (program) {
        List<SourceFile> inputs = new ArrayList<>();
        for (String rel : files) {
          inputs.add(
              SourceFile.fromCode(rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8)));
        }
        StringBuilder sb = new StringBuilder();
        dump(sb, "## program " + files.get(0) + " (+" + (files.size() - 1) + ")", inputs);
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
        return;
      }
      int index = 0;
      for (String rel : files) {
        StringBuilder sb = new StringBuilder();
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        dump(sb, "## file " + rel, ImmutableList.of(SourceFile.fromCode(rel, src)));
        index++;
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
        if (index % 200 == 0) {
          System.err.println("progress " + index + "/" + files.size());
        }
      }
    }
  }

  private static Compiler newCompiler(List<SourceFile> inputs) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    compiler.init(ImmutableList.<SourceFile>of(), inputs, options);
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

  private static void gatherModuleMetadata(Compiler compiler, Node externs, Node js) {
    // Force creation of the synthetic input so that we create metadata for it (as
    // DefaultPassConfig.gatherModuleMetadataPass does).
    compiler.getSynthesizedExternsInput();
    new GatherModuleMetadata(compiler, false, ResolutionMode.BROWSER).process(externs, js);
  }

  private static void createModuleMap(Compiler compiler, Node externs, Node js) {
    new ModuleMapCreator(compiler, compiler.getModuleMetadataMap()).process(externs, js);
  }

  private static void dump(StringBuilder sb, String header, List<SourceFile> inputs) {
    line(sb, header);
    Compiler parsed;
    try {
      parsed = newCompiler(inputs);
    } catch (RuntimeException e) {
      // e.g. CompilerInput's dependency scan rejects goog.loadModule('<string>').
      line(sb, "# parse exception " + exceptionText(e));
      return;
    }
    List<String> parseDiags = diagnostics(parsed);
    if (parsed.hasErrors()) {
      line(sb, "# parse errors");
      for (String d : parseDiags) {
        line(sb, d);
      }
      return;
    }
    line(sb, "# parse ok");
    for (String d : parseDiags) {
      line(sb, d);
    }

    pass(sb, inputs, parseDiags, "moduleMap", (compiler, externs, js, extra) -> {
      gatherModuleMetadata(compiler, externs, js);
      createModuleMap(compiler, externs, js);
      moduleMap(compiler.getModuleMap(), extra);
    });
    pass(sb, inputs, parseDiags, "checkClosureImports", (compiler, externs, js, extra) -> {
      gatherModuleMetadata(compiler, externs, js);
      createModuleMap(compiler, externs, js);
      new CheckClosureImports(compiler, compiler.getModuleMetadataMap()).process(externs, js);
    });
    pass(sb, inputs, parseDiags, "rewriteGoogJsImports", (compiler, externs, js, extra) -> {
      gatherModuleMetadata(compiler, externs, js);
      createModuleMap(compiler, externs, js);
      new RewriteGoogJsImports(
              compiler, RewriteGoogJsImports.Mode.LINT_AND_REWRITE, compiler.getModuleMap())
          .process(externs, js);
      moduleMap(compiler.getModuleMap(), extra);
    });
    pass(sb, inputs, parseDiags, "closureGoogScopeAliases", (compiler, externs, js, extra) -> {
      gatherModuleMetadata(compiler, externs, js);
      createModuleMap(compiler, externs, js);
      ScopedAliases.builder(compiler)
          .setPreprocessorSymbolTable(null)
          .setModuleMetadataMap(compiler.getModuleMetadataMap())
          .setInvalidModuleGetHandling(InvalidModuleGetHandling.GIVE_UNIQUE_NAME)
          .build()
          .process(externs, js);
    });
    pass(sb, inputs, parseDiags, "closureGoogScopeAliasesForIjs", (compiler, externs, js, extra) -> {
      gatherModuleMetadata(compiler, externs, js);
      createModuleMap(compiler, externs, js);
      ScopedAliases.builder(compiler).build().process(externs, js);
    });
    pass(sb, inputs, parseDiags, "closureRewriteModule", (compiler, externs, js, extra) -> {
      gatherModuleMetadata(compiler, externs, js);
      createModuleMap(compiler, externs, js);
      new ClosureRewriteModule(compiler, null, compiler.getTopScope()).process(externs, js);
    });
    pass(
        sb, inputs, parseDiags, "closureGoogScopeAliases+closureRewriteModule",
        (compiler, externs, js, extra) -> {
          gatherModuleMetadata(compiler, externs, js);
          createModuleMap(compiler, externs, js);
          ScopedAliases.builder(compiler)
              .setModuleMetadataMap(compiler.getModuleMetadataMap())
              .setInvalidModuleGetHandling(InvalidModuleGetHandling.GIVE_UNIQUE_NAME)
              .build()
              .process(externs, js);
          new ClosureRewriteModule(compiler, null, compiler.getTopScope()).process(externs, js);
        });
    for (ProcessDefines.Mode mode : ProcessDefines.Mode.values()) {
      pass(sb, inputs, parseDiags, "processDefines_" + mode.name(), (compiler, externs, js, extra) -> {
        new ProcessDefines.Builder(compiler)
            .putReplacements(defineReplacements())
            .setMode(mode)
            .setRecognizeClosureDefines(true)
            .setEnableZonesDefineName(compiler.getOptions().getEnableZonesDefineName())
            .setZoneInputPattern(compiler.getOptions().getZoneInputPattern())
            .setUnknownDefinesToIgnore(compiler.getOptions().getUnknownDefinesToIgnore())
            .build()
            .process(externs, js);
        extra.add("defineNames " + compiler.getDefineNames());
      });
    }
    pass(sb, inputs, parseDiags, "closureCodeRemoval", (compiler, externs, js, extra) -> {
      new ClosureCodeRemoval(compiler, true, true).process(externs, js);
    });
    for (boolean allowNonGlobalExports : new boolean[] {false, true}) {
      pass(
          sb, inputs, parseDiags, "generateExports nonGlobal=" + allowNonGlobalExports,
          (compiler, externs, js, extra) -> {
            GenerateExports g =
                new GenerateExports(
                    compiler, allowNonGlobalExports, "goog.exportSymbol", "goog.exportProperty");
            g.process(externs, js);
            extra.add("exported " + g.getExportedVariableNames());
          });
    }
    pass(sb, inputs, parseDiags, "externExports", (compiler, externs, js, extra) -> {
      new ExternExportsPass(compiler).process(externs, js);
      String result = compiler.getResult().externExport;
      if (result == null) {
        extra.add("externExport null");
      } else {
        for (String l : result.split("\n", -1)) {
          extra.add("| " + l);
        }
      }
    });
  }

  /** The define map: --define flags of every kind and one unknown define. */
  private static Map<String, Node> defineReplacements() {
    Map<String, Node> m = new LinkedHashMap<>();
    m.put("COMPILED", new Node(Token.TRUE));
    m.put("goog.DEBUG", new Node(Token.FALSE));
    m.put("goog.LOCALE", IR.string("de"));
    m.put("goog.TRUSTED_SITE", new Node(Token.FALSE));
    m.put("goog.dom.ASSUME_STANDARDS_MODE", new Node(Token.TRUE));
    m.put("goog.userAgent.ASSUME_ANY_VERSION", new Node(Token.TRUE));
    m.put("goog.json.USE_NATIVE_JSON", new Node(Token.TRUE));
    m.put("goog.asserts.ENABLE_ASSERTS", new Node(Token.FALSE));
    m.put("goog.FEATURESET_YEAR", IR.number(2021));
    m.put("UNKNOWN_DEFINE_X", IR.number(1));
    return m;
  }

  private static void moduleMap(ModuleMap map, List<String> extra) {
    if (map == null) {
      extra.add("moduleMap null");
      return;
    }
    for (Map.Entry<String, Module> e : map.getModulesByPath().entrySet()) {
      extra.add("module " + e.getKey() + " " + module(e.getValue()));
      moduleDetails(e.getValue(), extra);
    }
    for (Map.Entry<String, Module> e : map.getModulesByClosureNamespace().entrySet()) {
      extra.add("closureModule " + e.getKey() + " " + module(e.getValue()));
    }
  }

  private static String module(Module m) {
    return "path="
        + m.path()
        + " "
        + metadata(m.metadata())
        + " closureNamespace="
        + m.closureNamespace()
        + " namespace="
        + m.namespace().keySet()
        + " boundNames="
        + m.boundNames().keySet()
        + " localExports="
        + m.localNameToLocalExport().keySet();
  }

  private static void moduleDetails(Module m, List<String> extra) {
    for (Map.Entry<String, Binding> e : m.namespace().entrySet()) {
      extra.add("  ns " + e.getKey() + " " + binding(e.getValue()));
    }
    for (Map.Entry<String, Binding> e : m.boundNames().entrySet()) {
      extra.add("  bound " + e.getKey() + " " + binding(e.getValue()));
    }
    for (Map.Entry<String, Export> e : m.localNameToLocalExport().entrySet()) {
      extra.add("  local " + e.getKey() + " " + export(e.getValue()));
    }
  }

  private static String metadata(ModuleMetadata md) {
    if (md == null) {
      return "metadata=null";
    }
    return "type=" + md.moduleType() + " goog=" + md.googNamespaces() + " mpath=" + md.path();
  }

  private static String binding(Binding b) {
    return "{"
        + metadata(b.metadata())
        + " source="
        + nodeRef(b.sourceNode())
        + " nsObj="
        + b.isModuleNamespace()
        + " closureNs="
        + b.closureNamespace()
        + " createdBy="
        + b.createdBy()
        + " boundName="
        + (b.isModuleNamespace() ? "-" : b.boundName())
        + " mutated="
        + (b.isModuleNamespace() ? "-" : b.isMutated())
        + " export="
        + (b.originatingExport() == null ? "null" : export(b.originatingExport()))
        + "}";
  }

  private static String export(Export x) {
    return "["
        + x.exportName()
        + " req="
        + x.moduleRequest()
        + " import="
        + x.importName()
        + " local="
        + x.localName()
        + " path="
        + x.modulePath()
        + " exportNode="
        + nodeRef(x.exportNode())
        + " nameNode="
        + nodeRef(x.nameNode())
        + " closureNs="
        + x.closureNamespace()
        + " mutated="
        + x.mutated()
        + "]";
  }

  private static String nodeRef(Node n) {
    if (n == null) {
      return "null";
    }
    return n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
  }

  private static void pass(
      StringBuilder sb, List<SourceFile> inputs, List<String> parseDiags, String name, Run run) {
    line(sb, "## pass " + name);
    Compiler compiler = newCompiler(inputs);
    Node externs = compiler.getExternsRoot();
    Node js = compiler.getJsRoot();
    String before = js.toStringTree();
    String externsBefore = externs.toStringTree();
    List<String> extra = new ArrayList<>();
    String exception = null;
    try {
      run.run(compiler, externs, js, extra);
    } catch (RuntimeException e) {
      exception = "exception " + exceptionText(e);
    }
    List<String> remaining = new ArrayList<>(parseDiags);
    for (String d : diagnostics(compiler)) {
      if (!remaining.remove(d)) {
        line(sb, d);
      }
    }
    if (exception != null) {
      line(sb, exception);
    }
    tree(sb, "tree", before, js.toStringTree());
    tree(sb, "externs", externsBefore, externs.toStringTree());
    for (String e : extra) {
      line(sb, e);
    }
  }

  /**
   * Throwable.toString(), except that a NullPointerException prints its class only: the JDK's
   * helpful NullPointerException messages name JVM internals that a port cannot reproduce.
   */
  private static String exceptionText(RuntimeException e) {
    if (e instanceof NullPointerException) {
      return NullPointerException.class.getName();
    }
    return e.toString();
  }

  private static void tree(StringBuilder sb, String label, String before, String after) {
    if (after.equals(before)) {
      line(sb, label + " unchanged");
      return;
    }
    StringBuilder tree = new StringBuilder();
    for (String l : after.split("\n", -1)) {
      line(tree, l);
    }
    line(sb, label + " changed " + sha256(tree.toString()));
    if (full) {
      sb.append(tree);
    }
  }

  private static String sha256(String ascii) {
    try {
      byte[] d =
          java.security.MessageDigest.getInstance("SHA-256")
              .digest(ascii.getBytes(StandardCharsets.UTF_8));
      StringBuilder b = new StringBuilder();
      for (int i = 0; i < 8; i++) {
        b.append(String.format("%02x", d[i] & 0xff));
      }
      return b.toString();
    } catch (java.security.NoSuchAlgorithmException e) {
      throw new IllegalStateException(e);
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

  private static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }
}

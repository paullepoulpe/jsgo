// Differential oracle for the w5-es-modules passes (internal/jscomp): GatherModuleMetadata,
// ProcessCommonJSModules (with FindModuleDependencies, which Compiler.parse runs when
// processCommonJSModules is set), Es6RewriteScriptsToModules, Es6RewriteModules (with
// ModuleRenaming), Es6RewriteModulesToCommonJsModules, Es6RelativizeImportPaths,
// RewriteDynamicImports and ijs/ConvertToTypedInterface.
//
// For every file, each section runs on a fresh Compiler set up the way CompilerTestCase sets one
// up: default CompilerOptions with --language_in=ECMASCRIPT_NEXT and --language_out=NO_TRANSPILE,
// with the MODULE_LOAD diagnostic group and ModuleLoader.INVALID_MODULE_PATH at WARNING (a single
// file imports files that are not part of the compilation; as errors they would stop it at
// parsing), plus the section's own options; Compiler.init(no externs,
// [SourceFile.fromCode(<rel>, <src>)], options) and Compiler.parse(). Files with parse errors print
// their diagnostics and are skipped.
// The passes get (externsRoot, jsRoot), and run after the prerequisites DefaultPassConfig runs
// before them (the synthesized externs input that its gatherModuleMetadata factory forces,
// Es6RewriteScriptsToModules, GatherModuleMetadata, ModuleMapCreator):
//
//   gather              GatherModuleMetadata(processCommonJsModules=false, BROWSER) alone, as
//                       CompilerTestCase's enableCreateModuleMap does: the ModuleMetadataMap
//   gatherCjs           options.setProcessCommonJSModules(true); GatherModuleMetadata(true,
//                       BROWSER): every input's CompilerInput.getJsModuleType() after parsing,
//                       then the ModuleMetadataMap
//   rewriteCommonJs     options.setProcessCommonJSModules(true); synthesized externs,
//                       GatherModuleMetadata(true), ModuleMapCreator, ProcessCommonJSModules
//   es6RewriteModules   synthesized externs, Es6RewriteScriptsToModules, GatherModuleMetadata,
//                       ModuleMapCreator, Es6RewriteModules(no symbol table, no typed scope,
//                       GLOBAL_NAMESPACE) (rewriting before type checking)
//   es6ToCjs            synthesized externs, GatherModuleMetadata, ModuleMapCreator,
//                       Es6RewriteModulesToCommonJsModules
//   relativize          Es6RelativizeImportPaths
//   dynamicImports      options.setAllowDynamicImport(true); synthesized externs,
//                       GatherModuleMetadata, ModuleMapCreator, Es6RewriteModules,
//                       RewriteDynamicImports(alias null, GLOBAL_NAMESPACE)
//   ijs                 ConvertToTypedInterface
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside printable
// ASCII as \\uXXXX):
//   ## file <rel>
//   # parse ok|errors                then the parse diagnostics:
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   ## pass <name>
//   diag ...                         the diagnostics the section added (the parse ones are removed,
//                                    one occurrence each), in ErrorManager order
//   exception <Throwable.toString()> if the section threw
//   tree unchanged                   if jsRoot.toStringTree() did not change, else
//   tree changed <sha256>            the first 16 hex digits of the SHA-256 of the escaped lines
//                                    of jsRoot.toStringTree() (each followed by '\n'), and with
//                                    --full the escaped lines themselves
//   externs changed <sha256>         likewise for externsRoot, when it changed (ijs)
//   <section-specific lines>         for gather/gatherCjs, the ModuleMetadataMap:
//     input <name> <jsModuleType>      (gatherCjs only) per input of the chunk graph
//     module #<i> type=<T> root=<ref> path=<p|null> usesClosure=<b> testOnly=<b> ns=<ms>
//       strong=<ms> dynamic=<ms> maybe=<ms> weak=<ms> es6=<ms> toggles=<ms> nested=[#j, ...]
//                                      (one line) every module, numbered in
//                                      getAllModuleMetadata() order with nested modules after
//                                      their parent (depth first); <ms> is the multiset's
//                                      toString()
//     byPath <path> #<i>               getModulesByPath() in order
//     byNs <namespace> #<i>            getModulesByGoogNamespace() in order
//
// Usage: java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.EsModulesDump \
//          [--full] [--max-size=<B>] [--sections=a,b,...] <out.gz> <root> <rel-file>... |
//          --list=<listfile>
// --max-size=B prints "# skipped size <bytes>" instead of running anything on files larger than B
// bytes. --sections limits the sections (default: all, in the order above).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.ChunkOutputType;
import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.deps.ModuleLoader;
import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.ijs.ConvertToTypedInterface;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.modules.ModuleMapCreator;
import com.google.javascript.jscomp.modules.ModuleMetadataMap;
import com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.IdentityHashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Consumer;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class EsModulesDump {

  /** A section: options set up before init, and the run on (compiler, externs, js). */
  private interface Run {
    void run(Compiler compiler, Node externs, Node js, List<String> extra);
  }

  static final List<String> ALL_SECTIONS =
      Arrays.asList(
          "gather",
          "gatherCjs",
          "rewriteCommonJs",
          "es6RewriteModules",
          "es6ToCjs",
          "relativize",
          "dynamicImports",
          "ijs");

  private static boolean full = false;
  private static long maxSize = Long.MAX_VALUE;
  private static List<String> sections = ALL_SECTIONS;

  public static void main(String[] args) throws Exception {
    List<String> rest = new ArrayList<>();
    List<String> files = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.startsWith("--max-size=")) {
        maxSize = Long.parseLong(a.substring("--max-size=".length()));
      } else if (a.startsWith("--sections=")) {
        sections = Arrays.asList(a.substring("--sections=".length()).split(","));
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
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(out), true)) {
      int index = 0;
      for (String rel : files) {
        StringBuilder sb = new StringBuilder();
        long size = Files.size(root.resolve(rel));
        if (size > maxSize) {
          line(sb, "## file " + rel);
          line(sb, "# skipped size " + size);
        } else {
          String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
          dumpFile(sb, rel, src);
        }
        index++;
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
        os.flush();
        if (index % 500 == 0) {
          System.err.println("progress " + index + "/" + files.size());
        }
      }
    }
  }

  private static Compiler newCompiler(String rel, String src, Consumer<CompilerOptions> setup) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    // Single files import files that are not part of the compilation: keep the module loader's
    // errors, but as warnings, so that the file is not skipped.
    options.setWarningLevel(DiagnosticGroups.MODULE_LOAD, CheckLevel.WARNING);
    options.setWarningLevel(
        DiagnosticGroup.forType(ModuleLoader.INVALID_MODULE_PATH), CheckLevel.WARNING);
    setup.accept(options);
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

  private static final Consumer<CompilerOptions> NO_SETUP = (options) -> {};
  private static final Consumer<CompilerOptions> CJS = (o) -> o.setProcessCommonJSModules(true);

  private static void dumpFile(StringBuilder sb, String rel, String src) {
    line(sb, "## file " + rel);
    Compiler parsed = newCompiler(rel, src, NO_SETUP);
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

    for (String section : sections) {
      switch (section) {
        case "gather" ->
            pass(sb, rel, src, parseDiags, section, NO_SETUP, (compiler, externs, js, extra) -> {
              new GatherModuleMetadata(compiler, false, ResolutionMode.BROWSER)
                  .process(externs, js);
              metadata(compiler.getModuleMetadataMap(), extra);
            });
        case "gatherCjs" ->
            pass(sb, rel, src, null, section, CJS, (compiler, externs, js, extra) -> {
              for (CompilerInput input : compiler.getChunkGraph().getAllInputs()) {
                extra.add("input " + input.getName() + " " + input.getJsModuleType());
              }
              new GatherModuleMetadata(compiler, true, ResolutionMode.BROWSER)
                  .process(externs, js);
              metadata(compiler.getModuleMetadataMap(), extra);
            });
        case "rewriteCommonJs" ->
            pass(sb, rel, src, null, section, CJS, (compiler, externs, js, extra) -> {
              gatherAndCreateModuleMap(compiler, externs, js, true);
              new ProcessCommonJSModules(compiler).process(externs, js);
            });
        case "es6RewriteModules" ->
            pass(sb, rel, src, parseDiags, section, NO_SETUP, (compiler, externs, js, extra) -> {
              compiler.getSynthesizedExternsInput();
              new Es6RewriteScriptsToModules(compiler).process(externs, js);
              gatherAndCreateModuleMap(compiler, externs, js, false);
              es6RewriteModules(compiler, externs, js);
            });
        case "es6ToCjs" ->
            pass(sb, rel, src, parseDiags, section, NO_SETUP, (compiler, externs, js, extra) -> {
              gatherAndCreateModuleMap(compiler, externs, js, false);
              new Es6RewriteModulesToCommonJsModules(compiler).process(externs, js);
            });
        case "relativize" ->
            pass(sb, rel, src, parseDiags, section, NO_SETUP, (compiler, externs, js, extra) -> {
              new Es6RelativizeImportPaths(compiler).process(externs, js);
            });
        case "dynamicImports" ->
            pass(
                sb,
                rel,
                src,
                parseDiags,
                section,
                (o) -> o.setAllowDynamicImport(true),
                (compiler, externs, js, extra) -> {
                  gatherAndCreateModuleMap(compiler, externs, js, false);
                  es6RewriteModules(compiler, externs, js);
                  new RewriteDynamicImports(compiler, null, ChunkOutputType.GLOBAL_NAMESPACE)
                      .process(externs, js);
                });
        case "ijs" ->
            pass(sb, rel, src, parseDiags, section, NO_SETUP, (compiler, externs, js, extra) -> {
              new ConvertToTypedInterface(compiler).process(externs, js);
            });
        default -> throw new IllegalArgumentException("unknown section " + section);
      }
    }
  }

  /**
   * The synthesized externs input that DefaultPassConfig's gatherModuleMetadata factory forces,
   * GatherModuleMetadata and ModuleMapCreator.
   */
  private static void gatherAndCreateModuleMap(
      Compiler compiler, Node externs, Node js, boolean processCommonJsModules) {
    compiler.getSynthesizedExternsInput();
    new GatherModuleMetadata(compiler, processCommonJsModules, ResolutionMode.BROWSER)
        .process(externs, js);
    new ModuleMapCreator(compiler, compiler.getModuleMetadataMap()).process(externs, js);
  }

  private static void es6RewriteModules(Compiler compiler, Node externs, Node js) {
    new Es6RewriteModules(
            compiler,
            compiler.getModuleMetadataMap(),
            compiler.getModuleMap(),
            /* preprocessorSymbolTable= */ null,
            /* globalTypedScope= */ null,
            ChunkOutputType.GLOBAL_NAMESPACE)
        .process(externs, js);
  }

  private static void metadata(ModuleMetadataMap map, List<String> extra) {
    Map<ModuleMetadata, Integer> ids = new IdentityHashMap<>();
    List<ModuleMetadata> order = new ArrayList<>();
    for (ModuleMetadata m : map.getAllModuleMetadata()) {
      number(m, ids, order);
    }
    for (ModuleMetadata m : order) {
      StringBuilder nested = new StringBuilder("[");
      for (ModuleMetadata n : m.nestedModules()) {
        if (nested.length() > 1) {
          nested.append(", ");
        }
        nested.append('#').append(ids.get(n));
      }
      nested.append(']');
      extra.add(
          "module #"
              + ids.get(m)
              + " type="
              + m.moduleType()
              + " root="
              + nodeRef(m.rootNode())
              + " path="
              + m.path()
              + " usesClosure="
              + m.usesClosure()
              + " testOnly="
              + m.isTestOnly()
              + " ns="
              + m.googNamespaces()
              + " strong="
              + m.stronglyRequiredGoogNamespaces()
              + " dynamic="
              + m.dynamicallyRequiredGoogNamespaces()
              + " maybe="
              + m.maybeRequiredGoogNamespaces()
              + " weak="
              + m.weaklyRequiredGoogNamespaces()
              + " es6="
              + m.es6ImportSpecifiers()
              + " toggles="
              + m.readToggles()
              + " nested="
              + nested);
    }
    for (Map.Entry<String, ModuleMetadata> e : map.getModulesByPath().entrySet()) {
      extra.add("byPath " + e.getKey() + " #" + ids.get(e.getValue()));
    }
    for (Map.Entry<String, ModuleMetadata> e : map.getModulesByGoogNamespace().entrySet()) {
      extra.add("byNs " + e.getKey() + " #" + ids.get(e.getValue()));
    }
  }

  private static void number(
      ModuleMetadata m, Map<ModuleMetadata, Integer> ids, List<ModuleMetadata> order) {
    if (ids.containsKey(m)) {
      return;
    }
    ids.put(m, order.size());
    order.add(m);
    for (ModuleMetadata n : m.nestedModules()) {
      number(n, ids, order);
    }
  }

  /**
   * Runs one section on a fresh compiler. parseDiags are the diagnostics of the default parse;
   * null when the section's options change parsing, in which case the section's own parse
   * diagnostics are printed as "# parse" lines first.
   */
  private static void pass(
      StringBuilder sb,
      String rel,
      String src,
      List<String> parseDiags,
      String name,
      Consumer<CompilerOptions> setup,
      Run run) {
    line(sb, "## pass " + name);
    Compiler compiler = newCompiler(rel, src, setup);
    if (parseDiags == null) {
      parseDiags = diagnostics(compiler);
      for (String d : parseDiags) {
        line(sb, "# parse " + d);
      }
      if (compiler.hasErrors()) {
        return;
      }
    }
    Node externs = compiler.getExternsRoot();
    Node js = compiler.getJsRoot();
    String before = js.toStringTree();
    String externsBefore = externs.toStringTree();
    List<String> extra = new ArrayList<>();
    String exception = null;
    try {
      run.run(compiler, externs, js, extra);
    } catch (RuntimeException e) {
      exception = "exception " + e;
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
    tree(sb, "tree", js.toStringTree(), before);
    String externsAfter = externs.toStringTree();
    if (!externsAfter.equals(externsBefore)) {
      tree(sb, "externs", externsAfter, null);
    }
    for (String e : extra) {
      line(sb, e);
    }
  }

  private static void tree(StringBuilder sb, String what, String after, String before) {
    if (after.equals(before)) {
      line(sb, what + " unchanged");
      return;
    }
    StringBuilder tree = new StringBuilder();
    for (String l : after.split("\n", -1)) {
      line(tree, l);
    }
    line(sb, what + " changed " + sha256(tree.toString()));
    if (full) {
      sb.append(tree);
    }
  }

  private static String nodeRef(Node n) {
    if (n == null) {
      return "null";
    }
    return n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
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

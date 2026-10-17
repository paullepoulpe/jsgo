// Differential oracle for the property renaming passes of the w5-renaming-misc cluster
// (internal/jscomp/rename_properties.go, disambiguate_ambiguate_properties.go,
// disambiguate_disambiguate_properties.go, rescope_global_symbols.go and
// rescope_global_symbols_rewrite_callback.go), each set up as
// DefaultPassConfig creates it for ADVANCED_OPTIMIZATIONS.
//
// For every file, each variant starts from a fresh Compiler with TypeCheckDump.options()
// (CompilerTypeTestCase's options), the externs <file> as "externs.js" and the input as one
// SourceFile named <rel>: Compiler.init, Compiler.parse. Files with parse errors print their
// diagnostics and are skipped. Every variant then runs the same prerequisites, the steps
// DefaultPassConfig runs before these passes in ADVANCED mode:
//   type checking exactly as SerializationDump.typeCheck (TypeCheckDump's steps), RemoveCastNodes,
//   ConvertTypesToColors(includeDebugInfo false) + setLifeCycleStage(COLORS_AND_SIMPLIFIED_JSDOC),
//   Normalize.createNormalizeForOptimizations, GatherGetterAndSetterProperties.update and
//   GatherExternProperties(OPTIMIZE)
// and the pass:
//   disambiguate   new DisambiguateProperties(compiler, options.getPropertiesThatMustDisambiguate())
//   ambiguate      new AmbiguateProperties(compiler, options.getPropertyReservedNamingFirstChars(),
//                  options.getPropertyReservedNamingNonFirstChars(), compiler.getExternProperties()),
//                  with its renamingMap recorded (as makePassForTesting does, set reflectively)
//   rename         new RenameProperties(compiler, false, null, reserved first/non-first chars of the
//                  options, options.getNameGenerator()) and getPropertyMap()
//   rename-pseudo  the same with generatePseudoNames = true
//   rescope        new RescopeGlobalSymbols(compiler, "_", false, true, OptimizeLocalAccess.DISABLED)
//   rescope-extern new RescopeGlobalSymbols(compiler, "_", true, false, OptimizeLocalAccess.DISABLED)
//   rescope-local  new RescopeGlobalSymbols(compiler, "_", false, true, DEFINING_CHUNK_ONLY)
//   rescope-all    new RescopeGlobalSymbols(compiler, "_", false, true, ALL_CHUNKS)
//   rescope-wrapped new RescopeGlobalSymbols(compiler, "_", true, true,
//                  ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS)
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside printable
// ASCII as \\uXXXX):
//   ## file <rel>
//   # parse ok|errors                 then the parse diagnostics:
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   ## variant <name>
//   checks <n> <sha16>                the diagnostics reported by the prerequisites
//   exception <stage> <class>: <message>   if a step threw (the variant ends); stage is checks or pass
//   diags <n> <sha16>                 the diagnostics the pass reported
//   tree <n> <sha16>                  jsRoot after the pass: one line per node, 2 spaces per depth +
//                                     Node.toString(true, true, false)
//   externs <n> <sha16>               the same for externsRoot
//   map <n> <sha16>                   the renaming map ("<old> <new>" in its iteration order;
//                                     ambiguate and rename variants only)
//   ... with --full (or for every --full-every'th file) the lines of each section follow it.
// sha16 is the first 16 hex digits of the SHA-256 of the escaped lines (each followed by '\n').
//
// Usage: java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//          com.google.javascript.jscomp.RenamingDump --externs=<file> [--full]
//          [--sample=<N>] [--full-every=<M>] <out.gz> <root> <rel-file>... | --list=<listfile>
// (compile together with TypeCheckDump.java and SerializationDump.java). The Go side is
// internal/jscomp/renaming_oracle_test.go.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage;
import com.google.javascript.jscomp.disambiguate.AmbiguateProperties;
import com.google.javascript.jscomp.disambiguate.DisambiguateProperties;
import com.google.javascript.jscomp.serialization.ConvertTypesToColors;
import com.google.javascript.jscomp.serialization.SerializationOptions;
import com.google.javascript.rhino.Node;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Field;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class RenamingDump {

  private static boolean full = false;
  private static int sample = 1;
  private static int fullEvery = 0;

  static final String[] VARIANTS = {
    "disambiguate", "ambiguate", "rename", "rename-pseudo", "rescope", "rescope-extern",
    "rescope-local", "rescope-all", "rescope-wrapped"
  };

  public static void main(String[] args) throws Exception {
    String externsFile = null;
    List<String> rest = new ArrayList<>();
    List<String> files = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.startsWith("--sample=")) {
        sample = Integer.parseInt(a.substring("--sample=".length()));
      } else if (a.startsWith("--full-every=")) {
        fullEvery = Integer.parseInt(a.substring("--full-every=".length()));
      } else if (a.startsWith("--externs=")) {
        externsFile = a.substring("--externs=".length());
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
    if (externsFile == null) {
      throw new IllegalArgumentException("--externs=<file> is required");
    }
    String externs = Files.readString(Paths.get(externsFile), StandardCharsets.UTF_8);
    Path out = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    files.addAll(rest.subList(2, rest.size()));
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(out))) {
      int index = 0;
      int dumped = 0;
      for (String rel : files) {
        if (index++ % sample != 0) {
          continue;
        }
        boolean fullThis = full || (fullEvery > 0 && dumped % fullEvery == 0);
        dumped++;
        StringBuilder sb = new StringBuilder();
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        dumpFile(sb, externs, rel, src, fullThis);
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  private static final class Setup {
    final TypeCheckDump.RecordingErrorManager errors = new TypeCheckDump.RecordingErrorManager();
    final Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    final CompilerOptions options = TypeCheckDump.options();

    Setup(SourceFile externsFile, SourceFile inputFile) {
      compiler.setErrorManager(errors);
      compiler.init(ImmutableList.of(externsFile), ImmutableList.of(inputFile), options);
    }
  }

  static void dumpFile(StringBuilder sb, String externs, String rel, String src, boolean fullThis)
      throws Exception {
    line(sb, "## file " + rel);
    SourceFile externsFile = SourceFile.fromCode("externs.js", externs);
    SourceFile inputFile = SourceFile.fromCode(rel, src);
    {
      Setup s = new Setup(externsFile, inputFile);
      try {
        s.compiler.parse();
      } catch (RuntimeException | StackOverflowError e) {
        line(sb, "# parse exception " + e.getClass().getName() + ": " + e.getMessage());
        return;
      }
      boolean ok = s.errors.getErrorCount() == 0;
      line(sb, "# parse " + (ok ? "ok" : "errors"));
      if (!ok) {
        for (String d : s.errors.lines) {
          line(sb, "diag " + d);
        }
        return;
      }
    }
    for (String variant : VARIANTS) {
      line(sb, "## variant " + variant);
      dumpVariant(sb, variant, externsFile, inputFile, fullThis);
    }
  }

  private static void dumpVariant(
      StringBuilder sb, String variant, SourceFile externsFile, SourceFile inputFile,
      boolean fullThis)
      throws Exception {
    Setup s = new Setup(externsFile, inputFile);
    String stage = "checks";
    try {
      s.compiler.parse();
      Node externsRoot = s.compiler.getExternsRoot();
      Node jsRoot = s.compiler.getJsRoot();
      SerializationDump.typeCheck(s.compiler, s.options);
      new RemoveCastNodes(s.compiler).process(externsRoot, jsRoot);
      new ConvertTypesToColors(
              s.compiler, SerializationOptions.builder().setIncludeDebugInfo(false).build())
          .process(externsRoot, jsRoot);
      s.compiler.setLifeCycleStage(LifeCycleStage.COLORS_AND_SIMPLIFIED_JSDOC);
      Normalize.createNormalizeForOptimizations(s.compiler).process(externsRoot, jsRoot);
      GatherGetterAndSetterProperties.update(s.compiler, externsRoot, jsRoot);
      new GatherExternProperties(s.compiler, GatherExternProperties.Mode.OPTIMIZE)
          .process(externsRoot, jsRoot);
      int before = s.errors.lines.size();
      line(sb, "checks " + before + " " + sha16(prefixed(s.errors.lines)));

      stage = "pass";
      CompilerOptions options = s.options;
      List<String> map = null;
      switch (variant) {
        case "disambiguate" ->
            new DisambiguateProperties(s.compiler, options.getPropertiesThatMustDisambiguate())
                .process(externsRoot, jsRoot);
        case "ambiguate" -> {
          AmbiguateProperties ap =
              new AmbiguateProperties(
                  s.compiler,
                  options.getPropertyReservedNamingFirstChars(),
                  options.getPropertyReservedNamingNonFirstChars(),
                  s.compiler.getExternProperties());
          Map<String, String> renaming = new LinkedHashMap<>();
          Field f = AmbiguateProperties.class.getDeclaredField("renamingMap");
          f.setAccessible(true);
          f.set(ap, renaming);
          ap.process(externsRoot, jsRoot);
          map = new ArrayList<>();
          for (Map.Entry<String, String> e : renaming.entrySet()) {
            map.add(e.getKey() + " " + e.getValue());
          }
        }
        case "rename", "rename-pseudo" -> {
          RenameProperties rprop =
              new RenameProperties(
                  s.compiler,
                  variant.equals("rename-pseudo"),
                  null,
                  options.getPropertyReservedNamingFirstChars(),
                  options.getPropertyReservedNamingNonFirstChars(),
                  options.getNameGenerator());
          rprop.process(externsRoot, jsRoot);
          map = new ArrayList<>();
          for (Map.Entry<String, String> e :
              rprop.getPropertyMap().getOriginalNameToNewNameMap().entrySet()) {
            map.add(e.getKey() + " " + e.getValue());
          }
        }
        case "rescope" ->
            new RescopeGlobalSymbols(
                    s.compiler, "_", false, true, CompilerOptions.OptimizeLocalAccess.DISABLED)
                .process(externsRoot, jsRoot);
        case "rescope-extern" ->
            new RescopeGlobalSymbols(
                    s.compiler, "_", true, false, CompilerOptions.OptimizeLocalAccess.DISABLED)
                .process(externsRoot, jsRoot);
        case "rescope-local" ->
            new RescopeGlobalSymbols(
                    s.compiler,
                    "_",
                    false,
                    true,
                    CompilerOptions.OptimizeLocalAccess.DEFINING_CHUNK_ONLY)
                .process(externsRoot, jsRoot);
        case "rescope-all" ->
            new RescopeGlobalSymbols(
                    s.compiler, "_", false, true, CompilerOptions.OptimizeLocalAccess.ALL_CHUNKS)
                .process(externsRoot, jsRoot);
        case "rescope-wrapped" ->
            new RescopeGlobalSymbols(
                    s.compiler,
                    "_",
                    true,
                    true,
                    CompilerOptions.OptimizeLocalAccess.ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS)
                .process(externsRoot, jsRoot);
        default -> throw new IllegalArgumentException(variant);
      }
      List<String> diags = prefixed(s.errors.lines.subList(before, s.errors.lines.size()));
      section(sb, "diags", diags, fullThis);
      List<String> tree = new ArrayList<>();
      dumpTree(tree, jsRoot, 0);
      section(sb, "tree", tree, fullThis);
      List<String> ext = new ArrayList<>();
      dumpTree(ext, externsRoot, 0);
      section(sb, "externs", ext, fullThis);
      if (map != null) {
        section(sb, "map", map, fullThis);
      }
    } catch (RuntimeException | StackOverflowError e) {
      if (System.getenv("RENAMING_DUMP_TRACE") != null) {
        e.printStackTrace(); // debugging aid, on stderr only
      }
      line(sb, "exception " + stage + " " + e.getClass().getName() + ": " + e.getMessage());
    }
  }

  private static void dumpTree(List<String> out, Node n, int level) {
    out.add("  ".repeat(level) + n.toString(true, true, false));
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      dumpTree(out, c, level + 1);
    }
  }

  private static List<String> prefixed(List<String> diags) {
    List<String> out = new ArrayList<>();
    for (String d : diags) {
      out.add("diag " + d);
    }
    return out;
  }

  private static void section(StringBuilder sb, String name, List<String> lines, boolean fullThis)
      throws Exception {
    line(sb, name + " " + lines.size() + " " + sha16(lines));
    if (fullThis) {
      for (String l : lines) {
        line(sb, l);
      }
    }
  }

  private static String sha16(List<String> lines) throws Exception {
    MessageDigest md = MessageDigest.getInstance("SHA-256");
    for (String l : lines) {
      md.update((TypeCheckDump.escape(l) + "\n").getBytes(StandardCharsets.UTF_8));
    }
    StringBuilder hex = new StringBuilder();
    for (byte b : md.digest()) {
      hex.append(String.format("%02x", b));
    }
    return hex.substring(0, 16);
  }

  private static void line(StringBuilder sb, String s) {
    sb.append(TypeCheckDump.escape(s)).append('\n');
  }
}

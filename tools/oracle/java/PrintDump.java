// Oracle tool for the code printer: parses each file with ParserRunner.parse (the compiler's own
// parsing.Config, as ParsingDump.java obtains it) and prints the AST with CodePrinter in several
// configurations, with and without source maps, plus JSDocInfoPrinter.print for every JSDocInfo and
// (for small files) DotFormatter.toDot. The Go side is internal/jscomp/print_oracle_test.go.
//
// In package com.google.javascript.jscomp: it reads package-private API (SourceAndMappings,
// SourceMap.Mapping, SourceMap.Format.getInstance()).
//
// Build/run (see tools/oracle/README.md):
//   javac -cp closure-compiler.jar -d build PrintDump.java
//   java -Xss1g -cp closure-compiler.jar:build com.google.javascript.jscomp.PrintDump [--gzip]
//        [--every=N] [--only=<section>,...] [--no-mappings] [--dot-limit=BYTES] <outdir> <root>
//        <rel-file>... | @<listfile>
//
// For each <rel-file> (with --every=N: only every Nth file of the list, starting with the first)
// writes <outdir>/<rel-file>.print (or .print.gz), starting with the line
//   # options only=<sorted sections, or *> no-mappings=<bool> dot-limit=<n>
// --only keeps the named sections (configuration
// names, "jsdoc-types", "jsdoc-all", "dot"), in their usual order; --no-mappings leaves out the
// "# mappings" lines (the map JSON is still written). Format, every line escaped as in ParsingDump
// (backslash and every char outside printable ASCII become \\uXXXX):
//   # parse ok | # parse errors <n> | # parse exception <class>: <message>
//   then, only when the parse succeeded without errors, for each configuration <name>:
//   # print <name>
//   | <line>          one per line of the output: the output split at '\n' (so a trailing
//                     newline gives a last "| " line)
//   # mappings        (source map configurations) one line per mapping, in output order:
//                     <token> <lineno>:<charno> <original name or -> <start> <end>
//                     with positions <line>:<column>
//   # map             (source map configurations) the SourceMap V3 JSON (appendTo "out.js"), as
//                     "| " lines
//   # exception <class>: <message>   instead of the above when printing threw
//   # jsdoc <parse> <token> <lineno>:<charno>   for every node with JSDoc (pre-order, closure-unaware
//                     shadows included), followed by "| " lines of JSDocInfoPrinter.print; <parse>
//                     is "types" (default parse, printer (false)) or "all" (INCLUDE_ALL_COMMENTS
//                     parse, printer (false, true))
//   # dot             DotFormatter.toDot(root) as "| " lines, for sources up to --dot-limit bytes
//                     (default 4096)
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.jstype.JSTypeRegistry;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Consumer;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class PrintDump {
  private static int dotLimit = 4096;
  private static java.util.Set<String> only = null;
  private static boolean noMappings = false;

  /** Whether the section is selected (--only=a,b,... selects sections by name). */
  private static boolean selected(String name) {
    return only == null || only.contains(name);
  }

  public static void main(String[] args) throws Exception {
    boolean gzip = false;
    int every = 1;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--gzip")) {
        gzip = true;
      } else if (a.startsWith("--every=")) {
        every = Integer.parseInt(a.substring("--every=".length()));
      } else if (a.equals("--no-mappings")) {
        noMappings = true;
      } else if (a.startsWith("--only=")) {
        only = new java.util.HashSet<>(java.util.Arrays.asList(a.substring("--only=".length()).split(",")));
      } else if (a.startsWith("--dot-limit=")) {
        dotLimit = Integer.parseInt(a.substring("--dot-limit=".length()));
      } else {
        rest.add(a);
      }
    }
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
    Config typesConfig = compilerParserConfig(null);
    Config allConfig = compilerParserConfig(Config.JsDocParsing.INCLUDE_ALL_COMMENTS);
    for (int i = 0; i < files.size(); i += every) {
      String rel = files.get(i);
      StringBuilder sb = new StringBuilder();
      dumpOne(sb, typesConfig, allConfig, rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8));
      Path out = outDir.resolve(rel + (gzip ? ".print.gz" : ".print"));
      Files.createDirectories(out.getParent());
      try (OutputStream os =
          gzip ? new GZIPOutputStream(Files.newOutputStream(out)) : Files.newOutputStream(out)) {
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  /** The parsing.Config that CompilerInput passes to ParserRunner.parse for non-extern inputs. */
  private static Config compilerParserConfig(Config.JsDocParsing jsdocMode) throws Exception {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    if (jsdocMode != null) {
      options.setParseJsDocDocumentation(jsdocMode);
    }
    compiler.initOptions(options);
    Class<?> ctxClass = Class.forName("com.google.javascript.jscomp.AbstractCompiler$ConfigContext");
    Object defaultCtx = null;
    for (Object c : ctxClass.getEnumConstants()) {
      if (((Enum<?>) c).name().equals("DEFAULT")) {
        defaultCtx = c;
      }
    }
    Method getParserConfig = Compiler.class.getDeclaredMethod("getParserConfig", ctxClass);
    getParserConfig.setAccessible(true);
    return (Config) getParserConfig.invoke(compiler, defaultCtx);
  }

  /** Parses src; returns null (after writing the "# parse" line) unless the parse is clean. */
  private static Node parse(StringBuilder sb, Config config, String name, String src) {
    int[] errors = {0};
    ErrorReporter reporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int line, int lineOffset) {}

          @Override
          public void error(String message, String sourceName, int line, int lineOffset) {
            errors[0]++;
          }
        };
    ParserRunner.ParseResult result;
    try {
      result = ParserRunner.parse(SourceFile.fromCode(name, src), src, config, reporter);
    } catch (RuntimeException e) {
      if (sb != null) {
        line(sb, "# parse exception " + e.getClass().getName() + ": " + e.getMessage());
      }
      return null;
    }
    if (errors[0] > 0 || result.ast == null) {
      if (sb != null) {
        line(sb, "# parse errors " + errors[0]);
      }
      return null;
    }
    if (sb != null) {
      line(sb, "# parse ok");
    }
    return result.ast;
  }

  private static void dumpOne(
      StringBuilder sb, Config typesConfig, Config allConfig, String name, String src) {
    line(
        sb,
        "# options only="
            + (only == null ? "*" : String.join(",", new java.util.TreeSet<>(only)))
            + " no-mappings="
            + noMappings
            + " dot-limit="
            + dotLimit);
    Node root = parse(sb, typesConfig, name, src);
    if (root == null) {
      return;
    }

    // Plain configurations.
    print(sb, "compact", root, b -> {});
    print(sb, "pretty", root, b -> b.setPrettyPrint(true));
    print(sb, "linebreak", root, b -> b.setLineBreak(true));
    print(
        sb,
        "compact-threshold-60",
        root,
        b -> {
          CompilerOptions o = new CompilerOptions();
          o.setLineLengthThreshold(60);
          b.setCompilerOptions(o);
        });
    print(
        sb,
        "single-quotes",
        root,
        b -> {
          CompilerOptions o = new CompilerOptions();
          o.setPreferSingleQuotes(true);
          b.setCompilerOptions(o);
        });
    print(
        sb,
        "es3-untrusted",
        root,
        b -> {
          CompilerOptions o = new CompilerOptions();
          o.setLanguageOut(CompilerOptions.LanguageMode.ECMASCRIPT3);
          o.setTrustedStrings(false);
          b.setCompilerOptions(o);
        });
    print(
        sb,
        "es5-utf8",
        root,
        b -> {
          CompilerOptions o = new CompilerOptions();
          o.setLanguageOut(CompilerOptions.LanguageMode.ECMASCRIPT5);
          o.setOutputCharset(StandardCharsets.UTF_8);
          b.setCompilerOptions(o);
        });
    print(
        sb,
        "pretty-types",
        root,
        b -> {
          CompilerOptions o = new CompilerOptions();
          o.setPreserveTypeAnnotations(true);
          o.setPrettyPrint(true);
          b.setCompilerOptions(o);
        });
    print(
        sb,
        "typesummary-strict",
        root,
        b -> b.setPrettyPrint(true).setTagAsTypeSummary(true).setTagAsStrict(true));
    print(
        sb,
        "typed",
        root,
        b -> b.setOutputTypes(true).setTypeRegistry(new JSTypeRegistry(null)));

    // Source maps, after the original names are recorded as the compiler does.
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    compiler.initOptions(new CompilerOptions());
    NodeTraversal.traverse(compiler, root, SourceInformationAnnotator.create());
    printWithSourceMap(sb, "sourcemap-compact", root, false, false, 500, SourceMap.DetailLevel.ALL);
    printWithSourceMap(
        sb, "sourcemap-compact-threshold-60", root, false, false, 60, SourceMap.DetailLevel.ALL);
    printWithSourceMap(sb, "sourcemap-linebreak", root, false, true, 500, SourceMap.DetailLevel.ALL);
    printWithSourceMap(sb, "sourcemap-pretty", root, true, false, 500, SourceMap.DetailLevel.ALL);
    printWithSourceMap(
        sb, "sourcemap-symbols", root, false, false, 500, SourceMap.DetailLevel.SYMBOLS);
    print(
        sb,
        "original-names",
        root,
        b -> {
          CompilerOptions o = new CompilerOptions();
          o.setUseOriginalNamesInOutput(true);
          b.setCompilerOptions(o);
        });

    // Comments: a parse that keeps every comment.
    Node allRoot =
        selected("pretty-comments") || selected("compact-comments") || selected("jsdoc-all")
            ? parse(null, allConfig, name, src)
            : null;
    if (allRoot != null) {
      print(
          sb,
          "pretty-comments",
          allRoot,
          b -> {
            CompilerOptions o = new CompilerOptions();
            o.setPreserveTypeAnnotations(true);
            o.setPreserveNonJSDocComments(true);
            o.setPrettyPrint(true);
            b.setCompilerOptions(o);
          });
      print(
          sb,
          "compact-comments",
          allRoot,
          b -> {
            CompilerOptions o = new CompilerOptions();
            o.setPreserveNonJSDocComments(true);
            b.setCompilerOptions(o);
          });
    }

    // JSDocInfoPrinter on every JSDocInfo.
    if (selected("jsdoc-types")) {
      dumpJsDoc(sb, "types", root, new JSDocInfoPrinter(false));
    }
    if (allRoot != null && selected("jsdoc-all")) {
      dumpJsDoc(sb, "all", allRoot, new JSDocInfoPrinter(false, true));
    }

    if (src.length() <= dotLimit && selected("dot")) {
      line(sb, "# dot");
      try {
        lines(sb, DotFormatter.toDot(root));
      } catch (Exception e) {
        exception(sb, e);
      }
    }
  }

  private static void print(
      StringBuilder sb, String name, Node root, Consumer<CodePrinter.Builder> configure) {
    if (!selected(name)) {
      return;
    }
    line(sb, "# print " + name);
    try {
      CodePrinter.Builder b = new CodePrinter.Builder(root);
      configure.accept(b);
      lines(sb, b.build());
    } catch (RuntimeException | Error e) {
      exception(sb, e);
    }
  }

  private static void printWithSourceMap(
      StringBuilder sb,
      String name,
      Node root,
      boolean pretty,
      boolean lineBreak,
      int threshold,
      SourceMap.DetailLevel detail) {
    if (!selected(name)) {
      return;
    }
    line(sb, "# print " + name);
    try {
      CompilerOptions o = new CompilerOptions();
      o.setSourceMapOutputPath("out.map");
      o.setSourceMapDetailLevel(detail);
      o.setLineLengthThreshold(threshold);
      o.setPrettyPrint(pretty);
      o.setLineBreak(lineBreak);
      CodePrinter.SourceAndMappings result =
          new CodePrinter.Builder(root).setCompilerOptions(o).buildWithSourceMappings();
      lines(sb, result.source);
      if (!noMappings) {
        line(sb, "# mappings");
      }
      SourceMap sourceMap = SourceMap.Format.V3.getInstance();
      for (SourceMap.Mapping m : result.mappings) {
        String original = m.node.getOriginalName();
        if (!noMappings) {
          line(
            sb,
            m.node.getToken()
                + " "
                + m.node.getLineno()
                + ":"
                + m.node.getCharno()
                + " "
                + (original == null ? "-" : original)
                + " "
                + m.start.getLine()
                + ":"
                + m.start.getColumn()
                + " "
                + m.end.getLine()
                + ":"
                + m.end.getColumn());
        }
        sourceMap.addMapping(m);
      }
      line(sb, "# map");
      StringBuilder map = new StringBuilder();
      sourceMap.appendTo(map, "out.js");
      lines(sb, map.toString());
    } catch (Exception | Error e) {
      exception(sb, e);
    }
  }

  private static void dumpJsDoc(StringBuilder sb, String parse, Node n, JSDocInfoPrinter printer) {
    JSDocInfo info = n.getJSDocInfo();
    if (info != null) {
      line(sb, "# jsdoc " + parse + " " + n.getToken() + " " + n.getLineno() + ":" + n.getCharno());
      try {
        lines(sb, printer.print(info));
      } catch (RuntimeException | Error e) {
        exception(sb, e);
      }
    }
    Node shadow = n.getClosureUnawareShadow();
    if (shadow != null) {
      dumpJsDoc(sb, parse, shadow, printer);
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      dumpJsDoc(sb, parse, c, printer);
    }
  }

  private static void exception(StringBuilder sb, Throwable e) {
    line(sb, "# exception " + e.getClass().getName() + ": " + e.getMessage());
  }

  /** Writes s as "| " lines: s split at '\n', keeping a trailing empty piece. */
  private static void lines(StringBuilder sb, String s) {
    int start = 0;
    while (true) {
      int nl = s.indexOf('\n', start);
      if (nl < 0) {
        line(sb, "| " + s.substring(start));
        return;
      }
      line(sb, "| " + s.substring(start, nl));
      start = nl + 1;
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

  private PrintDump() {}
}

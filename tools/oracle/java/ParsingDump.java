// Oracle tool: dumps what com.google.javascript.jscomp.parsing.ParserRunner.parse returns for a
// file — the rhino.Node AST built by IRFactory (with JSDoc), the parse features, the source map
// URL, the comments — plus every error/warning reported to the rhino ErrorReporter. This is the
// exact boundary of the Go package internal/parsing; the compiler-level AstDump adds what
// CompilerInput and RhinoErrorReporter do on top.
//
// Build/run (see tools/oracle/README.md):
//   javac -cp closure-compiler.jar -d build ParsingDump.java
//   java -Xss1g -cp closure-compiler.jar:build ParsingDump [--gzip] [--lang=<LanguageMode>] [--jsdoc]
//        [--jsdoc-mode=<JsDocParsing>] [--keep-going] <outdir> <root> <rel-file>... | @<listfile>
//
// For each <rel-file> writes <outdir>/<rel-file>.parse (or .parse.gz). The parsing.Config is the
// one the compiler uses (Compiler.getParserConfig(ConfigContext.DEFAULT), called reflectively);
// --jsdoc selects JsDocParsing.INCLUDE_DESCRIPTIONS_NO_WHITESPACE instead of TYPES_ONLY,
// --jsdoc-mode=<JsDocParsing> any mode (CompilerOptions.setParseJsDocDocumentation), and
// --keep-going RunMode.KEEP_GOING (CompilerOptions.setContinueAfterErrors(true)).
//
// Format:
//   # config languageMode=<..> jsDocParsingMode=<..> runMode=<..> strictMode=<..>
//   # diagnostics                then  ERROR|WARNING <line>:<lineOffset> <quoted message>
//   # features <FeatureSet.toString()>
//   # sourceMapURL <quoted string or null>
//   # comments                   then  <Comment.Type> <range> <quoted value>
//   # tree                       then  one node per line: 4 spaces per depth +
//                                      Node.toString(true, true, false); a node with JSDoc is
//                                      followed by a line at the same indent: "  @jsdoc " +
//                                      JSDocInfo.toStringVerbose(); a closure-unaware shadow host
//                                      (Node.getClosureUnawareShadow() != null) is followed by
//                                      "  @shadow" at the same indent and the shadow's ROOT tree
//                                      one level deeper
//   # exception <class>: <message>   if ParserRunner.parse threw
// Every line is escaped: '\' -> "\\", and every char outside printable ASCII -> \\uXXXX
// (newlines included), so the output is ASCII and lone surrogates are unambiguous.
import com.google.javascript.jscomp.CompilerOptions;
import com.google.javascript.jscomp.SourceFile;
import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.jscomp.parsing.parser.trees.Comment;
import com.google.javascript.jscomp.parsing.parser.util.SourcePosition;
import com.google.javascript.jscomp.parsing.parser.util.SourceRange;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class ParsingDump {
  private static CompilerOptions.LanguageMode languageIn =
      CompilerOptions.LanguageMode.ECMASCRIPT_NEXT;
  private static Config.JsDocParsing jsdocMode = null;
  private static boolean keepGoing = false;

  public static void main(String[] args) throws Exception {
    boolean gzip = false;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--gzip")) {
        gzip = true;
      } else if (a.equals("--jsdoc")) {
        jsdocMode = Config.JsDocParsing.INCLUDE_DESCRIPTIONS_NO_WHITESPACE;
      } else if (a.startsWith("--jsdoc-mode=")) {
        jsdocMode = Config.JsDocParsing.valueOf(a.substring("--jsdoc-mode=".length()));
      } else if (a.equals("--keep-going")) {
        keepGoing = true;
      } else if (a.startsWith("--lang=")) {
        languageIn = CompilerOptions.LanguageMode.valueOf(a.substring("--lang=".length()));
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
    Config config = compilerParserConfig();
    for (String rel : files) {
      StringBuilder sb = new StringBuilder();
      dumpOne(sb, config, rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8));
      Path out = outDir.resolve(rel + (gzip ? ".parse.gz" : ".parse"));
      Files.createDirectories(out.getParent());
      try (OutputStream os =
          gzip ? new GZIPOutputStream(Files.newOutputStream(out)) : Files.newOutputStream(out)) {
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  /** The parsing.Config that CompilerInput passes to ParserRunner.parse for non-extern inputs. */
  private static Config compilerParserConfig() throws Exception {
    com.google.javascript.jscomp.Compiler compiler =
        new com.google.javascript.jscomp.Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(languageIn);
    if (jsdocMode != null) {
      options.setParseJsDocDocumentation(jsdocMode);
    }
    if (keepGoing) {
      options.setContinueAfterErrors(true);
    }
    compiler.initOptions(options);
    Class<?> ctxClass = Class.forName("com.google.javascript.jscomp.AbstractCompiler$ConfigContext");
    Object defaultCtx = null;
    for (Object c : ctxClass.getEnumConstants()) {
      if (((Enum<?>) c).name().equals("DEFAULT")) {
        defaultCtx = c;
      }
    }
    Method getParserConfig =
        com.google.javascript.jscomp.Compiler.class.getDeclaredMethod("getParserConfig", ctxClass);
    getParserConfig.setAccessible(true);
    return (Config) getParserConfig.invoke(compiler, defaultCtx);
  }

  private static void dumpOne(StringBuilder sb, Config config, String name, String src) {
    line(
        sb,
        "# config languageMode="
            + config.languageMode()
            + " jsDocParsingMode="
            + config.jsDocParsingMode()
            + " runMode="
            + config.runMode()
            + " strictMode="
            + config.strictMode());
    List<String> diags = new ArrayList<>();
    ErrorReporter reporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int line, int lineOffset) {
            diags.add("WARNING " + line + ":" + lineOffset + " " + quote(message));
          }

          @Override
          public void error(String message, String sourceName, int line, int lineOffset) {
            diags.add("ERROR " + line + ":" + lineOffset + " " + quote(message));
          }
        };
    ParserRunner.ParseResult result;
    try {
      result = ParserRunner.parse(SourceFile.fromCode(name, src), src, config, reporter);
    } catch (RuntimeException e) {
      line(sb, "# diagnostics");
      for (String d : diags) {
        line(sb, d);
      }
      line(sb, "# exception " + e.getClass().getName() + ": " + e.getMessage());
      return;
    }
    line(sb, "# diagnostics");
    for (String d : diags) {
      line(sb, d);
    }
    line(sb, "# features " + result.features);
    line(sb, "# sourceMapURL " + (result.sourceMapURL == null ? "null" : quote(result.sourceMapURL)));
    line(sb, "# comments");
    for (Comment c : result.comments) {
      line(sb, c.type + " " + range(c.location) + " " + quote(c.value));
    }
    line(sb, "# tree");
    if (result.ast == null) {
      line(sb, "null");
    } else {
      dumpTree(sb, result.ast, 0);
    }
  }

  private static void dumpTree(StringBuilder sb, Node n, int level) {
    String indent = "    ".repeat(level);
    line(sb, indent + n.toString(true, true, false));
    JSDocInfo info = n.getJSDocInfo();
    if (info != null) {
      line(sb, indent + "  @jsdoc " + info.toStringVerbose());
    }
    Node shadow = n.getClosureUnawareShadow();
    if (shadow != null) {
      line(sb, indent + "  @shadow");
      dumpTree(sb, shadow, level + 1);
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      dumpTree(sb, c, level + 1);
    }
  }

  private static String range(SourceRange r) {
    return pos(r.start) + "-" + pos(r.end);
  }

  private static String pos(SourcePosition p) {
    return p.line + ":" + p.column + "@" + p.offset;
  }

  /** Wraps in double quotes; escaping happens once, per line, in line(). */
  private static String quote(String s) {
    return "\"" + s + "\"";
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

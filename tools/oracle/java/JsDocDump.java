// Oracle tool: runs com.google.javascript.jscomp.parsing.JsDocInfoParser on every JSDoc and
// important comment of a file, the way IRFactory does, and dumps the results: the JSDocInfo
// (toStringVerbose() plus every marker), the license text, the file-overview state, and every
// warning/error sent to the rhino ErrorReporter. The Go side is
// internal/parsing/jsdocdump_test.go, which must reproduce the output byte for byte.
//
// Build/run (see tools/oracle/README.md; this class lives in the parsing package because the
// JsDocInfoParser entry points IRFactory uses are package-private):
//   javac -cp closure-compiler.jar -d build JsDocDump.java
//   java -Xss1g -cp closure-compiler.jar:build com.google.javascript.jscomp.parsing.JsDocDump \
//        [--gzip] [--lang=<LanguageMode>] [--jsdoc[=<JsDocParsing>]] <outdir> <root> <rel-file>...
//        | @<listfile>
//
// For each <rel-file> writes <outdir>/<rel-file>.jsdoc (or .jsdoc.gz).
//
// Configuration: the parsing.Config is the compiler's own (Compiler.getParserConfig(
// ConfigContext.DEFAULT), called reflectively, like ParsingDump), with --jsdoc selecting
// JsDocParsing.INCLUDE_DESCRIPTIONS_NO_WHITESPACE and --jsdoc=<mode> any mode (the default is
// TYPES_ONLY); the parser.Parser.Config is ParserRunner.newParserConfig(config), like
// ParseTreeDump. The comments are Parser.getComments() after parseProgram().
//
// Each JSDOC or IMPORTANT comment, in source order, is handled as IRFactory does:
//   - createJsDocInfoParser: a JsDocTokenStream over comment.value.substring(3) starting at
//     (line + 1, column + 3), the comment's start offset, a SCRIPT template node carrying the
//     source file, the Config, JsDocSourceKind.NORMAL, setFileOverviewJSDocInfo(the current
//     first file overview), then parseImportantComment() for a non-empty IMPORTANT comment and
//     parse() otherwise;
//   - handlePossibleFileOverviewJsDoc: the license text is added to the file's license set;
//     a new file overview becomes (or is merged into) the file's first file overview;
//   - parseJSDocInfoFrom: the JSDocInfo is retrieveAndResetParsedJSDocInfo(), or null when the
//     comment was a file overview;
//   - parseInlineJSDocAt: a JSDOC comment without '@' is additionally parsed by a fresh parser
//     with parseInlineTypeDoc() (IRFactory uses one or the other depending on the AST position).
// At the end, setFileOverviewJsDoc's JSDocInfo (first file overview + joined licenses,
// build(false)) is dumped.
//
// Format (every line escaped like ParsingDump: '\' -> "\\", chars outside printable ASCII ->
// \\uXXXX; JSDocInfo$Marker identity hashes are replaced by "JSDocInfo$Marker"):
//   # config <languageMode> <strictMode> <jsDocParsingMode> <runMode> parseInlineSourceMaps=<b>
//   # annotations <Config.annotations()>
//   # suppressions <Config.suppressionNames()>
//   # primitives <Config.closurePrimitiveNames()>
//   # parse-exception <class>: <message>        if parseProgram threw
//   # comment <type> <line>:<col>@<offset>-<line>:<col>@<offset>
//   WARNING|ERROR <line>:<charno> "<message>"   each diagnostic, in report order
//   info <JSDocInfo.toStringVerbose() | null>
//   marker annotation=<pos> name=<pos> description=<pos> type=<pos>   one per marker
//   license <"text" | null>
//   fileoverview <first file overview after this comment, when this comment was one | ->
//   inline ...diagnostics... then <toStringVerbose() | null>   (JSDOC comments without '@')
//   exception <class>: <message>               if the JsDocInfoParser threw
//   # file <JSDocInfo>
// A <pos> is "null" or "<item>@<startLine>:<startChar>-<endLine>:<endChar>" (a Node item is
// rendered with toStringTree(); a type position appends "[brackets]" when hasBrackets()).
package com.google.javascript.jscomp.parsing;

import com.google.javascript.jscomp.CompilerOptions;
import com.google.javascript.jscomp.parsing.parser.Parser;
import com.google.javascript.jscomp.parsing.parser.trees.Comment;
import com.google.javascript.jscomp.parsing.parser.util.SourcePosition;
import com.google.javascript.jscomp.parsing.parser.util.SourceRange;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.SimpleSourceFile;
import com.google.javascript.rhino.StaticSourceFile;
import com.google.javascript.rhino.Token;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.regex.Pattern;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class JsDocDump {
  private static CompilerOptions.LanguageMode languageIn =
      CompilerOptions.LanguageMode.ECMASCRIPT_NEXT;
  private static Config.JsDocParsing jsdocMode = null;

  private static final Pattern MARKER_HASH =
      Pattern.compile("com\\.google\\.javascript\\.rhino\\.JSDocInfo\\$Marker@[0-9a-f]+");

  public static void main(String[] args) throws Exception {
    boolean gzip = false;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--gzip")) {
        gzip = true;
      } else if (a.equals("--jsdoc")) {
        jsdocMode = Config.JsDocParsing.INCLUDE_DESCRIPTIONS_NO_WHITESPACE;
      } else if (a.startsWith("--jsdoc=")) {
        jsdocMode = Config.JsDocParsing.valueOf(a.substring("--jsdoc=".length()));
      } else if (a.startsWith("--lang=")) {
        languageIn = CompilerOptions.LanguageMode.valueOf(a.substring("--lang=".length()));
      } else if (a.startsWith("@")) {
        for (String l : Files.readAllLines(Paths.get(a.substring(1)), StandardCharsets.UTF_8)) {
          if (!l.isEmpty()) {
            rest.add(l);
          }
        }
      } else {
        rest.add(a);
      }
    }
    CorpusExclusions.removeFrom(rest.subList(2, rest.size())); // tools/oracle/gpl-excluded-corpus.txt
    final boolean gz = gzip;
    Thread t =
        new Thread(
            null,
            () -> {
              try {
                run(gz, rest);
              } catch (Exception e) {
                throw new RuntimeException(e);
              }
            },
            "dump",
            1L << 30);
    t.start();
    t.join();
  }

  private static void run(boolean gzip, List<String> rest) throws Exception {
    Config config = compilerParserConfig();
    Method newParserConfig =
        ParserRunner.class.getDeclaredMethod("newParserConfig", Config.class);
    newParserConfig.setAccessible(true);
    Parser.Config parserConfig = (Parser.Config) newParserConfig.invoke(null, config);
    Path outDir = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    for (String rel : rest.subList(2, rest.size())) {
      String src = new String(Files.readAllBytes(root.resolve(rel)), StandardCharsets.UTF_8);
      StringBuilder sb = new StringBuilder();
      dumpOne(sb, config, parserConfig, rel, src);
      Path out = outDir.resolve(rel + (gzip ? ".jsdoc.gz" : ".jsdoc"));
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

  private static final class RecordingReporter implements ErrorReporter {
    final List<String> lines = new ArrayList<>();

    @Override
    public void warning(String message, String sourceName, int line, int lineOffset) {
      lines.add("WARNING " + line + ":" + lineOffset + " " + quote(message));
    }

    @Override
    public void error(String message, String sourceName, int line, int lineOffset) {
      lines.add("ERROR " + line + ":" + lineOffset + " " + quote(message));
    }
  }

  /** The per-file state IRFactory keeps for file-level JSDoc. */
  private static final class FileState {
    final LinkedHashSet<String> licenseBuilder = new LinkedHashSet<>();
    JSDocInfo firstFileoverview = null;
  }

  static void dumpOne(
      StringBuilder sb, Config config, Parser.Config parserConfig, String name, String src) {
    line(
        sb,
        "# config "
            + config.languageMode()
            + " "
            + config.strictMode()
            + " "
            + config.jsDocParsingMode()
            + " "
            + config.runMode()
            + " parseInlineSourceMaps="
            + config.parseInlineSourceMaps());
    line(sb, "# annotations " + config.annotations());
    line(sb, "# suppressions " + config.suppressionNames());
    line(sb, "# primitives " + config.closurePrimitiveNames());

    com.google.javascript.jscomp.parsing.parser.SourceFile file =
        new com.google.javascript.jscomp.parsing.parser.SourceFile(name, src);
    Parser p =
        new Parser(
            parserConfig,
            new com.google.javascript.jscomp.parsing.parser.util.ErrorReporter() {
              @Override
              protected void reportError(SourcePosition location, String message) {}

              @Override
              protected void reportWarning(SourcePosition location, String message) {}
            },
            file);
    try {
      p.parseProgram();
    } catch (RuntimeException | Error e) {
      line(sb, "# parse-exception " + e.getClass().getName() + ": " + e.getMessage());
    }
    List<Comment> comments = p.getComments();

    StaticSourceFile sourceFile = new SimpleSourceFile(name, StaticSourceFile.SourceKind.STRONG);
    Node templateNode = new Node(Token.SCRIPT);
    templateNode.setStaticSourceFile(sourceFile);
    FileState state = new FileState();

    for (Comment c : comments) {
      if (c.type != Comment.Type.JSDOC && c.type != Comment.Type.IMPORTANT) {
        continue;
      }
      line(sb, "# comment " + c.type + " " + range(c.location));
      RecordingReporter reporter = new RecordingReporter();
      try {
        JsDocInfoParser parser = createJsDocInfoParser(c, templateNode, config, reporter, state);
        boolean isFileOverview = handlePossibleFileOverviewJsDoc(parser, state);
        JSDocInfo info = isFileOverview ? null : parser.retrieveAndResetParsedJSDocInfo();
        flush(sb, reporter);
        line(sb, "info " + render(info));
        markers(sb, info);
        String license = parser.getLicenseText();
        line(sb, "license " + (license == null ? "null" : quote(license)));
        line(
            sb,
            "fileoverview " + (isFileOverview ? render(state.firstFileoverview) : "-"));
        if (isFileOverview) {
          markers(sb, state.firstFileoverview);
        }
      } catch (RuntimeException e) {
        flush(sb, reporter);
        line(sb, "exception " + e.getClass().getName() + ": " + e.getMessage());
      }
      if (c.type == Comment.Type.JSDOC && !c.value.contains("@")) {
        RecordingReporter inlineReporter = new RecordingReporter();
        try {
          JSDocInfo inline = parseInlineTypeDoc(c, templateNode, config, inlineReporter);
          line(sb, "inline");
          flush(sb, inlineReporter);
          line(sb, "inline " + render(inline));
          markers(sb, inline);
        } catch (RuntimeException e) {
          line(sb, "inline");
          flush(sb, inlineReporter);
          line(sb, "exception " + e.getClass().getName() + ": " + e.getMessage());
        }
      }
    }

    // IRFactory.setFileOverviewJsDoc
    JSDocInfo.Builder fileoverview =
        (state.firstFileoverview == null)
            ? JSDocInfo.builder()
            : state.firstFileoverview.toBuilder();
    if (!state.licenseBuilder.isEmpty()) {
      fileoverview.recordLicense(String.join("", state.licenseBuilder));
    }
    JSDocInfo fileInfo = fileoverview.build(false);
    line(sb, "# file " + render(fileInfo));
    markers(sb, fileInfo);
  }

  /** IRFactory.createJsDocInfoParser(Comment, boolean), with useLicensesOnlyConfig = false. */
  private static JsDocInfoParser createJsDocInfoParser(
      Comment node,
      Node templateNode,
      Config config,
      ErrorReporter errorReporter,
      FileState state) {
    String comment = node.value;
    int lineno = node.location.start.line + 1;
    int charno = node.location.start.column;
    int position = node.location.start.offset;

    // The JsDocInfoParser expects the comment without the initial '/**'.
    int numOpeningChars = 3;
    JsDocInfoParser jsdocParser =
        new JsDocInfoParser(
            new JsDocTokenStream(
                comment.substring(numOpeningChars), lineno, charno + numOpeningChars),
            comment,
            position,
            templateNode,
            config,
            JsDocInfoParser.JsDocSourceKind.NORMAL,
            errorReporter);
    jsdocParser.setFileOverviewJSDocInfo(state.firstFileoverview);
    if (node.type == Comment.Type.IMPORTANT && node.value.length() > 0) {
      jsdocParser.parseImportantComment();
    } else {
      jsdocParser.parse();
    }

    return jsdocParser;
  }

  /** IRFactory.handlePossibleFileOverviewJsDoc(JsDocInfoParser). */
  private static boolean handlePossibleFileOverviewJsDoc(
      JsDocInfoParser jsDocParser, FileState state) {
    if (jsDocParser.getLicenseText() != null) {
      state.licenseBuilder.add(jsDocParser.getLicenseText());
    }

    JSDocInfo newFileoverview = jsDocParser.getFileOverviewJSDocInfo();
    if (newFileoverview == state.firstFileoverview) {
      return false;
    }

    if (state.firstFileoverview == null) {
      state.firstFileoverview = newFileoverview;
    } else {
      JSDocInfo.Builder merged = state.firstFileoverview.toBuilder();
      if (!newFileoverview.getSuppressions().isEmpty()) {
        merged.recordSuppressions(newFileoverview.getSuppressions());
      }
      if (newFileoverview.isExterns()) {
        merged.recordExterns();
      }
      if (newFileoverview.isNoCoverage()) {
        merged.recordNoCoverage();
      }
      state.firstFileoverview = merged.build();
    }

    return true;
  }

  /** IRFactory.parseInlineTypeDoc(Comment). */
  private static JSDocInfo parseInlineTypeDoc(
      Comment node, Node templateNode, Config config, ErrorReporter errorReporter) {
    String comment = node.value;
    int lineno = node.location.start.line + 1;
    int charno = node.location.start.column;

    // The JsDocInfoParser expects the comment without the initial '/**'.
    int numOpeningChars = 3;
    JsDocInfoParser parser =
        new JsDocInfoParser(
            new JsDocTokenStream(
                comment.substring(numOpeningChars), lineno, charno + numOpeningChars),
            comment,
            node.location.start.offset,
            templateNode,
            config,
            JsDocInfoParser.JsDocSourceKind.NORMAL,
            errorReporter);
    return parser.parseInlineTypeDoc();
  }

  private static void flush(StringBuilder sb, RecordingReporter reporter) {
    for (String l : reporter.lines) {
      line(sb, l);
    }
    reporter.lines.clear();
  }

  private static String render(JSDocInfo info) {
    if (info == null) {
      return "null";
    }
    return MARKER_HASH.matcher(info.toStringVerbose()).replaceAll("JSDocInfo\\$Marker");
  }

  private static void markers(StringBuilder sb, JSDocInfo info) {
    if (info == null) {
      return;
    }
    for (JSDocInfo.Marker m : info.getMarkers()) {
      String type =
          m.getType() == null
              ? "null"
              : pos(m.getType(), m.getType().getItem().toStringTree())
                  + (m.getType().hasBrackets() ? "[brackets]" : "");
      line(
          sb,
          "marker annotation="
              + (m.getAnnotation() == null
                  ? "null"
                  : pos(m.getAnnotation(), quote(m.getAnnotation().getItem())))
              + " name="
              + (m.getNameNode() == null
                  ? "null"
                  : pos(m.getNameNode(), m.getNameNode().getItem().toStringTree()))
              + " description="
              + (m.getDescription() == null
                  ? "null"
                  : pos(m.getDescription(), quote(m.getDescription().getItem())))
              + " type="
              + type);
    }
  }

  private static String pos(com.google.javascript.rhino.SourcePosition<?> p, String item) {
    return item
        + "@"
        + p.getStartLine()
        + ":"
        + p.getPositionOnStartLine()
        + "-"
        + p.getEndLine()
        + ":"
        + p.getPositionOnEndLine();
  }

  private static String range(SourceRange r) {
    return position(r.start) + "-" + position(r.end);
  }

  private static String position(SourcePosition p) {
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

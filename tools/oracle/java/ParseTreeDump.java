// Oracle tool: dumps the raw parse tree (com.google.javascript.jscomp.parsing.parser.trees)
// that the pinned Java Closure Compiler's parser produces for a file, before IRFactory turns it
// into rhino Nodes, together with every error and warning the parser reports.
//
// Build/run (see tools/oracle/README.md):
//   javac -cp closure-compiler.jar -d build ParseTreeDump.java
//   java -cp closure-compiler.jar:build ParseTreeDump [--gzip] [--lang=<LanguageMode>] <outdir> <root> <rel-file>... | @<listfile>
//
// For each <rel-file> (relative to <root>) writes <outdir>/<rel-file>.ptree (or .ptree.gz). The Go
// port's dumper (internal/parser/parsetreedump_test.go) must reproduce this format byte for byte.
//
// The parser is constructed exactly as ParserRunner.parse does for the compiler's default
// configuration with --language_in=ECMASCRIPT_NEXT (or the --lang given): the parsing.Config comes from
// Compiler.getParserConfig(ConfigContext.DEFAULT) and the parser.Parser.Config from
// ParserRunner.newParserConfig(config) (both called reflectively, so nothing is re-derived here).
//
// Format (one value per line, depth = number of leading spaces):
//   # config <Parser.Config fields>
//   # diagnostics             then ERROR|WARNING <pos> "<message>" in report order
//   # features <FeatureSet.toString()>
//   # sourceMapURL <string or null>
//   # comments                then the Parser.getComments() list
//   # tree                    then the ProgramTree (or "null")
//   # exception <class>: <message>   if parsing threw instead
// Objects (ParseTree, Token, Comment) print "<SimpleClassName> [<type> <range>]" and then each of
// their fields (all instance fields of the class and its superclasses below ParseTree/Token, any
// visibility, sorted by name) as "name=<value>". Lists print "[n]" then their elements one level
// deeper. Positions print "line:column@offset" (0-based); ranges "start-end". Strings are quoted
// with \" \\ and \\uXXXX escapes for everything outside printable ASCII. Enums print name().
import com.google.javascript.jscomp.CompilerOptions;
import com.google.javascript.jscomp.parsing.parser.FeatureSet;
import com.google.javascript.jscomp.parsing.parser.Parser;
import com.google.javascript.jscomp.parsing.parser.SourceFile;
import com.google.javascript.jscomp.parsing.parser.Token;
import com.google.javascript.jscomp.parsing.parser.trees.Comment;
import com.google.javascript.jscomp.parsing.parser.trees.ParseTree;
import com.google.javascript.jscomp.parsing.parser.trees.ProgramTree;
import com.google.javascript.jscomp.parsing.parser.util.ErrorReporter;
import com.google.javascript.jscomp.parsing.parser.util.SourcePosition;
import com.google.javascript.jscomp.parsing.parser.util.SourceRange;
import java.io.IOException;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.zip.GZIPOutputStream;

public final class ParseTreeDump {
  private static Parser.Config parserConfig;
  private static CompilerOptions.LanguageMode languageIn = CompilerOptions.LanguageMode.ECMASCRIPT_NEXT;

  public static void main(String[] args) throws Exception {
    boolean gzip = false;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--gzip")) {
        gzip = true;
      } else if (a.startsWith("--lang=")) {
        languageIn = CompilerOptions.LanguageMode.valueOf(a.substring("--lang=".length()));
      } else if (a.startsWith("@")) {
        for (String line : Files.readAllLines(Paths.get(a.substring(1)), StandardCharsets.UTF_8)) {
          if (!line.isEmpty()) {
            rest.add(line);
          }
        }
      } else {
        rest.add(a);
      }
    }
    // Run in a thread with a large stack: Java reports "Too deep recursion while parsing" on
    // StackOverflowError, which the Go port cannot reproduce (Go stack exhaustion is fatal).
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
    parserConfig = defaultParserConfig();
    Path outDir = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    for (String rel : rest.subList(2, rest.size())) {
      Path out = outDir.resolve(rel + (gzip ? ".ptree.gz" : ".ptree"));
      Files.createDirectories(out.getParent());
      String src = new String(Files.readAllBytes(root.resolve(rel)), StandardCharsets.UTF_8);
      StringBuilder sb = new StringBuilder();
      dumpOne(sb, rel, src);
      try (OutputStream os =
          gzip ? new GZIPOutputStream(Files.newOutputStream(out)) : Files.newOutputStream(out)) {
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  /** The Parser.Config that ParserRunner.parse uses for the default compiler configuration. */
  private static Parser.Config defaultParserConfig() throws Exception {
    com.google.javascript.jscomp.Compiler compiler =
        new com.google.javascript.jscomp.Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(languageIn);
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
    Object config = getParserConfig.invoke(compiler, defaultCtx);
    Method newParserConfig =
        com.google.javascript.jscomp.parsing.ParserRunner.class.getDeclaredMethod(
            "newParserConfig", com.google.javascript.jscomp.parsing.Config.class);
    newParserConfig.setAccessible(true);
    return (Parser.Config) newParserConfig.invoke(null, config);
  }

  private static final class RecordingErrorReporter extends ErrorReporter {
    final List<String> lines = new ArrayList<>();
    private final SourceFile file;

    RecordingErrorReporter(SourceFile file) {
      this.file = file;
    }

    @Override
    protected void reportError(SourcePosition location, String message) {
      lines.add("ERROR " + position(location, file) + " " + quote(message));
    }

    @Override
    protected void reportWarning(SourcePosition location, String message) {
      lines.add("WARNING " + position(location, file) + " " + quote(message));
    }
  }

  static void dumpOne(StringBuilder sb, String name, String src) throws Exception {
    SourceFile file = new SourceFile(name, src);
    RecordingErrorReporter reporter = new RecordingErrorReporter(file);
    sb.append("# config");
    for (Field f : instanceFields(Parser.Config.class, null)) {
      sb.append(' ').append(f.getName()).append('=').append(f.get(parserConfig));
    }
    sb.append('\n');
    Parser p = new Parser(parserConfig, reporter, file);
    ProgramTree tree;
    try {
      tree = p.parseProgram();
    } catch (RuntimeException | Error e) {
      sb.append("# diagnostics\n");
      for (String l : reporter.lines) {
        sb.append(l).append('\n');
      }
      sb.append("# exception ")
          .append(e.getClass().getName())
          .append(": ")
          .append(quote(String.valueOf(e.getMessage())))
          .append('\n');
      return;
    }
    sb.append("# diagnostics\n");
    for (String l : reporter.lines) {
      sb.append(l).append('\n');
    }
    FeatureSet features = p.getFeatures();
    sb.append("# features ").append(quote(features.toString())).append('\n');
    String url = p.getSourceMapURL();
    sb.append("# sourceMapURL ").append(url == null ? "null" : quote(url)).append('\n');
    sb.append("# comments\n");
    render(sb, p.getComments(), 0, "", file);
    sb.append("# tree\n");
    render(sb, tree, 0, "", file);
  }

  private static void indent(StringBuilder sb, int depth) {
    for (int i = 0; i < depth; i++) {
      sb.append(' ');
    }
  }

  static void render(StringBuilder sb, Object v, int depth, String label, SourceFile file)
      throws IllegalAccessException {
    indent(sb, depth);
    sb.append(label);
    if (v instanceof List<?> list) {
      if (list.isEmpty()) {
        sb.append("[]\n");
        return;
      }
      sb.append('[').append(list.size()).append("]\n");
      for (Object e : list) {
        render(sb, e, depth + 1, "", file);
      }
      return;
    }
    Class<?> base;
    if (v instanceof ParseTree t) {
      sb.append(v.getClass().getSimpleName())
          .append(' ')
          .append(t.type.name())
          .append(' ')
          .append(range(t.location, file))
          .append('\n');
      base = ParseTree.class;
    } else if (v instanceof Token t) {
      sb.append(v.getClass().getSimpleName())
          .append(' ')
          .append(t.type.name())
          .append(' ')
          .append(range(t.location, file))
          .append('\n');
      base = Token.class;
    } else if (v instanceof Comment) {
      sb.append("Comment\n");
      base = Object.class;
    } else {
      sb.append(scalar(v, file)).append('\n');
      return;
    }
    for (Field f : instanceFields(v.getClass(), base)) {
      render(sb, f.get(v), depth + 1, f.getName() + "=", file);
    }
  }

  private static final Map<Class<?>, List<Field>> FIELDS = new HashMap<>();

  /** Instance fields of c and its superclasses strictly below stop, sorted by name. */
  static List<Field> instanceFields(Class<?> c, Class<?> stop) {
    return FIELDS.computeIfAbsent(
        c,
        k -> {
          List<Field> out = new ArrayList<>();
          for (Class<?> x = k; x != null && x != stop && x != Object.class; x = x.getSuperclass()) {
            for (Field f : x.getDeclaredFields()) {
              if (Modifier.isStatic(f.getModifiers()) || f.isSynthetic()) {
                continue;
              }
              f.setAccessible(true);
              out.add(f);
            }
          }
          out.sort(Comparator.comparing(Field::getName));
          return out;
        });
  }

  static String scalar(Object v, SourceFile file) {
    if (v == null) {
      return "null";
    } else if (v instanceof String s) {
      return quote(s);
    } else if (v instanceof Enum<?> e) {
      return e.name();
    } else if (v instanceof SourceRange r) {
      return range(r, file);
    } else if (v instanceof SourcePosition p) {
      return position(p, file);
    } else if (v instanceof Double d) {
      return Double.toString(d);
    } else if (v instanceof Character c) {
      return quote(String.valueOf(c));
    } else if (v instanceof Boolean || v instanceof Integer || v instanceof Long) {
      return v.toString();
    }
    return "?" + v.getClass().getName() + " " + quote(String.valueOf(v));
  }

  static String range(SourceRange r, SourceFile file) {
    if (r == null) {
      return "null";
    }
    return position(r.start, file) + "-" + position(r.end, file);
  }

  static String position(SourcePosition p, SourceFile file) {
    if (p == null) {
      return "null";
    }
    return p.line + ":" + p.column + "@" + p.offset + (p.source == file ? "" : "!src");
  }

  static String quote(String s) {
    StringBuilder sb = new StringBuilder(s.length() + 2);
    sb.append('"');
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '"' || c == '\\') {
        sb.append('\\').append(c);
      } else if (c >= 0x20 && c < 0x7f) {
        sb.append(c);
      } else {
        sb.append(String.format("\\u%04X", (int) c));
      }
    }
    return sb.append('"').toString();
  }
}

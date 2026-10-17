import com.google.javascript.jscomp.CompilerOptions;
import com.google.javascript.jscomp.SourceFile;
import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.jscomp.parsing.parser.LiteralToken;
import com.google.javascript.jscomp.parsing.parser.Parser;
import com.google.javascript.jscomp.parsing.parser.TokenType;
import com.google.javascript.jscomp.parsing.parser.trees.ParseTree;
import com.google.javascript.jscomp.parsing.parser.util.SourcePosition;
import com.google.javascript.jscomp.regex.RegExpTree;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Deque;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;
import jsgo.oracle.CorpusExclusions;

/**
 * Oracle for jscomp.regex (internal/regex).
 *
 * <pre>
 * java -cp $CP RegExpDump extract <root> @<listfile> > patterns.txt
 *     Parses each file with ParserRunner (ECMASCRIPT_NEXT, as the compiler does) and prints every
 *     distinct REGEXP literal, one per line: escaped pattern, TAB, escaped flags. For a file without
 *     an AST (a syntax error), the regular expression tokens of parser.Parser's error-recovered
 *     parse tree are used instead, split at the last '/' (without IRFactory's escape
 *     normalization: they are extra inputs, not what the compiler would check).
 * java -cp $CP RegExpDump dump <patterns.txt> > dump.txt
 *     For each pattern line, RegExpTree.parseRegExp(pattern, flags) and a deterministic dump of the
 *     result (see dumpOne); internal/regex/reg_exp_tree_oracle_test.go produces the same format.
 * </pre>
 *
 * Escaping (both modes): backslash becomes two backslashes, chars outside printable ASCII become
 * backslash-u and four hex digits, so lone surrogates survive.
 */
public final class RegExpDump {
  public static void main(String[] args) throws Exception {
    PrintStream out =
        new PrintStream(new java.io.FileOutputStream(java.io.FileDescriptor.out), false, "UTF-8");
    if (args[0].equals("extract")) {
      extract(out, Paths.get(args[1]), readList(args[2]));
    } else if (args[0].equals("dump")) {
      for (String l : Files.readAllLines(Paths.get(args[1]), StandardCharsets.UTF_8)) {
        if (l.isEmpty()) {
          continue;
        }
        int tab = l.indexOf('\t');
        StringBuilder sb = new StringBuilder();
        dumpOne(sb, unescape(l.substring(0, tab)), unescape(l.substring(tab + 1)));
        out.print(sb);
      }
    } else {
      throw new IllegalArgumentException(args[0]);
    }
    out.flush();
  }

  private static List<String> readList(String arg) throws Exception {
    List<String> files = new ArrayList<>();
    if (arg.startsWith("@")) {
      for (String l : Files.readAllLines(Paths.get(arg.substring(1)), StandardCharsets.UTF_8)) {
        if (!l.isEmpty()) {
          files.add(l);
        }
      }
    } else {
      files.add(arg);
    }
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    return files;
  }

  private static void extract(PrintStream out, Path root, List<String> files) throws Exception {
    Config config = compilerParserConfig();
    ErrorReporter reporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int line, int lineOffset) {}

          @Override
          public void error(String message, String sourceName, int line, int lineOffset) {}
        };
    Set<String> seen = new LinkedHashSet<>();
    int nFiles = 0;
    int nFailed = 0;
    int nLiterals = 0;
    for (String rel : files) {
      nFiles++;
      String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
      ParserRunner.ParseResult result;
      try {
        result = ParserRunner.parse(SourceFile.fromCode(rel, src), src, config, reporter);
      } catch (RuntimeException | StackOverflowError e) {
        result = null;
      }
      if (result == null || result.ast == null) {
        nFailed++;
        for (String key : regexpTokens(rel, src)) {
          nLiterals++;
          if (seen.add(key)) {
            out.println(key);
          }
        }
        continue;
      }
      Deque<Node> stack = new ArrayDeque<>();
      stack.push(result.ast);
      while (!stack.isEmpty()) {
        Node n = stack.pop();
        if (n.isRegExp()) {
          nLiterals++;
          String pattern = n.getFirstChild().getString();
          String flags = n.hasTwoChildren() ? n.getLastChild().getString() : "";
          String key = escape(pattern) + "\t" + escape(flags);
          if (seen.add(key)) {
            out.println(key);
          }
        }
        for (Node c = n.getLastChild(); c != null; c = c.getPrevious()) {
          stack.push(c);
        }
      }
    }
    System.err.println(
        nFiles + " files, " + nFailed + " without AST, " + nLiterals + " regexp literals, "
            + seen.size() + " distinct");
  }

  private static Parser.Config parseTreeConfig;

  /** Regular expression tokens in parser.Parser's (error-recovered) parse tree of src. */
  private static List<String> regexpTokens(String name, String src) throws Exception {
    if (parseTreeConfig == null) {
      Method newParserConfig =
          ParserRunner.class.getDeclaredMethod("newParserConfig", Config.class);
      newParserConfig.setAccessible(true);
      parseTreeConfig = (Parser.Config) newParserConfig.invoke(null, compilerParserConfig());
    }
    com.google.javascript.jscomp.parsing.parser.util.ErrorReporter reporter =
        new com.google.javascript.jscomp.parsing.parser.util.ErrorReporter() {
          @Override
          protected void reportError(SourcePosition location, String message) {}

          @Override
          protected void reportWarning(SourcePosition location, String message) {}
        };
    List<String> out = new ArrayList<>();
    try {
      Object tree =
          new Parser(
                  parseTreeConfig,
                  reporter,
                  new com.google.javascript.jscomp.parsing.parser.SourceFile(name, src))
              .parseProgram();
      collectRegexpTokens(tree, out, new java.util.IdentityHashMap<>());
    } catch (RuntimeException | StackOverflowError e) {
      // Keep what was collected.
    }
    return out;
  }

  private static void collectRegexpTokens(
      Object v, List<String> out, java.util.IdentityHashMap<Object, Boolean> visited)
      throws IllegalAccessException {
    if (v == null || visited.put(v, true) != null) {
      return;
    }
    if (v instanceof List<?> list) {
      for (Object e : list) {
        collectRegexpTokens(e, out, visited);
      }
    } else if (v instanceof LiteralToken t) {
      if (t.type == TokenType.REGULAR_EXPRESSION) {
        String value = t.value;
        int lastSlash = value.lastIndexOf('/');
        out.add(escape(value.substring(1, lastSlash)) + "\t" + escape(value.substring(lastSlash + 1)));
      }
    } else if (v instanceof ParseTree) {
      for (Class<?> x = v.getClass(); x != Object.class; x = x.getSuperclass()) {
        for (Field f : x.getDeclaredFields()) {
          if (Modifier.isStatic(f.getModifiers()) || f.getType().isPrimitive()) {
            continue;
          }
          f.setAccessible(true);
          collectRegexpTokens(f.get(v), out, visited);
        }
      }
    }
  }

  /** The parsing.Config that CompilerInput passes to ParserRunner.parse for non-extern inputs. */
  private static Config compilerParserConfig() throws Exception {
    com.google.javascript.jscomp.Compiler compiler =
        new com.google.javascript.jscomp.Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
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

  private static final Method APPEND_SOURCE_CODE;
  private static final Method APPEND_DEBUG_INFO;

  static {
    try {
      APPEND_SOURCE_CODE =
          RegExpTree.class.getDeclaredMethod("appendSourceCode", StringBuilder.class);
      APPEND_SOURCE_CODE.setAccessible(true);
      APPEND_DEBUG_INFO = RegExpTree.class.getDeclaredMethod("appendDebugInfo", StringBuilder.class);
      APPEND_DEBUG_INFO.setAccessible(true);
    } catch (ReflectiveOperationException e) {
      throw new AssertionError(e);
    }
  }

  private static String sourceCode(RegExpTree t) throws Exception {
    StringBuilder sb = new StringBuilder();
    APPEND_SOURCE_CODE.invoke(t, sb);
    return sb.toString();
  }

  private static String debugInfo(RegExpTree t) throws Exception {
    StringBuilder sb = new StringBuilder();
    APPEND_DEBUG_INFO.invoke(t, sb);
    return sb.toString();
  }

  /**
   * One record per input:
   *
   * <pre>
   * === pattern TAB flags
   * parse: exception class: message      (or the tree, see dumpTree)
   * simplify(flags): ...                 (same shape, for t.simplify(flags))
   * simplify(flags+i): ...               (only when flags has no 'i')
   * </pre>
   */
  static void dumpOne(StringBuilder sb, String pattern, String flags) {
    sb.append("=== ").append(escape(pattern)).append('\t').append(escape(flags)).append('\n');
    RegExpTree t;
    try {
      t = RegExpTree.parseRegExp(pattern, flags);
    } catch (Throwable e) {
      line(sb, "parse: " + exception(e));
      return;
    }
    dumpAnalyzed(sb, "parse", t, flags);
    simplified(sb, "simplify(flags)", t, flags);
    if (flags.indexOf('i') < 0) {
      simplified(sb, "simplify(flags+i)", t, flags + "i");
    }
  }

  private static void simplified(StringBuilder sb, String label, RegExpTree t, String flags) {
    RegExpTree s;
    try {
      s = t.simplify(flags);
    } catch (Throwable e) {
      line(sb, label + ": " + exception(e));
      return;
    }
    dumpAnalyzed(sb, label, s, flags);
    line(sb, "  equals(parse)=" + s.equals(t));
  }

  private static void dumpAnalyzed(StringBuilder sb, String label, RegExpTree t, String flags) {
    try {
      line(sb, label + ": " + t);
      line(sb, "  source=" + sourceCode(t));
      line(
          sb,
          "  groups="
              + t.numCapturingGroups()
              + " hasCapturingGroup="
              + t.hasCapturingGroup()
              + " caseSensitive="
              + t.isCaseSensitive()
              + " anchor="
              + t.containsAnchor()
              + " wholeInput="
              + RegExpTree.matchesWholeInput(t, flags)
              + " hash="
              + t.hashCode());
      dumpTree(sb, t, 1);
    } catch (Throwable e) {
      line(sb, "  " + exception(e));
    }
  }

  private static void dumpTree(StringBuilder sb, RegExpTree t, int level) throws Exception {
    line(sb, "  ".repeat(level) + t.getClass().getSimpleName() + " [" + debugInfo(t) + "]");
    for (RegExpTree c : t.children()) {
      dumpTree(sb, c, level + 1);
    }
  }

  private static String exception(Throwable e) {
    if (e instanceof java.lang.reflect.InvocationTargetException ite) {
      e = ite.getCause();
    }
    return e.getClass().getName() + ": " + e.getMessage();
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

  static String unescape(String s) {
    StringBuilder b = new StringBuilder(s.length());
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        char d = s.charAt(++i);
        if (d == '\\') {
          b.append('\\');
        } else if (d == 'u') {
          b.append((char) Integer.parseInt(s.substring(i + 1, i + 5), 16));
          i += 4;
        } else {
          throw new IllegalArgumentException(s);
        }
      } else {
        b.append(c);
      }
    }
    return b.toString();
  }

  private static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }
}

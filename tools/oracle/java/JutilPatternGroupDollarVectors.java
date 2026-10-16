// Oracle tool: writes Java reference vectors for jutil.Pattern's `$` inside a group or before
// other alternatives (java.util.regex) into internal/jutil/testdata/pattern_group_dollar.tsv.
//   java -cp build JutilPatternGroupDollarVectors internal/jutil/testdata
//
// Same line format as JutilMatcherDollarVectors.java (find / matches lines; strings as "u" +
// 4-hex-digit UTF-16 units; a match is start,end,g1start:g1end,...).
import java.io.PrintStream;
import java.nio.file.*;
import java.util.regex.*;

public class JutilPatternGroupDollarVectors {
  static String units(String s) {
    StringBuilder sb = new StringBuilder("u");
    for (int i = 0; i < s.length(); i++) sb.append(String.format("%04x", (int) s.charAt(i)));
    return sb.toString();
  }

  static String match(Matcher m) {
    StringBuilder sb = new StringBuilder();
    sb.append(m.start()).append(',').append(m.end());
    for (int g = 1; g <= m.groupCount(); g++) {
      sb.append(',').append(m.start(g)).append(':').append(m.end(g));
    }
    return sb.toString();
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    String[] patterns = {
      "(code$)", "(a$|b)", "a$|b", "(.*_test\\.js$|foo)", "((x$))", "(^a$)", "(.+_bar$)",
      "(x|(y$))", "(?:ab$|c)", "(a|b$)|c", "(\\d+$)", "([^/]+\\.js$)",
    };
    String[] inputs = {
      "", "code", "xcode", "codex", "a", "b", "ab", "ba", "x", "y", "xy", "c", "foo", "foo_test.js",
      "a/foo_test.js", "foo_test.jsx", "foo_bar", "x_bar_baz", "12", "a12", "12a", "a/b.js", "é😀.js",
    };
    try (PrintStream ps =
        new PrintStream(Files.newOutputStream(out.resolve("pattern_group_dollar.tsv")), true, "UTF-8")) {
      for (String p : patterns) {
        Pattern pat = Pattern.compile(p);
        for (String in : inputs) {
          Matcher m = pat.matcher(in);
          StringBuilder sb = new StringBuilder();
          int n = 0;
          while (m.find() && n++ < 50) {
            if (sb.length() > 0) sb.append(';');
            sb.append(match(m));
          }
          ps.println("find\t" + units(p) + "\t" + units(in) + "\t" + sb);
          m = pat.matcher(in);
          ps.println("matches\t" + units(p) + "\t" + units(in) + "\t" + (m.matches() ? match(m) : "-"));
        }
        // An input ending with a line terminator: matches() is supported.
        for (String in : new String[] {"code\n", "a\r\n", "b\r", "x "}) {
          Matcher m = pat.matcher(in);
          ps.println("matches\t" + units(p) + "\t" + units(in) + "\t" + (m.matches() ? match(m) : "-"));
        }
      }
    }
  }
}

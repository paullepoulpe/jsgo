// Oracle tool: vectors for internal/jutil/pattern_syntax_exception.go — which regexes
// java.util.regex.Pattern.compile rejects with a PatternSyntaxException, and the exception's
// description, index and message. One TSV line per pattern: the escaped pattern, then "OK", or
// "PSE", the description, the index and the message (escaped as in ParsingDump.java).
//
//   javac -d build JutilPatternSyntaxVectors.java \
//     && java -cp build JutilPatternSyntaxVectors > internal/jutil/testdata/pattern_syntax.tsv
import java.util.regex.Pattern;
import java.util.regex.PatternSyntaxException;

public final class JutilPatternSyntaxVectors {
  static String esc(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c == '\t') {
        sb.append("\\t");
      } else if (c == '\n') {
        sb.append("\\n");
      } else if (c == '\r') {
        sb.append("\\r");
      } else if (c < 0x20 || c > 0x7e) {
        sb.append(String.format("\\u%04x", (int) c));
      } else {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  static final String[] PATTERNS = {
    "(", ")", "[", "]", "*", "+", "?", "a{", "a{1", "a{1,", "a{2}", "{", "}", "\\",
    "a**", "a++", "a*?", "a??", "x{2}{3}", "a{1001}", "a{2,1}", "(eval|(|x)", "(a))", "((a)",
    "[a-", "[a-]", "[]", "[^]", "[z-a]", "\\q", "\\y", "a\\", "(?:", "(?:a", "(?:a)", "|", "a|",
    "()", "(|)", "[\\]", "[\\]]", "^*", "$*", ".*\\.closure\\.js", ".*\\.tsx\\.cl\\.js",
    "(.*blah.*)", "(foo|bar)", "javascript/closure/.*", "\\.js$", "a{,2}", "a{ 1}", "\\d+",
    "[[a]]", "[a&&b]", "x*{2}", "(*)", "(+a)", "a|*", "[a-\\d]", "\\Qab", "\\Qa(\\E(",
  };

  public static void main(String[] args) {
    for (String p : PATTERNS) {
      StringBuilder line = new StringBuilder(esc(p));
      try {
        Pattern.compile(p);
        line.append("\tOK");
      } catch (PatternSyntaxException e) {
        line.append("\tPSE\t").append(esc(e.getDescription()))
            .append('\t').append(e.getIndex())
            .append('\t').append(esc(e.getMessage()));
      }
      System.out.println(line);
    }
  }
}

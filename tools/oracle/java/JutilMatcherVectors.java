// Oracle tool: vectors for internal/jutil/matcher.go (java.util.regex.Matcher over the patterns
// that pattern.go translates). Prints one TSV line per (op, pattern, input): the JDK result,
// with strings escaped as in ParsingDump.java (\\, \t, \n, \r, \\uXXXX for non-printable-ASCII).
//
//   javac -d build JutilMatcherVectors.java && java -cp build JutilMatcherVectors > matcher.tsv
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public final class JutilMatcherVectors {
  static String esc(String s) {
    if (s == null) {
      return "<null>";
    }
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

  static final String[][] FIND_CASES = {
    {"\\{([A-Z_0-9]+)\\}", "{NUM_PEOPLE, plural, =0 {I see {BEGIN_BOLD}no one{END_BOLD} in {X_2}.}}"},
    {"\\{([A-Z_0-9]+)\\}", "café {A} 😀 {B_1}{C}"},
    {"\\{([A-Z_0-9]+)\\}", "no placeholders"},
    {"grammatical_gender_case:\\s*(\\w+)", "grammatical_gender_case: FEMININE"},
    {"grammatical_gender_case:\\s*(\\w+)", "x grammatical_gender_case:\u000bOTHER"},
    {"grammatical_gender_case:\\s*(\\w+)", "grammatical_gender_case:\tMASCULINE,grammatical_gender_case:NEUTER"},
    {"a*", "baaacaa"},
    {"a*", ""},
    {"x?", "éx😀"},
    {"(a)|(b)", "abba"},
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "$jscomp$scope$12$0$MSG_HELLO"},
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "a$jscomp$scope$1$MSG_X$MSG_Y b$jscomp$scope$2$MSG_Z"},
    {"MSG_UNNAMED.*", "MSG_UNNAMED_123"},
  };

  static final String[][] MATCH_CASES = {
    {"MSG_UNNAMED.*", "MSG_UNNAMED_123"},
    {"MSG_UNNAMED.*", "MSG_UNNAMED\n"},
    {"MSG_UNNAMED.*", "XMSG_UNNAMED"},
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "$jscomp$scope$12$0$MSG_HELLO"},
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "$jscomp$scope$1 2$MSG_HELLO"},
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "$jscomp$scope$\u000b$MSG_HELLO"},
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "MSG_$jscomp$scope$1$MSG_"},
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "$jscomp$scope$1$MSG_A$MSG_B"},
    {"a*", "baaa"},
    {"(a)|(b)", "b"},
  };

  static final String[][] REPLACE_CASES = {
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "$jscomp$scope$12$0$MSG_HELLO", "MSG_"},
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "x$jscomp$scope$1$MSG_A $jscomp$scope$2$MSG_B", "MSG_"},
    {"\\$jscomp\\$scope\\$\\S+\\$MSG_", "MSG_HELLO", "MSG_"},
    {"(a)(b)?", "xaby a", "<$2|$1|\\$>"},
    {"(a)(b)?", "xaby a", "$10"},
    {"a*", "baaacaa", "-"},
  };

  public static void main(String[] args) {
    for (String[] c : FIND_CASES) {
      Matcher m = Pattern.compile(c[0]).matcher(c[1]);
      StringBuilder sb = new StringBuilder();
      while (m.find()) {
        sb.append(m.start()).append(',').append(m.end());
        for (int g = 1; g <= m.groupCount(); g++) {
          sb.append(',').append(esc(m.group(g))).append('@').append(m.start(g)).append('-')
              .append(m.end(g));
        }
        sb.append(';');
      }
      System.out.println("find\t" + esc(c[0]) + "\t" + esc(c[1]) + "\t" + sb);
    }
    for (String[] c : MATCH_CASES) {
      Matcher m = Pattern.compile(c[0]).matcher(c[1]);
      String matches = m.matches() ? "true:" + m.end() : "false";
      m.reset();
      String lookingAt = m.lookingAt() ? "true:" + m.end() : "false";
      System.out.println("match\t" + esc(c[0]) + "\t" + esc(c[1]) + "\t" + matches + "\t" + lookingAt);
    }
    for (String[] c : REPLACE_CASES) {
      Pattern p = Pattern.compile(c[0]);
      String first;
      String all;
      try {
        first = esc(p.matcher(c[1]).replaceFirst(c[2]));
      } catch (RuntimeException e) {
        first = "!" + esc(e.toString());
      }
      try {
        all = esc(p.matcher(c[1]).replaceAll(c[2]));
      } catch (RuntimeException e) {
        all = "!" + esc(e.toString());
      }
      System.out.println(
          "replace\t" + esc(c[0]) + "\t" + esc(c[1]) + "\t" + esc(c[2]) + "\t" + first + "\t" + all);
    }
  }
}

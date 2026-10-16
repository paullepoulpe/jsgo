// Oracle tool: writes internal/jutil/testdata/string_math.tsv, the JDK's
// Math.round(double), Math.max/min(double, double), Math.signum(double),
// String.lastIndexOf(String, int) and Matcher.replaceAll/replaceFirst (the
// replacement-string syntax of Matcher.appendReplacement), for
// internal/jutil/string_math_test.go.
//   java -cp .cache/oracle-build JutilStringMathVectors internal/jutil/testdata [seed [count]]
//
// Rows (tab-separated; doubles as 16 hex digits of their raw bits, strings
// as "u" followed by 4 hex digits per UTF-16 unit, an exception as "!"
// followed by its toString(), which is printable ASCII):
//   round    <a>  <Math.round(a)>
//   signum   <a>  <Math.signum(a)>
//   max      <a>  <b>  <Math.max(a, b)>
//   min      <a>  <b>  <Math.min(a, b)>
//   lastidx  <s>  <str>  <fromIndex>  <s.lastIndexOf(str, fromIndex)>
//   replace  <regex>  <input>  <replacement>  <replaceAll>  <replaceFirst>
import java.io.PrintStream;
import java.nio.file.*;
import java.util.*;
import java.util.regex.*;

public class JutilStringMathVectors {
  static String units(String s) {
    StringBuilder sb = new StringBuilder("u");
    for (int i = 0; i < s.length(); i++) sb.append(String.format("%04x", (int) s.charAt(i)));
    return sb.toString();
  }

  static String bits(double d) {
    return String.format("%016x", Double.doubleToRawLongBits(d));
  }

  // Doubles where rounding and comparisons have edge cases, and random ones.
  static double randomDouble(Random r) {
    switch (r.nextInt(10)) {
      case 0:
        return Double.longBitsToDouble(r.nextLong()); // any bits, NaNs included
      case 1: {
        double[] special = {0.0, -0.0, Double.NaN, Double.POSITIVE_INFINITY, Double.NEGATIVE_INFINITY,
          Double.MIN_VALUE, -Double.MIN_VALUE, Double.MIN_NORMAL, Double.MAX_VALUE, -Double.MAX_VALUE,
          0x1p63, -0x1p63, 0x1p63 - 1024, -0x1p63 + 1024, 0x1p62, 0x1p52, -0x1p52, 0x1p53, 0x1p53 + 2,
          0.49999999999999994, -0.49999999999999994, 0.5, -0.5, 1.5, -1.5, 1.0, -1.0,
          Double.longBitsToDouble(0x7ff8000000000001L), Double.longBitsToDouble(0xfff0000000000123L),
          Double.longBitsToDouble(0x7ff0000000000001L)};
        return special[r.nextInt(special.length)];
      }
      case 2:
      case 3: {
        // Halfway cases k + 0.5 and their neighbours, at every magnitude.
        long k = r.nextLong() >> r.nextInt(64);
        double h = k + 0.5;
        switch (r.nextInt(3)) {
          case 0: return h;
          case 1: return Math.nextUp(h);
          default: return Math.nextDown(h);
        }
      }
      case 4: {
        double d = Math.scalb(1.0, 50 + r.nextInt(15)) * (r.nextBoolean() ? 1 : -1);
        for (int i = r.nextInt(4); i > 0; i--) d = r.nextBoolean() ? Math.nextUp(d) : Math.nextDown(d);
        return d;
      }
      case 5:
        return (r.nextInt(41) - 20) / 4.0;
      case 6:
        return r.nextGaussian() * Math.pow(10, r.nextInt(40) - 20);
      default:
        return (r.nextDouble() - 0.5) * Math.scalb(1.0, r.nextInt(140) - 70);
    }
  }

  static final String[] STRING_PARTS = {"a", "b", "ab", "\u00e9", "\uD83D\uDE00", "\uD83D", "\uDE00", "\uDE00\uD83D", ""};

  static String randomParts(Random r, int max) {
    StringBuilder sb = new StringBuilder();
    for (int i = r.nextInt(max + 1); i > 0; i--) sb.append(STRING_PARTS[r.nextInt(STRING_PARTS.length)]);
    return sb.toString();
  }

  static final String[] PATTERNS = {
    "a", "(a)", "(a)(b)?", "(a)|(b)", "[ab]+", "(\u00e9)", "(\uD83D\uDE00)", "((a)(b))(c)?",
    "(a)(b)(c)(d)(e)(f)(g)(h)(i)(j)(k)(l)", "\\s+", "(\\d+)", "(a)(a)?(a)?", "(b)(x)?",
  };
  // Patterns that can match the empty string: their inputs stay in the BMP
  // (the Go Matcher cannot resume inside a surrogate pair).
  static final String[] EMPTY_PATTERNS = {"", "a*", "(x)?", "(a*)(b*)"};
  static final String INPUT_CHARS = "aabbcx \u00e9123";
  static final String[] REPLACEMENT_TOKENS = {
    "x", "-", "\u00e9", "\uD83D\uDE00", "\uDE00", "$0", "$1", "$2", "$3", "$9", "$10", "$11", "$12", "$13", "$01", "$00",
    "$123", "$1x", "\\$", "\\\\", "\\x", "\\\u00e9", "\\\uD83D\uDE00", "${a}", "${A1}", "${1a}", "${}", "${a", "${a-b}",
    "${\u00e9}", "${_a}", "$x", "$ ", "$$", "$-1", "$\u0661", "$\uff11", "{", "}", " ",
  };

  static String replacementOf(Random r) {
    StringBuilder sb = new StringBuilder();
    for (int i = r.nextInt(4); i >= 0; i--) {
      if (r.nextInt(3) == 0) {
        sb.append("$").append(r.nextInt(r.nextBoolean() ? 20 : 2000));
      } else {
        sb.append(REPLACEMENT_TOKENS[r.nextInt(REPLACEMENT_TOKENS.length)]);
      }
    }
    switch (r.nextInt(12)) {
      case 0: sb.append('$'); break;
      case 1: sb.append('\\'); break;
      case 2: sb.append("${"); break;
      default: break;
    }
    return sb.toString();
  }

  static String outcome(java.util.function.Supplier<String> f) {
    try {
      return units(f.get());
    } catch (Throwable t) {
      String m = t.toString();
      if (!m.chars().allMatch(c -> c >= 0x20 && c < 0x7f)) throw new IllegalStateException("not printable ASCII: " + m);
      return "!" + m;
    }
  }

  static void replace(PrintStream ps, String regex, String input, String repl) {
    Pattern p = Pattern.compile(regex);
    ps.println("replace\t" + units(regex) + "\t" + units(input) + "\t" + units(repl) + "\t"
        + outcome(() -> p.matcher(input).replaceAll(repl)) + "\t" + outcome(() -> p.matcher(input).replaceFirst(repl)));
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    long seed = args.length > 1 ? Long.parseLong(args[1]) : 20261016L;
    int count = args.length > 2 ? Integer.parseInt(args[2]) : 2000;
    Random r = new Random(seed);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("string_math.tsv")), false, "UTF-8")) {
      for (int i = 0; i < count; i++) {
        double a = randomDouble(r);
        ps.println("round\t" + bits(a) + "\t" + Math.round(a));
        ps.println("signum\t" + bits(a) + "\t" + bits(Math.signum(a)));
        double b = r.nextInt(4) == 0 ? (r.nextBoolean() ? a : -a) : randomDouble(r);
        ps.println("max\t" + bits(a) + "\t" + bits(b) + "\t" + bits(Math.max(a, b)));
        ps.println("min\t" + bits(a) + "\t" + bits(b) + "\t" + bits(Math.min(a, b)));
      }

      for (int i = 0; i < count; i++) {
        String s = randomParts(r, 6);
        String sub;
        if (r.nextBoolean() && s.length() > 0) {
          int from = r.nextInt(s.length() + 1);
          sub = s.substring(from, from + r.nextInt(Math.min(4, s.length() - from) + 1));
        } else {
          sub = randomParts(r, 2);
        }
        int fromIndex;
        switch (r.nextInt(12)) {
          case 0: fromIndex = Integer.MIN_VALUE; break;
          case 1: fromIndex = Integer.MAX_VALUE; break;
          default: fromIndex = r.nextInt(s.length() + 7) - 3; break;
        }
        ps.println("lastidx\t" + units(s) + "\t" + units(sub) + "\t" + fromIndex + "\t" + s.lastIndexOf(sub, fromIndex));
      }

      // Fixed cases, then random replacements on every pattern.
      String[] fixed = {"", "x", "$0", "$1", "$2", "$9", "$10", "$11", "$12", "$01", "$001", "$00", "$1x", "$", "x$", "$x",
        "$-1", "\\", "x\\", "\\$1", "\\\\", "\\x", "\\n", "${", "${}", "${a", "${a}", "${A1}", "${1a}", "${1}", "${1a", "${a-b}",
        "${a b}", "${\u00e9}", "${a}}", "$ {a}", "$${a}", "$1$", "\\$", "$\\1", "${_a}", "${a_b}", "$99999999999",
        "$4294967297", "$2147483648", "$1\\", "$\u0661", "$\uff11", "$1\u0661", "${a}$9", "$9${a}", "\\$9", "$\uD83D\uDE00",
        "\\\uD83D\uDE00", "x\\\\$", "$9\\", "[$1|$2]", "$0$0", "<$1>", "$13$1"};
      String[][] inputs = {{"(a)(b)?", "xaay"}, {"(a)", "a"}, {"a", "bab"}, {"x", "bab"}, {"(a)|(b)", "ab"}, {"", "ab"},
        {"(a)(b)(c)(d)(e)(f)(g)(h)(i)(j)(k)(l)", "abcdefghijkl"}, {"(\u00e9)", "\u00e9t\u00e9"}};
      for (String[] in : inputs) for (String repl : fixed) replace(ps, in[0], in[1], repl);
      for (int i = 0; i < count / 2; i++) {
        boolean empty = r.nextInt(4) == 0;
        String regex = empty ? EMPTY_PATTERNS[r.nextInt(EMPTY_PATTERNS.length)] : PATTERNS[r.nextInt(PATTERNS.length)];
        StringBuilder input = new StringBuilder();
        for (int j = r.nextInt(14); j > 0; j--) {
          if (!empty && r.nextInt(8) == 0) input.append("\uD83D\uDE00");
          else if (!empty && r.nextInt(12) == 0) input.append("abcdefghijkl");
          else input.append(INPUT_CHARS.charAt(r.nextInt(INPUT_CHARS.length())));
        }
        replace(ps, regex, input.toString(), replacementOf(r));
      }
    }
  }
}

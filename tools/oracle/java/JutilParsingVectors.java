// Oracle tool: writes the Java reference vectors for the internal/jutil helpers added for
// internal/parsing (parse_double.tsv, bigint_radix.tsv, properties.tsv).
//   javac -d .cache/oracle-build tools/oracle/java/JutilParsingVectors.java
//   java -cp .cache/oracle-build JutilParsingVectors internal/jutil/testdata \
//       internal/parsing/ParserConfig.properties [seed]
// Strings are written as "x" + the hex of their UTF-16 code units.
import java.io.PrintStream;
import java.io.StringReader;
import java.math.BigInteger;
import java.nio.file.*;
import java.util.*;

public class JutilParsingVectors {
  static String hex(String s) {
    StringBuilder sb = new StringBuilder("x");
    for (int i = 0; i < s.length(); i++) sb.append(String.format("%04x", (int) s.charAt(i)));
    return sb.toString();
  }

  static String outcome(Throwable t) {
    return "E:" + t.getClass().getSimpleName() + ":" + hex(String.valueOf(t.getMessage()));
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(args.length > 2 ? Long.parseLong(args[2]) : 4242);

    List<String> doubles = new ArrayList<>(List.of(
        "", " ", "0", "-0", "+0", "1", "1.", ".5", ".", "-", "+", "1e", "1e+", "1e-5", "1E5", "1e400",
        "1e-400", "4.9e-324", "2.4703282292062328e-324", "2.4703282292062327e-324", "1.7976931348623157e308",
        "1.7976931348623158e308", "1.7976931348623159e308", "NaN", "-NaN", "+NaN", "nan", "Infinity",
        "-Infinity", "+Infinity", "infinity", "Inf", "1f", "1F", "1d", "1D", "1.5fx", "1.5ff", " 12 ",
        "\t12\n", "\u000012\u0001", "1.2.3", "1..2", "00012", "0.000001", "123456789012345678901234567890",
        "0x1p0", "0x1P0", "0X1p0", "0x1.8p1", "0x.8p1", "0x1.p1", "0x.p1", "0xp1", "0x1", "0x1p", "0x1p+",
        "-0x1p-1074", "0x1p-1075", "0x1.0000000000001p-1075", "0x1.fffffffffffff8p1023", "0x1p1024",
        "0x1p0f", "0x1p0d", "1_000", "1e5.5", "e5", ".e5", "1.e5", "0.1", "0.2", "0.3", "9007199254740993",
        "1٠", "١", "12e2147483648", "12e-2147483649", "0e999999999999999999"));
    for (int k = 0; k < 3000; k++) {
      StringBuilder sb = new StringBuilder();
      if (r.nextInt(5) == 0) sb.append(r.nextBoolean() ? "-" : "+");
      int intDigits = r.nextInt(25);
      for (int i = 0; i < intDigits; i++) sb.append((char) ('0' + r.nextInt(10)));
      if (r.nextInt(3) > 0) {
        sb.append('.');
        int frac = r.nextInt(25);
        for (int i = 0; i < frac; i++) sb.append((char) ('0' + r.nextInt(10)));
      }
      if (r.nextInt(3) == 0) {
        sb.append(r.nextBoolean() ? 'e' : 'E');
        if (r.nextBoolean()) sb.append(r.nextBoolean() ? '-' : '+');
        sb.append(r.nextInt(r.nextBoolean() ? 30 : 400));
      }
      doubles.add(sb.toString());
    }
    for (int k = 0; k < 1000; k++) {
      double d = Double.longBitsToDouble(r.nextLong());
      if (Double.isNaN(d) || Double.isInfinite(d)) continue;
      doubles.add(Double.toString(d));
      doubles.add(Double.toHexString(d));
      doubles.add(String.format("%.17g", d));
    }
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("parse_double.tsv")))) {
      for (String s : doubles) {
        String res;
        try {
          res = Long.toHexString(Double.doubleToRawLongBits(Double.parseDouble(s)));
        } catch (RuntimeException e) {
          res = outcome(e);
        }
        ps.println(hex(s) + "\t" + res);
      }
    }

    List<String[]> bigints = new ArrayList<>();
    String[] fixed = {"", "-", "+", "0", "-0", "00", "1", "-1", "+1", "1-", "1+2", "--1", "-+1", "0x10", "ff",
        "FF", "fF", "1010", "777", "8", "z", "Z", "١٢", "1_0", " 1"};
    for (String s : fixed) {
      for (int radix : new int[] {2, 8, 10, 16, 36, 1, 37}) bigints.add(new String[] {s, "" + radix});
    }
    for (int k = 0; k < 600; k++) {
      int radix = new int[] {2, 8, 16, 10, 3, 36}[r.nextInt(6)];
      BigInteger x = new BigInteger(1 + r.nextInt(300), r);
      if (r.nextBoolean()) x = x.negate();
      String s = x.toString(radix);
      if (radix == 16 && r.nextBoolean()) s = s.toUpperCase(Locale.ROOT);
      if (r.nextInt(4) == 0) s = "000" + s;
      bigints.add(new String[] {s, "" + radix});
    }
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("bigint_radix.tsv")))) {
      for (String[] b : bigints) {
        String res;
        try {
          res = new BigInteger(b[0], Integer.parseInt(b[1])).toString();
        } catch (RuntimeException e) {
          res = outcome(e);
        }
        ps.println(hex(b[0]) + "\t" + b[1] + "\t" + res);
      }
    }

    List<String> texts = new ArrayList<>();
    texts.add(new String(Files.readAllBytes(Paths.get(args[1])), "UTF-8"));
    texts.addAll(List.of(
        "a=b", "a = b", "a:b", "a b", "a\tb", "a\fb", "  a=b  ", "a==b", "a=:b", "a :=b", "a", "a=", "=b",
        "# c\na=b", "! c\na=b", "a=b\\\n  c", "a=b\\\r\n  c", "a=b\\\rc", "a=b\\\\\nc=d", "a=b\\", "a=b\\\\",
        "a\\ b=c", "a\\=b=c", "a\\:b=c", "a=\\u0041\\t\\n\\r\\f\\x", "a=\\u004", "a=\\u00g1", "a=b\nb=c\na=d",
        "a=b\\\n# not a comment", "a=b\\\n\nc=d", "\n\n  \n", "#only", "a=b\r\nc=d\rd=e\n", "x=é€",
        "a=b\\\n", "a=1,\\\n   2,\\\n   3,\\\n\n# c\nb=x"));
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("properties.tsv")))) {
      for (String t : texts) {
        String res;
        try {
          Properties p = new Properties();
          p.load(new StringReader(t));
          List<String> keys = new ArrayList<>(p.stringPropertyNames());
          Collections.sort(keys);
          StringBuilder sb = new StringBuilder();
          for (String k : keys) sb.append(sb.length() == 0 ? "" : ",").append(hex(k)).append("=").append(hex(p.getProperty(k)));
          res = "M:" + sb;
        } catch (RuntimeException | java.io.IOException e) {
          res = outcome(e);
        }
        ps.println(hex(t) + "\t" + res);
      }
    }
  }
}

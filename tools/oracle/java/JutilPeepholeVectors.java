// Generates the JDK vectors for the jutil helpers the peephole passes need:
//   internal/jutil/testdata/java_math.tsv.gz  Math.pow/round/signum/max/min, double %,
//                                              BigInteger.mod/intValueExact
//   internal/jutil/testdata/string_ops.tsv    String.lastIndexOf(String, int),
//                                              Guava Ascii.equalsIgnoreCase
//   internal/jutil/testdata/string_case_root.tsv
//                                              Character.toUpperCase/toLowerCase(int) for
//                                              every code point whose mapping is not the
//                                              identity, and every code point whose
//                                              String.toUpperCase(Locale.ROOT) is not that
//                                              simple mapping
//   internal/jutil/jstring_case_table.go      the special-casing table itself
//
//   java -cp .cache/closure-compiler.jar:.cache/oracle-build JutilPeepholeVectors internal/jutil
//
// Doubles are written as their raw 64-bit patterns in hex; strings as "u" + UTF-16 code units
// (4 hex digits each).
import com.google.javascript.jscomp.jarjar.com.google.common.base.Ascii;
import java.io.PrintStream;
import java.math.BigInteger;
import java.nio.file.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;
import java.util.Random;

public class JutilPeepholeVectors {
  static String d(double x) {
    return Long.toHexString(Double.doubleToRawLongBits(x));
  }

  static String u(String s) {
    StringBuilder sb = new StringBuilder("u");
    for (int i = 0; i < s.length(); i++) {
      sb.append(String.format("%04x", (int) s.charAt(i)));
    }
    return sb.toString();
  }

  static final double[] SPECIAL = {
    0.0, -0.0, 1.0, -1.0, 2.0, -2.0, 0.5, -0.5, 3.0, -3.0, 1.5, -1.5, 10.0, 0.1, 1e-300, 1e300,
    Double.MIN_VALUE, -Double.MIN_VALUE, Double.MAX_VALUE, -Double.MAX_VALUE, Double.MIN_NORMAL,
    Double.POSITIVE_INFINITY, Double.NEGATIVE_INFINITY, Double.NaN, 0.49999999999999994,
    -0.49999999999999994, 4503599627370495.5, -4503599627370495.5, 4503599627370496.0,
    9007199254740993.0, 2.5, -2.5, 1e16, -1e16, 1 - 1e-16, 1 + 2.220446049250313e-16, 53, -53,
    1024, -1074, 0.25, 1.0 / 3, 2.7182818, 2.2
  };

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(20260415);
    try (PrintStream ps =
        new PrintStream(
            new java.util.zip.GZIPOutputStream(
                Files.newOutputStream(out.resolve("testdata/java_math.tsv.gz"))))) {
      // Binary double ops over the special grid.
      for (double x : SPECIAL) {
        for (double y : SPECIAL) {
          binary(ps, x, y);
        }
        unary(ps, x);
      }
      // Integer powers (exact cases and not), as constant folding sees them.
      for (int b = -40; b <= 40; b++) {
        for (int e = -40; e <= 40; e++) {
          ps.println("pow\t" + d(b) + "\t" + d(e) + "\t" + d(Math.pow(b, e)));
        }
      }
      for (int b : new int[] {2, 10, 16, 1024, 3, 7}) {
        for (int e = -1100; e <= 1100; e += 7) {
          ps.println("pow\t" + d(b) + "\t" + d(e) + "\t" + d(Math.pow(b, e)));
        }
      }
      // Decimal bases and exponents.
      for (int i = 0; i < 4000; i++) {
        double x = r.nextInt(1000) / 10.0;
        double y = r.nextInt(100) / 10.0 - 5;
        ps.println("pow\t" + d(x) + "\t" + d(y) + "\t" + d(Math.pow(x, y)));
      }
      // Random doubles.
      for (int i = 0; i < 2000; i++) {
        double x = r.nextDouble() * 100;
        double y = r.nextDouble() * 20 - 10;
        binary(ps, x, y);
        binary(ps, -x, Math.rint(y));
        unary(ps, x * (r.nextBoolean() ? 1 : -1) * Math.pow(10, r.nextInt(40) - 20));
        unary(ps, Double.longBitsToDouble(r.nextLong()));
        binary(ps, Double.longBitsToDouble(r.nextLong()), Double.longBitsToDouble(r.nextLong()));
      }
      // BigInteger.mod and intValueExact.
      List<BigInteger> bigs = new ArrayList<>();
      for (long v : new long[] {0, 1, -1, 2, -2, 7, -7, 2147483647L, -2147483648L, 2147483648L,
          -2147483649L, Long.MAX_VALUE, Long.MIN_VALUE}) {
        bigs.add(BigInteger.valueOf(v));
      }
      for (int i = 0; i < 200; i++) {
        BigInteger b = new BigInteger(1 + r.nextInt(100), r);
        bigs.add(r.nextBoolean() ? b : b.negate());
      }
      for (BigInteger a : bigs) {
        String ive;
        try {
          ive = Integer.toString(a.intValueExact());
        } catch (ArithmeticException e) {
          ive = "AE:" + e.getMessage();
        }
        ps.println("intValueExact\t" + a + "\t\t" + ive);
        for (int k = 0; k < 4; k++) {
          BigInteger m = bigs.get(r.nextInt(bigs.size()));
          String res;
          try {
            res = a.mod(m).toString();
          } catch (ArithmeticException e) {
            res = "AE:" + e.getMessage();
          }
          ps.println("mod\t" + a + "\t" + m + "\t" + res);
        }
      }
    }

    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("testdata/string_ops.tsv")))) {
      String[] alphabet = {"a", "b", "ab", "😀", "\uD83D", "\uDE00", "é", "x"};
      for (int i = 0; i < 3000; i++) {
        String s = randomString(r, alphabet, r.nextInt(8));
        String t = randomString(r, alphabet, r.nextInt(3));
        int from = r.nextInt(s.length() + 6) - 3;
        if (r.nextInt(20) == 0) {
          from = r.nextBoolean() ? Integer.MAX_VALUE : Integer.MIN_VALUE;
        }
        ps.println("lastIndexOf\t" + u(s) + "\t" + u(t) + "\t" + from + "\t" + s.lastIndexOf(t, from));
      }
      String[] ascii = {"0x", "0X", "0y", "x0", "", "a", "A", "à", "À", "@", "`", "[", "{"};
      for (String a : ascii) {
        for (String b : ascii) {
          ps.println("equalsIgnoreCase\t" + u(a) + "\t" + u(b) + "\t\t" + Ascii.equalsIgnoreCase(a, b));
        }
      }
      for (int c = 0; c < 0x10000; c += 1 + r.nextInt(40)) {
        String a = "0" + (char) c;
        String b = "0" + (char) (c ^ 0x20);
        ps.println("equalsIgnoreCase\t" + u(a) + "\t" + u(b) + "\t\t" + Ascii.equalsIgnoreCase(a, b));
      }
    }

    StringBuilder upperTable = new StringBuilder();
    try (PrintStream ps =
        new PrintStream(Files.newOutputStream(out.resolve("testdata/string_case_root.tsv")))) {
      for (int cp = 0; cp <= Character.MAX_CODE_POINT; cp++) {
        int su = Character.toUpperCase(cp);
        int sl = Character.toLowerCase(cp);
        if (su != cp || sl != cp) {
          ps.println("simple\t" + Integer.toHexString(cp) + "\t" + Integer.toHexString(su) + "\t"
              + Integer.toHexString(sl));
        }
        if (Character.isSurrogate((char) cp) && cp < 0x10000) {
          continue;
        }
        String s = new String(Character.toChars(cp));
        String up = s.toUpperCase(Locale.ROOT);
        if (!up.equals(new String(Character.toChars(su)))) {
          ps.println("upper\t" + Integer.toHexString(cp) + "\t" + u(up));
          upperTable.append("\t0x").append(Integer.toHexString(cp)).append(": ")
              .append(goString(up)).append(",\n");
        }
      }
      // Whole strings, including context-free ones mixing planes and lone surrogates.
      String[] alphabet = {"a", "Z", "ß", "İ", "ŉ", "ﬀ", "𐐀",
          "𐐨", "\uD83D", "\uDE00", "ǅ", "ẞ", "ς", "1", " ", "և"};
      for (int i = 0; i < 2000; i++) {
        String s = randomString(r, alphabet, r.nextInt(10));
        ps.println("string\t" + u(s) + "\t" + u(s.toUpperCase(Locale.ROOT)));
      }
    }

    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("jstring_case_table.go")))) {
      ps.println("// Code generated by tools/oracle/java/JutilPeepholeVectors.java from JDK 21; DO NOT EDIT.");
      ps.println();
      ps.println("package jutil");
      ps.println();
      ps.println("// stringUpperCaseRootSpecial holds, for every code point whose");
      ps.println("// String.toUpperCase(Locale.ROOT) is not the simple Character.toUpperCase(int)");
      ps.println("// mapping, the full uppercase form (from SpecialCasing.txt).");
      ps.println("var stringUpperCaseRootSpecial = map[rune]string{");
      ps.print(upperTable);
      ps.println("}");
    }
  }

  static String goString(String s) {
    StringBuilder sb = new StringBuilder("\"");
    s.codePoints().forEach(cp -> sb.append(String.format("\\U%08x", cp)));
    return sb.append('"').toString();
  }

  static String randomString(Random r, String[] alphabet, int n) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < n; i++) {
      sb.append(alphabet[r.nextInt(alphabet.length)]);
    }
    return sb.toString();
  }

  static void binary(PrintStream ps, double x, double y) {
    ps.println("pow\t" + d(x) + "\t" + d(y) + "\t" + d(Math.pow(x, y)));
    ps.println("max\t" + d(x) + "\t" + d(y) + "\t" + d(Math.max(x, y)));
    ps.println("min\t" + d(x) + "\t" + d(y) + "\t" + d(Math.min(x, y)));
    ps.println("rem\t" + d(x) + "\t" + d(y) + "\t" + d(x % y));
  }

  static void unary(PrintStream ps, double x) {
    ps.println("round\t" + d(x) + "\t\t" + Math.round(x));
    ps.println("signum\t" + d(x) + "\t\t" + d(Math.signum(x)));
  }
}

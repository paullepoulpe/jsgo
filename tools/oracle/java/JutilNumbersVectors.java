// Oracle tool: writes the Java reference vectors for the internal/jutil number helpers
// (Integer.parseInt, Long.parseLong, new BigDecimal(String) and its accessors,
// Double.parseDouble, new BigInteger(String, int), BigInteger.pow,
// String.format("%.<n>f", double), Arrays.binarySearch(int[], int) and
// Random.nextInt(int)) → internal/jutil/testdata/numbers.tsv.gz, replayed by
// internal/jutil/numbers_vectors_test.go. Only the public JDK API is called.
//   javac -d .cache/oracle-build tools/oracle/java/JutilNumbersVectors.java
//   java -cp .cache/oracle-build JutilNumbersVectors internal/jutil/testdata/numbers.tsv.gz [seed]
// Strings are written as "x" + the hex of their UTF-16 code units; an exception as
// "E:" + the hex of its toString() (class name and message, if any).
import java.io.*;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.function.Supplier;
import java.util.zip.GZIPOutputStream;

public class JutilNumbersVectors {
  static Random r;
  static PrintWriter out;

  static String hex(String s) {
    StringBuilder sb = new StringBuilder("x");
    for (int i = 0; i < s.length(); i++) sb.append(String.format("%04x", (int) s.charAt(i)));
    return sb.toString();
  }

  static String run(Supplier<Object> f) {
    try {
      return "=" + f.get();
    } catch (RuntimeException e) {
      return "E:" + hex(e.toString());
    }
  }

  static <T> T pick(T[] a) {
    return a[r.nextInt(a.length)];
  }

  static final String[] ODD_CHARS = {
    "٠", "٩", "०", "０", "９", "Ａ", "ａ", "ｚ", "²", "Ⅷ",
    "𝟎", "\ud835", "\udfce", " ", "\t", "\u0000", " ", "_", ",", "x", "X", "z", "Z", "g", "G",
    "e", "E", "p", "P", "f", "F", "d", "D", ".", "-", "+", "a", "A", "N", "I"
  };

  static String digits(int n, String alphabet) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < n; i++) sb.append(alphabet.charAt(r.nextInt(alphabet.length())));
    return sb.toString();
  }

  // Inserts, replaces or deletes a few chars of s at random places.
  static String mutate(String s, String[] chars) {
    StringBuilder sb = new StringBuilder(s);
    int k = 1 + r.nextInt(3);
    for (int i = 0; i < k; i++) {
      int at = sb.length() == 0 ? 0 : r.nextInt(sb.length() + 1);
      switch (r.nextInt(3)) {
        case 0 -> sb.insert(at, pick(chars));
        case 1 -> {
          if (at < sb.length()) sb.replace(at, at + 1, pick(chars));
        }
        default -> {
          if (at < sb.length()) sb.deleteCharAt(at);
        }
      }
    }
    return sb.toString();
  }

  static void integers() {
    List<String> in = new ArrayList<>(List.of(
        "", "-", "+", "0", "-0", "+0", "00", "1", "-1", "+1", "--1", "+-1", "-+1", "1-", "1+", " 1", "1 ",
        "2147483647", "2147483648", "-2147483648", "-2147483649", "+2147483647", "4294967295", "4294967296",
        "9223372036854775807", "9223372036854775808", "-9223372036854775808", "-9223372036854775809",
        "18446744073709551615", "18446744073709551616", "99999999999999999999999", "000000000000000000000000001",
        "7fffffff", "80000000", "-80000000", "-80000001", "ffffffff", "7fffffffffffffff", "8000000000000000",
        "-8000000000000000", "-8000000000000001", "zzzzzz", "ZZZZZZ", "1y2p0ij32e8e7", "-1y2p0ij32e8e8",
        "1y2p0ij32e8e8", "1111111111111111111111111111111", "10000000000000000000000000000000",
        "-10000000000000000000000000000000", "١٢", "１２", "०१", "Ａｂ",
        "𝟎", "1_000", "0x10", "1e3", "1.0", "\u0000", "12\u0000"));
    int[] radixes = {2, 3, 8, 10, 16, 36, 1, 37, 0, -10};
    for (int k = 0; k < 1500; k++) {
      int radix = 2 + r.nextInt(35);
      String s;
      switch (r.nextInt(5)) {
        case 0 -> s = Integer.toString(r.nextInt(), radix);
        case 1 -> s = Long.toString(r.nextLong(), radix);
        case 2 -> s = Long.toString(r.nextBoolean() ? Long.MAX_VALUE - r.nextInt(3) : Long.MIN_VALUE + r.nextInt(3), radix)
            + (r.nextInt(4) == 0 ? "0" : "");
        case 3 -> s = Long.toString((r.nextBoolean() ? Integer.MAX_VALUE : Integer.MIN_VALUE) + (long) (r.nextInt(5) - 2), radix);
        default -> s = mutate(Long.toString(r.nextLong() >> r.nextInt(64), radix), ODD_CHARS);
      }
      if (r.nextInt(5) == 0) s = s.toUpperCase(Locale.ROOT);
      if (r.nextInt(6) == 0) s = (r.nextBoolean() ? "+" : "-") + s.replace("-", "");
      String t = s;
      out.println("int\t" + hex(t) + "\t" + radix + "\t" + run(() -> Integer.parseInt(t, radix)));
      out.println("long\t" + hex(t) + "\t" + radix + "\t" + run(() -> Long.parseLong(t, radix)));
    }
    for (String s : in) {
      for (int radix : radixes) {
        out.println("int\t" + hex(s) + "\t" + radix + "\t" + run(() -> Integer.parseInt(s, radix)));
        out.println("long\t" + hex(s) + "\t" + radix + "\t" + run(() -> Long.parseLong(s, radix)));
      }
    }
  }

  static String describe(BigDecimal d) {
    String tbi;
    int scale = d.scale();
    if ((scale >= -2000 && scale <= 2000) || scale >= 536870920 || scale <= -536870920) {
      tbi = run(() -> d.toBigInteger());
    } else {
      tbi = "-"; // too expensive for the JDK to compute
    }
    return d.unscaledValue() + "|" + scale + "|" + d.precision() + "|" + d.signum() + "|" + d.intValue() + "|"
        + d.longValue() + "|" + tbi;
  }

  static void bigDecimals() {
    List<String> in = new ArrayList<>(List.of(
        "", " ", "-", "+", ".", "-.", "+.", ".e5", "e", "E", "e5", "1e", "1E", "1e+", "1e-", "1e+-5", "1ee5",
        "1.5.e5", "1e5.5", "1e5e5", "--1", "+-1", "1-", "1+", " 1", "1 ", "1e٠٠٠٠٠٠٠٠٠٠٠1", "1e00000000000001",
        "1e0000000000000", "1e1234567890", "1e12345678901", "1e1234567890x", "1e000x", "1e2147483647",
        "1e2147483648", "1e2147483649", "1e-2147483647", "1e-2147483648", "1.5e-2147483647", "1.5e2147483648",
        "0e2147483648", "0e-2147483647", "0e536870920", "0e-536870920", "0.0e-2147483647", "1e536870919",
        "1e-536870919", "1e536870920", "1e-536870920", "-7e63", "-7e64", "7e63", "1e19", "1e20", "9e18",
        "12345678901234567890", "-12345678901234567890", "123456789012345678901234567890.123456789",
        "000000000000000000000000000001.000000000000000000000", "-0.00", "0.", "00", "1.", ".5", "١٢", "1١",
        "1e١", "１.２", "𝟎", "²", "4294967295.9", "18446744073709551617.5", "-1.5e10",
        "1e-64", "1e63", "1e64", "-1e64", "5e63", "3e9", "-3e9"));
    String[] bdChars = {".", "e", "E", "+", "-", "x", " ", "١", "１", "\ud835", "0", "5", "9"};
    for (int k = 0; k < 2500; k++) {
      StringBuilder sb = new StringBuilder();
      if (r.nextInt(4) == 0) sb.append(r.nextBoolean() ? "-" : "+");
      int intDigits = r.nextInt(r.nextBoolean() ? 6 : 30);
      sb.append(digits(intDigits, "0123456789"));
      if (r.nextInt(3) > 0) {
        sb.append('.');
        sb.append(digits(r.nextInt(r.nextBoolean() ? 6 : 30), "0123456789"));
      }
      if (r.nextInt(2) == 0) {
        sb.append(r.nextBoolean() ? 'e' : 'E');
        if (r.nextBoolean()) sb.append(r.nextBoolean() ? '-' : '+');
        switch (r.nextInt(4)) {
          case 0 -> sb.append(r.nextInt(30));
          case 1 -> sb.append(r.nextInt(3000));
          case 2 -> sb.append(2147483600L + r.nextInt(100));
          default -> sb.append(digits(1 + r.nextInt(14), "0123456789"));
        }
      }
      String s = sb.toString();
      in.add(r.nextInt(3) == 0 ? mutate(s, bdChars) : s);
    }
    for (String s : in) out.println("bd\t" + hex(s) + "\t" + run(() -> describe(new BigDecimal(s))));
  }

  static void doubles() {
    List<String> in = new ArrayList<>(List.of(
        "0.1f", "1.1F", "0.1d", "3.4028235e38f", "1e39f", "0x1.0p0", "0x1..1p1", "0x1.2.3p1", "-0x..p1", "0x.8.p1",
        "1..e5", "..", ". .", "1 .5", "0x1p1 ", " \u0000 1.5e3d \u001f", "Infinity ", " -Infinity", "\u0000NaN",
        "+-NaN", "NaN.", "1e+", "1e-+1", "1e2147483647", "1e-2147483648", "0x1p-2147483648", "0x1p+2147483647",
        "4.9e-324", "2.4703282292062328e-324", "2.4703282292062327e-324", "1.7976931348623158e308",
        "1.797693134862315807e308", "0x1.fffffffffffff8p1023", "0x1.fffffffffffff7ffffp1023",
        "0x0.0000000000001p-1022", "0x0.00000000000008p-1022", "0x0.000000000000080001p-1022",
        "0x0.00000000000018p-1022", "0x0.00000000000028p-1022", "0x1.00000000000008p0", "0x1.00000000000018p0",
        "0x1.000000000000080000000000000000001p0"));
    String alphabet = "0123456789..eE+-xXpPfFdDNaIn abAB";
    for (int k = 0; k < 3000; k++) {
      in.add(digits(1 + r.nextInt(9), alphabet));
    }
    for (int k = 0; k < 1500; k++) {
      double d = Double.longBitsToDouble(r.nextLong());
      if (Double.isNaN(d) || Double.isInfinite(d)) continue;
      double next = Math.nextUp(d);
      if (Double.isInfinite(next)) continue;
      // The exact midpoint between two adjacent doubles, and values just around it.
      BigDecimal mid = new BigDecimal(d).add(new BigDecimal(next)).divide(BigDecimal.valueOf(2));
      String m = mid.toString();
      in.add(m);
      BigDecimal tiny = BigDecimal.ONE.movePointLeft(mid.scale() + 5);
      in.add(mid.add(tiny).toString());
      in.add(mid.subtract(tiny).toString());
      in.add(Double.toHexString(d));
      // A hex significand with extra digits, some of them beyond double precision.
      String hx = Long.toHexString(r.nextLong()) + Long.toHexString(r.nextLong() >>> r.nextInt(64));
      int point = r.nextInt(hx.length() + 1);
      String zeros = "0".repeat(r.nextInt(4));
      in.add((r.nextBoolean() ? "-" : "") + "0x" + zeros + hx.substring(0, point) + "." + hx.substring(point) + "p"
          + (r.nextInt(2400) - 1200) + (r.nextInt(5) == 0 ? "d" : ""));
    }
    for (int k = 0; k < 600; k++) {
      // Long digit strings with exponents that push the value to the edges.
      String sig = digits(1 + r.nextInt(r.nextBoolean() ? 40 : 900), "0123456789");
      int point = r.nextInt(sig.length() + 1);
      String lead = "0".repeat(r.nextInt(3) == 0 ? r.nextInt(400) : 0);
      int exp = switch (r.nextInt(4)) {
        case 0 -> r.nextInt(700) - 350;
        case 1 -> -320 - r.nextInt(30) - point;
        case 2 -> 300 + r.nextInt(20) - point;
        default -> r.nextInt(40000) - 20000;
      };
      in.add(lead + sig.substring(0, point) + "." + sig.substring(point) + "e" + exp);
    }
    for (String s : in) {
      out.println("pd\t" + hex(s) + "\t" + run(() -> Long.toHexString(Double.doubleToRawLongBits(Double.parseDouble(s)))));
    }
  }

  static void bigIntegerRadix() {
    String[] chars = {"-", "+", "x", " ", "_", "١", "１", "Ａ", "\ud835", "z", "Z", "0", "1", "9"};
    List<String[]> in = new ArrayList<>();
    for (String s : new String[] {"x-1", "1-x", "-", "+", "--", "-x", "0x10", "000", "-000", "0٠０x", "٠١",
        "+١", "1𝟎", "𝟎1"}) {
      for (int radix : new int[] {2, 10, 16, 36, 1, 37}) in.add(new String[] {s, "" + radix});
    }
    for (int k = 0; k < 2000; k++) {
      int radix = 2 + r.nextInt(35);
      BigInteger x = new BigInteger(1 + r.nextInt(400), r);
      if (r.nextBoolean()) x = x.negate();
      String s = x.toString(radix);
      if (r.nextInt(4) == 0) s = s.toUpperCase(Locale.ROOT);
      if (r.nextInt(4) == 0) s = "0".repeat(1 + r.nextInt(12)) + s;
      if (r.nextInt(2) == 0) s = mutate(s, chars);
      in.add(new String[] {s, "" + radix});
    }
    for (String[] b : in) {
      int radix = Integer.parseInt(b[1]);
      out.println("bir\t" + hex(b[0]) + "\t" + radix + "\t" + run(() -> new BigInteger(b[0], radix)));
    }
  }

  static void pow() {
    List<String[]> in = new ArrayList<>();
    String[][] fixed = {{"0", "0"}, {"0", "5"}, {"0", "-1"}, {"5", "-1"}, {"1", "2147483647"}, {"-1", "2147483647"},
        {"-1", "2147483646"}, {"2", "0"}, {"-2", "3"}, {"5", "800000000"}, {"5", "715827894"}, {"3", "1073741840"},
        {"10", "536870920"}, {"12", "536870920"}, {"255", "268435460"}, {"1099511627777", "52377651"},
        {"-5", "715827894"}, {"4", "1073741824"}, {"-4", "1073741824"}, {"6", "900000000"}, {"7", "2147483647"},
        {"2", "2147483647"}, {"8", "715827883"}, {"1024", "214748365"}, {"-1024", "214748365"}};
    in.addAll(Arrays.asList(fixed));
    for (int k = 0; k < 400; k++) {
      BigInteger b = new BigInteger(r.nextInt(70), r);
      if (r.nextInt(4) == 0) b = b.shiftLeft(r.nextInt(40));
      if (r.nextInt(6) == 0) b = BigInteger.ONE.shiftLeft(r.nextInt(70));
      if (r.nextBoolean()) b = b.negate();
      in.add(new String[] {b.toString(), "" + r.nextInt(b.bitLength() > 40 ? 30 : 200)});
    }
    for (String[] p : in) {
      BigInteger b = new BigInteger(p[0]);
      int e = Integer.parseInt(p[1]);
      out.println("pow\t" + p[0] + "\t" + e + "\t" + run(() -> b.pow(e)));
    }
  }

  static void formats() {
    List<Double> ds = new ArrayList<>(Arrays.asList(0.0, -0.0, 0.5, 1.5, 2.5, -0.5, -2.5, 0.05, 0.15, 0.25, 0.35,
        0.45, 0.95, 9.95, 99.95, 999.9995, 0.125, 0.375, 1.005, 2.675, 1e-5, 5e-5, 4.9e-324, Double.MIN_NORMAL,
        Double.MAX_VALUE, 1e22, 1e23, 2e23, 8.41e21, 9.999999999999999e22, 123456789.125, 0.1 + 0.2,
        Double.NaN, Double.POSITIVE_INFINITY, Double.NEGATIVE_INFINITY, -0.04, -1e-300, 1e7, 1e-3, 9.5e-4));
    for (int k = 0; k < 2500; k++) {
      switch (k % 6) {
        case 0 -> ds.add(Double.longBitsToDouble(r.nextLong()));
        case 1 -> ds.add(r.nextDouble() * Math.pow(10, r.nextInt(30) - 15));
        case 2 -> ds.add((r.nextInt(2000001) - 1000000) / Math.pow(10, r.nextInt(9)));
        case 3 -> ds.add((r.nextInt(20000) * 2 + 1) / Math.pow(10, 1 + r.nextInt(8)) / 2 * 10);
        case 4 -> ds.add((double) r.nextLong() / (1L << r.nextInt(63)));
        default -> ds.add((double) r.nextInt(1000) / (1 + r.nextInt(1000)) * 100);
      }
    }
    int[] precs = {0, 1, 2, 3, 4, 6, 9, 12, 17, 20, 25, 40, 330};
    for (double d : ds) {
      if (Double.isNaN(d) && Double.doubleToRawLongBits(d) != Double.doubleToRawLongBits(Double.NaN)) continue;
      int prec = r.nextInt(8) == 0 ? precs[r.nextInt(precs.length)] : r.nextInt(8);
      String bits = Long.toHexString(Double.doubleToRawLongBits(d));
      out.println("fmt\t" + bits + "\t" + prec + "\t" + hex(String.format("%." + prec + "f", d)));
    }
  }

  static void binarySearch() {
    for (int k = 0; k < 1500; k++) {
      int n = r.nextInt(k < 100 ? 6 : 70);
      int[] a = new int[n];
      boolean extreme = r.nextInt(10) == 0;
      int range = 1 + r.nextInt(r.nextBoolean() ? 6 : 200);
      for (int i = 0; i < n; i++) a[i] = extreme ? (r.nextBoolean() ? Integer.MAX_VALUE - r.nextInt(3) : Integer.MIN_VALUE + r.nextInt(3)) : r.nextInt(range) - range / 2;
      Arrays.sort(a);
      List<Integer> keys = new ArrayList<>();
      for (int key = -range / 2 - 2; key <= range - range / 2 + 1; key++) keys.add(key);
      keys.addAll(List.of(Integer.MIN_VALUE, Integer.MIN_VALUE + 1, Integer.MAX_VALUE, Integer.MAX_VALUE - 1));
      StringBuilder as = new StringBuilder(), ks = new StringBuilder(), rs = new StringBuilder();
      for (int v : a) as.append(as.length() == 0 ? "" : ",").append(v);
      for (int key : keys) {
        ks.append(ks.length() == 0 ? "" : ",").append(key);
        rs.append(rs.length() == 0 ? "" : ",").append(Arrays.binarySearch(a, key));
      }
      out.println("bs\t" + (n == 0 ? "-" : as) + "\t" + ks + "\t" + rs);
    }
  }

  static void nextIntBound() {
    int[] special = {1, 2, 3, 7, 1 << 30, (1 << 30) + 1, (1 << 30) + 7, Integer.MAX_VALUE, Integer.MAX_VALUE - 1,
        1431655765, 1431655766, 0x60000000, 0, -1, Integer.MIN_VALUE};
    for (int k = 0; k < 300; k++) {
      long seed = r.nextLong();
      Random x = new Random(seed);
      StringBuilder bs = new StringBuilder(), rs = new StringBuilder();
      for (int i = 0; i < 40; i++) {
        int bound = switch (r.nextInt(4)) {
          case 0 -> special[r.nextInt(special.length - 3)];
          case 1 -> 1 << r.nextInt(31);
          case 2 -> 1 + r.nextInt(1000);
          default -> 1 + r.nextInt(Integer.MAX_VALUE);
        };
        bs.append(i == 0 ? "" : ",").append(bound);
        rs.append(i == 0 ? "" : ",").append(x.nextInt(bound));
      }
      out.println("rnd\t" + seed + "\t" + bs + "\t" + rs);
    }
    for (int bound : new int[] {0, -1, Integer.MIN_VALUE}) {
      out.println("rnd\t1\t" + bound + "\t" + run(() -> new Random(1).nextInt(bound)));
    }
  }

  public static void main(String[] args) throws IOException {
    long seed = args.length > 1 ? Long.parseLong(args[1]) : 20261016L;
    r = new Random(seed);
    try (PrintWriter pw = new PrintWriter(new OutputStreamWriter(
        new GZIPOutputStream(new FileOutputStream(args[0])), StandardCharsets.UTF_8))) {
      out = pw;
      out.println("seed\t" + seed);
      integers();
      bigDecimals();
      doubles();
      bigIntegerRadix();
      pow();
      formats();
      binarySearch();
      nextIntBound();
    }
  }
}

// Oracle tool: writes the Java reference vectors for the jutil helpers added
// for the jscomp port (random.tsv, murmur3.tsv, utf8.tsv, caseformat.tsv,
// charset.tsv, bitset.tsv, bimap.tsv) into internal/jutil/testdata.
//   java -cp .cache/closure-compiler.jar:.cache/oracle-build JutilExtraVectors internal/jutil/testdata
import com.google.javascript.jscomp.jarjar.com.google.common.base.CaseFormat;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.BiMap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.HashBiMap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableBiMap;
import com.google.javascript.jscomp.jarjar.com.google.common.hash.Hasher;
import com.google.javascript.jscomp.jarjar.com.google.common.hash.Hashing;
import java.io.PrintStream;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.*;

public class JutilExtraVectors {
  static String hex(byte[] b) {
    StringBuilder sb = new StringBuilder("x");
    for (byte x : b) sb.append(String.format("%02x", x & 0xff));
    return sb.toString();
  }

  static String units(String s) {
    StringBuilder sb = new StringBuilder("u");
    for (int i = 0; i < s.length(); i++) sb.append(String.format("%04x", (int) s.charAt(i)));
    return sb.toString();
  }

  static String randomString(Random r, int len, String alphabet) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < len; i++) sb.append(alphabet.charAt(r.nextInt(alphabet.length())));
    return sb.toString();
  }

  static String randomUtf16(Random r, int len) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < len; i++) {
      switch (r.nextInt(6)) {
        case 0 -> sb.append((char) r.nextInt(0x80));
        case 1 -> sb.append((char) (0x80 + r.nextInt(0x780)));
        case 2 -> sb.append((char) (0xD800 + r.nextInt(0x800))); // any surrogate
        case 3 -> sb.appendCodePoint(0x10000 + r.nextInt(0x100000));
        default -> sb.append((char) r.nextInt(0x10000));
      }
    }
    return sb.toString();
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(4242);

    // java.util.Random: per seed, a fixed sequence of calls.
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("random.tsv")))) {
      long[] seeds = {0, 1, -1, 42, 123456789L, Long.MIN_VALUE, Long.MAX_VALUE, 0x5DEECE66DL};
      int[] bounds = {1, 2, 3, 7, 16, 52, 64, 100, 1000, 1 << 30, Integer.MAX_VALUE, 0x40000001};
      for (long seed : seeds) {
        Random x = new Random(seed);
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < 5; i++) sb.append(x.nextInt()).append(',');
        for (int b : bounds) sb.append(x.nextInt(b)).append(',');
        for (int i = 0; i < 3; i++) sb.append(x.nextLong()).append(',');
        sb.append(x.nextBoolean()).append(',');
        sb.append(Double.doubleToRawLongBits(x.nextDouble())).append(',');
        sb.append(Float.floatToRawIntBits(x.nextFloat()));
        List<Integer> list = new ArrayList<>();
        for (int i = 0; i < 64; i++) list.add(i);
        Collections.shuffle(list, x);
        ps.println(seed + "\t" + sb + "\t" + list.toString().replace(" ", ""));
      }
    }

    // Guava murmur3_128: raw bytes, and the putInt + putUnencodedChars
    // sequence that RandomNameGenerator uses.
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("murmur3.tsv")))) {
      for (int len = 0; len <= 70; len++) {
        for (int k = 0; k < 3; k++) {
          byte[] b = new byte[len];
          r.nextBytes(b);
          ps.println("B\t" + hex(b) + "\t" + Hashing.murmur3_128().hashBytes(b));
        }
      }
      for (int k = 0; k < 300; k++) {
        int n = r.nextInt(8);
        String s = k % 3 == 0 ? randomUtf16(r, n) : randomString(r, n, "abcXYZ_$09");
        int len = r.nextInt(10);
        Hasher h = Hashing.murmur3_128().newHasher();
        h.putInt(len);
        h.putUnencodedChars(s);
        ps.println("C\t" + len + "," + units(s) + "\t" + h.hash().asInt());
      }
    }

    // String.getBytes(UTF_8) and new String(bytes, UTF_8).
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("utf8.tsv")))) {
      for (int k = 0; k < 2000; k++) {
        String s = randomUtf16(r, r.nextInt(8));
        ps.println("E\t" + units(s) + "\t" + hex(s.getBytes(StandardCharsets.UTF_8)));
      }
      byte[] interesting = {
        (byte) 0x41, (byte) 0x80, (byte) 0xbf, (byte) 0xc0, (byte) 0xc1, (byte) 0xc2, (byte) 0xdf,
        (byte) 0xe0, (byte) 0xe1, (byte) 0xed, (byte) 0xee, (byte) 0xef, (byte) 0xf0, (byte) 0xf1,
        (byte) 0xf4, (byte) 0xf5, (byte) 0xff, (byte) 0x9f, (byte) 0xa0, (byte) 0x8f, (byte) 0x90
      };
      for (int k = 0; k < 6000; k++) {
        byte[] b = new byte[r.nextInt(7)];
        for (int i = 0; i < b.length; i++) {
          b[i] = r.nextInt(4) == 0 ? (byte) r.nextInt(256) : interesting[r.nextInt(interesting.length)];
        }
        ps.println("D\t" + hex(b) + "\t" + units(new String(b, StandardCharsets.UTF_8)));
      }
    }

    // CaseFormat.LOWER_CAMEL -> UPPER_UNDERSCORE.
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("caseformat.tsv")))) {
      String[] fixed = {"", "a", "A", "AB", "aB", "abc", "Abc", "fooBar", "fooBarBaz", "a4Letter",
          "startSpan_1_23", "ABC", "aBC", "x_Y", "__", "_a", "a_", "éA", "éa"};
      for (String s : fixed) {
        ps.println(units(s) + "\t" + units(CaseFormat.LOWER_CAMEL.to(CaseFormat.UPPER_UNDERSCORE, s)));
      }
      for (int k = 0; k < 1000; k++) {
        String s = randomString(r, r.nextInt(9), "aAbBzZ_09éÉ");
        ps.println(units(s) + "\t" + units(CaseFormat.LOWER_CAMEL.to(CaseFormat.UPPER_UNDERSCORE, s)));
      }
    }

    // Standard charsets: aliases and canEncode(char) ranges.
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("charset.tsv")))) {
      Charset[] cs = {StandardCharsets.US_ASCII, StandardCharsets.ISO_8859_1, StandardCharsets.UTF_8,
          StandardCharsets.UTF_16, StandardCharsets.UTF_16BE, StandardCharsets.UTF_16LE};
      for (Charset c : cs) {
        List<String> aliases = new ArrayList<>(c.aliases());
        Collections.sort(aliases);
        StringBuilder ranges = new StringBuilder();
        var enc = c.newEncoder();
        int start = -1;
        for (int ch = 0; ch <= 0x10000; ch++) {
          boolean ok = ch < 0x10000 && enc.canEncode((char) ch);
          if (ok && start < 0) start = ch;
          if (!ok && start >= 0) {
            ranges.append(Integer.toHexString(start)).append('-').append(Integer.toHexString(ch - 1)).append(',');
            start = -1;
          }
        }
        ps.println(c.name() + "\t" + String.join(",", aliases) + "\t" + ranges);
      }
    }

    // BitSet: random set/clear sequences.
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("bitset.tsv")))) {
      for (int k = 0; k < 300; k++) {
        BitSet bs = new BitSet();
        StringBuilder ops = new StringBuilder();
        int n = r.nextInt(12);
        int range = 1 + r.nextInt(k % 2 == 0 ? 70 : 300);
        for (int i = 0; i < n; i++) {
          int bit = r.nextInt(range);
          if (r.nextInt(4) == 0) {
            bs.clear(bit);
            ops.append('c').append(bit).append(',');
          } else {
            bs.set(bit);
            ops.append('s').append(bit).append(',');
          }
        }
        int from = r.nextInt(range + 1);
        ps.println(ops + "\t" + bs + "\t" + bs.length() + "\t" + bs.cardinality() + "\t"
            + hex(bs.toByteArray()) + "\t" + bs.hashCode() + "\t" + bs.nextSetBit(from) + "\t"
            + bs.nextClearBit(from) + "\t" + from);
      }
    }

    // HashBiMap put/forcePut sequences, and ImmutableBiMap.copyOf conflicts.
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("bimap.tsv")))) {
      for (int k = 0; k < 300; k++) {
        BiMap<String, String> m = HashBiMap.create();
        StringBuilder ops = new StringBuilder();
        StringBuilder results = new StringBuilder();
        int n = r.nextInt(10);
        for (int i = 0; i < n; i++) {
          String key = randomString(r, 1, "abcde");
          String val = randomString(r, 1, "vwxyz");
          int op = r.nextInt(5);
          String res;
          try {
            if (op == 0) {
              ops.append('F');
              res = String.valueOf(m.forcePut(key, val));
            } else if (op == 1) {
              ops.append('R');
              res = String.valueOf(m.remove(key));
            } else {
              ops.append('P');
              res = String.valueOf(m.put(key, val));
            }
          } catch (IllegalArgumentException e) {
            res = "!" + e.getMessage();
          }
          ops.append(key).append(val).append(',');
          results.append(res).append(';');
        }
        ps.println("H\t" + ops + "\t" + results + "\t" + m + "\t" + m.inverse());
      }
      for (int k = 0; k < 100; k++) {
        LinkedHashMap<String, String> src = new LinkedHashMap<>();
        int n = r.nextInt(6);
        StringBuilder ops = new StringBuilder();
        for (int i = 0; i < n; i++) {
          String key = randomString(r, 1, "abcde");
          String val = randomString(r, 1, "vwxy");
          src.put(key, val);
          ops.append(key).append(val).append(',');
        }
        String res;
        try {
          ImmutableBiMap<String, String> m = ImmutableBiMap.copyOf(src);
          res = m + "\t" + m.inverse();
        } catch (IllegalArgumentException e) {
          res = "!" + e.getMessage() + "\t-";
        }
        ps.println("I\t" + ops + "\t" + res);
      }
    }
  }
}

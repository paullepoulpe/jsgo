// Oracle tool: writes internal/jutil/testdata/string_case.tsv, the JDK's
// Character.toLowerCase(int) and String.toLowerCase(Locale.ROOT), for
// internal/jutil/string_case_test.go.
//   java -cp .cache/oracle-build JutilStringCaseVectors internal/jutil/testdata [seed [count]]
//
// Rows (tab-separated):
//   cp     <cp>  <lower>    Character.toLowerCase(int) for every int in
//                           0..0x10FFFF whose mapping is not the identity,
//                           and for a few ints that are not code points;
//   str1   <cp>  <units>    toLowerCase(Locale.ROOT) of the string holding
//                           just that code point (a lone surrogate for
//                           U+D800..U+DFFF), for every code point where it
//                           is not the string of Character.toLowerCase(cp);
//   sigma  <ctx> <runs>     for every code point c, whether the U+03A3 of a
//                           context string becomes U+03C2 (final sigma):
//                           pre  = "A" c "\u03A3",  pre0  = c "\u03A3",
//                           post = "A\u03A3" c "a", post0 = "A\u03A3" c,
//                           as run lengths of false/true bits over
//                           c = 0..0x10FFFF, starting with false;
//   mixed  <units> <units>  random strings (seeded) drawn from pools of
//                           letters, marks, format characters, punctuation,
//                           supplementary characters and lone surrogates,
//                           and their lowercase;
//   text   <units> <units>  random Greek/Latin text (seeded): words of
//                           capital and small letters (sigma often last),
//                           with spaces, punctuation, digits, hyphens,
//                           apostrophes and accents.
// Code points are hex, strings are "u" followed by 4 hex digits per UTF-16
// unit.
import java.io.PrintStream;
import java.nio.file.*;
import java.util.*;

public class JutilStringCaseVectors {
  static String units(String s) {
    StringBuilder sb = new StringBuilder("u");
    for (int i = 0; i < s.length(); i++) sb.append(String.format("%04x", (int) s.charAt(i)));
    return sb.toString();
  }

  static String str(int cp) {
    return new String(Character.toChars(cp));
  }

  static final int[][] POOLS = {
    // Greek capital sigma is drawn from its own slot (below).
    {'A', 'Z'}, {'a', 'z'}, {'0', '9'}, {' ', ' '}, {'\'', '\''}, {'.', '.'}, {':', ':'}, {'^', '^'}, {'`', '`'},
    {'-', '-'}, {',', ','}, {'!', '!'}, {'_', '_'},
    {0x00C0, 0x00FF}, {0x0130, 0x0131}, {0x0391, 0x03A9}, {0x03AC, 0x03CE}, {0x1F08, 0x1F0F},
    {0x0300, 0x036F}, {0x0483, 0x0489}, {0x0903, 0x0903}, {0x093E, 0x0940}, {0x20DD, 0x20E0},
    {0x00AD, 0x00AD}, {0x200B, 0x200F}, {0x2060, 0x2064}, {0xFEFF, 0xFEFF}, {0x061C, 0x061C}, {0xE0001, 0xE0001},
    {0x02B0, 0x02FF}, {0x3005, 0x3005}, {0x00A8, 0x00A8}, {0x00B4, 0x00B4}, {0x1FBD, 0x1FC1},
    {0x00B7, 0x00B7}, {0x2018, 0x2019}, {0x2024, 0x2024}, {0x2027, 0x2027}, {0xFF07, 0xFF07}, {0xFF0E, 0xFF0E}, {0x0387, 0x0387},
    {0x24B6, 0x24E9}, {0x2160, 0x217F}, {0x10400, 0x1044F}, {0x1D400, 0x1D433}, {0x1F600, 0x1F64F},
    {0x05D0, 0x05EA}, {0x0627, 0x064A}, {0x4E00, 0x4E10}, {0x3041, 0x3096},
    {0xD800, 0xDFFF},
  };

  static String randomString(Random r) {
    StringBuilder sb = new StringBuilder();
    int len = 1 + r.nextInt(10);
    for (int i = 0; i < len; i++) {
      int k = r.nextInt(POOLS.length + 8);
      int cp;
      if (k < POOLS.length) {
        int[] p = POOLS[k];
        cp = p[0] + r.nextInt(p[1] - p[0] + 1);
      } else if (k < POOLS.length + 5) {
        cp = 0x03A3;
      } else if (k == POOLS.length + 5) {
        cp = r.nextInt(0x10000);
      } else {
        cp = r.nextInt(0x110000);
      }
      sb.appendCodePoint(cp);
    }
    return sb.toString();
  }

  static final String[] SEPARATORS = {
    " ", " ", " ", ", ", ". ", "; ", "! ", "? ", "\n", "-", "_", ":", "'", "\u2019", "\"", "(", ")", "1", "42", "/", " - ", ".",
  };
  static final String GREEK_CAPITALS = "\u0391\u0392\u0393\u0394\u0395\u0396\u0397\u0398\u0399\u039A\u039B\u039C\u039D\u039E\u039F\u03A0\u03A1\u03A3\u03A4\u03A5\u03A6\u03A7\u03A8\u03A9\u0386\u0388\u0389\u038A\u038C\u038E\u038F\u03AA\u03AB";
  static final String GREEK_SMALL = "\u03B1\u03B2\u03B3\u03B4\u03B5\u03B6\u03B7\u03B8\u03B9\u03BA\u03BB\u03BC\u03BD\u03BE\u03BF\u03C0\u03C1\u03C3\u03C4\u03C5\u03C6\u03C7\u03C8\u03C9\u03AC\u03AD\u03AE\u03AF\u03CC\u03CD\u03CE\u03C2";

  static String randomText(Random r) {
    StringBuilder sb = new StringBuilder();
    int words = 1 + r.nextInt(5);
    for (int w = 0; w < words; w++) {
      if (w > 0 || r.nextInt(8) == 0) sb.append(SEPARATORS[r.nextInt(SEPARATORS.length)]);
      int len = 1 + r.nextInt(7);
      int style = r.nextInt(6); // 0-2 capitals, 3 small, 4 capitalised, 5 Latin
      for (int i = 0; i < len; i++) {
        char c;
        if (i == len - 1 && r.nextInt(3) == 0) {
          c = '\u03A3';
        } else if (style == 5) {
          c = (char) ((r.nextBoolean() ? 'A' : 'a') + r.nextInt(26));
        } else if (style == 3 || (style == 4 && i > 0)) {
          c = GREEK_SMALL.charAt(r.nextInt(GREEK_SMALL.length()));
        } else {
          c = GREEK_CAPITALS.charAt(r.nextInt(GREEK_CAPITALS.length()));
        }
        sb.append(c);
        if (r.nextInt(12) == 0) sb.append(r.nextBoolean() ? '\u0301' : '\u0308');
      }
    }
    if (r.nextInt(3) == 0) sb.append(SEPARATORS[r.nextInt(SEPARATORS.length)]);
    return sb.toString();
  }

  static void runs(PrintStream ps, String ctx, boolean[] bits) {
    StringBuilder sb = new StringBuilder();
    boolean cur = false;
    int n = 0;
    for (boolean b : bits) {
      if (b != cur) {
        sb.append(n).append(',');
        cur = b;
        n = 0;
      }
      n++;
    }
    sb.append(n);
    ps.println("sigma\t" + ctx + "\t" + sb);
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    long seed = args.length > 1 ? Long.parseLong(args[1]) : 20261016L;
    int count = args.length > 2 ? Integer.parseInt(args[2]) : 2000;
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("string_case.tsv")), false, "UTF-8")) {
      for (int cp = 0; cp <= 0x10FFFF; cp++) {
        int l = Character.toLowerCase(cp);
        if (l != cp) ps.printf("cp\t%X\t%X%n", cp, l);
      }
      for (int cp : new int[] {-1, -0x41, Integer.MIN_VALUE, Integer.MAX_VALUE, 0x110000, 0x110041, 0x1F0000}) {
        ps.printf("cp\t%X\t%X%n", cp, Character.toLowerCase(cp));
      }

      for (int cp = 0; cp <= 0x10FFFF; cp++) {
        String got = str(cp).toLowerCase(Locale.ROOT);
        if (!got.equals(str(Character.toLowerCase(cp)))) ps.printf("str1\t%X\t%s%n", cp, units(got));
      }

      int n = 0x110000;
      boolean[] pre = new boolean[n], pre0 = new boolean[n], post = new boolean[n], post0 = new boolean[n];
      for (int cp = 0; cp < n; cp++) {
        String c = str(cp);
        String a = ("A" + c + "\u03A3").toLowerCase(Locale.ROOT);
        pre[cp] = a.charAt(a.length() - 1) == '\u03C2';
        String b = (c + "\u03A3").toLowerCase(Locale.ROOT);
        pre0[cp] = b.charAt(b.length() - 1) == '\u03C2';
        post[cp] = ("A\u03A3" + c + "a").toLowerCase(Locale.ROOT).charAt(1) == '\u03C2';
        post0[cp] = ("A\u03A3" + c).toLowerCase(Locale.ROOT).charAt(1) == '\u03C2';
      }
      runs(ps, "pre", pre);
      runs(ps, "pre0", pre0);
      runs(ps, "post", post);
      runs(ps, "post0", post0);

      Random r = new Random(seed);
      for (int i = 0; i < count; i++) {
        String s = randomString(r);
        ps.println("mixed\t" + units(s) + "\t" + units(s.toLowerCase(Locale.ROOT)));
      }
      for (int i = 0; i < count; i++) {
        String s = randomText(r);
        ps.println("text\t" + units(s) + "\t" + units(s.toLowerCase(Locale.ROOT)));
      }
    }
  }
}

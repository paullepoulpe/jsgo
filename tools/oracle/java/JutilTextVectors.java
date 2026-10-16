// Oracle tool: writes black-box reference vectors for the clean-room text
// helpers of internal/jutil (java.net.URI, java.nio.file.Path,
// java.text.MessageFormat, Properties.load, String.split, regionMatches,
// Pattern.quote, UTF-16 decoding) into <dir>/text_*.tsv, from random inputs
// drawn with the given seed. Only the public API of the JDK is called.
//   javac -d .cache/oracle-build tools/oracle/java/JutilTextVectors.java
//   LC_ALL=C.UTF-8 java -cp .cache/oracle-build JutilTextVectors internal/jutil/testdata 20261016
// (a UTF-8 locale, so that sun.jnu.encoding, which java.nio.file.Path uses,
// is UTF-8; the generator sets the default locale to en_US itself).
// The Go side is internal/jutil/text_vectors_test.go; JUTIL_TEXT_VECTORS=<dir>
// points it at vectors written with another seed.
import java.io.PrintStream;
import java.io.StringReader;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import java.nio.charset.CodingErrorAction;
import java.nio.charset.MalformedInputException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.text.MessageFormat;
import java.util.*;
import java.util.regex.Pattern;

public class JutilTextVectors {
  static Random r;

  static String esc(String s) {
    if (s == null) return "\\N";
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      switch (c) {
        case '\\': sb.append("\\\\"); break;
        case '\t': sb.append("\\t"); break;
        case '\n': sb.append("\\n"); break;
        case '\r': sb.append("\\r"); break;
        default:
          if (c < 0x20 || (c >= 0xD800 && c <= 0xDFFF) || c > 0x7e) {
            sb.append(String.format("\\u%04x", (int) c));
          } else {
            sb.append(c);
          }
      }
    }
    return sb.toString();
  }

  static String exc(Throwable t) {
    return "EXC:" + t.getClass().getName() + ":" + esc(t.getMessage());
  }

  interface Op { String run() throws Throwable; }

  static String guard(Op op) {
    try {
      return esc(op.run());
    } catch (Throwable t) {
      return exc(t);
    }
  }

  static <T> T pick(T[] a) { return a[r.nextInt(a.length)]; }

  static String pickChars(String alphabet, int max) {
    int n = r.nextInt(max + 1);
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < n; i++) sb.append(alphabet.charAt(r.nextInt(alphabet.length())));
    return sb.toString();
  }

  // Inserts, deletes or replaces a few characters.
  static String mutate(String s, String alphabet) {
    StringBuilder sb = new StringBuilder(s);
    int k = r.nextInt(3);
    for (int i = 0; i < k; i++) {
      int op = r.nextInt(3);
      int pos = sb.length() == 0 ? 0 : r.nextInt(sb.length() + 1);
      char c = alphabet.charAt(r.nextInt(alphabet.length()));
      if (op == 0 || sb.length() == 0) sb.insert(pos, c);
      else if (op == 1 && pos < sb.length()) sb.deleteCharAt(pos);
      else if (pos < sb.length()) sb.setCharAt(pos, c);
    }
    return sb.toString();
  }

  // ---------------------------------------------------------------- URI
  static final String URI_NOISE = "aZ09:/?#[]@!$&'()*+,;=%-._~ \"<>{}|\\^`é \u0085 　😀\u0000\u007fAfF.:%%";

  static String hex4() {
    String h = "0123456789abcdefABCDEF";
    int n = 1 + r.nextInt(r.nextInt(8) == 0 ? 6 : 4);
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < n; i++) sb.append(h.charAt(r.nextInt(h.length())));
    return sb.toString();
  }

  static String ipv4() {
    StringBuilder sb = new StringBuilder();
    int parts = r.nextInt(6) == 0 ? 2 + r.nextInt(4) : 4;
    for (int i = 0; i < parts; i++) {
      if (i > 0) sb.append('.');
      int k = r.nextInt(10);
      if (k == 0) sb.append(256 + r.nextInt(1000));
      else if (k == 1) sb.append("0").append(r.nextInt(300));
      else if (k == 2) sb.append("");
      else sb.append(r.nextInt(256));
    }
    return sb.toString();
  }

  static String ipv6() {
    StringBuilder sb = new StringBuilder();
    int groups = r.nextInt(10);
    int compressAt = r.nextInt(3) == 0 ? -1 : r.nextInt(groups + 1);
    for (int i = 0; i < groups; i++) {
      if (i == compressAt) sb.append("::");
      else if (i > 0) sb.append(':');
      sb.append(hex4());
    }
    if (compressAt == groups) sb.append("::");
    if (r.nextInt(4) == 0) {
      if (sb.length() > 0 && sb.charAt(sb.length() - 1) != ':') sb.append(':');
      sb.append(ipv4());
    }
    if (r.nextInt(5) == 0) sb.append('%').append(pickChars("abz09_.-%:", 4));
    String s = sb.toString();
    if (r.nextInt(3) == 0) s = mutate(s, "0123456789abcdefg:.%x");
    return s;
  }

  static String hostname() {
    StringBuilder sb = new StringBuilder();
    int labels = 1 + r.nextInt(4);
    for (int i = 0; i < labels; i++) {
      if (i > 0) sb.append('.');
      sb.append(pickChars("abcxyz019-", 1 + r.nextInt(5)));
    }
    if (r.nextInt(5) == 0) sb.append('.');
    return sb.toString();
  }

  static String authority() {
    StringBuilder sb = new StringBuilder();
    if (r.nextInt(4) == 0) sb.append(pickChars("uU:;&=+$,%41!~*'()é", 5)).append('@');
    switch (r.nextInt(4)) {
      case 0: sb.append('[').append(ipv6()).append(r.nextInt(8) == 0 ? "" : "]"); break;
      case 1: sb.append(ipv4()); break;
      default: sb.append(hostname());
    }
    if (r.nextInt(3) == 0) {
      sb.append(':');
      int k = r.nextInt(6);
      if (k == 0) sb.append("99999999999");
      else if (k == 1) sb.append(pickChars("0123456789x-+", 4));
      else if (k == 2) sb.append("2147483647");
      else if (k == 3) sb.append("2147483648");
      else sb.append(r.nextInt(70000));
    }
    String s = sb.toString();
    if (r.nextInt(4) == 0) s = mutate(s, URI_NOISE);
    return s;
  }

  static String pickTokens(String[] tokens, int max) {
    int n = r.nextInt(max + 1);
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < n; i++) sb.append(pick(tokens));
    return sb.toString();
  }

  static final String[] SEGMENT_TOKENS = {"a", "b", "Z", "0", "9", "-", "_", ".", "!", "~", "*", "'", "(", ")", ";", ":", "@",
      "&", "=", "+", "$", ",", "%41", "%4a", "%e9", "\u00e9", "\ud83d\ude00", "ab", "..", "x"};

  static String segment() {
    return pickTokens(SEGMENT_TOKENS, 4);
  }

  static String uri() {
    StringBuilder sb = new StringBuilder();
    int shape = r.nextInt(14);
    if (shape >= 10) {
      // A server-based authority, the part that has the most rules.
      sb.append(pick(new String[] {"s://", "//", "HTTP://"}));
      sb.append(shape == 10 ? "[" + ipv6() + "]" + (r.nextBoolean() ? ":" + r.nextInt(100) : "") : authority());
      if (r.nextBoolean()) sb.append('/').append(segment());
      return sb.toString();
    }
    if (shape < 7) {
      if (r.nextInt(3) > 0) sb.append(pick(new String[] {"http", "s", "x+y-z.1", "HTTP", "file", "a", "1a", "a_b", "é", ""})).append(':');
      if (r.nextInt(4) > 0) sb.append("//").append(r.nextInt(8) == 0 ? "" : authority());
      int segs = r.nextInt(4);
      for (int i = 0; i < segs; i++) {
        if (i > 0 || r.nextInt(3) > 0) sb.append('/');
        sb.append(segment());
      }
    } else if (shape < 9) {
      sb.append(pick(new String[] {"mailto", "urn", "x", "jar", "data"})).append(':').append(pickChars("abc:/@?%41[];,=", 10));
    } else {
      sb.append(pickChars(URI_NOISE, 12));
    }
    if (r.nextInt(3) == 0) sb.append('?').append(pickTokens(new String[] {"a", "b", "=", "&", "/", "?", ":", "@", "[", "]", "%4a", "%A0", "\u00e9"}, 6));
    if (r.nextInt(3) == 0) sb.append('#').append(pickTokens(new String[] {"a", "b", "=", "&", "/", "?", ":", "@", "[", "]", "%4a", "%", "#"}, 6));
    String s = sb.toString();
    if (r.nextInt(5) == 0) s = mutate(s, URI_NOISE);
    return s;
  }

  static String describe(URI x) {
    return x.isAbsolute() + "|" + x.isOpaque() + "|" + x.hashCode() + "|" + x + "|" + x.getScheme() + "|"
        + x.getRawSchemeSpecificPart() + "|" + x.getRawAuthority() + "|" + x.getRawUserInfo() + "|" + x.getHost() + "|"
        + x.getPort() + "|" + x.getRawPath() + "|" + x.getRawQuery() + "|" + x.getRawFragment();
  }

  // Changes the case of letters (equal URIs differ only by the case of the
  // scheme, the host and escaped octets).
  static String flipCase(String s) {
    StringBuilder sb = new StringBuilder();
    for (char c : s.toCharArray()) {
      if (r.nextInt(3) == 0) sb.append(Character.isUpperCase(c) ? Character.toLowerCase(c) : Character.toUpperCase(c));
      else sb.append(c);
    }
    return sb.toString();
  }

  static void uris(PrintStream ps, int n) {
    List<String> valid = new ArrayList<>();
    for (int i = 0; i < n; i++) {
      String u = uri();
      String res = guard(() -> describe(new URI(u)));
      if (!res.startsWith("EXC:")) valid.add(u);
      ps.println("p\t" + esc(u) + "\t" + res);
    }
    for (int i = 0; i < n / 4 && !valid.isEmpty(); i++) {
      String a = valid.get(r.nextInt(valid.size()));
      String b = r.nextBoolean() ? flipCase(a) : valid.get(r.nextInt(valid.size()));
      ps.println("eq\t" + esc(a) + "\t" + esc(b) + "\t" + guard(() -> {
        URI x = new URI(a), y = new URI(b);
        return x.equals(y) + "|" + (x.hashCode() == y.hashCode());
      }));
    }
  }

  // ---------------------------------------------------------------- Path
  static String path() {
    String[] parts = {"a", "b", "c.js", ".", "..", "", "é", "ab", "a\u0000", "\ud800", "x/y", "/"};
    StringBuilder sb = new StringBuilder();
    if (r.nextInt(3) == 0) sb.append('/');
    int n = r.nextInt(6);
    for (int i = 0; i < n; i++) {
      if (i > 0) sb.append(r.nextInt(6) == 0 ? "//" : "/");
      String p = pick(parts);
      if (p.contains("\u0000") || p.contains("\ud800")) {
        if (r.nextInt(8) != 0) p = "q";
      }
      sb.append(p);
    }
    if (r.nextInt(6) == 0) sb.append('/');
    return sb.toString();
  }

  static String describe(Path x) {
    StringBuilder names = new StringBuilder();
    for (int i = 0; i < x.getNameCount(); i++) names.append(i == 0 ? "" : ",").append(x.getName(i));
    return x + "|" + x.getParent() + "|" + x.getFileName() + "|" + x.getRoot() + "|" + x.getNameCount() + "|" + names + "|"
        + x.normalize() + "|" + x.hashCode() + "|" + x.isAbsolute();
  }

  static void paths(PrintStream ps, int n) {
    for (int i = 0; i < n; i++) {
      String p = path(), q = path();
      ps.println("of\t" + esc(p) + "\t" + guard(() -> describe(Path.of(p))));
      String more = r.nextBoolean() ? "" : path();
      ps.println("of2\t" + esc(p) + "\t" + esc(q) + "\t" + esc(more) + "\t" + guard(() -> describe(Path.of(p, q, more))));
      ps.println("rel\t" + esc(p) + "\t" + esc(q) + "\t" + guard(() -> Path.of(p).relativize(Path.of(q)).toString()));
      ps.println("res\t" + esc(p) + "\t" + esc(q) + "\t" + guard(() -> Path.of(p).resolve(q) + "|" + Path.of(p).resolveSibling(q) + "|" + Path.of(p).resolve(Path.of(q))));
      ps.println("cmp\t" + esc(p) + "\t" + esc(q) + "\t" + guard(() -> Path.of(p).startsWith(Path.of(q)) + "|" + Path.of(p).endsWith(Path.of(q)) + "|"
          + Path.of(p).startsWith(q) + "|" + Path.of(p).endsWith(q) + "|" + Path.of(p).compareTo(Path.of(q)) + "|" + Path.of(p).equals(Path.of(q))));
      int b = r.nextInt(5) - 1, e = r.nextInt(6) - 1;
      ps.println("sub\t" + esc(p) + "\t" + b + "\t" + e + "\t" + guard(() -> Path.of(p).subpath(b, e).toString()) + "\t" + guard(() -> Path.of(p).getName(b).toString()));
    }
  }


  // ---------------------------------------------------------------- MessageFormat
  // Argument sets; the Go test has the same ones (by number).
  static Object[][] argSets() {
    return new Object[][] {
        {}, {"A", "B", "C"}, {7, -1234567L, 2.5}, {null, "x'y", 0.0005}, {1.5, 0, -0.0},
        {Double.NaN, Double.POSITIVE_INFINITY, 1234567.891}};
  }

  static final String[] MF_INDEX = {"0", "1", "2", "3", "00", "+1", "-1", " 0", "x", "", "10000", "9999", "١", "1 "};
  static final String[] MF_LITERAL = {"a", " ", "'", "''", "}", "x'{'y", "'{0}'", "é", "#", "|", ","};
  static final String[] MF_CHOICE = {"0", "1", "2", "1.5", "-1", "∞", "-∞", "NaN", "1e1", "+1", " 1", "1f", "0x1p0", "#", "#", "<",
      "≤", "|", "|", "a", "b", "'", "''", "{0}", "{1}", " ", "{2,number}"};

  static String mfPattern() {
    StringBuilder sb = new StringBuilder();
    int n = r.nextInt(5);
    for (int i = 0; i < n; i++) {
      if (r.nextInt(3) == 0) {
        sb.append(pick(MF_LITERAL));
        continue;
      }
      sb.append('{').append(r.nextInt(3) == 0 ? pick(MF_INDEX) : String.valueOf(r.nextInt(3)));
      int t = r.nextInt(12);
      if (t < 4) {
        // no type
        if (t == 1) sb.append(',');
        if (t == 2) sb.append(", ,x");
      } else if (t < 7) {
        sb.append(',').append(pick(new String[] {"number", " NUMBER ", "number", "Number"}));
        String[] styles = {"", "integer", "percent", "currency", " Integer", "PERCENT ", "#.##", "'{'#", "0.0", "#,##0.#;(#)"};
        String st = pick(styles);
        if (r.nextBoolean() || !st.isEmpty()) sb.append(',').append(st);
      } else if (t < 8) {
        sb.append(',').append(pick(new String[] {"date", "time", "DATE", " time"}));
        if (r.nextBoolean()) sb.append(',').append(pick(new String[] {"short", " MEDIUM ", "long", "full", "yyyy", "H'h'", "abc", "q", "''", "é", "G"}));
      } else if (t < 11) {
        sb.append(",choice,");
        int k = r.nextInt(8);
        for (int j = 0; j < k; j++) sb.append(pick(MF_CHOICE));
      } else {
        sb.append(',').append(pick(new String[] {"foo", "'x'", "{}", " ", " number", "tıme"}));
      }
      if (r.nextInt(10) > 0) sb.append('}');
    }
    String p = sb.toString();
    if (r.nextInt(6) == 0) p = mutate(p, "{}',a0#|");
    return p;
  }

  static String[] NUM_PATTERNS = {"{0}", "{0,number}", "{0,number,integer}", "{0,number,percent}", "{0,number,currency}"};

  static double randomDouble() {
    switch (r.nextInt(6)) {
      case 0: return Double.longBitsToDouble(r.nextLong());
      case 1: return (r.nextInt(2000000) - 1000000) / Math.pow(10, r.nextInt(9));
      case 2: return (10 * r.nextInt(100) + 5) / Math.pow(10, 1 + r.nextInt(8));
      case 3: return r.nextDouble() * Math.pow(10, r.nextInt(40) - 15);
      case 4: return (r.nextInt(100000) + 0.5) * (r.nextBoolean() ? 1 : -1);
      default: return 5 * Math.pow(10, -r.nextInt(12));
    }
  }

  // DecimalFormat pattern errors, which the Go side does not check (see
  // message_format.go).
  static boolean decimalFormatPatternError(String p) {
    try {
      new MessageFormat(p);
      return false;
    } catch (IllegalArgumentException e) {
      String m = String.valueOf(e.getMessage());
      return m.contains(" pattern \"") || m.startsWith("Malformed pattern") || m.startsWith("Malformed exponential");
    }
  }

  // Doubles of 2^63 or more whose DecimalFormat digits are not the shortest
  // ones (rare; not reproduced, see message_format.go).
  static boolean longDigits(double d) {
    if (Math.abs(d) < 0x1p63 || Double.isInfinite(d) || Double.isNaN(d)) return false;
    String sci = new java.text.DecimalFormat("0.#####################E0", java.text.DecimalFormatSymbols.getInstance(Locale.US)).format((Object) Math.abs(d));
    String digits = sci.substring(0, sci.indexOf('E')).replace(".", "").replaceAll("0+$", "");
    String shortest = new java.math.BigDecimal(Double.toString(Math.abs(d))).unscaledValue().toString().replaceAll("0+$", "");
    return !digits.equals(shortest);
  }

  // Whether every format of the pattern can format numbers on the Go side:
  // no date or time formats, no DecimalFormat patterns other than the four
  // standard ones (also inside ChoiceFormat results that are patterns).
  static boolean numbersSupported(String p, int depth) {
    MessageFormat mf;
    try {
      mf = new MessageFormat(p);
    } catch (IllegalArgumentException e) {
      return depth == 0; // the error itself is reproduced
    }
    java.text.Format[] standard = {java.text.NumberFormat.getInstance(), java.text.NumberFormat.getIntegerInstance(),
        java.text.NumberFormat.getPercentInstance(), java.text.NumberFormat.getCurrencyInstance()};
    for (java.text.Format f : mf.getFormats()) {
      if (f == null) continue;
      if (f instanceof java.text.ChoiceFormat) {
        for (Object s : ((java.text.ChoiceFormat) f).getFormats()) {
          if (((String) s).indexOf('{') >= 0 && (depth > 3 || !numbersSupported((String) s, depth + 1))) return false;
        }
        continue;
      }
      if (!Arrays.asList(standard).contains(f)) return false;
    }
    return true;
  }

  static void messageFormats(PrintStream ps, int n) {
    Object[][] sets = argSets();
    for (int i = 0; i < n; i++) {
      String p = mfPattern();
      while (decimalFormatPatternError(p)) p = mfPattern();
      // Numbers only where the Go side can format them (see
      // numbersSupported); otherwise no arguments or text ones.
      int set = numbersSupported(p, 0) ? r.nextInt(sets.length) : r.nextInt(2);
      Object[] args = sets[set];
      String pattern = p;
      ps.println("mf\t" + esc(pattern) + "\t" + set + "\t" + guard(() -> MessageFormat.format(pattern, args)));
    }
    for (int i = 0; i < n; i++) {
      String p = pick(NUM_PATTERNS);
      if (r.nextInt(4) == 0) {
        long v = r.nextInt(3) == 0 ? r.nextLong() : r.nextInt(2000001) - 1000000;
        ps.println("ml\t" + esc(p) + "\t" + v + "\t" + guard(() -> MessageFormat.format(p, v)));
      } else {
        double d0 = randomDouble();
        while (longDigits(d0) || longDigits(d0 * 100)) d0 = randomDouble();
        double d = d0;
        if (r.nextInt(8) == 0) {
          float f = (float) d;
          ps.println("mF\t" + esc(p) + "\t" + Integer.toHexString(Float.floatToRawIntBits(f)) + "\t" + guard(() -> MessageFormat.format(p, f)));
        } else {
          ps.println("md\t" + esc(p) + "\t" + Long.toHexString(Double.doubleToRawLongBits(d)) + "\t" + guard(() -> MessageFormat.format(p, d)));
        }
      }
    }
  }

  // ---------------------------------------------------------------- Properties
  static final String[] PROP_TOKENS = {"a", "b", "k", "=", ":", " ", " ", "\t", "\f", "\\", "\\", "\n", "\n", "\r", "\r\n", "#", "!",
      "\\u0041", "\\u00", "\\uzzzz", "\\n", "\\t", "é", "\\\n", "\\\r\n", "x", "\\=", "\\ ", "😀", "\u000b"};

  static void properties(PrintStream ps, int n) {
    for (int i = 0; i < n; i++) {
      String text = pickTokens(PROP_TOKENS, 14);
      ps.println("p\t" + esc(text) + "\t" + guard(() -> {
        Properties p = new Properties();
        p.load(new StringReader(text));
        List<String> keys = new ArrayList<>(p.stringPropertyNames());
        Collections.sort(keys);
        StringBuilder sb = new StringBuilder();
        for (String k : keys) sb.append(sb.length() == 0 ? "" : "\u0001").append(k).append('\u0002').append(p.getProperty(k));
        return sb.toString();
      }));
    }
  }

  // ---------------------------------------------------------------- split
  // Regexes that mean the same in java.util.regex and Go's RE2 on these
  // inputs; the first ones never match the empty string.
  static final String[] SPLIT_REGEX = {"\r?\n", ",", ",+", ";|,", "\\s+", "[,;]", "(,)", "é", "😀|,", "a|b",
      "a*", "", "b?", ",|", "x*", ",*"};
  static final int SPLIT_NULLABLE = 10;

  static void splits(PrintStream ps, int n) {
    for (int i = 0; i < n; i++) {
      int ri = r.nextInt(SPLIT_REGEX.length);
      String[] tokens = ri >= SPLIT_NULLABLE ? new String[] {"a", "b", ",", ";", " ", "\n", "\r\n", "é", "aa", ",,"}
          : new String[] {"a", "b", ",", ";", " ", "\n", "\r\n", "é", "😀", ",,", "\r"};
      String s = pickTokens(tokens, 8);
      int limit = r.nextInt(5) - 1;
      ps.println("s\t" + esc(s) + "\t" + ri + "\t" + limit + "\t" + guard(() -> {
        String[] parts = s.split(SPLIT_REGEX[ri], limit);
        StringBuilder sb = new StringBuilder().append(parts.length);
        for (String x : parts) sb.append('\u0001').append(x);
        return sb.toString();
      }));
    }
  }

  // ---------------------------------------------------------------- regionMatches, quote
  static final String[] CASE_TOKENS = {"a", "A", "b", "B", "ı", "I", "İ", "i", "ſ", "s", "S", "K", "k", "K",
      "µ", "Μ", "μ", "ß", "ẞ", "𐐀", "𐐨", "\ud801", "\udc00", "\udc28", "ǅ",
      "Ǆ", "ǆ", "ς", "Σ", "σ", "1"};

  static void regionMatches(PrintStream ps, int n) {
    for (int i = 0; i < n; i++) {
      String a = pickTokens(CASE_TOKENS, 5), b = r.nextBoolean() ? pickTokens(CASE_TOKENS, 5) : flipCase(a);
      int to = r.nextInt(a.length() + 3) - 1, oo = r.nextInt(b.length() + 3) - 1, len = r.nextInt(Math.max(a.length(), b.length()) + 3) - 1;
      boolean ic = r.nextInt(4) > 0;
      ps.println("r\t" + esc(a) + "\t" + esc(b) + "\t" + ic + "\t" + to + "\t" + oo + "\t" + len + "\t" + a.regionMatches(ic, to, b, oo, len));
    }
    // Every supplementary code point with a case mapping.
    StringBuilder sb = new StringBuilder();
    for (int cp = 0x10000; cp < 0x110000; cp++) {
      int u = Character.toUpperCase(cp), l = Character.toLowerCase(cp);
      if (u != cp || l != cp) ps.println("c\t" + Integer.toHexString(cp) + "\t" + Integer.toHexString(u) + "\t" + Integer.toHexString(l));
    }
  }

  static void quotes(PrintStream ps, int n) {
    for (int i = 0; i < n; i++) {
      String s = pickTokens(new String[] {"a", "\\", "\\E", "\\Q", "E", "Q", ".", "*", "é", "\\\\E", "$", "("}, 6);
      String q = Pattern.quote(s);
      ps.println("q\t" + esc(s) + "\t" + esc(q) + "\t" + Pattern.compile(q).matcher(s).matches());
    }
  }

  // ---------------------------------------------------------------- UTF-16
  static void utf16(PrintStream ps, int n) {
    String[] css = {"UTF-16", "UTF-16BE", "UTF-16LE"};
    for (int i = 0; i < n; i++) {
      byte[] b = new byte[r.nextInt(13)];
      for (int k = 0; k < b.length; k++) {
        int x = r.nextInt(8);
        b[k] = (byte) (x == 0 ? 0xD8 + r.nextInt(4) : x == 1 ? 0xDC + r.nextInt(4) : x == 2 ? 0xFE + r.nextInt(2) : x == 3 ? 0 : r.nextInt(256));
      }
      String csName = pick(css);
      Charset cs = Charset.forName(csName);
      StringBuilder hx = new StringBuilder();
      for (byte x : b) hx.append(String.format("%02x", x & 0xff));
      String strict;
      try {
        strict = esc(cs.newDecoder().onMalformedInput(CodingErrorAction.REPORT).onUnmappableCharacter(CodingErrorAction.REPORT)
            .decode(ByteBuffer.wrap(b)).toString());
      } catch (MalformedInputException e) {
        strict = "MALFORMED:" + e.getInputLength();
      } catch (Exception e) {
        strict = exc(e);
      }
      ps.println(csName + "\t" + hx + "\t" + esc(new String(b, cs)) + "\t" + strict);
    }
  }

  // ---------------------------------------------------------------- BufferedWriter
  // Records what reaches the underlying writer.
  static class LogWriter extends java.io.Writer {
    final StringBuilder log = new StringBuilder();
    public void write(char[] c, int off, int len) { log.append(" w").append(len); }
    public void write(String s, int off, int len) { log.append(" w").append(len); }
    public void flush() { log.append(" f"); }
    public void close() { log.append(" c"); }
  }

  static void bufferedWriters(PrintStream ps, int n) {
    for (int i = 0; i < n; i++) {
      LogWriter out = new LogWriter();
      java.io.BufferedWriter w = new java.io.BufferedWriter(out);
      StringBuilder ops = new StringBuilder();
      int k = 1 + r.nextInt(8);
      try {
        for (int j = 0; j < k; j++) {
          int op = r.nextInt(10);
          if (op < 6) {
            int len = r.nextInt(4) == 0 ? r.nextInt(20000) : r.nextInt(10);
            ops.append(" w").append(len);
            w.write("x".repeat(len));
          } else if (op < 7) {
            ops.append(" a");
            w.append('y');
          } else if (op < 9) {
            ops.append(" f");
            w.flush();
          } else {
            ops.append(" c");
            w.close();
          }
        }
        out.log.append(" .");
      } catch (java.io.IOException e) {
        out.log.append(" !" + e.getMessage());
      }
      ps.println("b\t" + ops.toString().trim() + "\t" + out.log.toString().trim());
    }
  }

  // ---------------------------------------------------------------- main
  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    long seed = Long.parseLong(args[1]);
    Locale.setDefault(Locale.US); // MessageFormat's numbers are en_US
    // Each file has its own random stream, so that adding a section does not
    // change the others.
    r = new Random(seed * 31 + 1);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("text_uri.tsv")), true, "UTF-8")) {
      uris(ps, 6000);
    }
    r = new Random(seed * 31 + 2);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("text_path.tsv")), true, "UTF-8")) {
      paths(ps, 1000);
    }
    r = new Random(seed * 31 + 3);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("text_messageformat.tsv")), true, "UTF-8")) {
      messageFormats(ps, 3000);
    }
    r = new Random(seed * 31 + 4);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("text_properties.tsv")), true, "UTF-8")) {
      properties(ps, 3000);
    }
    r = new Random(seed * 31 + 5);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("text_split.tsv")), true, "UTF-8")) {
      splits(ps, 3000);
    }
    r = new Random(seed * 31 + 6);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("text_case.tsv")), true, "UTF-8")) {
      regionMatches(ps, 3000);
      quotes(ps, 500);
    }
    r = new Random(seed * 31 + 7);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("text_utf16.tsv")), true, "UTF-8")) {
      utf16(ps, 3000);
    }
    r = new Random(seed * 31 + 8);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("text_bufferedwriter.tsv")), true, "UTF-8")) {
      bufferedWriters(ps, 500);
    }
  }
}

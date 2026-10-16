// Oracle tool: writes black-box reference vectors for the clean-room text
// helpers of internal/jutil (java.net.URI, java.nio.file.Path,
// java.text.MessageFormat, Properties.load, String.split, regionMatches,
// Pattern.quote, UTF-16 decoding) into <dir>/text_*.tsv, from random inputs
// drawn with the given seed. Only the public API of the JDK is called.
//   javac -d .cache/oracle-build tools/oracle/java/JutilTextVectors.java
//   LC_ALL=C.UTF-8 java -cp .cache/oracle-build JutilTextVectors internal/jutil/testdata 20261016
// (a UTF-8 locale, so that sun.jnu.encoding, which java.nio.file.Path uses,
// is UTF-8; the default locale must be en_US for MessageFormat's numbers).
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

  // ---------------------------------------------------------------- main
  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    long seed = Long.parseLong(args[1]);
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
  }
}

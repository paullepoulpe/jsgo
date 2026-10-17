// Oracle tool: writes Java reference vectors for the jutil helpers added for
// the jscomp diagnostics/deps port (MessageFormat, Path, URI, %.Nf,
// String.split, BaseEncoding, PercentEscaper, TreeSet, LoadingCache,
// ImmutableMultiset, Pattern) into internal/jutil/testdata/jscompa_*.tsv.
//   LC_ALL=C.UTF-8 java -cp .cache/$CLOSURE_JAR:.cache/oracle-build JutilJscompAVectors internal/jutil/testdata
// (a UTF-8 locale, so that sun.jnu.encoding, which java.nio.file.Path uses, is UTF-8)
import com.google.javascript.jscomp.jarjar.com.google.common.cache.CacheBuilder;
import com.google.javascript.jscomp.jarjar.com.google.common.cache.CacheLoader;
import com.google.javascript.jscomp.jarjar.com.google.common.cache.LoadingCache;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableMultiset;
import com.google.javascript.jscomp.jarjar.com.google.common.io.BaseEncoding;
import com.google.javascript.jscomp.jarjar.com.google.common.net.PercentEscaper;
import java.io.PrintStream;
import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.text.MessageFormat;
import java.util.*;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;
import java.util.regex.Pattern;

public class JutilJscompAVectors {
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

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(4242);

    // MessageFormat: every DiagnosticType format in the jar, plus edge cases.
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("jscompa_messageformat.tsv")), true, "UTF-8")) {
      TreeSet<String> formats = new TreeSet<>();
      String jar = System.getProperty("java.class.path").split(":")[0];
      ClassLoader cl = JutilJscompAVectors.class.getClassLoader();
      Class<?> dtClass = Class.forName("com.google.javascript.jscomp.DiagnosticType");
      Field formatField = dtClass.getField("format");
      try (JarFile jf = new JarFile(jar)) {
        for (Enumeration<JarEntry> e = jf.entries(); e.hasMoreElements(); ) {
          String n = e.nextElement().getName();
          if (!n.startsWith("com/google/javascript/jscomp/") || !n.endsWith(".class") || n.contains("jarjar")) continue;
          String cn = n.substring(0, n.length() - 6).replace('/', '.');
          Class<?> c;
          try {
            c = Class.forName(cn, true, cl);
          } catch (Throwable t) {
            continue;
          }
          Field[] fields;
          try {
            fields = c.getDeclaredFields();
          } catch (Throwable t) {
            continue;
          }
          for (Field f : fields) {
            if (!Modifier.isStatic(f.getModifiers()) || f.getType() != dtClass) continue;
            try {
              f.setAccessible(true);
              Object dt = f.get(null);
              if (dt != null) formats.add((String) formatField.get(dt));
            } catch (Throwable t) {
              // ignore
            }
          }
        }
      }
      formats.addAll(Arrays.asList(
          "", "plain", "{0}", "{1}{0}", "'{0}'", "''{0}''", "it''s {0}", "'quoted {0} text", "{0,}", "{0,,}",
          "{ 0}", "{0 }", "{-1}", "{x}", "{0", "0}", "{0}}", "{{0}}", "{0,number}", "{0,number,#.##}",
          "{0,date}", "{0,time,short}", "{0,choice,0#a|1#b}", "{0,foo}", "{0, NUMBER }", "{0,number,'{'}",
          "a {0} b {1} c {2} d {3} e {9}", "{10000}", "{9999}", "{00}", "{+1}", "'", "''", "'''", "x'{'y",
          "{0,'x'}", "{0,{}}", "\u00e9{0}\ud83d\ude00{1}"));
      String[][] argSets = {{}, {"A", "B", "C", "D"}, {"x'y", "{0}", "\u00fc"}};
      for (String f : formats) {
        for (String[] a : argSets) {
          ps.println(esc(f) + "\t" + a.length + "\t" + guard(() -> MessageFormat.format(f, (Object[]) a)));
        }
      }
    }

    // Path
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("jscompa_path.tsv")), true, "UTF-8")) {
      String[] paths = {"", "/", "//", "a", "a/", "a//b", "/a/b/c.js", "./a.js", "../a.js", "a/./b/../c", "/../a",
          "a/b/../../..", "foo/bar/baz.js", "/foo/./bar//", "m0.js", "../../x/y", ".", "..", "./", "a/..", "/a/b/../../../c",
          "x/y/z", "\u00e9/f", "a\u0000b"};
      for (String p : paths) {
        ps.println("of\t" + esc(p) + "\t" + guard(() -> {
          Path x = Path.of(p);
          return x + "|" + x.getParent() + "|" + x.getFileName() + "|" + x.getNameCount() + "|" + x.normalize()
              + "|" + x.hashCode() + "|" + x.isAbsolute();
        }));
        for (String q : paths) {
          ps.println("rel\t" + esc(p) + "\t" + esc(q) + "\t" + guard(() -> Path.of(p).relativize(Path.of(q)).toString()));
          ps.println("sib\t" + esc(p) + "\t" + esc(q) + "\t" + guard(() -> Path.of(p).resolveSibling(q).normalize().toString()));
          ps.println("sw\t" + esc(p) + "\t" + esc(q) + "\t" + guard(() -> Path.of(p).startsWith(Path.of(q)) + "," + Path.of(p).endsWith(Path.of(q)) + "," + Integer.signum(Path.of(p).compareTo(Path.of(q)))));
        }
      }
    }

    // URI
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("jscompa_uri.tsv")), true, "UTF-8")) {
      String[] uris = {"", "a", "a.js", "/a/b.js", "http://example.com/x", "HTTP://Example.COM/x", "file:///foo/bar",
          "mailto:x@y.z", "urn:isbn:123", "data:application/json;base64,eyJ9", "//host/path", "http://[::1]:80/",
          "http://[::1%eth0]/", "http://[1:2:3:4:5:6:7:8]/", "http://[1:2]/", "http://1.2.3.4:8080/p?q#f",
          "http://256.1.1.1/", "http://a_b.com/", "http://-a.com/", "http://a-.com/", "http://user@host:99999999999/",
          "a b", "a%20b", "a%2", "a%zz", "%41", "?q", "#f", ":x", "1a:b", "a:", "a:#", "http:", "http://", "x://y:/",
          "foo/bar#frag", "\u00e9", "a\u00a0b", "jar:file:/x.zip!/y.js", "C:\\path", "a?b?c", "a#b#c",
          "http://host/%41", "http://host/%61", "../x/y.js", "./m.js", "http://x.y.z./", "http://1.2.3/"};
      for (String u : uris) {
        ps.println("p\t" + esc(u) + "\t" + guard(() -> {
          URI x = new URI(u);
          return x.isAbsolute() + "|" + x.isOpaque() + "|" + x.hashCode() + "|" + x + "|" + x.getScheme() + "|" + x.getHost() + "|" + x.getPort() + "|" + x.getRawPath();
        }));
      }
      String[][] pairs = {{"http://example.com/x", "HTTP://EXAMPLE.com/x"}, {"http://example.com/X", "http://example.com/x"},
          {"a%41", "a%41"}, {"a%4a", "a%4A"}, {"a%4a", "a%4b"}, {"mailto:x", "MAILTO:x"}, {"a", "a#f"}, {"//h/p", "//H/p"},
          {"file:///a", "file:/a"}, {"x:y", "x:/y"}};
      for (String[] p : pairs) {
        ps.println("eq\t" + esc(p[0]) + "\t" + esc(p[1]) + "\t" + guard(() -> String.valueOf(new URI(p[0]).equals(new URI(p[1])))));
      }
    }

    // %.Nf
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("jscompa_fixed.tsv")), true, "UTF-8")) {
      List<Double> ds = new ArrayList<>(Arrays.asList(0.0, -0.0, 0.05, 0.15, 0.25, 0.35, 0.45, 1.0, 9.95, 99.95, 99.96, 100.0,
          -0.04, 1e-10, 1e10, 123456789.125, 4.9e-324, Double.MAX_VALUE, Double.NaN, Double.POSITIVE_INFINITY,
          Double.NEGATIVE_INFINITY, 0.5, 1.5, 2.5, 50.0, 66.66666666666667, 33.333333333333336, 12.345));
      for (int i = 0; i < 1500; i++) {
        switch (i % 4) {
          case 0 -> ds.add(r.nextDouble() * 100);
          case 1 -> ds.add(Math.round(r.nextDouble() * 10000) / 100.0);
          case 2 -> ds.add(Double.longBitsToDouble(r.nextLong()));
          default -> ds.add((double) r.nextInt(1000) / (1 + r.nextInt(1000)) * 100);
        }
      }
      for (double d : ds) {
        if (Double.isNaN(d) && Double.doubleToRawLongBits(d) != Double.doubleToRawLongBits(Double.NaN)) continue;
        ps.println(Long.toHexString(Double.doubleToRawLongBits(d)) + "\t" + String.format("%.1f", d) + "\t" + String.format("%.0f", d) + "\t" + String.format("%.3f", d));
      }
    }

    // String.split, base64, PercentEscaper
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("jscompa_strings.tsv")), true, "UTF-8")) {
      String[] ss = {"", "/", "//", "a", "/a", "a/", "a//b", "/a/b/", "x/node_modules/y", "/node_modules/", "/node_modules/a/node_modules/b.js", "a/node_modules"};
      for (String s : ss) {
        ps.println("split/\t" + esc(s) + "\t" + esc(String.join("|", s.split("/"))) + "\t" + s.split("/").length);
        ps.println("splitnm\t" + esc(s) + "\t" + esc(String.join("|", s.split("/node_modules/"))) + "\t" + s.split("/node_modules/").length);
      }
      String[] b64 = {"", "YQ", "YQ==", "YQ=", "YQ=====", "YWI", "YWJj", "Y", "Y===", "eyJ2ZXJzaW9uIjozfQ==", "a b", "a\nb", "\u00e9", "@@@@", "YW\u007fJ", "YW!J"};
      for (String s : b64) {
        ps.println("b64d\t" + esc(s) + "\t" + guard(() -> {
          byte[] b = BaseEncoding.base64().decode(s);
          StringBuilder sb = new StringBuilder();
          for (byte x : b) sb.append(String.format("%02x", x & 0xff));
          return sb.toString();
        }));
      }
      PercentEscaper pe = new PercentEscaper("-_.*", false);
      String[] pes = {"", "abc", "a b", "{\"version\":3,\"file\":\"x.js\"}", "\u00e9\u4e2d\ud83d\ude00", "-_.*~!", "a\ud800", "\udc00a", "\ud800a", "%+&=/:;"};
      for (String s : pes) {
        ps.println("pe\t" + esc(s) + "\t" + guard(() -> pe.escape(s)));
      }
      ps.println("ms\t" + esc(ImmutableMultiset.of("a", "b", "a", "c", "a", "b").toString()));
      ps.println("ms\t" + esc(ImmutableMultiset.of().toString()));
    }

    // Charset decoding (REPLACE and REPORT).
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("jscompa_decode.tsv")), true, "UTF-8")) {
      String[] css = {"UTF-8", "ISO-8859-1", "US-ASCII", "UTF-16", "UTF-16BE", "UTF-16LE"};
      byte[][] special = {{}, {(byte) 0xfe, (byte) 0xff, 0, 0x41}, {(byte) 0xff, (byte) 0xfe, 0x41, 0}, {(byte) 0xd8, 0, (byte) 0xdc, 0},
          {(byte) 0xdc, 0, 0, 0x41}, {(byte) 0xd8, 0}, {0x41}, {(byte) 0xed, (byte) 0xa0, (byte) 0x80}, {(byte) 0xe0, (byte) 0x80},
          {(byte) 0xf0, (byte) 0x90, (byte) 0x80}, {(byte) 0xf4, (byte) 0x90}, {(byte) 0xc0, (byte) 0xaf}, {(byte) 0xe2, (byte) 0x82, 0x41}};
      List<byte[]> inputs = new ArrayList<>(Arrays.asList(special));
      for (int i = 0; i < 300; i++) {
        byte[] b = new byte[r.nextInt(12)];
        r.nextBytes(b);
        if (i % 3 == 0) for (int k = 0; k < b.length; k++) if (r.nextBoolean()) b[k] &= 0x7f;
        inputs.add(b);
      }
      for (String csName : css) {
        java.nio.charset.Charset cs = java.nio.charset.Charset.forName(csName);
        for (byte[] b : inputs) {
          StringBuilder hx = new StringBuilder();
          for (byte x : b) hx.append(String.format("%02x", x & 0xff));
          String strict;
          try {
            strict = esc(cs.newDecoder().onMalformedInput(java.nio.charset.CodingErrorAction.REPORT)
                .onUnmappableCharacter(java.nio.charset.CodingErrorAction.REPORT).decode(java.nio.ByteBuffer.wrap(b)).toString());
          } catch (java.nio.charset.MalformedInputException e) {
            strict = "MALFORMED:" + e.getInputLength();
          } catch (Exception e) {
            strict = exc(e);
          }
          ps.println(csName + "\t" + hx + "\t" + esc(new String(b, cs)) + "\t" + strict);
        }
      }
    }

    // TreeSet with a comparator that changes (ComposeWarningsGuard's pattern).
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("jscompa_treeset.tsv")), true, "UTF-8")) {
      for (int trial = 0; trial < 300; trial++) {
        Map<Integer, Integer> order = new HashMap<>();
        int[] prio = new int[12];
        for (int i = 0; i < prio.length; i++) prio[i] = r.nextInt(4);
        TreeSet<Integer> set = new TreeSet<>((a, b) -> {
          int d = prio[a] - prio[b];
          if (d != 0) return d;
          return order.get(b) - order.get(a);
        });
        StringBuilder ops = new StringBuilder();
        int adds = 0;
        int n = 1 + r.nextInt(30);
        for (int k = 0; k < n; k++) {
          int g = r.nextInt(prio.length);
          if (r.nextInt(5) == 0 && order.containsKey(g)) {
            ops.append("r").append(g).append(' ');
            set.remove(g);
          } else {
            ops.append("a").append(g).append(' ');
            adds++;
            order.put(g, adds);
            set.remove(g);
            set.add(g);
          }
        }
        StringBuilder asc = new StringBuilder();
        for (int x : set) asc.append(x).append(',');
        StringBuilder desc = new StringBuilder();
        for (int x : set.descendingSet()) desc.append(x).append(',');
        StringBuilder pr = new StringBuilder();
        for (int p : prio) pr.append(p).append(',');
        ps.println(pr + "\t" + ops.toString().trim() + "\t" + asc + "\t" + desc + "\t" + set.size());
      }
    }

    // LoadingCache eviction order.
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("jscompa_cache.tsv")), true, "UTF-8")) {
      for (int trial = 0; trial < 150; trial++) {
        int max = r.nextInt(30);
        int keys = 1 + r.nextInt(60);
        List<String> loads = new ArrayList<>();
        String spec = "maximumSize=" + max;
        LoadingCache<String, String> cache = CacheBuilder.from(spec).build(new CacheLoader<String, String>() {
          @Override public String load(String key) {
            loads.add(key);
            return "v" + key;
          }
        });
        StringBuilder ops = new StringBuilder();
        for (int k = 0; k < 60; k++) {
          String key = "k" + r.nextInt(keys);
          ops.append(key).append(' ');
          cache.getUnchecked(key);
        }
        ps.println(spec + "\t" + ops.toString().trim() + "\t" + String.join(" ", loads) + "\t" + cache.size());
      }
    }
  }
}

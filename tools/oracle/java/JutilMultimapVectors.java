// Generates internal/jutil/testdata/multimap.tsv: Guava ImmutableSetMultimap
// built from random puts and its inverse(), and ImmutableMultimap.builder()
// with and without orderKeysBy(naturalOrder()).
//   java -cp .cache/$CLOSURE_JAR:.cache/oracle-build JutilMultimapVectors internal/jutil/testdata
//
// Rows: "set\t<puts>\t<toString>\t<inverse().toString()>" and
// "list\t<sorted 0|1>\t<puts>\t<toString>", where <puts> is k=v;k=v;... with
// keys and values as hex UTF-16 code units (4 digits each).
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSetMultimap;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.Comparator;
import java.util.Random;

public class JutilMultimapVectors {
  static final String[] WORDS = {
    "a", "b", "c", "B", "aa", "", "é", "￿", "😀", "z", "中", "a\u0000"
  };

  static String hex(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      sb.append(String.format("%04x", (int) s.charAt(i)));
    }
    return sb.toString();
  }

  static String esc(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c < 0x20 || c >= 0x7f || c == '\\') {
        sb.append(String.format("\\u%04x", (int) c));
      } else {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(42);
    try (PrintStream ps =
        new PrintStream(
            Files.newOutputStream(out.resolve("multimap.tsv")), false, StandardCharsets.UTF_8)) {
      for (int k = 0; k < 300; k++) {
        int n = r.nextInt(12);
        ImmutableSetMultimap.Builder<String, String> b = ImmutableSetMultimap.builder();
        StringBuilder puts = new StringBuilder();
        for (int i = 0; i < n; i++) {
          String key = WORDS[r.nextInt(6)];
          String value = WORDS[r.nextInt(WORDS.length)];
          b.put(key, value);
          puts.append(i == 0 ? "" : ";").append(hex(key)).append('=').append(hex(value));
        }
        ImmutableSetMultimap<String, String> m = b.build();
        ps.println("set\t" + puts + "\t" + esc(m.toString()) + "\t" + esc(m.inverse().toString()));
      }
      for (int k = 0; k < 300; k++) {
        int n = r.nextInt(12);
        boolean sorted = r.nextBoolean();
        ImmutableMultimap.Builder<String, String> b = ImmutableMultimap.builder();
        if (sorted) {
          b.orderKeysBy(Comparator.naturalOrder());
        }
        StringBuilder puts = new StringBuilder();
        for (int i = 0; i < n; i++) {
          String key = WORDS[r.nextInt(WORDS.length)];
          String value = WORDS[r.nextInt(6)];
          b.put(key, value);
          puts.append(i == 0 ? "" : ";").append(hex(key)).append('=').append(hex(value));
        }
        ps.println("list\t" + (sorted ? 1 : 0) + "\t" + puts + "\t" + esc(b.build().toString()));
      }
    }
  }
}

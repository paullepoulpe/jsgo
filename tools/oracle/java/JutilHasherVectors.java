// Generates internal/jutil/testdata/farmhash_hasher.tsv: Guava's
// Hashing.farmHashFingerprint64().newHasher() fed a random sequence of
// putInt / putLong / putString(UTF_8) / putBytes calls.
//   java -cp .cache/closure-compiler.jar:.cache/oracle-build JutilHasherVectors internal/jutil/testdata
//
// Each row is "<op>;<op>;...\t<hash().asLong()>", where an op is I<int>,
// L<long>, B<hex bytes> or S<hex UTF-16 code units, 4 digits each>.
import com.google.javascript.jscomp.jarjar.com.google.common.hash.Hasher;
import com.google.javascript.jscomp.jarjar.com.google.common.hash.Hashing;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.Random;

public class JutilHasherVectors {
  static final char[] CHARS = {
    'a', 'Z', '0', '_', '$', '.', ' ', 'é', 'Δ', '中', '￿', '\ud83d', '\ude00',
    '\ud800', '\udfff', '\u0000'
  };

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(20260415);
    try (PrintStream ps =
        new PrintStream(Files.newOutputStream(out.resolve("farmhash_hasher.tsv")))) {
      for (int k = 0; k < 400; k++) {
        int n = 1 + r.nextInt(6);
        Hasher hs = Hashing.farmHashFingerprint64().newHasher();
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < n; i++) {
          if (i > 0) {
            sb.append(';');
          }
          switch (r.nextInt(4)) {
            case 0 -> {
              int v = r.nextBoolean() ? r.nextInt() : r.nextInt(100);
              hs.putInt(v);
              sb.append('I').append(v);
            }
            case 1 -> {
              long v = r.nextLong();
              hs.putLong(v);
              sb.append('L').append(v);
            }
            case 2 -> {
              byte[] b = new byte[r.nextInt(40)];
              r.nextBytes(b);
              hs.putBytes(b);
              sb.append('B');
              for (byte x : b) {
                sb.append(String.format("%02x", x & 0xff));
              }
            }
            default -> {
              int len = r.nextInt(30);
              StringBuilder s = new StringBuilder();
              for (int j = 0; j < len; j++) {
                s.append(CHARS[r.nextInt(CHARS.length)]);
              }
              hs.putString(s, StandardCharsets.UTF_8);
              sb.append('S');
              for (int j = 0; j < s.length(); j++) {
                sb.append(String.format("%04x", (int) s.charAt(j)));
              }
            }
          }
        }
        ps.println(sb + "\t" + hs.hash().asLong());
      }
    }
  }
}

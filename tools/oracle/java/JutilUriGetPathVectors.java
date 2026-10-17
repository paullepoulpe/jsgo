// Generates internal/jutil/testdata/uri_get_path.tsv: java.net.URI.getPath() (the decoded path)
// of seeded random URI strings whose paths mix plain ASCII, non-ASCII characters and escaped
// octets (valid and invalid UTF-8 sequences), plus a few fixed cases.
//   java -cp .cache/oracle-build JutilUriGetPathVectors internal/jutil/testdata
//
// Rows: "<input hex>\t<getPath() hex or 'null'>\t<getRawPath() hex or 'null'>", strings as hex
// UTF-16 code units (4 digits each), "-" for the empty string; rows whose input does not parse
// are "<input hex>\terror".
import java.io.PrintStream;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;

public class JutilUriGetPathVectors {
  static String hex(String s) {
    if (s == null) {
      return "null";
    }
    if (s.isEmpty()) {
      return "-";
    }
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      sb.append(String.format("%04x", (int) s.charAt(i)));
    }
    return sb.toString();
  }

  static final String[] PIECES = {
    "a", "b", "/", "/", "-", ".", "_", "~", "x.js", "é", "中", "😀", "%41", "%2F", "%2f", "%20",
    "%C3%A9", "%E4%B8%AD", "%F0%9F%98%80", "%C3", "%A9", "%FF", "%80", "%E4%B8", "%ED%A0%80",
    "%F0%9F%98", "%C0%80", "%25", "%00", "%7E", ":", "@", "!", "$", "&", "'", "(", ")", "*", "+",
    ",", ";", "="
  };

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(11);
    List<String> inputs = new ArrayList<>();
    inputs.add("");
    inputs.add("/a/b.js");
    inputs.add("http://host/a%20b/c");
    inputs.add("mailto:a@b");
    inputs.add("a?q=%41#f%42");
    inputs.add("//host");
    inputs.add("/%");
    for (int k = 0; k < 600; k++) {
      StringBuilder sb = new StringBuilder();
      int kind = r.nextInt(4);
      if (kind == 1) {
        sb.append("http://h.example");
        sb.append('/');
      } else if (kind == 2) {
        sb.append("file:/");
      }
      int n = 1 + r.nextInt(10);
      for (int i = 0; i < n; i++) {
        sb.append(PIECES[r.nextInt(PIECES.length)]);
      }
      if (r.nextInt(5) == 0) {
        sb.append("?x=%41");
      }
      inputs.add(sb.toString());
    }
    try (PrintStream ps =
        new PrintStream(
            Files.newOutputStream(out.resolve("uri_get_path.tsv")), false, StandardCharsets.UTF_8)) {
      for (String in : inputs) {
        try {
          URI u = new URI(in);
          ps.println(hex(in) + "\t" + hex(u.getPath()) + "\t" + hex(u.getRawPath()));
        } catch (URISyntaxException e) {
          ps.println(hex(in) + "\terror");
        }
      }
    }
  }
}

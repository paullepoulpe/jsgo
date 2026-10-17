// Generates internal/jutil/testdata/char_sink.tsv: Guava Files.asCharSink(File, UTF_8)
// .write(CharSequence) over a file that already holds longer content (so that truncation shows),
// for fixed strings (ASCII, non-ASCII, surrogate pairs, unpaired surrogates), plus the exception a
// write into a directory and into a missing directory throws.
//   java -cp .cache/closure-compiler.jar:.cache/oracle-build JutilCharSinkVectors internal/jutil/testdata
//
// Rows: "<input hex>\t<file bytes hex>" (input as hex UTF-16 code units, 4 digits each; "-" for
// an empty string or empty file), and "dir\t<exception class>: <message without the temp dir>"
// and "missing\t<exception class>: <message without the temp dir>".
import com.google.javascript.jscomp.jarjar.com.google.common.io.Files;
import java.io.File;
import java.io.IOException;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.nio.file.Paths;

public class JutilCharSinkVectors {
  static String hex16(String s) {
    if (s.isEmpty()) {
      return "-";
    }
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      sb.append(String.format("%04x", (int) s.charAt(i)));
    }
    return sb.toString();
  }

  static String hex8(byte[] b) {
    if (b.length == 0) {
      return "-";
    }
    StringBuilder sb = new StringBuilder();
    for (byte x : b) {
      sb.append(String.format("%02x", x & 0xff));
    }
    return sb.toString();
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    File dir = java.nio.file.Files.createTempDirectory("charsink").toFile();
    String[] inputs = {
      "", "a", "abc\n", "xé中\n", "😀", "a\ud800b", "a\udc00b", "end\ud800",
      "\udc00\ud800", "\r\n\t\u0000", "  \u0085", "￿￾"
    };
    try (PrintStream ps = new PrintStream(out.resolve("char_sink.tsv").toFile(), "UTF-8")) {
      int k = 0;
      for (String in : inputs) {
        File f = new File(dir, "f" + (k++));
        Files.asCharSink(f, StandardCharsets.UTF_8).write("previous longer content\n");
        Files.asCharSink(f, StandardCharsets.UTF_8).write(in);
        ps.println(hex16(in) + "\t" + hex8(java.nio.file.Files.readAllBytes(f.toPath())));
      }
      String prefix = dir.getPath() + File.separator;
      try {
        Files.asCharSink(dir, StandardCharsets.UTF_8).write("x");
        ps.println("dir\tnone");
      } catch (IOException e) {
        ps.println("dir\t" + e.getClass().getName() + ": " + e.getMessage().replace(dir.getPath(), "<dir>"));
      }
      try {
        Files.asCharSink(new File(dir, "missing/f"), StandardCharsets.UTF_8).write("x");
        ps.println("missing\tnone");
      } catch (IOException e) {
        ps.println("missing\t" + e.getClass().getName() + ": " + e.getMessage().replace(prefix, "<dir>/"));
      }
    }
  }
}

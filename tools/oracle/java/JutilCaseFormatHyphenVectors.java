// Oracle tool: writes the Guava reference vectors for
// jutil.CaseFormat_LowerHyphenToUpperCamel (CaseFormat.LOWER_HYPHEN.to(
// CaseFormat.UPPER_CAMEL, s)) into internal/jutil/testdata/caseformat_hyphen.tsv.
//   java -cp .cache/$CLOSURE_JAR:.cache/oracle-build JutilCaseFormatHyphenVectors internal/jutil/testdata
import com.google.javascript.jscomp.jarjar.com.google.common.base.CaseFormat;
import java.io.PrintStream;
import java.nio.file.*;
import java.util.*;

public class JutilCaseFormatHyphenVectors {
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

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(20260415);
    try (PrintStream ps =
        new PrintStream(Files.newOutputStream(out.resolve("caseformat_hyphen.tsv")))) {
      String[] fixed = {"", "-", "--", "a", "A", "x-foo", "x-Foo", "X-FOO", "foo--bar", "-foo",
          "foo-", "my-element", "paper-input", "iron-a11y-keys", "a$b", "é-éa", "É-a", "😀-a",
          "a-😀b", "\ud800-\udc00"};
      for (String s : fixed) {
        ps.println(units(s) + "\t" + units(CaseFormat.LOWER_HYPHEN.to(CaseFormat.UPPER_CAMEL, s)));
      }
      for (int k = 0; k < 2000; k++) {
        String s = randomString(r, r.nextInt(10), "aAbBzZ-09_$éÉ😀\ud800");
        ps.println(units(s) + "\t" + units(CaseFormat.LOWER_HYPHEN.to(CaseFormat.UPPER_CAMEL, s)));
      }
    }
  }
}

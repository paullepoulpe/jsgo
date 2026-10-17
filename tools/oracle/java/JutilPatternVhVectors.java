// Oracle tool: writes internal/jutil/testdata/pattern_vh.tsv, which characters
// java.util.regex matches for the \v, \V, \h and \H classes (alone and inside
// a character class), as jutil.Pattern translates them.
//   java -cp .cache/oracle-build JutilPatternVhVectors > internal/jutil/testdata/pattern_vh.tsv
import java.util.regex.Pattern;

public class JutilPatternVhVectors {
  public static void main(String[] args) {
    String[] regexes = {"\\v", "\\V", "\\h", "\\H", "[^}\\t\\n\\v\\f\\r ]", "[a\\h]", "x\\v*y"};
    StringBuilder chars = new StringBuilder();
    for (int c = 0; c < 0x80; c++) chars.append((char) c);
    for (int c : new int[] {0x85, 0xa0, 0x1680, 0x180e, 0x200b, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000, 0xfeff, 0xe9}) {
      chars.append((char) c);
    }
    for (int c = 0x2000; c <= 0x200a; c++) chars.append((char) c);
    for (String regex : regexes) {
      Pattern p = Pattern.compile(regex);
      StringBuilder bits = new StringBuilder();
      for (int i = 0; i < chars.length(); i++) {
        String in = regex.startsWith("x") ? "x" + chars.charAt(i) + "y" : String.valueOf(chars.charAt(i));
        bits.append(p.matcher(in).matches() ? '1' : '0');
      }
      StringBuilder hex = new StringBuilder();
      for (int i = 0; i < chars.length(); i++) {
        if (i > 0) hex.append(',');
        hex.append(Integer.toHexString(chars.charAt(i)));
      }
      System.out.println(regex + "\t" + hex + "\t" + bits);
    }
  }
}

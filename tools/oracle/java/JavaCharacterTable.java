// Generates internal/jutil/testdata/character.tsv (the JDK running it must be 21, the oracle JDK):
//   java -cp build JavaCharacterTable > internal/jutil/testdata/character.tsv
// Dumps Java 21 Character classification for every UTF-16 code unit as
// run-length-encoded lines: "<startHex> <endHex> <getType> <isLetter> <isDigit> <digit16> <digit10> <digit36>".
public class JavaCharacterTable {
  public static void main(String[] args) {
    StringBuilder sb = new StringBuilder();
    int start = 0;
    String prev = null;
    for (int c = 0; c <= 0x10000; c++) {
      String cur = c == 0x10000 ? null : key((char) c);
      if (prev != null && !prev.equals(cur)) {
        sb.append(String.format("%04X %04X %s%n", start, c - 1, prev));
        start = c;
      }
      if (prev == null) start = c;
      prev = cur;
    }
    System.out.print(sb);
  }
  static String key(char c) {
    return Character.getType(c) + " " + (Character.isLetter(c) ? 1 : 0) + " " + (Character.isDigit(c) ? 1 : 0)
        + " " + Character.digit(c, 16) + " " + Character.digit(c, 10) + " " + Character.digit(c, 36);
  }
}

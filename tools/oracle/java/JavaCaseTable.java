/**
 * Prints Character.toUpperCase(char) / toLowerCase(char) for every char whose
 * case mapping is not the identity ("CCCC UUUU LLLL" in hex), and the
 * String.regionMatches(true, ...) results the sourcemap escaper relies on.
 * Output: internal/jutil/testdata/character_case.tsv.
 */
public final class JavaCaseTable {
  public static void main(String[] args) {
    StringBuilder sb = new StringBuilder();
    for (int c = 0; c < 0x10000; c++) {
      char u = Character.toUpperCase((char) c);
      char l = Character.toLowerCase((char) c);
      if (u != c || l != c) {
        sb.append(String.format("%04X %04X %04X%n", c, (int) u, (int) l));
      }
    }
    System.out.print(sb);
  }
}

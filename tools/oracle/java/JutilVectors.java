// Oracle tool: writes the Java reference vectors for internal/jutil/testdata
// (farmhash.tsv, character_codepoint.tsv, bigint.tsv).
//   java -cp .cache/$CLOSURE_JAR:.cache/oracle-build JutilVectors internal/jutil/testdata
import com.google.javascript.jscomp.jarjar.com.google.common.hash.Hashing;
import com.google.javascript.jscomp.jarjar.com.google.common.hash.Hasher;
import java.io.PrintStream;
import java.math.BigInteger;
import java.nio.file.*;
import java.util.Random;

public class JutilVectors {
  static String hex(byte[] b) {
    StringBuilder sb = new StringBuilder("x");
    for (byte x : b) sb.append(String.format("%02x", x & 0xff));
    return sb.toString();
  }
  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(12345);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("farmhash.tsv")))) {
      for (int len = 0; len <= 300; len++) {
        for (int k = 0; k < 3; k++) {
          byte[] b = new byte[len];
          r.nextBytes(b);
          long h = Hashing.farmHashFingerprint64().hashBytes(b).asLong();
          ps.println(hex(b) + "\t" + h);
        }
      }
      for (int k = 0; k < 50; k++) {
        int n = 2 + r.nextInt(6);
        Hasher hs = Hashing.farmHashFingerprint64().newHasher();
        StringBuilder sb = new StringBuilder("L");
        for (int i = 0; i < n; i++) { long v = r.nextLong(); if (k % 3 == 0) v = r.nextInt(1000); hs.putLong(v); sb.append(i == 0 ? "" : ",").append(v); }
        ps.println(sb + "\t" + hs.hash().asLong());
      }
    }
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("character_codepoint.tsv")))) {
      int start = 0;
      String prev = null;
      for (int cp = 0; cp <= 0x110000; cp++) {
        String cur = cp == 0x110000 ? "END" :
            Character.getType(cp) + "," + (Character.isWhitespace(cp) ? 1 : 0) + (Character.isJavaIdentifierStart(cp) ? 1 : 0)
            + (Character.isJavaIdentifierPart(cp) ? 1 : 0) + (Character.isIdentifierIgnorable(cp) ? 1 : 0)
            + (Character.isLetter(cp) ? 1 : 0) + (Character.isDigit(cp) ? 1 : 0);
        if (prev != null && !cur.equals(prev)) {
          ps.println(Integer.toHexString(start) + "\t" + Integer.toHexString(cp - 1) + "\t" + prev);
          start = cp;
        }
        if (prev == null) start = cp;
        prev = cur;
      }
    }
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("bigint.tsv")))) {
      for (int k = 0; k < 500; k++) {
        BigInteger x = new BigInteger(1 + r.nextInt(200), r);
        if (k % 2 == 1) x = x.negate();
        if (k == 0) x = BigInteger.ZERO;
        if (k == 2) x = BigInteger.valueOf(-128);
        if (k == 4) x = BigInteger.valueOf(-129);
        if (k == 6) x = BigInteger.valueOf(255);
        ps.println(x + "\t" + hex(x.toByteArray()) + "\t" + x.intValue());
      }
    }
  }
}

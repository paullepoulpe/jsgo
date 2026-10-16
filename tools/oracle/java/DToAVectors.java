// Oracle tool: writes internal/dtoa/testdata/dtoa.tsv, the Java reference
// output of DToA.numberToString for many doubles (hex bits <TAB> string).
//   java -cp .cache/closure-compiler.jar:.cache/oracle-build DToAVectors internal/jutil/testdata/double_tostring.tsv internal/dtoa/testdata/dtoa.tsv
import com.google.javascript.rhino.dtoa.DToA;
import java.io.PrintStream;
import java.nio.file.*;
import java.util.*;

public class DToAVectors {
  static boolean isSmallSubnormal(long b) {
    long m = b & 0x000fffffffffffffL;
    return (b & 0x7ff0000000000000L) == 0 && m != 0 && m < (1L << 32);
  }

  public static void main(String[] args) throws Exception {
    LinkedHashSet<Long> bits = new LinkedHashSet<>();
    for (String line : Files.readAllLines(Paths.get(args[0]))) {
      bits.add(Long.parseUnsignedLong(line.split("\t")[0], 16));
    }
    // Rhino's DToA mis-estimates the exponent of small subnormals and then
    // emits (and quadratically builds) strings of up to ~10^6 characters,
    // taking minutes in Java. Keep only a few hand-picked small subnormals.
    bits.removeIf(b -> isSmallSubnormal(b));
    for (long b : new long[] {1L, 2L, 3L, 5L, 0xffL, 0xfffL, 0x1234L, 0x10000L, 0x316a2L, 0x7fffffffL, 0x80000001L}) {
      bits.add(b);
      bits.add(b | 0x8000000000000000L);
    }
    Random r = new Random(4242);
    double[] specials = {0.0, -0.0, 1.0, -1.0, 0.5, 0.1, 0.2, 0.3, 1e21, 1e-7, 1e-6, 123e-20, 1e300, 1.7976931348623157e308,
        Double.MIN_VALUE, Double.MIN_NORMAL, Double.NaN, Double.POSITIVE_INFINITY, Double.NEGATIVE_INFINITY,
        9007199254740992.0, 9007199254740993.0, 4294967296.0, 123456789012345680000.0, 1e22, 1e23, 5e-324, 2.5, 1.5};
    for (double d : specials) bits.add(Double.doubleToRawLongBits(d));
    for (int i = -1022; i <= 1023; i++) { bits.add(Double.doubleToRawLongBits(Math.pow(2, i))); }
    for (int i = -330; i <= 310; i++) { bits.add(Double.doubleToRawLongBits(Double.parseDouble("1e" + i))); bits.add(Double.doubleToRawLongBits(Double.parseDouble("7.77e" + i))); }
    for (int i = 0; i < 20000; i++) bits.add(r.nextLong());
    for (int i = 0; i < 5000; i++) bits.add(Double.doubleToRawLongBits((double) r.nextInt()));
    for (int i = 0; i < 5000; i++) bits.add(Double.doubleToRawLongBits(r.nextInt(1000000) / Math.pow(10, r.nextInt(12))));
    for (int i = 0; i < 3000; i++) bits.add((r.nextLong() & 0x000fffffffffffffL) | (1L << (32 + r.nextInt(20)))); // subnormals
    for (int i = 0; i < 3000; i++) bits.add(Double.doubleToRawLongBits(r.nextDouble() * Math.pow(10, r.nextInt(40) - 20)));
    bits.removeIf(b -> isSmallSubnormal(b) && !java.util.Set.of(1L, 2L, 3L, 5L, 0xffL, 0xfffL, 0x1234L, 0x10000L, 0x316a2L, 0x7fffffffL, 0x80000001L).contains(b & 0x7fffffffffffffffL));
    try (PrintStream ps = new PrintStream(Files.newOutputStream(Paths.get(args[1])), false, "UTF-8")) {
      for (long b : bits) {
        double d = Double.longBitsToDouble(b);
        String out;
        try {
          out = DToA.numberToString(d);
        } catch (RuntimeException e) {
          out = "!" + e.getClass().getName() + ": " + e.getMessage();
        }
        ps.println(Long.toHexString(b) + "\t" + out);
      }
    }
  }
}

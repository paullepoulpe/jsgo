// Vectors for java.util.BitSet.flip(int) and flip(int, int)
// (internal/jutil/bit_set_flip.go): random set/clear/flip sequences, then
// the resulting toString(), length() and cardinality() (and the
// IndexOutOfBoundsException messages for bad indices).
// java -cp .cache/oracle-build JutilBitSetFlipVectors > internal/jutil/testdata/bitset_flip.tsv
import java.util.BitSet;
import java.util.Random;

public class JutilBitSetFlipVectors {
  public static void main(String[] args) {
    Random rnd = new Random(20260415);
    StringBuilder out = new StringBuilder();
    for (int c = 0; c < 200; c++) {
      BitSet bs = new BitSet();
      StringBuilder ops = new StringBuilder();
      StringBuilder res = new StringBuilder();
      int n = 1 + rnd.nextInt(14);
      int span = 1 + rnd.nextInt(c % 3 == 0 ? 140 : 40);
      for (int i = 0; i < n; i++) {
        int k = rnd.nextInt(span);
        int op = rnd.nextInt(6);
        String r = "";
        try {
          if (op == 0) {
            bs.clear(k);
            ops.append('c').append(k);
          } else if (op == 1) {
            bs.set(k);
            ops.append('s').append(k);
          } else if (op == 2) {
            int f = k - 1 - rnd.nextInt(3);
            ops.append('f').append(f);
            bs.flip(f);
          } else {
            int to = k + rnd.nextInt(span) - (c % 5 == 0 ? 20 : 0);
            int from = k - (c % 7 == 0 ? 2 : 0);
            ops.append('F').append(from).append(':').append(to);
            bs.flip(from, to);
          }
          r = bs.toString() + "/" + bs.length() + "/" + bs.cardinality();
        } catch (IndexOutOfBoundsException e) {
          r = "!" + e.getMessage();
        }
        ops.append(',');
        res.append(r).append(';');
      }
      out.append(ops).append('\t').append(res).append('\n');
    }
    System.out.print(out);
  }
}

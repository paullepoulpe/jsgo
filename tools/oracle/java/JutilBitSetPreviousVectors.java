// Vectors for java.util.BitSet.previousSetBit (internal/jutil/bitset.go):
// random set/clear sequences, then previousSetBit at several indices
// (including -1 and the IndexOutOfBoundsException for fromIndex < -1).
// java -cp .cache/oracle-build JutilBitSetPreviousVectors > internal/jutil/testdata/bitset_previous.tsv
import java.util.BitSet;
import java.util.Random;

public class JutilBitSetPreviousVectors {
  public static void main(String[] args) {
    Random rnd = new Random(20260415);
    StringBuilder out = new StringBuilder();
    for (int c = 0; c < 400; c++) {
      BitSet bs = new BitSet();
      StringBuilder ops = new StringBuilder();
      int n = rnd.nextInt(14);
      int span = 1 + rnd.nextInt(c % 3 == 0 ? 300 : 70);
      for (int i = 0; i < n; i++) {
        int k = rnd.nextInt(span);
        if (rnd.nextInt(4) == 0) {
          bs.clear(k);
          ops.append('c').append(k).append(',');
        } else {
          bs.set(k);
          ops.append('s').append(k).append(',');
        }
      }
      StringBuilder res = new StringBuilder();
      int[] froms = {-2, -1, 0, rnd.nextInt(span + 70), rnd.nextInt(span + 70), 63, 64, 127, 128, span + 200};
      for (int from : froms) {
        String r;
        try {
          r = Integer.toString(bs.previousSetBit(from));
        } catch (IndexOutOfBoundsException e) {
          r = "!" + e.getMessage();
        }
        res.append(from).append('=').append(r).append(';');
      }
      out.append(ops).append('\t').append(res).append('\n');
    }
    System.out.print(out);
  }
}

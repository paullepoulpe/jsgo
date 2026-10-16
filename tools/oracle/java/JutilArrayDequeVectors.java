// Vectors for java.util.ArrayDeque (internal/jutil/array_deque.go): seeded
// random op scripts on ArrayDeque<Integer> (push/pop/peek, add/addFirst/
// addLast, poll/peekFirst/peekLast/pollFirst/pollLast, removeFirst/
// removeLast, getFirst/getLast, copy construction, addAll, clear, contains,
// descending iteration), each op followed by its result (or exception
// class), the size and toString().
// java -cp .cache/oracle-build JutilArrayDequeVectors > internal/jutil/testdata/array_deque.tsv
import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.Random;

public class JutilArrayDequeVectors {
  public static void main(String[] args) {
    Random rnd = new Random(20260415);
    StringBuilder out = new StringBuilder();
    for (int c = 0; c < 300; c++) {
      ArrayDeque<Integer> d = new ArrayDeque<>();
      StringBuilder ops = new StringBuilder();
      StringBuilder res = new StringBuilder();
      int n = 1 + rnd.nextInt(30);
      for (int i = 0; i < n; i++) {
        int op = rnd.nextInt(21);
        int v = rnd.nextInt(100);
        String r;
        try {
          switch (op) {
            case 0: case 1: ops.append("push").append(v); d.push(v); r = "-"; break;
            case 2: ops.append("pop"); r = String.valueOf(d.pop()); break;
            case 3: ops.append("peek"); r = String.valueOf(d.peek()); break;
            case 4: case 5: ops.append("add").append(v); r = String.valueOf(d.add(v)); break;
            case 6: ops.append("addFirst").append(v); d.addFirst(v); r = "-"; break;
            case 7: ops.append("addLast").append(v); d.addLast(v); r = "-"; break;
            case 8: ops.append("poll"); r = String.valueOf(d.poll()); break;
            case 9: ops.append("peekLast"); r = String.valueOf(d.peekLast()); break;
            case 10: ops.append("pollLast"); r = String.valueOf(d.pollLast()); break;
            case 11: ops.append("removeFirst"); r = String.valueOf(d.removeFirst()); break;
            case 12: ops.append("removeLast"); r = String.valueOf(d.removeLast()); break;
            case 13: ops.append("getFirst"); r = String.valueOf(d.getFirst()); break;
            case 14: ops.append("getLast"); r = String.valueOf(d.getLast()); break;
            case 15: ops.append("copy"); d = new ArrayDeque<>(d); r = "-"; break;
            case 16: {
              List<Integer> l = new ArrayList<>();
              int k = rnd.nextInt(4);
              ops.append("addAll");
              for (int j = 0; j < k; j++) { int x = rnd.nextInt(100); l.add(x); ops.append(j == 0 ? "" : ":").append(x); }
              r = String.valueOf(d.addAll(l));
              break;
            }
            case 17: ops.append("clear"); d.clear(); r = "-"; break;
            case 18: ops.append("contains").append(v); r = String.valueOf(d.contains(v)); break;
            case 19: {
              StringBuilder sb = new StringBuilder();
              for (Iterator<Integer> it = d.descendingIterator(); it.hasNext(); ) sb.append(it.next()).append(' ');
              r = sb.toString().trim(); ops.append("desc");
              break;
            }
            default: ops.append("pollFirst"); r = String.valueOf(d.pollFirst()); break;
          }
        } catch (RuntimeException e) {
          r = "!" + e.getClass().getSimpleName();
        }
        ops.append(',');
        res.append(r).append('/').append(d.size()).append('/').append(d).append(';');
      }
      out.append(ops).append('\t').append(res).append('\n');
    }
    System.out.print(out);
  }
}

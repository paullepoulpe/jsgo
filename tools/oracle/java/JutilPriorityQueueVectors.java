// Generates internal/jutil/testdata/priority_queue.tsv: java.util.PriorityQueue
// with a comparator, driven by random op scripts.
//   java -cp .cache/$CLOSURE_JAR:.cache/oracle-build JutilPriorityQueueVectors internal/jutil/testdata [seed]
//
// Elements are objects {key, id} (ids 0, 1, 2, ... in insertion order); the
// comparator looks at the key only, so many elements tie and the order in
// which ties leave the queue depends on the heap's exact sift sequence.
// Each row is "<initialCapacity>\t<op>;<op>;...\t<result>;<result>;...\t<ids in
// iteration order>". Ops: A<key> (add -> "t"), P (poll -> id or "null"),
// R (remove() -> id or "NSEE"), K (peek -> id or "null"), D<id> (remove(Object)
// of the element with that id, which may have left already -> "t"/"f"),
// C<id> (contains -> "t"/"f"), S (size -> n), X (clear -> "").
import java.io.PrintStream;
import java.nio.file.*;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.NoSuchElementException;
import java.util.PriorityQueue;
import java.util.Random;

public class JutilPriorityQueueVectors {
  static final class E {
    final int key;
    final int id;

    E(int key, int id) {
      this.key = key;
      this.id = id;
    }
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(args.length > 1 ? Long.parseLong(args[1]) : 20260415L);
    try (PrintStream ps =
        new PrintStream(Files.newOutputStream(out.resolve("priority_queue.tsv")))) {
      for (int k = 0; k < 600; k++) {
        int cap = 1 + r.nextInt(12);
        PriorityQueue<E> q = new PriorityQueue<>(cap, Comparator.comparingInt((E e) -> e.key));
        List<E> all = new ArrayList<>();
        int nops = 1 + r.nextInt(k < 300 ? 25 : 120);
        int keyRange = 1 + r.nextInt(k % 3 == 0 ? 4 : 50);
        StringBuilder ops = new StringBuilder();
        StringBuilder res = new StringBuilder();
        for (int i = 0; i < nops; i++) {
          if (i > 0) {
            ops.append(';');
            res.append(';');
          }
          int c = r.nextInt(20);
          if (c < 9) {
            E e = new E(r.nextInt(keyRange), all.size());
            all.add(e);
            q.add(e);
            ops.append('A').append(e.key);
            res.append('t');
          } else if (c < 12) {
            E e = q.poll();
            ops.append('P');
            res.append(e == null ? "null" : String.valueOf(e.id));
          } else if (c < 13) {
            ops.append('R');
            try {
              res.append(q.remove().id);
            } catch (NoSuchElementException ex) {
              res.append("NSEE");
            }
          } else if (c < 14) {
            E e = q.peek();
            ops.append('K');
            res.append(e == null ? "null" : String.valueOf(e.id));
          } else if (c < 17) {
            int id = all.isEmpty() ? 0 : r.nextInt(all.size());
            ops.append('D').append(id);
            res.append(!all.isEmpty() && q.remove(all.get(id)) ? 't' : 'f');
          } else if (c < 18) {
            int id = all.isEmpty() ? 0 : r.nextInt(all.size());
            ops.append('C').append(id);
            res.append(!all.isEmpty() && q.contains(all.get(id)) ? 't' : 'f');
          } else if (c < 19 || r.nextInt(8) != 0) {
            ops.append('S');
            res.append(q.size());
          } else {
            q.clear();
            ops.append('X');
          }
        }
        StringBuilder order = new StringBuilder();
        for (E e : q) {
          if (order.length() > 0) {
            order.append(',');
          }
          order.append(e.id);
        }
        ps.println(cap + "\t" + ops + "\t" + res + "\t" + order);
      }
    }
  }
}

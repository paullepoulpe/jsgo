// Generates internal/jutil/testdata/infra_collections.tsv: random operation
// scripts on java.util.PriorityQueue (with a comparator that has ties),
// Guava HashMultiset, ArrayListMultimap and
// MultimapBuilder.hashKeys().hashSetValues() multimaps, with the result of
// every operation.
//   java -cp .cache/$CLOSURE_JAR:.cache/oracle-build JutilInfraVectors internal/jutil/testdata
//
// Rows: "<kind>\t<op>;<op>;...\t<result>;<result>;...", kind one of
//   pq  ops a<key>.<id> (add) | p (poll: "<key>.<id>" or "null") | k (peek)
//   ms  ops a<w> (add) | r<w> (remove) | R<w>,<w>.. (removeAll) | c<w> (count) | h<w> (contains)
//   lm  ops p<k>=<v> (put) | R<k> (removeAll, in order) | g<k> (get) | c<k> (containsKey)
//   sm  ops p<k>=<v> (put) | P<k>=<v>,<v>.. (putAll) | R<k> (removeAll, sorted) | g<k> (get, sorted)
//       | c<k> (containsKey)
// and a final result "size=<n>" (plus, for pq, the remaining polls).
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ArrayListMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.HashMultiset;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ListMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.MultimapBuilder;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.Multiset;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.SetMultimap;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.*;

public class JutilInfraVectors {
  static final String[] W = {"a", "b", "c", "d", "e", "f"};

  static String sorted(Collection<String> c) {
    List<String> l = new ArrayList<>(c);
    Collections.sort(l);
    return l.toString();
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(4711);
    try (PrintStream ps =
        new PrintStream(
            Files.newOutputStream(out.resolve("infra_collections.tsv")),
            false,
            StandardCharsets.UTF_8)) {
      for (int k = 0; k < 200; k++) {
        PriorityQueue<int[]> q = new PriorityQueue<>((x, y) -> Integer.compare(x[0], y[0]));
        List<String> ops = new ArrayList<>();
        List<String> res = new ArrayList<>();
        int n = r.nextInt(60);
        int range = 1 + r.nextInt(8);
        for (int i = 0; i < n; i++) {
          int c = r.nextInt(10);
          if (c < 6) {
            int key = r.nextInt(range);
            ops.add("a" + key + "." + i);
            q.add(new int[] {key, i});
            res.add("true");
          } else if (c < 9) {
            ops.add("p");
            int[] e = q.poll();
            res.add(e == null ? "null" : e[0] + "." + e[1]);
          } else {
            ops.add("k");
            int[] e = q.peek();
            res.add(e == null ? "null" : e[0] + "." + e[1]);
          }
        }
        res.add("size=" + q.size());
        while (!q.isEmpty()) {
          int[] e = q.poll();
          res.add(e[0] + "." + e[1]);
        }
        ps.println("pq\t" + String.join(";", ops) + "\t" + String.join(";", res));
      }
      for (int k = 0; k < 150; k++) {
        Multiset<String> m = HashMultiset.create();
        List<String> ops = new ArrayList<>();
        List<String> res = new ArrayList<>();
        int n = r.nextInt(30);
        for (int i = 0; i < n; i++) {
          int c = r.nextInt(10);
          String w = W[r.nextInt(4)];
          if (c < 4) {
            ops.add("a" + w);
            res.add("" + m.add(w));
          } else if (c < 5) {
            ops.add("r" + w);
            res.add("" + m.remove(w));
          } else if (c < 6) {
            List<String> ws = new ArrayList<>();
            int cnt = r.nextInt(3);
            for (int j = 0; j < cnt; j++) {
              ws.add(W[r.nextInt(4)]);
            }
            ops.add("R" + String.join(",", ws));
            res.add("" + m.removeAll(ws));
          } else if (c < 8) {
            ops.add("c" + w);
            res.add("" + m.count(w));
          } else {
            ops.add("h" + w);
            res.add("" + m.contains(w));
          }
        }
        res.add("size=" + m.size());
        ps.println("ms\t" + String.join(";", ops) + "\t" + String.join(";", res));
      }
      for (int k = 0; k < 150; k++) {
        ListMultimap<String, String> m =
            k % 2 == 0 ? ArrayListMultimap.create() : MultimapBuilder.hashKeys().arrayListValues().build();
        List<String> ops = new ArrayList<>();
        List<String> res = new ArrayList<>();
        int n = r.nextInt(30);
        for (int i = 0; i < n; i++) {
          int c = r.nextInt(10);
          String key = W[r.nextInt(3)];
          if (c < 5) {
            String v = W[r.nextInt(4)];
            ops.add("p" + key + "=" + v);
            res.add("" + m.put(key, v));
          } else if (c < 7) {
            ops.add("R" + key);
            res.add(m.removeAll(key).toString());
          } else if (c < 9) {
            ops.add("g" + key);
            res.add(m.get(key).toString());
          } else {
            ops.add("c" + key);
            res.add("" + m.containsKey(key));
          }
        }
        res.add("size=" + m.size());
        ps.println("lm\t" + String.join(";", ops) + "\t" + String.join(";", res));
      }
      for (int k = 0; k < 150; k++) {
        SetMultimap<String, String> m = MultimapBuilder.hashKeys().hashSetValues().build();
        List<String> ops = new ArrayList<>();
        List<String> res = new ArrayList<>();
        int n = r.nextInt(30);
        for (int i = 0; i < n; i++) {
          int c = r.nextInt(10);
          String key = W[r.nextInt(3)];
          if (c < 4) {
            String v = W[r.nextInt(4)];
            ops.add("p" + key + "=" + v);
            res.add("" + m.put(key, v));
          } else if (c < 5) {
            List<String> vs = new ArrayList<>();
            int cnt = r.nextInt(3);
            for (int j = 0; j < cnt; j++) {
              vs.add(W[r.nextInt(4)]);
            }
            ops.add("P" + key + "=" + String.join(",", vs));
            res.add("" + m.putAll(key, vs));
          } else if (c < 7) {
            ops.add("R" + key);
            res.add(sorted(m.removeAll(key)));
          } else if (c < 9) {
            ops.add("g" + key);
            res.add(sorted(m.get(key)));
          } else {
            ops.add("c" + key);
            res.add("" + m.containsKey(key));
          }
        }
        res.add("size=" + m.size());
        ps.println("sm\t" + String.join(";", ops) + "\t" + String.join(";", res));
      }
    }
  }
}

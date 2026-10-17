// Generates internal/jutil/testdata/tree_set.tsv: java.util.TreeSet with a
// comparator, driven by random op scripts. Consumed by
// internal/jutil/tree_set_test.go.
//   java -cp .cache/oracle-build JutilTreeSetVectors internal/jutil/testdata [seed]
//
// Two kinds of rows:
//   v  a comparator that changes over time (ComposeWarningsGuard's pattern:
//      prio[a] - prio[b], then the element added most recently first). Which
//      elements remove() and contains() still find depends on the exact shape
//      of the red-black tree, so these rows pin down its algorithms.
//   n  Integer order over a wider range, with many removals.
// Row: kind \t nElems \t prio ("-" for n) \t ops \t results \t ascending \t
// descending \t size. Ops (space separated): a<g> (v: remove(g) then add(g),
// result "<removed><added>" as t/f; n: add(g) -> t/f), r<g> (remove -> t/f),
// c<g> (contains -> t/f), f (first -> element or NSEE), l (last -> element
// or NSEE), x (clear -> "-").
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.Random;
import java.util.TreeSet;
import java.util.function.Supplier;

public class JutilTreeSetVectors {
  static String tf(boolean b) {
    return b ? "t" : "f";
  }

  static String orNsee(Supplier<Integer> s) {
    try {
      return String.valueOf(s.get());
    } catch (NoSuchElementException e) {
      return "NSEE";
    }
  }

  static String join(Iterable<Integer> xs) {
    StringBuilder sb = new StringBuilder();
    for (int x : xs) {
      if (sb.length() > 0) {
        sb.append(',');
      }
      sb.append(x);
    }
    return sb.toString();
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    long seed = args.length > 1 ? Long.parseLong(args[1]) : 20260415L;
    Random r = new Random(seed);
    try (PrintStream ps =
        new PrintStream(Files.newOutputStream(out.resolve("tree_set.tsv")), false, StandardCharsets.UTF_8)) {
      for (int trial = 0; trial < 500; trial++) {
        boolean v = trial % 2 == 0;
        int nElems = v ? 4 + r.nextInt(40) : 2 + r.nextInt(r.nextBoolean() ? 30 : 400);
        int[] prio = new int[nElems];
        for (int i = 0; i < nElems; i++) {
          prio[i] = r.nextInt(1 + r.nextInt(6));
        }
        Map<Integer, Integer> order = new HashMap<>();
        TreeSet<Integer> set =
            v
                ? new TreeSet<>(
                    (a, b) -> {
                      int d = prio[a] - prio[b];
                      if (d != 0) {
                        return d;
                      }
                      return order.get(b) - order.get(a);
                    })
                : new TreeSet<>(Integer::compare);
        List<String> ops = new ArrayList<>();
        List<String> res = new ArrayList<>();
        int adds = 0;
        int nOps = 1 + r.nextInt(v ? 80 : 300);
        for (int k = 0; k < nOps; k++) {
          int g = r.nextInt(nElems);
          int c = r.nextInt(20);
          if (c < 4 && (!v || order.containsKey(g))) {
            ops.add("r" + g);
            res.add(tf(set.remove(g)));
          } else if (c < 6 && (!v || order.containsKey(g))) {
            ops.add("c" + g);
            res.add(tf(set.contains(g)));
          } else if (c == 6) {
            ops.add("f");
            res.add(orNsee(set::first));
          } else if (c == 7) {
            ops.add("l");
            res.add(orNsee(set::last));
          } else if (c == 8 && r.nextInt(10) == 0) {
            ops.add("x");
            set.clear();
            res.add("-");
          } else if (v) {
            ops.add("a" + g);
            adds++;
            order.put(g, adds);
            boolean removed = set.remove(g);
            res.add(tf(removed) + tf(set.add(g)));
          } else {
            ops.add("a" + g);
            res.add(tf(set.add(g)));
          }
        }
        StringBuilder pr = new StringBuilder();
        if (v) {
          for (int p : prio) {
            if (pr.length() > 0) {
              pr.append(',');
            }
            pr.append(p);
          }
        } else {
          pr.append('-');
        }
        ps.println(
            (v ? "v" : "n")
                + "\t"
                + nElems
                + "\t"
                + pr
                + "\t"
                + String.join(" ", ops)
                + "\t"
                + String.join(" ", res)
                + "\t"
                + join(set)
                + "\t"
                + join(set.descendingSet())
                + "\t"
                + set.size());
      }
    }
  }
}

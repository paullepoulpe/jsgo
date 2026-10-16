import com.google.javascript.jscomp.jarjar.com.google.common.collect.HashMultiset;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.LinkedHashMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ListMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.Multimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.MultimapBuilder;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.Multiset;
import java.io.PrintStream;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;

/**
 * Vectors for jutil.LinkedHashMultimap, jutil.ArrayListMultimap_Create and jutil.ObjectHashMultiset:
 * seeded random
 * operation scripts run against Guava (as shaded in the pinned closure-compiler.jar), each row
 * holding the script and the result of every operation.
 *
 * <p>Usage: java -cp closure-compiler.jar:oracle-build JutilMultiCollectionVectors > out.tsv
 *
 * <p>Row: kind TAB op;op;... TAB result|result|... where an op is name[:arg[:arg]].
 */
public final class JutilMultiCollectionVectors {
  private static final String[] KEYS = {"a", "b", "c", "d"};
  private static final String[] VALUES = {"x", "y", "z", "w"};

  public static void main(String[] args) {
    PrintStream out = System.out;
    Random rnd = new Random(20260415L);
    for (int i = 0; i < 300; i++) {
      linkedHashMultimap(out, rnd);
      listMultimap(out, rnd);
      hashMultiset(out, rnd);
    }
    out.flush();
  }

  private static String pick(Random rnd, String[] from) {
    return from[rnd.nextInt(from.length)];
  }

  private static void linkedHashMultimap(PrintStream out, Random rnd) {
    LinkedHashMultimap<String, String> mm = LinkedHashMultimap.create();
    List<String> ops = new ArrayList<>();
    List<String> results = new ArrayList<>();
    int n = 5 + rnd.nextInt(40);
    for (int i = 0; i < n; i++) {
      String k = pick(rnd, KEYS);
      String v = pick(rnd, VALUES);
      int op = rnd.nextInt(16);
      if (op < 6) {
        ops.add("put:" + k + ":" + v);
        results.add(String.valueOf(mm.put(k, v)));
      } else if (op == 6) {
        ops.add("rm:" + k + ":" + v);
        results.add(String.valueOf(mm.remove(k, v)));
      } else if (op == 7) {
        ops.add("rmAll:" + k);
        results.add(String.valueOf(mm.removeAll(k)));
      } else if (op == 8) {
        ops.add("get:" + k);
        results.add(String.valueOf(mm.get(k)));
      } else if (op == 9) {
        ops.add("ck:" + k);
        results.add(String.valueOf(mm.containsKey(k)));
      } else if (op == 10) {
        ops.add("ce:" + k + ":" + v);
        results.add(String.valueOf(mm.containsEntry(k, v)));
      } else if (op == 11) {
        ops.add("size");
        results.add(String.valueOf(mm.size()));
      } else if (op == 12) {
        ops.add("vals");
        results.add(String.valueOf(new ArrayList<>(mm.values())));
      } else if (op == 13) {
        ops.add("ents");
        results.add(entries(mm));
      } else if (op == 14) {
        ops.add("keys");
        results.add(String.valueOf(new ArrayList<>(mm.keySet())));
      } else {
        ops.add("str");
        results.add(mm.toString());
      }
    }
    row(out, "lhm", ops, results);
  }

  private static void listMultimap(PrintStream out, Random rnd) {
    ListMultimap<String, String> mm = MultimapBuilder.linkedHashKeys().arrayListValues().build();
    List<String> ops = new ArrayList<>();
    List<String> results = new ArrayList<>();
    int n = 5 + rnd.nextInt(40);
    for (int i = 0; i < n; i++) {
      String k = pick(rnd, KEYS);
      String v = pick(rnd, VALUES);
      int op = rnd.nextInt(12);
      if (op < 5) {
        ops.add("put:" + k + ":" + v);
        results.add(String.valueOf(mm.put(k, v)));
      } else if (op == 5) {
        ops.add("rmAll:" + k);
        results.add(String.valueOf(mm.removeAll(k)));
      } else if (op == 6) {
        ops.add("get:" + k);
        results.add(String.valueOf(mm.get(k)));
      } else if (op == 7) {
        ops.add("ck:" + k);
        results.add(String.valueOf(mm.containsKey(k)));
      } else if (op == 8) {
        ops.add("size");
        results.add(String.valueOf(mm.size()));
      } else if (op == 9) {
        ops.add("ents");
        results.add(entries(mm));
      } else if (op == 10) {
        ops.add("keys");
        results.add(String.valueOf(new ArrayList<>(mm.keySet())));
      } else {
        ops.add("str");
        results.add(mm.toString());
      }
    }
    row(out, "list", ops, results);
  }

  private static void hashMultiset(PrintStream out, Random rnd) {
    Multiset<String> ms = HashMultiset.create();
    List<String> ops = new ArrayList<>();
    List<String> results = new ArrayList<>();
    int n = 5 + rnd.nextInt(40);
    for (int i = 0; i < n; i++) {
      String e = pick(rnd, KEYS);
      int op = rnd.nextInt(7);
      if (op < 2) {
        ops.add("add:" + e);
        results.add(String.valueOf(ms.add(e)));
      } else if (op == 2) {
        int count = rnd.nextInt(3);
        ops.add("addn:" + e + ":" + count);
        results.add(String.valueOf(ms.add(e, count)));
      } else if (op == 3) {
        ops.add("rm:" + e);
        results.add(String.valueOf(ms.remove(e)));
      } else if (op == 4) {
        ops.add("cnt:" + e);
        results.add(String.valueOf(ms.count(e)));
      } else if (op == 5) {
        ops.add("has:" + e);
        results.add(String.valueOf(ms.contains(e)));
      } else {
        ops.add("size");
        results.add(String.valueOf(ms.size()));
      }
    }
    row(out, "ms", ops, results);
  }

  private static String entries(Multimap<String, String> mm) {
    List<String> parts = new ArrayList<>();
    mm.entries().forEach(e -> parts.add(e.getKey() + "=" + e.getValue()));
    return String.valueOf(parts);
  }

  private static void row(PrintStream out, String kind, List<String> ops, List<String> results) {
    out.println(kind + "\t" + String.join(";", ops) + "\t" + String.join("|", results));
  }

  private JutilMultiCollectionVectors() {}
}

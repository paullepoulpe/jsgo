// Generates internal/jutil/testdata/java_hash_{map,set,multimap}.tsv: random
// operation sequences on java.util.HashMap / HashSet (as the JDK running
// closure-compiler.jar implements them), the shaded Guava hash factories
// (Maps.newHashMap*, Sets.newHashSet*) and HashMultimap / MultimapBuilder
// .hashKeys().hashSetValues(), with the full iteration order dumped every
// few operations. Consumed by internal/jutil/java_hash_vectors_test.go,
// which replays the operations on JavaHashMap / JavaHashSet / HashMultimap.
//   java -cp .cache/closure-compiler.jar:.cache/oracle-build JutilJavaHashVectors internal/jutil/testdata [seed]
//
// Key pools mix short, long, non-ASCII and surrogate strings, families of
// equal-hash strings (concatenations of "Aa"/"BB") placed so that their tree
// bins treeify, untreeify and split on resize, strings whose hashes differ
// only in the high bits, Integer/Long keys (negative, bin-colliding, equal
// folded Long hashes), a Comparable value class with colliding hashes (CP), a
// non-Comparable value class with distinct bin-colliding hashes (NP), and
// mixed String/Integer/Long keys with equal hashes (tieBreakOrder by class
// name). No two distinct keys of one non-Comparable class share a hash, so no
// order depends on System.identityHashCode.
//
// Line format (tab-separated); keys and values are pool indices, Map values
// are Integers ("n" is null):
//   B <type> <pool> [<value pool>]   begin a scenario; pool entries are comma
//                                    separated: S hex UTF-16, I/L decimal,
//                                    P "a:b", N "h", O s<hex>|i<int>|l<long>
//   C <ctor> [args]                  construct the collection under test
//   <op> args... <result>            one operation and its result
//   D <dump>                         the iteration order
//   E <size> <hashCode> <toString>   end ("-" when toString is not checked)
import com.google.javascript.jscomp.jarjar.com.google.common.collect.HashMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.LinkedHashMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.Maps;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.MultimapBuilder;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.SetMultimap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.Sets;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collection;
import java.util.ConcurrentModificationException;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.Set;

public class JutilJavaHashVectors {
  /** A Comparable value class whose hashCode collides easily. */
  static final class CP implements Comparable<CP> {
    final int a;
    final int b;

    CP(int a, int b) {
      this.a = a;
      this.b = b;
    }

    @Override
    public boolean equals(Object o) {
      return o instanceof CP c && c.a == a && c.b == b;
    }

    @Override
    public int hashCode() {
      return 31 * a + b;
    }

    @Override
    public int compareTo(CP o) {
      return a != o.a ? Integer.compare(a, o.a) : Integer.compare(b, o.b);
    }

    @Override
    public String toString() {
      return "CP(" + a + "," + b + ")";
    }
  }

  /** A non-Comparable value class whose hashCode is its field. */
  static final class NP {
    final int h;

    NP(int h) {
      this.h = h;
    }

    @Override
    public boolean equals(Object o) {
      return o instanceof NP n && n.h == h;
    }

    @Override
    public int hashCode() {
      return h;
    }

    @Override
    public String toString() {
      return "NP(" + h + ")";
    }
  }

  static Random r;

  static String hex(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      sb.append(String.format("%04x", (int) s.charAt(i)));
    }
    return sb.toString();
  }

  static String esc(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c < 0x20 || c >= 0x7f || c == '\\') {
        sb.append(String.format("\\u%04x", (int) c));
      } else {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  static int spread(int h) {
    return h ^ (h >>> 16);
  }

  // ---------------------------------------------------------------------------
  // Key pools

  /** A pool of distinct keys and their encodings. */
  static final class Pool {
    final String type;
    final List<Object> keys = new ArrayList<>();
    final List<String> enc = new ArrayList<>();
    final Map<Object, Integer> index = new HashMap<>();

    Pool(String type) {
      this.type = type;
    }

    void add(Object k) {
      if (index.containsKey(k)) {
        return;
      }
      index.put(k, keys.size());
      keys.add(k);
      enc.add(encode(type, k));
    }

    /** A fresh instance equal to key i (value classes), or key i itself. */
    Object key(int i) {
      Object k = keys.get(i);
      if (k instanceof CP c) {
        return new CP(c.a, c.b);
      }
      if (k instanceof NP n) {
        return new NP(n.h);
      }
      if (k instanceof String s) {
        return new String(s);
      }
      return k;
    }

    int idx(Object k) {
      return index.get(k);
    }

    int pick() {
      // Skewed toward a hot subset, so that keys recur.
      int n = keys.size();
      return r.nextInt(4) == 0 ? r.nextInt(Math.min(n, 8)) : r.nextInt(n);
    }

    void shuffle() {
      List<Integer> order = new ArrayList<>();
      for (int i = 0; i < keys.size(); i++) {
        order.add(i);
      }
      java.util.Collections.shuffle(order, r);
      List<Object> ks = new ArrayList<>();
      for (int i : order) {
        ks.add(keys.get(i));
      }
      keys.clear();
      enc.clear();
      index.clear();
      for (Object k : ks) {
        add(k);
      }
    }

    String encoded() {
      return String.join(",", enc);
    }
  }

  static String encode(String type, Object k) {
    switch (type) {
      case "S":
        return hex((String) k);
      case "I":
      case "L":
        return k.toString();
      case "P":
        return ((CP) k).a + ":" + ((CP) k).b;
      case "N":
        return Integer.toString(((NP) k).h);
      case "O":
        if (k instanceof String s) {
          return "s" + hex(s);
        } else if (k instanceof Integer i) {
          return "i" + i;
        } else {
          return "l" + k;
        }
      default:
        throw new AssertionError(type);
    }
  }

  static final String ALPHA = "abcxyzAB019_$";

  static String randomString(int minLen, int maxLen, int flavor) {
    int len = minLen + r.nextInt(maxLen - minLen + 1);
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < len; i++) {
      int f = flavor == 0 ? 0 : r.nextInt(6);
      switch (f) {
        case 1 -> sb.append((char) (0xa0 + r.nextInt(0x60)));
        case 2 -> sb.append((char) (0x4e00 + r.nextInt(0x100)));
        case 3 -> sb.appendCodePoint(0x1f600 + r.nextInt(0x40));
        case 4 -> sb.append((char) (0xd800 + r.nextInt(0x800)));
        default -> sb.append(ALPHA.charAt(r.nextInt(ALPHA.length())));
      }
    }
    return sb.toString();
  }

  static String blocks(int n, int bitsSel) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < n; i++) {
      sb.append(((bitsSel >> i) & 1) == 0 ? "Aa" : "BB");
    }
    return sb.toString();
  }

  /** count members of the equal-hash family prefix + {Aa,BB}^n. */
  static void addFamily(Pool p, String prefix, int count) {
    int n = 1;
    while ((1 << n) < count) {
      n++;
    }
    List<Integer> sel = new ArrayList<>();
    for (int i = 0; i < (1 << n); i++) {
      sel.add(i);
    }
    java.util.Collections.shuffle(sel, r);
    for (int i = 0; i < count; i++) {
      p.add(prefix + blocks(n, sel.get(i)));
    }
  }

  /**
   * Prefixes whose families (with n blocks) share the bin of a 64-slot table but are
   * separated at 128 (and maybe 256) slots.
   */
  static List<String> splitPrefixes(int n, int want) {
    Map<Integer, List<String>> byBin = new HashMap<>();
    while (true) {
      String p = randomString(1, 3, 0);
      int h = spread((p + blocks(n, 0)).hashCode());
      List<String> l = byBin.computeIfAbsent(h & 63, x -> new ArrayList<>());
      if (l.contains(p)) {
        continue;
      }
      boolean lo = false;
      boolean hi = false;
      for (String q : l) {
        int qh = spread((q + blocks(n, 0)).hashCode());
        if ((qh & 64) == 0) {
          lo = true;
        } else {
          hi = true;
        }
        if (qh == h) {
          lo = hi = false;
          break;
        }
      }
      l.add(p);
      if (l.size() >= want && lo && hi) {
        return l;
      }
    }
  }

  /** Strings whose hashCodes share the low 16 bits. */
  static List<String> highBitStrings(int want) {
    Map<Integer, List<String>> byLow = new HashMap<>();
    while (true) {
      String s = randomString(2, 6, 0);
      List<String> l = byLow.computeIfAbsent(s.hashCode() & 0xffff, x -> new ArrayList<>());
      boolean dup = false;
      for (String q : l) {
        dup |= q.hashCode() == s.hashCode();
      }
      if (!dup) {
        l.add(s);
      }
      if (l.size() >= want) {
        return l;
      }
    }
  }

  static Pool stringPool(int design) {
    Pool p = new Pool("S");
    switch (design) {
      case 0 -> {
        // Short and long, ASCII, non-ASCII and surrogate strings.
        int n = 10 + r.nextInt(120);
        for (int i = 0; i < n; i++) {
          p.add(r.nextInt(6) == 0 ? randomString(8, 30, 1) : randomString(0, 4, r.nextInt(3)));
        }
      }
      case 1 -> {
        // One equal-hash family (one tree bin) plus filler.
        addFamily(p, randomString(0, 2, 0), 9 + r.nextInt(24));
        int n = r.nextInt(80);
        for (int i = 0; i < n; i++) {
          p.add(randomString(1, 5, r.nextInt(2)));
        }
      }
      case 2 -> {
        // Families sharing a bin at 64 slots and splitting at 128/256, plus
        // enough filler to grow the table there.
        int n = 4;
        for (String pre : splitPrefixes(n, 2 + r.nextInt(3))) {
          addFamily(p, pre, 1 + r.nextInt(14));
        }
        int f = 40 + r.nextInt(120);
        for (int i = 0; i < f; i++) {
          p.add(randomString(1, 6, 0));
        }
      }
      case 3 -> {
        // Hashes that differ only in their high 16 bits.
        for (String s : highBitStrings(3 + r.nextInt(5))) {
          p.add(s);
        }
        for (String s : highBitStrings(3 + r.nextInt(5))) {
          p.add(s);
        }
        int f = r.nextInt(40);
        for (int i = 0; i < f; i++) {
          p.add(randomString(1, 4, 0));
        }
      }
      default -> {
        // Several small families and everything else mixed.
        for (int i = 0; i < 1 + r.nextInt(4); i++) {
          addFamily(p, randomString(0, 3, 0), 2 + r.nextInt(12));
        }
        for (String s : highBitStrings(4)) {
          p.add(s);
        }
        int f = r.nextInt(100);
        for (int i = 0; i < f; i++) {
          p.add(randomString(0, 12, r.nextInt(2)));
        }
      }
    }
    p.shuffle();
    return p;
  }

  static Pool intPool() {
    Pool p = new Pool("I");
    int n = 5 + r.nextInt(120);
    for (int i = 0; i < n; i++) {
      switch (r.nextInt(5)) {
        case 0 -> p.add(r.nextInt());
        case 1 -> p.add(r.nextInt(64) - 32);
        case 2 -> p.add(r.nextInt(64) << 20); // bin-colliding, distinct hashes
        case 3 -> p.add((r.nextInt(16) << 16) + 7); // differ in high bits only
        default -> p.add(-(r.nextInt(1000) << 6));
      }
    }
    p.shuffle();
    return p;
  }

  static Pool longPool() {
    Pool p = new Pool("L");
    int n = 5 + r.nextInt(100);
    long c = r.nextInt(64);
    for (int i = 0; i < n; i++) {
      switch (r.nextInt(4)) {
        case 0 -> p.add(r.nextLong());
        case 1 -> {
          // Equal folded hash: high ^ low == c.
          long a = r.nextInt() & 0xffffffffL;
          p.add((a << 32) | ((a ^ c) & 0xffffffffL));
        }
        case 2 -> p.add((long) (r.nextInt(200) - 100));
        default -> p.add(-((long) r.nextInt(50) << 32));
      }
    }
    p.shuffle();
    return p;
  }

  static Pool cpPool() {
    Pool p = new Pool("P");
    int groups = 1 + r.nextInt(4);
    for (int g = 0; g < groups; g++) {
      int t = r.nextInt(2000) - 1000;
      int m = 2 + r.nextInt(16);
      for (int a = 0; a < m; a++) {
        p.add(new CP(a - 4, t - 31 * (a - 4)));
      }
    }
    int f = r.nextInt(60);
    for (int i = 0; i < f; i++) {
      p.add(new CP(r.nextInt(100), r.nextInt(100)));
    }
    p.shuffle();
    return p;
  }

  static Pool npPool() {
    Pool p = new Pool("N");
    int n = 5 + r.nextInt(100);
    for (int i = 0; i < n; i++) {
      switch (r.nextInt(3)) {
        case 0 -> p.add(new NP(r.nextInt(256) << 22)); // all in bin 0 below 1024 slots
        case 1 -> p.add(new NP(r.nextInt()));
        default -> p.add(new NP((r.nextInt(64) << 16) | 3));
      }
    }
    p.shuffle();
    return p;
  }

  static Pool mixedPool() {
    Pool p = new Pool("O");
    int fams = 1 + r.nextInt(3);
    for (int f = 0; f < fams; f++) {
      String pre = randomString(0, 2, 0);
      Pool tmp = new Pool("S");
      addFamily(tmp, pre, 4 + r.nextInt(10));
      int h = ((String) tmp.keys.get(0)).hashCode();
      for (Object k : tmp.keys) {
        p.add(k);
      }
      p.add(h);
      p.add((long) h & 0xffffffffL);
    }
    int n = r.nextInt(50);
    for (int i = 0; i < n; i++) {
      switch (r.nextInt(3)) {
        case 0 -> p.add(randomString(1, 3, 0));
        case 1 -> p.add(r.nextInt(100));
        default -> p.add((long) r.nextInt(100));
      }
    }
    p.shuffle();
    return p;
  }

  static Pool anyPool() {
    switch (r.nextInt(10)) {
      case 0:
        return intPool();
      case 1:
        return longPool();
      case 2:
        return cpPool();
      case 3:
        return npPool();
      case 4:
        return mixedPool();
      default:
        return stringPool(r.nextInt(5));
    }
  }

  // ---------------------------------------------------------------------------
  // Output

  static PrintStream ps;

  static void line(Object... cols) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < cols.length; i++) {
      if (i > 0) {
        sb.append('\t');
      }
      sb.append(cols[i]);
    }
    ps.println(sb);
  }

  static String v(Object o) {
    return o == null ? "n" : o.toString();
  }

  static String b(boolean x) {
    return x ? "1" : "0";
  }

  static String mapDump(Pool p, Map<Object, Integer> m) {
    StringBuilder sb = new StringBuilder();
    for (Map.Entry<Object, Integer> e : m.entrySet()) {
      if (sb.length() > 0) {
        sb.append(',');
      }
      sb.append(p.idx(e.getKey())).append('=').append(v(e.getValue()));
    }
    return sb.toString();
  }

  static String setDump(Pool p, Collection<Object> s) {
    StringBuilder sb = new StringBuilder();
    for (Object e : s) {
      if (sb.length() > 0) {
        sb.append(',');
      }
      sb.append(p.idx(e));
    }
    return sb.toString();
  }

  /** k=v,... over distinct random keys, with the map to build from it. */
  static String srcList(Pool p, int n, LinkedHashMap<Object, Integer> out, int[] val) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < n; i++) {
      int k = p.pick();
      if (out.containsKey(p.keys.get(k))) {
        continue;
      }
      int x = val[0]++;
      out.put(p.key(k), x);
      if (sb.length() > 0) {
        sb.append(',');
      }
      sb.append(k).append('=').append(x);
    }
    return sb.toString();
  }

  // ---------------------------------------------------------------------------
  // HashMap scenarios

  static Integer mergeFn(Integer a, Integer b) {
    int s = a + b;
    return s % 5 == 0 ? null : s;
  }

  static HashMap<Object, Integer> newMap(Pool p, int[] val) {
    switch (r.nextInt(11)) {
      case 0: {
        int n = r.nextInt(r.nextBoolean() ? 4 : 200);
        line("C", "c", n);
        return new HashMap<>(n);
      }
      case 1: {
        int n = r.nextInt(70);
        float lf = new float[] {0.25f, 0.5f, 1f, 2f, 0.9f}[r.nextInt(5)];
        line("C", "cl", n, lf);
        return new HashMap<>(n, lf);
      }
      case 2: {
        LinkedHashMap<Object, Integer> src = new LinkedHashMap<>();
        String s = srcList(p, r.nextInt(60), src, val);
        line("C", "m", s);
        return new HashMap<>(src);
      }
      case 3:
        line("C", "gd");
        return Maps.newHashMap();
      case 4: {
        LinkedHashMap<Object, Integer> src = new LinkedHashMap<>();
        String s = srcList(p, r.nextInt(60), src, val);
        line("C", "gm", s);
        return Maps.newHashMap(src);
      }
      case 5:
      case 6: {
        int n = r.nextInt(r.nextBoolean() ? 5 : 120);
        line("C", "ge", n);
        return Maps.newHashMapWithExpectedSize(n);
      }
      case 7: {
        int n = r.nextInt(100);
        line("C", "nh", n);
        return HashMap.newHashMap(n);
      }
      default:
        line("C", "d");
        return new HashMap<>();
    }
  }

  @SuppressWarnings("unchecked")
  static void mapScenario(Pool p, int nOps, int phase) {
    line("B", p.type, p.encoded());
    int[] val = {1};
    HashMap<Object, Integer> m = newMap(p, val);
    int untilDump = 1 + r.nextInt(6);
    for (int op = 0; op < nOps; op++) {
      int k = p.pick();
      Object key = p.key(k);
      Integer x = val[0]++;
      int w = r.nextInt(100);
      // phase 0: grow; 1: churn; 2: shrink.
      int putW = phase == 0 ? 60 : phase == 1 ? 35 : 15;
      if (w < putW) {
        if (r.nextInt(40) == 0) {
          x = null;
        }
        line("p", k, v(x), v(m.put(key, x)));
      } else if (w < putW + 25) {
        line("r", k, v(m.remove(key)));
      } else {
        switch (r.nextInt(20)) {
          case 0 -> line("pi", k, x, v(m.putIfAbsent(key, x)));
          case 1 -> {
            Integer y = r.nextInt(4) == 0 ? null : x;
            line("ca", k, v(y), v(m.computeIfAbsent(key, kk -> y)));
          }
          case 2 -> {
            Integer y = r.nextInt(3) == 0 ? null : x;
            line("cp", k, v(y), v(m.computeIfPresent(key, (kk, old) -> y)));
          }
          case 3 -> {
            Integer y = r.nextInt(3) == 0 ? null : x;
            line("co", k, v(y), v(m.compute(key, (kk, old) -> y)));
          }
          case 4 -> line("mg", k, x, v(m.merge(key, x, JutilJavaHashVectors::mergeFn)));
          case 5 -> line("g", k, v(m.get(key)));
          case 6 -> line("ck", k, b(m.containsKey(key)));
          case 7 -> {
            Integer y = m.get(key);
            if (r.nextBoolean() || y == null) {
              y = x;
            }
            line("re", k, v(y), b(m.remove(key, y)));
          }
          case 8 -> line("kr", k, b(m.keySet().remove(key)));
          case 9 -> {
            LinkedHashMap<Object, Integer> src = new LinkedHashMap<>();
            String s = srcList(p, r.nextInt(r.nextBoolean() ? 8 : 80), src, val);
            m.putAll(src);
            line("pa", s);
          }
          case 10 -> {
            int mod = 1 + r.nextInt(4);
            int rem = r.nextInt(mod);
            StringBuilder vis = new StringBuilder();
            Iterator<Map.Entry<Object, Integer>> it = m.entrySet().iterator();
            for (int i = 0; it.hasNext(); i++) {
              Map.Entry<Object, Integer> e = it.next();
              if (vis.length() > 0) {
                vis.append(',');
              }
              vis.append(p.idx(e.getKey()));
              if (i % mod == rem) {
                it.remove();
              }
            }
            line("ir", mod, rem, vis);
          }
          case 11 -> {
            int pos = r.nextInt(m.size() + 1);
            String res = "ok";
            int i = 0;
            try {
              for (Object kk : m.keySet()) {
                if (i == pos) {
                  m.put(key, -1);
                }
                i++;
              }
            } catch (ConcurrentModificationException e) {
              res = "cme" + i;
            }
            line("cme", k, pos, res);
          }
          case 12 -> {
            line("cpy");
            m = new HashMap<>(m);
          }
          case 13 -> {
            line("cln");
            m = (HashMap<Object, Integer>) m.clone();
          }
          case 14 -> {
            if (r.nextInt(4) == 0) {
              line("x");
              m.clear();
            }
          }
          default -> line("p", k, x, v(m.put(key, x)));
        }
      }
      if (--untilDump == 0) {
        line("D", mapDump(p, m));
        untilDump = 1 + r.nextInt(6) + m.size() / 12;
      }
    }
    line("D", mapDump(p, m));
    line("E", m.size(), m.hashCode(), m.size() <= 16 ? esc(m.toString()) : "-");
  }

  /** Inserts the whole pool, churns, then removes most of it again. */
  static void treeScenario(Pool p) {
    line("B", p.type, p.encoded());
    HashMap<Object, Integer> m = r.nextBoolean() ? new HashMap<>() : Maps.newHashMap();
    line("C", "d");
    int val = 1;
    int n = p.keys.size();
    for (int i = 0; i < n; i++) {
      Integer x = val++;
      line("p", i, x, v(m.put(p.key(i), x)));
      if (i % (4 + n / 8) == 3) {
        line("D", mapDump(p, m));
      }
    }
    line("D", mapDump(p, m));
    for (int j = 0; j < n; j++) {
      int i = r.nextInt(n);
      if (r.nextInt(3) == 0) {
        Integer x = val++;
        line("ca", i, x, v(m.computeIfAbsent(p.key(i), kk -> x)));
      } else {
        line("r", i, v(m.remove(p.key(i))));
      }
      if (j % (3 + m.size() / 10) == 2) {
        line("D", mapDump(p, m));
      }
    }
    if (r.nextBoolean()) {
      int mod = 2 + r.nextInt(3);
      int rem = r.nextInt(mod);
      StringBuilder vis = new StringBuilder();
      Iterator<Object> it = m.keySet().iterator();
      for (int i = 0; it.hasNext(); i++) {
        Object e = it.next();
        if (vis.length() > 0) {
          vis.append(',');
        }
        vis.append(p.idx(e));
        if (i % mod == rem) {
          it.remove();
        }
      }
      line("ir", mod, rem, vis);
      line("D", mapDump(p, m));
      for (int i = 0; i < n; i++) {
        Integer x = val++;
        line("p", i, x, v(m.put(p.key(i), x)));
      }
      line("D", mapDump(p, m));
    }
    line("E", m.size(), m.hashCode(), "-");
  }

  // ---------------------------------------------------------------------------
  // HashSet scenarios

  static String list(Pool p, int n, List<Object> out) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < n; i++) {
      int k = p.pick();
      out.add(p.key(k));
      if (sb.length() > 0) {
        sb.append(',');
      }
      sb.append(k);
    }
    return sb.toString();
  }

  static HashSet<Object> newSet(Pool p) {
    List<Object> src = new ArrayList<>();
    switch (r.nextInt(11)) {
      case 0: {
        int n = r.nextInt(100);
        line("C", "c", n);
        return new HashSet<>(n);
      }
      case 1: {
        int n = r.nextInt(60);
        float lf = new float[] {0.5f, 1f, 3f}[r.nextInt(3)];
        line("C", "cl", n, lf);
        return new HashSet<>(n, lf);
      }
      case 2: {
        line("C", "col", list(p, r.nextInt(60), src));
        return new HashSet<>(src);
      }
      case 3:
        line("C", "g0");
        return Sets.newHashSet();
      case 4: {
        line("C", "gv", list(p, r.nextInt(40), src));
        return Sets.newHashSet(src.toArray());
      }
      case 5: {
        line("C", "gi", list(p, r.nextInt(40), src));
        Iterable<Object> it = () -> src.iterator();
        return Sets.newHashSet(it);
      }
      case 6: {
        line("C", "gc", list(p, r.nextInt(40), src));
        return Sets.newHashSet((Iterable<Object>) src);
      }
      case 7: {
        int n = r.nextInt(r.nextBoolean() ? 5 : 100);
        line("C", "ge", n);
        return Sets.newHashSetWithExpectedSize(n);
      }
      case 8: {
        int n = r.nextInt(100);
        line("C", "nh", n);
        return HashSet.newHashSet(n);
      }
      default:
        line("C", "d");
        return new HashSet<>();
    }
  }

  static void setScenario(Pool p, int nOps) {
    line("B", p.type, p.encoded());
    HashSet<Object> s = newSet(p);
    int untilDump = 1 + r.nextInt(6);
    for (int op = 0; op < nOps; op++) {
      int k = p.pick();
      Object key = p.key(k);
      int w = r.nextInt(100);
      if (w < 50) {
        line("a", k, b(s.add(key)));
      } else if (w < 70) {
        line("r", k, b(s.remove(key)));
      } else {
        List<Object> c = new ArrayList<>();
        switch (r.nextInt(9)) {
          case 0 -> line("h", k, b(s.contains(key)));
          case 1 -> {
            String l = list(p, r.nextInt(30), c);
            line("aa", l, b(s.addAll(c)));
          }
          case 2 -> {
            String l = list(p, r.nextInt(r.nextBoolean() ? 6 : 60), c);
            line("ra", l, b(s.removeAll(c)));
          }
          case 3 -> {
            String l = list(p, r.nextInt(60), c);
            line("ta", l, b(s.retainAll(c)));
          }
          case 4 -> {
            int mod = 1 + r.nextInt(4);
            int rem = r.nextInt(mod);
            StringBuilder vis = new StringBuilder();
            Iterator<Object> it = s.iterator();
            for (int i = 0; it.hasNext(); i++) {
              Object e = it.next();
              if (vis.length() > 0) {
                vis.append(',');
              }
              vis.append(p.idx(e));
              if (i % mod == rem) {
                it.remove();
              }
            }
            line("ir", mod, rem, vis);
          }
          case 5 -> {
            line("cpy");
            s = new HashSet<>(s);
          }
          case 6 -> {
            line("gcp");
            s = Sets.newHashSet(s);
          }
          case 7 -> {
            if (r.nextInt(3) == 0) {
              line("x");
              s.clear();
            }
          }
          default -> line("a", k, b(s.add(key)));
        }
      }
      if (--untilDump == 0) {
        line("D", setDump(p, s));
        untilDump = 1 + r.nextInt(6);
      }
    }
    line("D", setDump(p, s));
    line("E", s.size(), s.hashCode(), s.size() <= 16 ? esc(s.toString()) : "-");
  }

  // ---------------------------------------------------------------------------
  // HashMultimap scenarios

  static String entriesDump(Pool kp, Pool vp, SetMultimap<Object, Object> mm) {
    StringBuilder sb = new StringBuilder();
    for (Map.Entry<Object, Object> e : mm.entries()) {
      if (sb.length() > 0) {
        sb.append(',');
      }
      sb.append(kp.idx(e.getKey())).append(':').append(vp.idx(e.getValue()));
    }
    return sb.toString();
  }

  static void multimapScenario(Pool kp, Pool vp, int nOps) {
    line("B", kp.type, kp.encoded(), vp.type, vp.encoded());
    SetMultimap<Object, Object> mm;
    switch (r.nextInt(6)) {
      case 0 -> {
        int ek = r.nextInt(40);
        int ev = r.nextInt(6);
        line("C", "hk", ek, ev);
        mm = HashMultimap.create(ek, ev);
      }
      case 1 -> {
        line("C", "mb");
        mm = MultimapBuilder.hashKeys().hashSetValues().build();
      }
      case 2 -> {
        int ek = r.nextInt(40);
        int ev = r.nextInt(6);
        line("C", "mbn", ek, ev);
        mm = MultimapBuilder.hashKeys(ek).hashSetValues(ev).build();
      }
      case 3 -> {
        LinkedHashMultimap<Object, Object> src = LinkedHashMultimap.create();
        StringBuilder sb = new StringBuilder();
        int n = r.nextInt(40);
        for (int i = 0; i < n; i++) {
          int k = kp.pick();
          int x = vp.pick();
          src.put(kp.key(k), vp.key(x));
          if (sb.length() > 0) {
            sb.append(',');
          }
          sb.append(k).append(':').append(x);
        }
        line("C", "hm", sb);
        mm = HashMultimap.create(src);
      }
      default -> {
        line("C", "hc");
        mm = HashMultimap.create();
      }
    }
    int untilDump = 1 + r.nextInt(6);
    for (int op = 0; op < nOps; op++) {
      int k = kp.pick();
      int x = vp.pick();
      Object key = kp.key(k);
      Object value = vp.key(x);
      int w = r.nextInt(100);
      if (w < 55) {
        line("p", k, x, b(mm.put(key, value)));
      } else if (w < 70) {
        line("r", k, x, b(mm.remove(key, value)));
      } else {
        switch (r.nextInt(8)) {
          case 0 -> line("ra", k, setDump(vp, mm.removeAll(key)));
          case 1 -> {
            List<Object> c = new ArrayList<>();
            String l = list(vp, r.nextInt(8), c);
            line("pa", k, l, b(mm.putAll(key, c)));
          }
          case 2 -> line("ce", k, x, b(mm.containsEntry(key, value)));
          case 3 -> line("ck", k, b(mm.containsKey(key)));
          case 4 -> line("g", k, setDump(vp, mm.get(key)));
          case 5 -> {
            if (r.nextInt(4) == 0) {
              line("x");
              mm.clear();
            }
          }
          default -> line("p", k, x, b(mm.put(key, value)));
        }
      }
      if (--untilDump == 0) {
        line("D", entriesDump(kp, vp, mm));
        untilDump = 1 + r.nextInt(6);
      }
    }
    line("D", entriesDump(kp, vp, mm));
    line("E", mm.size(), mm.hashCode(), mm.size() <= 16 ? esc(mm.toString()) : "-");
  }

  static Pool valuePool() {
    if (r.nextInt(3) == 0) {
      return intPool();
    }
    Pool p = new Pool("S");
    if (r.nextBoolean()) {
      addFamily(p, randomString(0, 2, 0), 2 + r.nextInt(14));
    }
    int n = 3 + r.nextInt(40);
    for (int i = 0; i < n; i++) {
      p.add(randomString(0, 6, r.nextInt(2)));
    }
    p.shuffle();
    return p;
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    long seed = args.length > 1 ? Long.parseLong(args[1]) : 20260415L;
    r = new Random(seed);

    try (PrintStream p =
        new PrintStream(
            Files.newOutputStream(out.resolve("java_hash_map.tsv")), false, StandardCharsets.UTF_8)) {
      ps = p;
      for (int i = 0; i < 170; i++) {
        mapScenario(anyPool(), 10 + r.nextInt(110), r.nextInt(3));
      }
      // Tree bins: every string design with families, and the other key kinds.
      for (int i = 0; i < 49; i++) {
        Pool pool =
            switch (i % 7) {
              case 0, 1 -> stringPool(1);
              case 2, 3 -> stringPool(2);
              case 4 -> cpPool();
              case 5 -> npPool();
              default -> i % 2 == 0 ? mixedPool() : longPool();
            };
        treeScenario(pool);
      }
    }

    try (PrintStream p =
        new PrintStream(
            Files.newOutputStream(out.resolve("java_hash_set.tsv")), false, StandardCharsets.UTF_8)) {
      ps = p;
      for (int i = 0; i < 100; i++) {
        setScenario(r.nextInt(4) == 0 ? intPool() : stringPool(r.nextInt(5)), 10 + r.nextInt(100));
      }
    }

    try (PrintStream p =
        new PrintStream(
            Files.newOutputStream(out.resolve("java_hash_multimap.tsv")),
            false,
            StandardCharsets.UTF_8)) {
      ps = p;
      for (int i = 0; i < 90; i++) {
        Pool kp = r.nextInt(5) == 0 ? intPool() : stringPool(r.nextInt(5));
        multimapScenario(kp, valuePool(), 10 + r.nextInt(100));
      }
    }
  }
}

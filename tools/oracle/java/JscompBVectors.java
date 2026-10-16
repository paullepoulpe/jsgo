// Oracle tool: writes the Java reference vectors for the jscomp-b port
// (name generators, graph algorithms, union-find, VariableMap, message ids,
// placeholder names, Xid, PathUtil, SourceCodeEscapers, CodeConsumer.addNumber)
// into internal/jscomp/testdata/jscomp_b.tsv. Every case records its inputs, so
// internal/jscomp/jscomp_b_differential_test.go replays it without a seeded
// generator of its own.
//
// It lives in package com.google.javascript.jscomp to reach package-private
// API (DefaultNameGenerator.favors, JsMessage.Hash, PlaceholderReference
// .createForJsName, CodeConsumer):
//   java -cp .cache/closure-compiler.jar:.cache/oracle-build \
//       com.google.javascript.jscomp.JscompBVectors internal/jscomp/testdata/jscomp_b.tsv
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.deps.PathUtil;
import com.google.javascript.jscomp.deps.SourceCodeEscapers;
import com.google.javascript.jscomp.graph.DiGraph;
import com.google.javascript.jscomp.graph.DominatorTree;
import com.google.javascript.jscomp.graph.Graph;
import com.google.javascript.jscomp.graph.GraphColoring;
import com.google.javascript.jscomp.graph.LinkedDirectedGraph;
import com.google.javascript.jscomp.graph.LinkedUndirectedGraph;
import com.google.javascript.jscomp.graph.LowestCommonAncestorFinder;
import com.google.javascript.jscomp.graph.StandardUnionFind;
import java.io.PrintStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.Set;

public class JscompBVectors {
  static PrintStream out;

  // Strings are written as "u" + the UTF-16 code units, with every unit outside
  // SAFE written as a backslash and four hex digits; null is "N".
  static final String SAFE =
      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 ._-[]()'\"<>#@?+*$%&^~{}";

  static String u(String s) {
    if (s == null) {
      return "N";
    }
    StringBuilder sb = new StringBuilder("u");
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (SAFE.indexOf(c) >= 0) {
        sb.append(c);
      } else {
        sb.append(String.format("\\%04x", (int) c));
      }
    }
    return sb.toString();
  }

  static String hex(byte[] b) {
    StringBuilder sb = new StringBuilder("x");
    for (byte x : b) {
      sb.append(String.format("%02x", x & 0xff));
    }
    return sb.toString();
  }

  static String exc(Throwable e) {
    return "!" + u(e.toString());
  }

  static void row(Object... cols) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < cols.length; i++) {
      if (i > 0) {
        sb.append('\t');
      }
      sb.append(cols[i]);
    }
    out.println(sb);
  }

  static String pick(Random r, String alphabet, int len) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < len; i++) {
      sb.append(alphabet.charAt(r.nextInt(alphabet.length())));
    }
    return sb.toString();
  }

  static String randomUtf16(Random r, int len) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < len; i++) {
      switch (r.nextInt(8)) {
        case 0 -> sb.append((char) (0xD800 + r.nextInt(0x800))); // lone surrogate
        case 1 -> sb.appendCodePoint(0x10000 + r.nextInt(0x100000));
        case 2 -> sb.append((char) (0x80 + r.nextInt(0x780)));
        case 3 -> sb.append(":\\\n\r".charAt(r.nextInt(4)));
        default -> sb.append((char) (0x20 + r.nextInt(0x5f)));
      }
    }
    return sb.toString();
  }

  // ---------------------------------------------------------------- names

  static final String IDENT = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_$";

  static Set<String> parseNames(String s) {
    Set<String> result = new LinkedHashSet<>();
    if (!s.isEmpty()) {
      for (String n : s.split(",", -1)) {
        result.add(n);
      }
    }
    return result;
  }

  static Set<Character> parseChars(String s) {
    Set<Character> result = new LinkedHashSet<>();
    for (char c : s.toCharArray()) {
      result.add(c);
    }
    return result;
  }

  static String fromU(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 1; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append((char) Integer.parseInt(s.substring(i + 1, i + 5), 16));
        i += 4;
      } else {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  static String randomReserved(Random r) {
    List<String> names = new ArrayList<>();
    int n = r.nextInt(4) == 0 ? 0 : r.nextInt(30);
    for (int i = 0; i < n; i++) {
      // Short names so that they collide with generated ones.
      names.add(pick(r, r.nextBoolean() ? "abcdefgh" : IDENT, 1 + r.nextInt(2)));
    }
    return String.join(",", names);
  }

  static String randomPrefix(Random r) {
    switch (r.nextInt(6)) {
      case 0, 1, 2:
        return "";
      case 3:
        return pick(r, IDENT, 1 + r.nextInt(3));
      case 4:
        return pick(r, "abc$_", 1 + r.nextInt(2));
      default:
        return pick(r, "0123-a+", 1 + r.nextInt(2)); // often invalid
    }
  }

  static String randomChars(Random r) {
    return r.nextInt(3) == 0 ? "" : pick(r, IDENT, r.nextInt(8));
  }

  static void nameGenerators(Random r) {
    for (int c = 0; c < 400; c++) {
      List<String> ops = new ArrayList<>();
      boolean isDefault = r.nextBoolean();
      if (isDefault) {
        switch (r.nextInt(3)) {
          case 0 -> ops.add("D0");
          case 1 -> ops.add("D3:" + randomReserved(r) + ":" + randomPrefix(r) + ":" + randomChars(r));
          default ->
              ops.add(
                  "D4:" + randomReserved(r) + ":" + randomPrefix(r) + ":" + randomChars(r) + ":"
                      + randomChars(r));
        }
      } else {
        long seed = r.nextInt(5) == 0 ? r.nextLong() : r.nextInt(1000);
        switch (r.nextInt(3)) {
          case 0 -> ops.add("R1:" + seed);
          case 1 ->
              ops.add(
                  "R3:" + randomReserved(r) + ":" + randomPrefix(r) + ":" + randomChars(r) + ":"
                      + seed);
          default ->
              ops.add(
                  "R4:" + randomReserved(r) + ":" + randomPrefix(r) + ":" + randomChars(r) + ":"
                      + randomChars(r) + ":" + seed);
        }
      }
      int steps = 1 + r.nextInt(5);
      for (int s = 0; s < steps; s++) {
        switch (r.nextInt(isDefault ? 6 : 5)) {
          case 0 -> ops.add("r3:" + randomReserved(r) + ":" + randomPrefix(r) + ":" + randomChars(r));
          case 1 ->
              ops.add(
                  "r4:" + randomReserved(r) + ":" + randomPrefix(r) + ":" + randomChars(r) + ":"
                      + randomChars(r));
          case 2 -> ops.add("c:" + randomReserved(r) + ":" + randomPrefix(r) + ":" + randomChars(r));
          case 5 -> ops.add("f:" + u(randomUtf16(r, r.nextInt(40))));
          default -> ops.add("g:" + (1 + r.nextInt(r.nextInt(20) == 0 ? 1500 : 100)));
        }
      }
      ops.add("g:" + (1 + r.nextInt(60)));
      row("ng", String.join(";", ops), runNameOps(ops));
    }
  }

  static String runNameOps(List<String> ops) {
    List<String> names = new ArrayList<>();
    NameGenerator g = null;
    try {
      for (String op : ops) {
        String[] f = op.split(":", -1);
        switch (f[0]) {
          case "D0" -> g = new DefaultNameGenerator();
          case "D3" -> g = new DefaultNameGenerator(parseNames(f[1]), f[2], parseChars(f[3]));
          case "D4" ->
              g =
                  new DefaultNameGenerator(
                      parseNames(f[1]), f[2], parseChars(f[3]), parseChars(f[4]));
          case "R1" -> g = new RandomNameGenerator(new Random(Long.parseLong(f[1])));
          case "R3" ->
              g =
                  new RandomNameGenerator(
                      parseNames(f[1]), f[2], parseChars(f[3]), new Random(Long.parseLong(f[4])));
          case "R4" ->
              g =
                  new RandomNameGenerator(
                      parseNames(f[1]),
                      f[2],
                      parseChars(f[3]),
                      parseChars(f[4]),
                      new Random(Long.parseLong(f[5])));
          case "r3" -> g.reset(parseNames(f[1]), f[2], parseChars(f[3]));
          case "r4" -> g.reset(parseNames(f[1]), f[2], parseChars(f[3]), parseChars(f[4]));
          case "c" -> g = g.clone(parseNames(f[1]), f[2], parseChars(f[3]));
          case "f" -> ((DefaultNameGenerator) g).favors(fromU(f[1]));
          case "g" -> {
            int n = Integer.parseInt(f[1]);
            for (int i = 0; i < n; i++) {
              names.add(g.generateNextName());
            }
          }
          default -> throw new AssertionError(op);
        }
      }
    } catch (RuntimeException e) {
      names.add(exc(e));
    }
    return String.join(",", names);
  }

  // ---------------------------------------------------------------- graphs

  record RandomGraph(List<String> nodes, List<String[]> edges) {
    String nodesCol() {
      return String.join(",", nodes);
    }

    String edgesCol() {
      List<String> es = new ArrayList<>();
      for (String[] e : edges) {
        es.add(e[0] + ">" + e[1]);
      }
      return String.join(",", es);
    }
  }

  static RandomGraph randomGraph(Random r, int maxNodes, boolean acyclic) {
    int n = 1 + r.nextInt(maxNodes);
    List<Integer> order = new ArrayList<>();
    for (int i = 0; i < n; i++) {
      order.add(i);
    }
    java.util.Collections.shuffle(order, r);
    List<String> nodes = new ArrayList<>();
    for (int i : order) {
      nodes.add(String.valueOf(i));
    }
    int m = r.nextInt(3 * n + 1);
    List<String[]> edges = new ArrayList<>();
    for (int i = 0; i < m; i++) {
      int a = r.nextInt(n);
      int b = r.nextInt(n);
      if (acyclic) {
        if (a == b) {
          continue;
        }
        // Edges point from larger to smaller values (parents above children).
        if (a < b) {
          int t = a;
          a = b;
          b = t;
        }
      }
      edges.add(new String[] {String.valueOf(a), String.valueOf(b)});
    }
    return new RandomGraph(nodes, edges);
  }

  static void graphColoring(Random r) {
    for (int c = 0; c < 300; c++) {
      boolean directed = r.nextInt(4) == 0;
      String tie = "NAD".substring(r.nextInt(3)).substring(0, 1);
      RandomGraph rg = randomGraph(r, r.nextInt(5) == 0 ? 40 : 12, false);
      Graph<String, Object> g =
          directed ? LinkedDirectedGraph.create() : LinkedUndirectedGraph.create();
      for (String n : rg.nodes()) {
        g.createNode("n" + n);
      }
      for (String[] e : rg.edges()) {
        g.connect("n" + e[0], null, "n" + e[1]);
      }
      Comparator<String> cmp =
          switch (tie) {
            case "A" -> Comparator.naturalOrder();
            case "D" -> Comparator.<String>reverseOrder();
            default -> null;
          };
      GraphColoring<String, Object> coloring =
          cmp == null
              ? new GraphColoring.GreedyGraphColoring<>(g)
              : new GraphColoring.GreedyGraphColoring<>(g, cmp);
      String result;
      try {
        int count = coloring.color();
        List<String> supers = new ArrayList<>();
        for (String n : rg.nodes()) {
          supers.add(coloring.getPartitionSuperNode("n" + n));
        }
        result = count + ";" + String.join(",", supers);
      } catch (RuntimeException e) {
        result = exc(e);
      }
      row("color", directed ? 1 : 0, tie, rg.nodesCol(), rg.edgesCol(), result);
    }
  }

  static void unionFind(Random r) {
    for (int c = 0; c < 300; c++) {
      int n = 1 + r.nextInt(16);
      List<String> ops = new ArrayList<>();
      int count = 1 + r.nextInt(40);
      for (int i = 0; i < count; i++) {
        String a = "e" + r.nextInt(n + 1);
        String b = "e" + r.nextInt(n + 1);
        switch (r.nextInt(12)) {
          case 0, 1 -> ops.add("a:" + a);
          case 2, 3, 4, 5 -> ops.add("u:" + a + ":" + b);
          case 6 -> ops.add("f:" + a);
          case 7 -> ops.add("e:" + a + ":" + b);
          case 8 -> ops.add("F:" + a);
          case 9 -> ops.add(r.nextBoolean() ? "C" : "R");
          case 10 -> ops.add("E");
          default -> ops.add("K");
        }
      }
      ops.add("C");
      StandardUnionFind<String> uf = new StandardUnionFind<>();
      List<String> results = new ArrayList<>();
      for (String op : ops) {
        String[] f = op.split(":", -1);
        try {
          switch (f[0]) {
            case "a" -> {
              uf.add(f[1]);
              results.add("");
            }
            case "u" -> results.add(uf.union(f[1], f[2]));
            case "f" -> results.add(uf.find(f[1]));
            case "e" -> results.add(String.valueOf(uf.areEquivalent(f[1], f[2])));
            case "F" -> results.add(String.join(",", uf.findAll(f[1])) + "/" + uf.findAll(f[1]).size());
            case "C" -> {
              List<String> classes = new ArrayList<>();
              for (Set<String> s : uf.allEquivalenceClasses()) {
                classes.add(String.join(",", s));
              }
              results.add(String.join("|", classes));
            }
            case "R" -> results.add(String.join(",", uf.allRepresentatives()));
            case "E" -> results.add(String.join(",", uf.elements()));
            case "K" -> {
              uf = new StandardUnionFind<>(uf);
              results.add("");
            }
            default -> throw new AssertionError(op);
          }
        } catch (RuntimeException e) {
          results.add(exc(e));
        }
      }
      row("uf", String.join(";", ops), String.join(";", results));
    }
  }

  static LinkedDirectedGraph<Integer, Object> intGraph(RandomGraph rg) {
    LinkedDirectedGraph<Integer, Object> g = LinkedDirectedGraph.create();
    for (String n : rg.nodes()) {
      g.createNode(Integer.valueOf(n));
    }
    for (String[] e : rg.edges()) {
      g.connect(Integer.valueOf(e[0]), null, Integer.valueOf(e[1]));
    }
    return g;
  }

  static void dominators(Random r) {
    for (int c = 0; c < 300; c++) {
      RandomGraph rg = randomGraph(r, r.nextInt(5) == 0 ? 60 : 14, false);
      int n = rg.nodes().size();
      int entry = r.nextInt(20) == 0 ? n : r.nextInt(n);
      DiGraph<Integer, Object> g = intGraph(rg);
      String result;
      try {
        DominatorTree<Integer> tree = DominatorTree.compute(g, entry);
        List<String> idoms = new ArrayList<>();
        List<String> sizes = new ArrayList<>();
        for (int i = 0; i <= n; i++) {
          Integer idom = tree.getImmediateDominator(i);
          idoms.add(idom == null ? "-" : idom.toString());
          sizes.add(String.valueOf(tree.getDominatedSubtreeSize(i)));
        }
        List<String> all = new ArrayList<>();
        for (Map.Entry<Integer, Integer> e : tree.getAllSubtreeSizes().entrySet()) {
          all.add(e.getKey() + "=" + e.getValue());
        }
        result = String.join(",", idoms) + ";" + String.join(",", sizes) + ";" + String.join(",", all);
      } catch (RuntimeException e) {
        result = exc(e);
      }
      row("dom", rg.nodesCol(), rg.edgesCol(), entry, result);
    }
  }

  static void lowestCommonAncestors(Random r) {
    for (int c = 0; c < 300; c++) {
      RandomGraph rg = randomGraph(r, r.nextInt(5) == 0 ? 50 : 14, r.nextInt(3) != 0);
      int n = rg.nodes().size();
      List<String> roots = new ArrayList<>();
      int k = r.nextInt(30) == 0 ? 32 : 1 + r.nextInt(Math.min(n, 5));
      Set<Integer> rootSet = new LinkedHashSet<>();
      for (int i = 0; i < k; i++) {
        int root = r.nextInt(40) == 0 ? n + r.nextInt(3) : r.nextInt(k == 32 ? 64 : n);
        if (rootSet.add(root)) {
          roots.add(String.valueOf(root));
        }
      }
      LinkedDirectedGraph<Integer, Object> g = intGraph(rg);
      if (k == 32) {
        for (int root : rootSet) {
          g.createNode(root);
          rg.nodes().add(String.valueOf(root));
        }
        // Keep the node column in creation order.
      }
      String result;
      try {
        List<String> found = new ArrayList<>();
        LowestCommonAncestorFinder<Integer, Object> finder = new LowestCommonAncestorFinder<>(g);
        for (Integer x : finder.findAll(rootSet)) {
          found.add(x.toString());
        }
        // A second search on the same finder must give the same answer.
        List<String> again = new ArrayList<>();
        for (Integer x : finder.findAll(rootSet)) {
          again.add(x.toString());
        }
        result = String.join(",", found) + ";" + String.join(",", again);
      } catch (RuntimeException e) {
        result = exc(e);
      }
      row("lca", rg.nodesCol(), rg.edgesCol(), String.join(",", roots), result);
    }
  }

  // ---------------------------------------------------------------- VariableMap

  static String entries(Map<String, String> m) {
    List<String> es = new ArrayList<>();
    for (Map.Entry<String, String> e : m.entrySet()) {
      es.add(u(e.getKey()) + "=" + u(e.getValue()));
    }
    return String.join(",", es);
  }

  static void variableMaps(Random r) {
    for (int c = 0; c < 300; c++) {
      Map<String, String> m = new LinkedHashMap<>();
      int n = r.nextInt(12);
      for (int i = 0; i < n; i++) {
        String k = r.nextInt(3) == 0 ? randomUtf16(r, 1 + r.nextInt(6)) : pick(r, "abcde:\\", 1 + r.nextInt(4));
        String v = r.nextInt(3) == 0 ? randomUtf16(r, r.nextInt(6)) : pick(r, "xyz:\\", r.nextInt(4)) + (r.nextInt(6) == 0 ? "" : String.valueOf(i));
        m.put(k, v);
      }
      String result;
      try {
        VariableMap vm = new VariableMap(m);
        byte[] bytes = vm.toBytes();
        String back;
        try {
          VariableMap parsed = VariableMap.fromBytes(bytes);
          back = entries(parsed.toMap()) + ";" + entries(parsed.getOriginalNameToNewNameMap())
              + ";" + entries(parsed.getNewNameToOriginalNameMap());
        } catch (Exception e) {
          back = exc(e);
        }
        List<String> lookups = new ArrayList<>();
        for (Map.Entry<String, String> e : m.entrySet()) {
          lookups.add(u(vm.lookupNewName(e.getKey())) + "/" + u(vm.lookupSourceName(e.getValue())));
        }
        result = hex(bytes) + ";" + String.join(",", lookups) + ";" + back;
      } catch (RuntimeException e) {
        result = exc(e);
      }
      row("vmap", entries(m), result);
    }
    // fromBytes on fuzzed input.
    for (int c = 0; c < 300; c++) {
      byte[] b;
      if (r.nextBoolean()) {
        b = new byte[r.nextInt(24)];
        String alphabet = "ab:\\\n\r=";
        for (int i = 0; i < b.length; i++) {
          b[i] = r.nextInt(6) == 0 ? (byte) r.nextInt(256) : (byte) alphabet.charAt(r.nextInt(alphabet.length()));
        }
      } else {
        StringBuilder sb = new StringBuilder();
        int lines = r.nextInt(5);
        for (int i = 0; i < lines; i++) {
          sb.append(pick(r, "ab\\", 1 + r.nextInt(3))).append(':').append(pick(r, "xy\\:", r.nextInt(3)));
          sb.append(r.nextInt(4) == 0 ? "\r\n" : "\n");
        }
        b = sb.toString().getBytes(java.nio.charset.StandardCharsets.UTF_8);
      }
      String result;
      try {
        result = entries(VariableMap.fromBytes(b).toMap());
      } catch (Exception e) {
        result = exc(e);
      }
      row("vbytes", hex(b), result);
    }
  }

  // ---------------------------------------------------------------- messages

  static final String PH = "abcXYZ_019AZaz٣é";

  static void messages(Random r) {
    for (int c = 0; c < 300; c++) {
      String projectId =
          switch (r.nextInt(4)) {
            case 0 -> null;
            case 1 -> "";
            case 2 -> "MyProject";
            default -> randomUtf16(r, r.nextInt(8));
          };
      String meaning = r.nextInt(25) == 0 ? null : (r.nextBoolean() ? "" : randomUtf16(r, r.nextInt(12)));
      List<String> partsCol = new ArrayList<>();
      String result;
      try {
        List<JsMessage.Part> parts = new ArrayList<>();
        int n = r.nextInt(6);
        for (int i = 0; i < n; i++) {
          switch (r.nextInt(3)) {
            case 0 -> {
              String s = randomUtf16(r, r.nextInt(10));
              partsCol.add("s" + u(s));
              parts.add(JsMessage.StringPart.create(s));
            }
            case 1 -> {
              String s = r.nextInt(8) != 0 ? pick(r, "abc", 1) + pick(r, "abcXY019", r.nextInt(6)) + pick(r, "_19", r.nextInt(3)) : pick(r, PH, r.nextInt(6));
              partsCol.add("j" + u(s));
              parts.add(JsMessage.PlaceholderReference.createForJsName(s));
            }
            default -> {
              String s = r.nextInt(8) != 0 ? pick(r, "ABC_019", r.nextInt(8)) : pick(r, PH, r.nextInt(6));
              partsCol.add("c" + u(s));
              parts.add(JsMessage.PlaceholderReference.createForCanonicalName(s));
            }
          }
        }
        result = new GoogleJsMessageIdGenerator(projectId).generateId(meaning, parts);
      } catch (RuntimeException e) {
        result = exc(e);
      }
      row("msgid", u(projectId), u(meaning), String.join(",", partsCol), result);
    }
    for (int c = 0; c < 200; c++) {
      String s = r.nextInt(30) == 0 ? null : randomUtf16(r, r.nextInt(r.nextInt(5) == 0 ? 60 : 14));
      row("hash", u(s), JsMessage.Hash.hash64(s));
    }
    for (int c = 0; c < 300; c++) {
      String s = pick(r, r.nextBoolean() ? PH : "AB_1a", r.nextInt(9));
      String js;
      try {
        js = u(JsMessage.PlaceholderReference.createForJsName(s).getCanonicalPlaceholderName());
      } catch (RuntimeException e) {
        js = exc(e);
      }
      String canon;
      try {
        canon = u(JsMessage.PlaceholderReference.createForCanonicalName(s).getJsPlaceholderName());
      } catch (RuntimeException e) {
        canon = exc(e);
      }
      String lower;
      try {
        lower = u(JsMessage.toLowerCamelCaseWithNumericSuffixes(s));
      } catch (RuntimeException e) {
        lower = exc(e);
      }
      row(
          "ph",
          u(s),
          JsMessage.isLowerCamelCaseWithNumericSuffixes(s),
          JsMessage.isCanonicalPlaceholderNameFormat(s),
          lower,
          js,
          canon);
    }
  }

  // ---------------------------------------------------------------- misc

  static void misc(Random r) {
    Xid xid = new Xid();
    for (int c = 0; c < 200; c++) {
      String s = randomUtf16(r, r.nextInt(16));
      row("xid", u(s), xid.get(s), xid.getAsInt(s));
    }
    String[] segs = {"a", "b", "..", ".", "", "c.js", "...", "x y"};
    for (int c = 0; c < 300; c++) {
      String[] paths = new String[2];
      for (int p = 0; p < 2; p++) {
        StringBuilder sb = new StringBuilder();
        if (r.nextInt(3) == 0) {
          sb.append('/');
        }
        int n = r.nextInt(6);
        for (int i = 0; i < n; i++) {
          if (i > 0) {
            sb.append('/');
          }
          sb.append(segs[r.nextInt(segs.length)]);
        }
        paths[p] = sb.toString();
      }
      // (makeAbsolute(String) depends on the working directory, so it is not used.)
      String op = "CAR".substring(r.nextInt(3)).substring(0, 1);
      String result;
      try {
        result =
            u(
                switch (op) {
                  case "C" -> PathUtil.collapseDots(paths[0]);
                  case "A" -> PathUtil.makeAbsolute(paths[0], paths[1]);
                  default -> PathUtil.makeRelative(paths[0], paths[1]);
                });
      } catch (RuntimeException e) {
        result = exc(e);
      }
      row("path", op, u(paths[0]), u(paths[1]), result);
    }
    for (int c = 0; c < 200; c++) {
      String s = randomUtf16(r, r.nextInt(16));
      if (r.nextInt(4) == 0) {
        s += "</script><!--  \u0000\u007f\u0085'\"";
      }
      row("esc", u(s), u(SourceCodeEscapers.javascriptEscaper().escape(s)));
    }
    List<Double> nums = new ArrayList<>();
    double[] fixed = {
      0, -0.0, 1, 0.5, 0.1, 100, 1000, 1e3, 1e21, 1e22, 123e5, 1e-7, 1.5e-7, 1e12, 1e12 + 1,
      1000000000001L, 0x7fffffffffffffffL, 9007199254740992.0, 9007199254740993.0, 1e300,
      Double.MAX_VALUE, Double.MIN_VALUE, Double.POSITIVE_INFINITY, Double.NaN, -1, -1e10,
      Math.PI, 255, 4294967295.0, 4294967296.0, 1e15, 1e16, 1e17, 1e18, 1e19, 1e20, 12345678912345.0,
      0x123456789abcL, 3e-5, 12.5, 0.000001, 0.0000001
    };
    for (double d : fixed) {
      nums.add(d);
    }
    for (int c = 0; c < 300; c++) {
      switch (r.nextInt(5)) {
        case 0 -> nums.add(Double.longBitsToDouble(r.nextLong()));
        case 1 -> nums.add((double) r.nextInt(Integer.MAX_VALUE));
        case 2 -> nums.add((double) (r.nextLong() >>> r.nextInt(64)));
        case 3 -> nums.add(r.nextInt(1000) * Math.pow(10, r.nextInt(30)));
        default -> nums.add(r.nextDouble() * Math.pow(10, r.nextInt(40) - 20));
      }
    }
    for (double d : nums) {
      StringBuilder buf = new StringBuilder();
      CodeConsumer cc =
          new CodeConsumer() {
            @Override
            char getLastChar() {
              return buf.length() == 0 ? '\0' : buf.charAt(buf.length() - 1);
            }

            @Override
            void append(String str) {
              buf.append(str);
            }
          };
      String result;
      try {
        cc.addNumber(d, null);
        result = u(buf.toString());
      } catch (RuntimeException e) {
        result = exc(e);
      }
      row("num", Long.toHexString(Double.doubleToRawLongBits(d)), result);
    }
  }

  public static void main(String[] args) throws Exception {
    out = new PrintStream(Files.newOutputStream(Paths.get(args[0])), false, "UTF-8");
    Random r = new Random(20260415);
    nameGenerators(r);
    graphColoring(r);
    unionFind(r);
    dominators(r);
    lowestCommonAncestors(r);
    variableMaps(r);
    messages(r);
    misc(r);
    out.close();
  }
}

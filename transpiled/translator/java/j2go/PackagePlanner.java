package j2go;

import com.sun.source.tree.ClassTree;
import com.sun.source.tree.CompilationUnitTree;
import com.sun.source.tree.Tree;
import com.sun.source.util.TreePath;
import com.sun.source.util.TreePathScanner;
import com.sun.source.util.Trees;
import java.net.URI;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.Deque;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeMap;
import java.util.TreeSet;
import javax.lang.model.element.Element;
import javax.lang.model.element.NestingKind;
import javax.lang.model.element.PackageElement;
import javax.lang.model.element.TypeElement;
import javax.lang.model.util.Elements;

/**
 * Computes Strongly Connected Components (SCCs) of the Java package dependency graph,
 * maps cyclic Java packages into minimal acyclic Go packages, and assigns deterministic
 * collision-free Go type names and 1:1 {@code .go} file names.
 */
public final class PackagePlanner {

  public static final class UnitPlan {
    public final CompilationUnitTree unit;
    public final Path sourcePath;
    public final String javaPkg;
    public final String goPkg;
    public final String goFileName;
    public final List<TypeElement> topLevelTypes;

    UnitPlan(
        CompilationUnitTree unit,
        Path sourcePath,
        String javaPkg,
        String goPkg,
        String goFileName,
        List<TypeElement> topLevelTypes) {
      this.unit = unit;
      this.sourcePath = sourcePath;
      this.javaPkg = javaPkg;
      this.goPkg = goPkg;
      this.goFileName = goFileName;
      this.topLevelTypes = topLevelTypes;
    }
  }

  public static final class SccGroup {
    public final int id;
    public final String goPkgName;
    public final String goImportPath;
    public final List<String> javaPackages;
    public final Set<String> depGoPackages = new TreeSet<>();

    SccGroup(int id, String goPkgName, String goImportPath, List<String> javaPackages) {
      this.id = id;
      this.goPkgName = goPkgName;
      this.goImportPath = goImportPath;
      this.javaPackages = javaPackages;
    }
  }

  /** Tracks per-file Go package imports so unused imports are never emitted. */
  public static final class ImportTracker {
    private final String currentGoPkg;
    // importPath -> alias (null if default basename matches)
    private final Map<String, String> stdImports = new TreeMap<>();
    private final Map<String, String> extImports = new TreeMap<>();

    public ImportTracker(String currentGoPkg) {
      this.currentGoPkg = currentGoPkg;
    }

    public String std(String pkgPath) {
      stdImports.putIfAbsent(pkgPath, null);
      int slash = pkgPath.lastIndexOf('/');
      return slash >= 0 ? pkgPath.substring(slash + 1) : pkgPath;
    }

    public String ext(String importPath, String alias) {
      String base = importPath.substring(importPath.lastIndexOf('/') + 1);
      if (base.equals(currentGoPkg) && importPath.endsWith("/" + currentGoPkg)) {
        return "";
      }
      if (!importPath.startsWith("github.com/paullepoulpe/jsgo/transpiled/shims/")
          && !base.equals("j2gort")
          && !alias.startsWith("pkg_")) {
        alias = "pkg_" + alias;
      }
      extImports.putIfAbsent(importPath, base.equals(alias) ? null : alias);
      return alias;
    }

    public String emitImportBlock(String bodyText) {
      List<Map.Entry<String, String>> usedStd = new ArrayList<>();
      for (Map.Entry<String, String> e : stdImports.entrySet()) {
        String alias = e.getValue() != null ? e.getValue() : lastSegment(e.getKey());
        if (bodyText.contains(alias + ".")) {
          usedStd.add(e);
        }
      }
      List<Map.Entry<String, String>> usedExt = new ArrayList<>();
      for (Map.Entry<String, String> e : extImports.entrySet()) {
        String alias = e.getValue() != null ? e.getValue() : lastSegment(e.getKey());
        if (bodyText.contains(alias + ".")) {
          usedExt.add(e);
        }
      }
      if (usedStd.isEmpty() && usedExt.isEmpty()) {
        return "";
      }
      StringBuilder sb = new StringBuilder();
      sb.append("import (\n");
      for (Map.Entry<String, String> e : usedStd) {
        sb.append('\t');
        if (e.getValue() != null) {
          sb.append(e.getValue()).append(' ');
        }
        sb.append('"').append(e.getKey()).append("\"\n");
      }
      if (!usedStd.isEmpty() && !usedExt.isEmpty()) {
        sb.append('\n');
      }
      for (Map.Entry<String, String> e : usedExt) {
        sb.append('\t');
        if (e.getValue() != null) {
          sb.append(e.getValue()).append(' ');
        }
        sb.append('"').append(e.getKey()).append("\"\n");
      }
      sb.append(")\n\n");
      return sb.toString();
    }

    private static String lastSegment(String importPath) {
      int slash = importPath.lastIndexOf('/');
      return slash >= 0 ? importPath.substring(slash + 1) : importPath;
    }
  }

  private final String goModuleBase;
  private final List<UnitPlan> unitPlans = new ArrayList<>();
  private final List<SccGroup> sccGroups = new ArrayList<>();
  private final Map<String, SccGroup> javaPkgToScc = new LinkedHashMap<>();
  private final Map<String, SccGroup> goPkgToScc = new LinkedHashMap<>();
  private final Set<TypeElement> allSourceTypes = new LinkedHashSet<>();
  private final Map<TypeElement, String> goTypeNames = new LinkedHashMap<>();
  private final Map<TypeElement, String> typeToGoPkg = new LinkedHashMap<>();

  private PackagePlanner(String goModuleBase) {
    this.goModuleBase = goModuleBase;
  }

  public static PackagePlanner plan(
      List<CompilationUnitTree> units, Trees trees, Elements elements, String goModuleBase) {
    PackagePlanner planner = new PackagePlanner(goModuleBase);
    planner.compute(units, trees, elements);
    return planner;
  }

  public String getGoModuleBase() {
    return goModuleBase;
  }

  public List<UnitPlan> getUnitPlans() {
    return Collections.unmodifiableList(unitPlans);
  }

  public List<SccGroup> getSccGroups() {
    return Collections.unmodifiableList(sccGroups);
  }

  public Set<TypeElement> getAllSourceTypes() {
    return Collections.unmodifiableSet(allSourceTypes);
  }

  public boolean isSourceType(TypeElement te) {
    return allSourceTypes.contains(te);
  }

  public String getGoPkgForJavaPkg(String javaPkg) {
    SccGroup g = javaPkgToScc.get(javaPkg);
    return g != null ? g.goPkgName : sanitizeGoPkgName(javaPkg);
  }

  public String getGoPkg(TypeElement te) {
    String p = typeToGoPkg.get(te);
    if (p != null) {
      return p;
    }
    Element cur = te.getEnclosingElement();
    while (cur != null) {
      if (cur instanceof TypeElement enc && typeToGoPkg.containsKey(enc)) {
        return typeToGoPkg.get(enc);
      }
      if (cur instanceof PackageElement pe) {
        return getGoPkgForJavaPkg(pe.getQualifiedName().toString());
      }
      cur = cur.getEnclosingElement();
    }
    return "pkg";
  }

  public String getGoImportPath(String goPkg) {
    SccGroup g = goPkgToScc.get(goPkg);
    return g != null ? g.goImportPath : (goModuleBase + "/" + goPkg);
  }

  public String getGoTypeName(TypeElement te) {
    String existing = goTypeNames.get(te);
    if (existing != null) {
      return existing;
    }
    return GoNames.exported(te.getSimpleName().toString());
  }

  public void registerLocalOrAnonType(TypeElement te, TypeElement enclosingTop, int seq, boolean isAnon) {
    allSourceTypes.add(te);
    String goPkg = getGoPkg(enclosingTop);
    typeToGoPkg.put(te, goPkg);
    String topName = getGoTypeName(enclosingTop);
    String name =
        isAnon
            ? topName + "__Anon_" + seq
            : topName + "__Local_" + seq + "_" + GoNames.exported(te.getSimpleName().toString());
    goTypeNames.put(te, name);
  }

  private void compute(List<CompilationUnitTree> units, Trees trees, Elements elements) {
    Set<String> sourceJavaPkgs = new TreeSet<>();
    Map<CompilationUnitTree, List<TypeElement>> unitTopTypes = new LinkedHashMap<>();

    // Pass 1: Discover all source TypeElements (top-level, member, local, anonymous)
    for (CompilationUnitTree cu : units) {
      String pkg = cu.getPackageName() == null ? "" : cu.getPackageName().toString();
      sourceJavaPkgs.add(pkg);
      List<TypeElement> tops = new ArrayList<>();
      unitTopTypes.put(cu, tops);

      new TreePathScanner<Void, Void>() {
        @Override
        public Void visitClass(ClassTree node, Void p) {
          Element el = trees.getElement(getCurrentPath());
          if (el instanceof TypeElement te) {
            allSourceTypes.add(te);
            if (te.getNestingKind() == NestingKind.TOP_LEVEL) {
              tops.add(te);
            }
          }
          return super.visitClass(node, p);
        }
      }.scan(cu, null);
    }

    // Pass 2: Build Java package dependency graph across all AST symbol references
    Map<String, Set<String>> pkgDeps = new TreeMap<>();
    for (String p : sourceJavaPkgs) {
      pkgDeps.put(p, new TreeSet<>());
    }

    for (CompilationUnitTree cu : units) {
      String curPkg = cu.getPackageName() == null ? "" : cu.getPackageName().toString();
      new TreePathScanner<Void, Void>() {
        @Override
        public Void scan(Tree tree, Void p) {
          if (tree != null) {
            TreePath path = new TreePath(getCurrentPath(), tree);
            recordElement(trees.getElement(path), curPkg, elements, pkgDeps);
          }
          return super.scan(tree, p);
        }
      }.scan(cu, null);
    }

    // Pass 3: Run Tarjan's SCC algorithm in deterministic order (post-order = leaves first)
    List<List<String>> rawSccs = computeTarjanSccs(sourceJavaPkgs, pkgDeps);

    Set<String> usedGoPkgNames = new LinkedHashSet<>();
    usedGoPkgNames.add("j2gort");
    usedGoPkgNames.add("javalang");
    usedGoPkgNames.add("javautil");
    usedGoPkgNames.add("javaio");
    usedGoPkgNames.add("guava");
    usedGoPkgNames.add("nodestore");

    for (int i = 0; i < rawSccs.size(); i++) {
      List<String> sccPkgs = new ArrayList<>(rawSccs.get(i));
      Collections.sort(sccPkgs);
      String candidate = chooseGoPkgName(sccPkgs);
      String uniqueName = candidate;
      int suffix = 2;
      while (!usedGoPkgNames.add(uniqueName)) {
        uniqueName = candidate + suffix++;
      }
      SccGroup group = new SccGroup(i, uniqueName, goModuleBase + "/" + uniqueName, sccPkgs);
      sccGroups.add(group);
      goPkgToScc.put(uniqueName, group);
      for (String jp : sccPkgs) {
        javaPkgToScc.put(jp, group);
      }
    }

    // Populate inter-SCC Go package dependencies
    for (SccGroup g : sccGroups) {
      for (String jp : g.javaPackages) {
        for (String depJp : pkgDeps.getOrDefault(jp, Set.of())) {
          SccGroup depGroup = javaPkgToScc.get(depJp);
          if (depGroup != null && depGroup != g) {
            g.depGoPackages.add(depGroup.goPkgName);
          }
        }
      }
    }

    // Pass 4: Assign collision-free Go type names within each SCC
    // First collect top-level TypeElements per SCC
    Map<SccGroup, List<TypeElement>> sccTopTypes = new LinkedHashMap<>();
    for (SccGroup g : sccGroups) {
      sccTopTypes.put(g, new ArrayList<>());
    }
    for (CompilationUnitTree cu : units) {
      String jp = cu.getPackageName() == null ? "" : cu.getPackageName().toString();
      SccGroup g = javaPkgToScc.get(jp);
      for (TypeElement te : unitTopTypes.get(cu)) {
        sccTopTypes.get(g).add(te);
        typeToGoPkg.put(te, g.goPkgName);
      }
    }

    for (Map.Entry<SccGroup, List<TypeElement>> entry : sccTopTypes.entrySet()) {
      SccGroup g = entry.getKey();
      List<TypeElement> tops = entry.getValue();
      String rootJavaPkg = shortestPkg(g.javaPackages);

      // Group by exported simple name to detect within-SCC name collisions
      Map<String, List<TypeElement>> bySimple = new LinkedHashMap<>();
      for (TypeElement te : tops) {
        String exp = GoNames.exported(te.getSimpleName().toString());
        bySimple.computeIfAbsent(exp, k -> new ArrayList<>()).add(te);
      }

      Set<String> assignedInScc = new LinkedHashSet<>();
      for (Map.Entry<String, List<TypeElement>> se : bySimple.entrySet()) {
        String expSimple = se.getKey();
        List<TypeElement> colliding = se.getValue();
        if (colliding.size() == 1) {
          TypeElement te = colliding.get(0);
          goTypeNames.put(te, expSimple);
          assignedInScc.add(expSimple);
        } else {
          // Sort deterministically by qualified name; root package keeps unprefixed name if unique
          colliding.sort(Comparator.comparing(t -> t.getQualifiedName().toString()));
          for (TypeElement te : colliding) {
            String pkg = elements.getPackageOf(te).getQualifiedName().toString();
            String candidate;
            if (pkg.equals(rootJavaPkg) && !assignedInScc.contains(expSimple)) {
              candidate = expSimple;
            } else {
              String rel = relativePkgPrefix(rootJavaPkg, pkg);
              candidate = rel + "_" + expSimple;
            }
            int k = 2;
            String unique = candidate;
            while (!assignedInScc.add(unique)) {
              unique = candidate + "_" + k++;
            }
            goTypeNames.put(te, unique);
          }
        }
      }
    }

    // Assign member TypeElement Go names recursively: <EnclosingGoType>_<SimpleName>
    for (SccGroup g : sccGroups) {
      for (TypeElement top : sccTopTypes.get(g)) {
        assignMemberNamesRecursive(top, g.goPkgName);
      }
    }

    // Pass 5: Assign 1:1 .go filenames per CompilationUnitTree
    Map<String, Set<String>> usedFilesPerGoPkg = new LinkedHashMap<>();
    for (CompilationUnitTree cu : units) {
      URI uri = cu.getSourceFile().toUri();
      Path srcPath = Paths.get(uri);
      String javaPkg = cu.getPackageName() == null ? "" : cu.getPackageName().toString();
      SccGroup g = javaPkgToScc.get(javaPkg);
      String rootJavaPkg = shortestPkg(g.javaPackages);

      String rawBase = srcPath.getFileName().toString();
      if (rawBase.endsWith(".java")) {
        rawBase = rawBase.substring(0, rawBase.length() - 5);
      }
      String snake = GoNames.toSnakeCase(rawBase);
      Set<String> usedFiles =
          usedFilesPerGoPkg.computeIfAbsent(g.goPkgName, k -> new LinkedHashSet<>());
      String candidate = snake;
      if (usedFiles.contains(candidate) && !javaPkg.equals(rootJavaPkg)) {
        String rel = GoNames.toSnakeCase(relativePkgPrefix(rootJavaPkg, javaPkg));
        candidate = rel + "_" + snake;
      }
      int idx = 2;
      String uniqueFile = candidate;
      while (!usedFiles.add(uniqueFile)) {
        uniqueFile = candidate + "_" + idx++;
      }
      unitPlans.add(
          new UnitPlan(
              cu, srcPath, javaPkg, g.goPkgName, uniqueFile + ".go", unitTopTypes.get(cu)));
    }
  }

  private void assignMemberNamesRecursive(TypeElement enclosing, String goPkg) {
    String encName = goTypeNames.get(enclosing);
    for (Element enclosed : enclosing.getEnclosedElements()) {
      if (enclosed instanceof TypeElement member && member.getNestingKind() == NestingKind.MEMBER) {
        allSourceTypes.add(member);
        typeToGoPkg.put(member, goPkg);
        String memberGoName = encName + "_" + GoNames.exported(member.getSimpleName().toString());
        goTypeNames.put(member, memberGoName);
        assignMemberNamesRecursive(member, goPkg);
      }
    }
  }

  private void recordElement(
      Element el, String curPkg, Elements elements, Map<String, Set<String>> pkgDeps) {
    if (el == null) {
      return;
    }
    Element cur = el;
    while (cur != null && !(cur instanceof TypeElement)) {
      cur = cur.getEnclosingElement();
    }
    if (cur instanceof TypeElement te) {
      // Find top-level enclosing TypeElement
      TypeElement top = te;
      while (top.getEnclosingElement() instanceof TypeElement p) {
        top = p;
      }
      if (allSourceTypes.contains(top)) {
        PackageElement pe = elements.getPackageOf(top);
        if (pe != null) {
          String depPkg = pe.getQualifiedName().toString();
          if (!depPkg.equals(curPkg) && pkgDeps.containsKey(depPkg)) {
            pkgDeps.get(curPkg).add(depPkg);
          }
        }
      }
    }
  }

  private static List<List<String>> computeTarjanSccs(
      Set<String> nodes, Map<String, Set<String>> edges) {
    Map<String, Integer> indexMap = new LinkedHashMap<>();
    Map<String, Integer> lowlink = new LinkedHashMap<>();
    Deque<String> stack = new ArrayDeque<>();
    Set<String> onStack = new LinkedHashSet<>();
    List<List<String>> result = new ArrayList<>();
    int[] index = {0};

    class Dfs {
      void visit(String v) {
        indexMap.put(v, index[0]);
        lowlink.put(v, index[0]);
        index[0]++;
        stack.push(v);
        onStack.add(v);

        for (String w : edges.getOrDefault(v, Set.of())) {
          if (!nodes.contains(w)) {
            continue;
          }
          if (!indexMap.containsKey(w)) {
            visit(w);
            lowlink.put(v, Math.min(lowlink.get(v), lowlink.get(w)));
          } else if (onStack.contains(w)) {
            lowlink.put(v, Math.min(lowlink.get(v), indexMap.get(w)));
          }
        }

        if (lowlink.get(v).equals(indexMap.get(v))) {
          List<String> scc = new ArrayList<>();
          while (true) {
            String w = stack.pop();
            onStack.remove(w);
            scc.add(w);
            if (w.equals(v)) {
              break;
            }
          }
          result.add(scc);
        }
      }
    }

    Dfs dfs = new Dfs();
    for (String n : nodes) {
      if (!indexMap.containsKey(n)) {
        dfs.visit(n);
      }
    }
    return result;
  }

  private static String shortestPkg(List<String> pkgs) {
    String best = pkgs.get(0);
    for (String p : pkgs) {
      if (p.length() < best.length() || (p.length() == best.length() && p.compareTo(best) < 0)) {
        best = p;
      }
    }
    return best;
  }

  private static String relativePkgPrefix(String rootPkg, String pkg) {
    String rel = pkg;
    if (!rootPkg.isEmpty() && pkg.startsWith(rootPkg + ".")) {
      rel = pkg.substring(rootPkg.length() + 1);
    } else if (pkg.contains(".")) {
      rel = pkg.substring(pkg.lastIndexOf('.') + 1);
    }
    StringBuilder sb = new StringBuilder();
    for (String part : rel.split("\\.")) {
      if (!part.isEmpty()) {
        if (sb.length() > 0) {
          sb.append('_');
        }
        sb.append(GoNames.exported(part));
      }
    }
    return sb.length() == 0 ? "Pkg" : sb.toString();
  }

  private static String chooseGoPkgName(List<String> sccPkgs) {
    String root = shortestPkg(sccPkgs);
    if (root.equals("com.google.javascript.jscomp.base.format")) {
      return "baseformat";
    }
    if (root.equals("com.google.javascript.jscomp.type")) {
      return "jscomptype";
    }
    return sanitizeGoPkgName(root);
  }

  private static String sanitizeGoPkgName(String javaPkg) {
    if (javaPkg == null || javaPkg.isEmpty()) {
      return "mainpkg";
    }
    String last = javaPkg.substring(javaPkg.lastIndexOf('.') + 1).toLowerCase();
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < last.length(); i++) {
      char c = last.charAt(i);
      if ((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
        sb.append(c);
      }
    }
    String res = sb.length() == 0 ? "pkg" : sb.toString();
    if (res.charAt(0) >= '0' && res.charAt(0) <= '9') {
      res = "p" + res;
    }
    if (res.equals("type")
        || res.equals("map")
        || res.equals("range")
        || res.equals("func")
        || res.equals("chan")
        || res.equals("select")
        || res.equals("struct")
        || res.equals("interface")
        || res.equals("package")
        || res.equals("import")
        || res.equals("default")
        || res.equals("switch")
        || res.equals("case")
        || res.equals("go")
        || res.equals("var")
        || res.equals("const")) {
      res = "pkg" + res;
    }
    return res;
  }
}

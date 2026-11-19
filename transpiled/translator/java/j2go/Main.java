package j2go;

import com.sun.source.tree.CompilationUnitTree;
import com.sun.source.util.JavacTask;
import com.sun.source.util.Trees;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.Locale;
import java.util.stream.Stream;
import javax.tools.Diagnostic;
import javax.tools.DiagnosticCollector;
import javax.tools.JavaCompiler;
import javax.tools.JavaFileObject;
import javax.tools.StandardJavaFileManager;
import javax.tools.ToolProvider;

/**
 * CLI entry point for {@code j2go}, the clean-room mechanical Java-to-Go translator
 * built on OpenJDK's public {@code javax.tools.JavaCompiler} and {@code com.sun.source.*} APIs.
 */
public final class Main {

  private Main() {}

  public static void main(String[] args) throws Exception {
    List<Path> srcDirs = new ArrayList<>();
    Path outDir = null;
    Path shadedJar = null;
    Path unshadedJarCache = null;
    String explicitCp = null;
    String goModuleBase = "github.com/paullepoulpe/jsgo/transpiled/gen";
    List<String> excludes = new ArrayList<>();
    excludes.add("com/google/javascript/jscomp/ant/");
    boolean printSccs = false;

    for (int i = 0; i < args.length; i++) {
      switch (args[i]) {
        case "--src" -> srcDirs.add(Paths.get(args[++i]).toAbsolutePath().normalize());
        case "--out" -> outDir = Paths.get(args[++i]).toAbsolutePath().normalize();
        case "--shaded-jar" -> shadedJar = Paths.get(args[++i]).toAbsolutePath().normalize();
        case "--unshaded-jar" -> unshadedJarCache = Paths.get(args[++i]).toAbsolutePath().normalize();
        case "--cp" -> explicitCp = args[++i];
        case "--go-module-base" -> goModuleBase = args[++i];
        case "--exclude" -> excludes.add(args[++i]);
        case "--print-sccs" -> printSccs = true;
        default -> throw new IllegalArgumentException("Unknown flag: " + args[i]);
      }
    }

    if (srcDirs.isEmpty() || outDir == null) {
      System.err.println(
          "Usage: j2go.Main --src <dir> --out <dir> [--shaded-jar <jar>] [--unshaded-jar <jar>] "
              + "[--cp <classpath>] [--go-module-base <prefix>] [--exclude <substr>] [--print-sccs]");
      System.exit(2);
    }

    long t0 = System.currentTimeMillis();
    List<Path> javaFiles = new ArrayList<>();
    for (Path sd : srcDirs) {
      try (Stream<Path> walk = Files.walk(sd)) {
        walk.filter(p -> Files.isRegularFile(p) && p.toString().endsWith(".java"))
            .filter(
                p -> {
                  String norm = p.toString().replace('\\', '/');
                  if (norm.endsWith("/package-info.java") || norm.endsWith("/module-info.java")) {
                    return false;
                  }
                  for (String ex : excludes) {
                    if (!ex.isEmpty() && norm.contains(ex)) {
                      return false;
                    }
                  }
                  return true;
                })
            .forEach(javaFiles::add);
      }
    }
    javaFiles.sort(Comparator.comparing(Path::toString));

    if (javaFiles.isEmpty()) {
      throw new IllegalStateException("No .java files found under " + srcDirs);
    }

    Files.createDirectories(outDir);

    String cp = explicitCp;
    if (shadedJar != null) {
      Path targetUnshaded =
          unshadedJarCache != null
              ? unshadedJarCache
              : outDir.getParent().resolve(".j2go-closure-deps-unshaded.jar");
      Path builtJar = ClasspathUnshader.ensureUnshadedJar(shadedJar, javaFiles, targetUnshaded);
      cp = (cp == null || cp.isEmpty()) ? builtJar.toString() : (cp + ":" + builtJar);
    }

    JavaCompiler compiler = ToolProvider.getSystemJavaCompiler();
    if (compiler == null) {
      throw new IllegalStateException("No system JavaCompiler available (run with full JDK 21+).");
    }

    DiagnosticCollector<JavaFileObject> diagnostics = new DiagnosticCollector<>();
    try (StandardJavaFileManager fm =
        compiler.getStandardFileManager(diagnostics, Locale.ROOT, StandardCharsets.UTF_8)) {
      List<String> javacOpts = new ArrayList<>();
      javacOpts.add("--release");
      javacOpts.add("21");
      javacOpts.add("-proc:none");
      javacOpts.add("-implicit:none");
      if (cp != null && !cp.isEmpty()) {
        javacOpts.add("-cp");
        javacOpts.add(cp);
      }

      Iterable<? extends JavaFileObject> fileObjects = fm.getJavaFileObjectsFromPaths(javaFiles);
      JavacTask task =
          (JavacTask) compiler.getTask(null, fm, diagnostics, javacOpts, null, fileObjects);

      List<CompilationUnitTree> units = new ArrayList<>();
      for (CompilationUnitTree cu : task.parse()) {
        units.add(cu);
      }
      task.analyze();

      long errCount =
          diagnostics.getDiagnostics().stream()
              .filter(d -> d.getKind() == Diagnostic.Kind.ERROR)
              .count();
      if (errCount > 0) {
        for (Diagnostic<? extends JavaFileObject> d : diagnostics.getDiagnostics()) {
          if (d.getKind() == Diagnostic.Kind.ERROR) {
            System.err.println("JAVAC ERROR: " + d);
          }
        }
        throw new IllegalStateException(
            "Aborting j2go: " + errCount + " javac type-resolution error(s).");
      }

      long tAnalyze = System.currentTimeMillis();
      Trees trees = Trees.instance(task);

      // Step 1: SCC Package Planning
      PackagePlanner planner = PackagePlanner.plan(units, trees, task.getElements(), goModuleBase);
      if (printSccs) {
        System.out.println("=== SCC Go Package Layout (" + planner.getSccGroups().size() + " packages) ===");
        for (PackagePlanner.SccGroup g : planner.getSccGroups()) {
          System.out.printf(
              "  [%02d] %-12s <- %s (deps: %s)%n",
              g.id, g.goPkgName, g.javaPackages, g.depGoPackages);
        }
      }

      // Step 2: Whole-Program Type, Override Union-Find, Capture & External Analysis
      TypeResolver resolver = new TypeResolver(task, planner);
      resolver.analyze(units);

      // Step 3: Emit 1:1 Go Source Files
      ClassEmitter emitter = new ClassEmitter(resolver);
      Path primarySrcRoot = srcDirs.get(0);
      int emittedCount = 0;
      for (PackagePlanner.UnitPlan plan : planner.getUnitPlans()) {
        String goCode = emitter.emitCompilationUnit(plan, primarySrcRoot);
        Path pkgDir = outDir.resolve(plan.goPkg);
        Files.createDirectories(pkgDir);
        Files.writeString(pkgDir.resolve(plan.goFileName), goCode, StandardCharsets.UTF_8);
        emittedCount++;
      }

      // Step 4: Emit Self-Contained `j2gort` Bridge Package
      resolver.runtime.emitRuntimePackage(outDir);

      long tDone = System.currentTimeMillis();
      System.out.printf(
          "j2go completed: %d .java -> %d .go files across %d acyclic SCC Go packages (+ j2gort) "
              + "[javac=%dms, emit=%dms, total=%dms]%n",
          javaFiles.size(),
          emittedCount,
          planner.getSccGroups().size(),
          (tAnalyze - t0),
          (tDone - tAnalyze),
          (tDone - t0));
    }
  }
}

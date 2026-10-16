// Benchmark: the pinned Java Closure Compiler's parser (parsing.parser.Parser.parseProgram, the
// layer ported in internal/parser) over a corpus, K times in one JVM.
//
//   javac -cp closure-compiler.jar -d build ParseBench.java
//   java -Xss1g -cp closure-compiler.jar:build ParseBench <iters> <root> @<listfile> <per-file-tsv>
//
// All files are read into memory first. Each iteration parses every file once, single-threaded,
// with the parser configuration the compiler uses for --language_in=ECMASCRIPT_NEXT. Per iteration
// it prints wall time, this thread's CPU time, process CPU time (includes GC and JIT threads),
// bytes allocated by this thread, and GC count/time. Per-file wall times of the last iteration go
// to <per-file-tsv> (rel-path, bytes, nanos). The Go counterpart is tools/bench/parsebench.
import com.google.javascript.jscomp.CompilerOptions;
import com.google.javascript.jscomp.parsing.parser.Parser;
import com.google.javascript.jscomp.parsing.parser.SourceFile;
import com.google.javascript.jscomp.parsing.parser.util.ErrorReporter;
import com.google.javascript.jscomp.parsing.parser.util.SourcePosition;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.management.GarbageCollectorMXBean;
import java.lang.management.ManagementFactory;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

public final class ParseBench {
  public static void main(String[] args) throws Exception {
    int iters = Integer.parseInt(args[0]);
    Path root = Paths.get(args[1]);
    List<String> names = new ArrayList<>();
    for (String l : Files.readAllLines(Paths.get(args[2].substring(1)), StandardCharsets.UTF_8)) {
      if (!l.isEmpty()) {
        names.add(l);
      }
    }
    Path perFile = Paths.get(args[3]);
    List<String> contents = new ArrayList<>();
    long totalChars = 0;
    for (String n : names) {
      String s = Files.readString(root.resolve(n), StandardCharsets.UTF_8);
      contents.add(s);
      totalChars += s.length();
    }
    Parser.Config config = defaultParserConfig();
    com.sun.management.ThreadMXBean threads =
        (com.sun.management.ThreadMXBean) ManagementFactory.getThreadMXBean();
    com.sun.management.OperatingSystemMXBean os =
        (com.sun.management.OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();
    long tid = Thread.currentThread().getId();
    System.out.println("# java files=" + names.size() + " utf16_chars=" + totalChars);
    System.out.println("iter\twall_ms\tthread_cpu_ms\tprocess_cpu_ms\talloc_mb\tgc_count\tgc_ms");
    long[] last = new long[names.size()];
    int errors = 0;
    for (int it = 1; it <= iters; it++) {
      long gcCount0 = 0, gcTime0 = 0;
      for (GarbageCollectorMXBean gc : ManagementFactory.getGarbageCollectorMXBeans()) {
        gcCount0 += gc.getCollectionCount();
        gcTime0 += gc.getCollectionTime();
      }
      long alloc0 = threads.getThreadAllocatedBytes(tid);
      long cpu0 = threads.getCurrentThreadCpuTime();
      long pcpu0 = os.getProcessCpuTime();
      long wall0 = System.nanoTime();
      errors = 0;
      for (int i = 0; i < names.size(); i++) {
        long t0 = System.nanoTime();
        SourceFile file = new SourceFile(names.get(i), contents.get(i));
        int[] errs = new int[1];
        ErrorReporter reporter =
            new ErrorReporter() {
              @Override
              protected void reportError(SourcePosition location, String message) {
                errs[0]++;
              }

              @Override
              protected void reportWarning(SourcePosition location, String message) {}
            };
        new Parser(config, reporter, file).parseProgram();
        errors += errs[0];
        last[i] = System.nanoTime() - t0;
      }
      long wall = System.nanoTime() - wall0;
      long cpu = threads.getCurrentThreadCpuTime() - cpu0;
      long pcpu = os.getProcessCpuTime() - pcpu0;
      long alloc = threads.getThreadAllocatedBytes(tid) - alloc0;
      long gcCount = -gcCount0, gcTime = -gcTime0;
      for (GarbageCollectorMXBean gc : ManagementFactory.getGarbageCollectorMXBeans()) {
        gcCount += gc.getCollectionCount();
        gcTime += gc.getCollectionTime();
      }
      System.out.printf(
          "%d\t%.0f\t%.0f\t%.0f\t%.0f\t%d\t%d%n",
          it, wall / 1e6, cpu / 1e6, pcpu / 1e6, alloc / 1048576.0, gcCount, gcTime);
    }
    System.out.println("# files_with_errors_reported=" + errors);
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < names.size(); i++) {
      sb.append(names.get(i)).append('\t').append(contents.get(i).length()).append('\t').append(last[i]).append('\n');
    }
    Files.writeString(perFile, sb.toString(), StandardCharsets.UTF_8);
  }

  /** Same configuration as tools/oracle/java/ParseTreeDump (what ParserRunner.parse uses). */
  private static Parser.Config defaultParserConfig() throws Exception {
    com.google.javascript.jscomp.Compiler compiler =
        new com.google.javascript.jscomp.Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    compiler.initOptions(options);
    Class<?> ctxClass = Class.forName("com.google.javascript.jscomp.AbstractCompiler$ConfigContext");
    Object defaultCtx = null;
    for (Object c : ctxClass.getEnumConstants()) {
      if (((Enum<?>) c).name().equals("DEFAULT")) {
        defaultCtx = c;
      }
    }
    Method getParserConfig =
        com.google.javascript.jscomp.Compiler.class.getDeclaredMethod("getParserConfig", ctxClass);
    getParserConfig.setAccessible(true);
    Object config = getParserConfig.invoke(compiler, defaultCtx);
    Method newParserConfig =
        com.google.javascript.jscomp.parsing.ParserRunner.class.getDeclaredMethod(
            "newParserConfig", com.google.javascript.jscomp.parsing.Config.class);
    newParserConfig.setAccessible(true);
    return (Parser.Config) newParserConfig.invoke(null, config);
  }
}

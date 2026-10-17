// Benchmark driver for the "Java warm" rows of the end-to-end benchmark (PLAN.md §6,
// tools/bench/bench.py): repeats the same CommandLineRunner invocations in one JVM, which models a
// build system's persistent worker and separates JIT warm-up from steady-state speed. It is the
// benchmark counterpart of tools/oracle/java/CliBatch.java (which runs many different invocations
// once each, for oracles).
//
//   javac -cp closure-compiler.jar -d build CliBatch.java
//   java -cp closure-compiler.jar:build CliBatch [--keep-stderr] <warmup> <iters> <jobfile> <outdir>
//
// <jobfile> holds one invocation per line, its arguments tab-separated. One iteration runs every
// job once, in order; <warmup> unmeasured iterations come first, then <iters> measured ones. Each
// job runs as `java -jar closure-compiler.jar <args>` would (CommandLineRunner.main): relative
// paths resolve against this JVM's working directory, the PhaseOptimizer logger is off, stdout and
// stderr are encoded as System.out and System.err would encode them, and the exit status is what
// the shell would see. For job i, the last iteration's stdout, stderr and exit status are written
// to <outdir>/<i>.stdout, <i>.stderr and <i>.exit. With --keep-stderr, each measured iteration's
// stderr of job i also goes to <outdir>/<i>.<iter>.stderr (to read --tracer_mode reports).
//
// Per iteration it prints one tab-separated line:
//   iter kind wall_ms process_cpu_ms gc_count gc_ms rss_mb start_epoch_ms end_epoch_ms
// where kind is "warmup" or "measure", process CPU includes the JIT and GC threads, and rss_mb is
// VmRSS after the iteration. The last line is "# peak_rss_mb=<VmHWM>", the peak over the JVM's life.
import com.google.javascript.jscomp.CommandLineRunner;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.PrintStream;
import java.io.UncheckedIOException;
import java.lang.management.GarbageCollectorMXBean;
import java.lang.management.ManagementFactory;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;

public final class CliBatch extends CommandLineRunner {
  private CliBatch(String[] args, PrintStream out, PrintStream err) {
    super(args, InputStream.nullInputStream(), out, err);
  }

  /** The exit status CommandLineRunner.main would end the process with. */
  private int runAsMain() {
    int[] status = {0};
    if (shouldRunCompiler()) {
      setExitCodeReceiver(
          code -> {
            status[0] = code;
            return null;
          });
      run();
    } else if (hasErrors()) {
      status[0] = -1;
    }
    return status[0] & 0xff;
  }

  private static Charset streamCharset(String property) {
    String name = System.getProperty(property, System.getProperty("native.encoding", "UTF-8"));
    try {
      return Charset.forName(name);
    } catch (RuntimeException e) {
      return Charset.defaultCharset();
    }
  }

  private static long statusKb(String field) throws IOException {
    for (String line : Files.readAllLines(Paths.get("/proc/self/status"), StandardCharsets.UTF_8)) {
      if (line.startsWith(field + ":")) {
        return Long.parseLong(line.substring(field.length() + 1).replace("kB", "").trim());
      }
    }
    return -1;
  }

  public static void main(String[] argv) {
    try {
      batch(argv);
    } catch (IOException e) {
      throw new UncheckedIOException(e);
    }
  }

  private static void batch(String[] argv) throws IOException {
    List<String> a = new ArrayList<>(Arrays.asList(argv));
    boolean keepStderr = a.remove("--keep-stderr");
    int warmup = Integer.parseInt(a.get(0));
    int iters = Integer.parseInt(a.get(1));
    List<String[]> jobs = new ArrayList<>();
    for (String line : Files.readAllLines(Paths.get(a.get(2)), StandardCharsets.UTF_8)) {
      if (!line.isEmpty()) {
        jobs.add(line.split("\t", -1));
      }
    }
    Path outDir = Paths.get(a.get(3));
    Files.createDirectories(outDir);

    // As CommandLineRunner.main.
    Logger.getLogger("com.google.javascript.jscomp.PhaseOptimizer").setLevel(Level.OFF);
    Charset outCharset = streamCharset("stdout.encoding");
    Charset errCharset = streamCharset("stderr.encoding");
    com.sun.management.OperatingSystemMXBean os =
        (com.sun.management.OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();
    PrintStream report = System.out;
    report.println("iter\tkind\twall_ms\tprocess_cpu_ms\tgc_count\tgc_ms\trss_mb\tstart_epoch_ms\tend_epoch_ms");

    for (int it = 1; it <= warmup + iters; it++) {
      boolean measured = it > warmup;
      long gcCount0 = 0, gcTime0 = 0;
      for (GarbageCollectorMXBean gc : ManagementFactory.getGarbageCollectorMXBeans()) {
        gcCount0 += gc.getCollectionCount();
        gcTime0 += gc.getCollectionTime();
      }
      long startEpoch = System.currentTimeMillis();
      long cpu0 = os.getProcessCpuTime();
      long wall0 = System.nanoTime();
      for (int j = 0; j < jobs.size(); j++) {
        ByteArrayOutputStream outBuf = new ByteArrayOutputStream();
        ByteArrayOutputStream errBuf = new ByteArrayOutputStream();
        PrintStream out = new PrintStream(outBuf, true, outCharset);
        PrintStream err = new PrintStream(errBuf, true, errCharset);
        int status;
        try {
          status = new CliBatch(jobs.get(j), out, err).runAsMain();
        } catch (Throwable t) {
          // An exception escaping main: the JVM's default handler, then status 1.
          err.print("Exception in thread \"main\" ");
          t.printStackTrace(err);
          status = 1;
        }
        out.flush();
        err.flush();
        if (it == warmup + iters) {
          Files.write(outDir.resolve(j + ".stdout"), outBuf.toByteArray());
          Files.write(outDir.resolve(j + ".stderr"), errBuf.toByteArray());
          Files.writeString(outDir.resolve(j + ".exit"), status + "\n");
        }
        if (keepStderr && measured) {
          Files.write(outDir.resolve(j + "." + (it - warmup) + ".stderr"), errBuf.toByteArray());
        }
      }
      long wall = System.nanoTime() - wall0;
      long cpu = os.getProcessCpuTime() - cpu0;
      long endEpoch = System.currentTimeMillis();
      long gcCount = -gcCount0, gcTime = -gcTime0;
      for (GarbageCollectorMXBean gc : ManagementFactory.getGarbageCollectorMXBeans()) {
        gcCount += gc.getCollectionCount();
        gcTime += gc.getCollectionTime();
      }
      report.printf(
          "%d\t%s\t%.1f\t%.1f\t%d\t%d\t%.0f\t%d\t%d%n",
          measured ? it - warmup : it - warmup - 1,
          measured ? "measure" : "warmup",
          wall / 1e6,
          cpu / 1e6,
          gcCount,
          gcTime,
          statusKb("VmRSS") / 1024.0,
          startEpoch,
          endEpoch);
      report.flush();
    }
    report.printf("# peak_rss_mb=%.0f%n", statusKb("VmHWM") / 1024.0);
  }
}

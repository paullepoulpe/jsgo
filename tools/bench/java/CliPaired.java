// Original benchmark helper. Uses the unmodified pinned Closure Compiler API.
// One fresh CommandLineRunner/Compiler and normal extern parsing per request.
import com.google.javascript.jscomp.CommandLineRunner;
import java.io.*;
import java.lang.management.*;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.*;
import java.util.logging.*;

public final class CliPaired extends CommandLineRunner {
  private CliPaired(String[] args, PrintStream out, PrintStream err) {
    super(args, InputStream.nullInputStream(), out, err);
  }
  private int runAsMain() {
    int[] status = {0};
    if (shouldRunCompiler()) {
      setExitCodeReceiver(code -> { status[0] = code; return null; });
      run();
    } else if (hasErrors()) { status[0] = -1; }
    return status[0] & 255;
  }
  private static String cg(String name) throws IOException {
    return Files.readString(Path.of("/sys/fs/cgroup", name)).trim().replace('\n', ';');
  }
  private static long status(String key) throws IOException {
    for (String line : Files.readAllLines(Path.of("/proc/self/status"))) {
      if (line.startsWith(key + ":")) return Long.parseLong(line.substring(key.length()+1).replace("kB", "").trim());
    }
    throw new IOException("missing process metric: " + key);
  }
  private static long[] gc() {
    long count = 0, ms = 0;
    for (GarbageCollectorMXBean bean : ManagementFactory.getGarbageCollectorMXBeans()) {
      count += Math.max(0, bean.getCollectionCount()); ms += Math.max(0, bean.getCollectionTime());
    }
    return new long[]{count, ms};
  }
  private static long cpu() {
    return ((com.sun.management.OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean()).getProcessCpuTime();
  }
  private static long jit() {
    CompilationMXBean bean = ManagementFactory.getCompilationMXBean();
    return bean != null && bean.isCompilationTimeMonitoringSupported() ? bean.getTotalCompilationTime() : -1;
  }
  public static void main(String[] args) {
    try { serve(args); } catch (IOException error) { throw new UncheckedIOException(error); }
  }
  private static void serve(String[] args) throws IOException {
    if (!cg("memory.max").equals("4294967296") || !cg("memory.swap.max").equals("0")) {
      throw new IOException("requires kernel cgroup memory.max=4GiB and memory.swap.max=0");
    }
    Logger.getLogger("com.google.javascript.jscomp.PhaseOptimizer").setLevel(Level.OFF);
    List<String> job = Files.readAllLines(Path.of(args[0]), StandardCharsets.UTF_8);
    // One argument per line; embedded newlines are rejected by the controller.
    String[] compilerArgs = job.toArray(String[]::new);
    Charset stdout = Charset.forName(System.getProperty("stdout.encoding", System.getProperty("native.encoding", "UTF-8")));
    Charset stderr = Charset.forName(System.getProperty("stderr.encoding", System.getProperty("native.encoding", "UTF-8")));
    System.out.println("READY\t"+cg("memory.max")+"\t"+cg("memory.swap.max")+"\t"+cg("cpu.max")+"\t"+cg("cpuset.cpus.effective"));
    System.out.flush();
    BufferedReader input = new BufferedReader(new InputStreamReader(System.in, StandardCharsets.UTF_8));
    for (String command; (command = input.readLine()) != null;) {
      if (command.equals("STOP")) break;
      String[] parts = command.split("\\t", -1);
      if (parts.length != 2 || !parts[0].matches("[0-9]+") || !parts[1].matches("warmup|measure")) {
        throw new IOException("invalid benchmark request");
      }
      Path directory = Path.of("protocol", parts[0]); Files.createDirectories(directory);
      long[] gcBefore = gc(); long jitBefore = jit(), cpuBefore = cpu(), start = System.nanoTime();
      int code;
      // Both runtimes write compiler stdout/stderr to real files inside their
      // timed interval. The protocol uses separate System.out, never a result
      // cache or a cheaper Java-only in-memory diagnostic sink.
      try (PrintStream output = new PrintStream(Files.newOutputStream(directory.resolve("stdout")), true, stdout);
           PrintStream errors = new PrintStream(Files.newOutputStream(directory.resolve("stderr")), true, stderr)) {
        code = new CliPaired(compilerArgs, output, errors).runAsMain();
      }
      long wall = System.nanoTime()-start, processCpu = cpu()-cpuBefore, jitAfter = jit(); long[] gcAfter = gc();
      Files.writeString(directory.resolve("exit"), code+"\n", StandardCharsets.US_ASCII);
      System.out.println(String.join("\t", "RESULT", parts[0], parts[1], Long.toString(wall), Long.toString(processCpu),
          Long.toString(gcAfter[0]-gcBefore[0]), Long.toString(gcAfter[1]-gcBefore[1]), Long.toString(jitBefore), Long.toString(jitAfter),
          Long.toString(status("VmRSS")), Long.toString(status("VmHWM")), Integer.toString(code), cg("memory.peak"), cg("memory.events")));
      System.out.flush();
    }
  }
}

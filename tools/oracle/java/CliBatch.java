// Oracle tool: runs the pinned Java Closure Compiler CLI (CommandLineRunner)
// in-process for many invocations, so oracle generation pays JVM startup once
// instead of once per fixture.
//
// stdin: one job per line, tab-separated: <outfile> <arg> <arg> ...
// Each <outfile> receives the combined stdout+stderr of that run, followed by
//   ---EXIT---
//   <exit code as the shell would see it>
// which is the same format scripts/gen-oracles.sh produces with `java -jar`.
import com.google.javascript.jscomp.CommandLineRunner;
import java.io.BufferedReader;
import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.logging.Level;
import java.util.logging.Logger;

public final class CliBatch extends CommandLineRunner {
  private CliBatch(String[] args, PrintStream out, PrintStream err) {
    super(args, InputStream.nullInputStream(), out, err);
  }

  private int runOnce(PrintStream err) {
    if (!shouldRunCompiler()) {
      return hasErrors() ? -1 : 0;
    }
    int result;
    try {
      result = doRun();
    } catch (com.google.javascript.jscomp.FlagUsageException e) {
      err.println(e.getMessage());
      result = -1;
    } catch (Throwable t) {
      t.printStackTrace(err);
      result = -2;
    }
    return result;
  }

  public static void main(String[] argv) {
    try {
      batch();
    } catch (java.io.IOException e) {
      throw new java.io.UncheckedIOException(e);
    }
  }

  private static void batch() throws java.io.IOException {
    Logger.getLogger("com.google.javascript.jscomp.PhaseOptimizer").setLevel(Level.OFF);
    BufferedReader in = new BufferedReader(new InputStreamReader(System.in, StandardCharsets.UTF_8));
    String line;
    while ((line = in.readLine()) != null) {
      if (line.isEmpty()) {
        continue;
      }
      String[] parts = line.split("\t");
      Path outFile = Paths.get(parts[0]);
      String[] args = Arrays.copyOfRange(parts, 1, parts.length);
      ByteArrayOutputStream buf = new ByteArrayOutputStream();
      PrintStream ps = new PrintStream(buf, true, StandardCharsets.UTF_8);
      int code;
      try {
        CliBatch runner = new CliBatch(args, ps, ps);
        code = runner.runOnce(ps);
      } catch (Throwable t) {
        t.printStackTrace(ps);
        code = -2;
      }
      ps.flush();
      if (outFile.getParent() != null) {
        Files.createDirectories(outFile.getParent());
      }
      String body = buf.toString(StandardCharsets.UTF_8);
      Files.writeString(outFile, body + "---EXIT---\n" + (code & 0xff) + "\n", StandardCharsets.UTF_8);
    }
  }
}

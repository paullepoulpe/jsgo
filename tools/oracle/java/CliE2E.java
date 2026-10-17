// Oracle tool for the end-to-end CLI comparison (tools/clie2e): runs the
// pinned CommandLineRunner exactly as CommandLineRunner.main does under
// `java -jar`, for many invocations in one JVM.
//
// stdin: one job per line, tab-separated: <resultdir> <arg> <arg> ...
// The working directory is the JVM's, so jobs name their inputs and outputs
// relative to it. If <resultdir>/stdin exists it is the job's System.in
// (otherwise System.in is empty). Each job writes <resultdir>/stdout,
// <resultdir>/stderr (encoded like System.err: with stderr.encoding) and
// <resultdir>/exit (the status the shell would see).
// A job that names a file of tools/oracle/gpl-excluded-corpus.txt is skipped
// (CorpusExclusions; run from the repository or pass -Djsgo.excludedCorpus).
import com.google.javascript.jscomp.CommandLineRunner;
import java.io.BufferedReader;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.PrintStream;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.logging.Level;
import java.util.logging.Logger;
import jsgo.oracle.CorpusExclusions;

public final class CliE2E extends CommandLineRunner {
  private CliE2E(String[] args, InputStream in, PrintStream out, PrintStream err) {
    super(args, in, out, err);
  }

  /** Thrown by the exit code receiver to stop the job like System.exit. */
  private static final class Exit extends RuntimeException {
    final int status;

    Exit(int status) {
      super(null, null, false, false);
      this.status = status;
    }
  }

  /** SystemExitCodeReceiver's mapping of the code passed to System.exit. */
  private static int exitStatus(int exitCodeValue) {
    byte exitCodeByte = (byte) exitCodeValue;
    if (exitCodeByte == 0 && exitCodeValue != 0) {
      exitCodeByte = (byte) -1;
    }
    return exitCodeByte & 0xff;
  }

  public static void main(String[] argv) {
    try {
      batch();
    } catch (java.io.IOException e) {
      throw new java.io.UncheckedIOException(e);
    }
  }

  private static void batch() throws java.io.IOException {
    Charset errCharset = Charset.forName(System.getProperty("stderr.encoding", "UTF-8"));
    BufferedReader jobs =
        new BufferedReader(new InputStreamReader(System.in, StandardCharsets.UTF_8));
    String line;
    while ((line = jobs.readLine()) != null) {
      if (line.isEmpty()) {
        continue;
      }
      String[] parts = line.split("\t", -1);
      Path dir = Paths.get(parts[0]);
      Files.createDirectories(dir);
      String[] args = Arrays.copyOfRange(parts, 1, parts.length);
      if (Arrays.stream(args).anyMatch(CorpusExclusions::isExcluded)) {
        // tools/oracle/gpl-excluded-corpus.txt: no output is made from these files.
        System.err.println(
            "CorpusExclusions: skipped job " + parts[0] + ": it names a file listed in "
                + CorpusExclusions.LIST);
        continue;
      }
      Path stdinFile = dir.resolve("stdin");
      InputStream in =
          Files.exists(stdinFile)
              ? new ByteArrayInputStream(Files.readAllBytes(stdinFile))
              : new ByteArrayInputStream(new byte[0]);
      ByteArrayOutputStream outBuf = new ByteArrayOutputStream();
      ByteArrayOutputStream errBuf = new ByteArrayOutputStream();
      PrintStream out = new PrintStream(outBuf, true, StandardCharsets.UTF_8);
      PrintStream err = new PrintStream(errBuf, true, errCharset);
      int status;
      try {
        // CommandLineRunner.main:
        Logger.getLogger("com.google.javascript.jscomp.PhaseOptimizer").setLevel(Level.OFF);
        CliE2E runner = new CliE2E(args, in, out, err);
        runner.setExitCodeReceiver(
            code -> {
              throw new Exit(exitStatus(code));
            });
        status = 0;
        try {
          if (runner.shouldRunCompiler()) {
            runner.run();
          }
          if (runner.hasErrors()) {
            status = exitStatus(-1);
          }
        } catch (Exit e) {
          status = e.status;
        }
      } catch (Throwable t) {
        // The default uncaught-exception handler.
        err.print("Exception in thread \"main\" ");
        t.printStackTrace(err);
        status = 1;
      }
      out.flush();
      err.flush();
      Files.write(dir.resolve("stdout"), outBuf.toByteArray());
      Files.write(dir.resolve("stderr"), errBuf.toByteArray());
      Files.writeString(dir.resolve("exit"), status + "\n");
    }
  }
}

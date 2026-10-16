import com.google.javascript.jscomp.diagnostic.LogFile;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Vectors for internal/diagnostic (WritingLogFile, NoOpLogFile): runs fixed scenarios against
 * the pinned jar's LogFile in a temporary directory and prints, per scenario, every call's
 * result (or exception) and the bytes of the log file.
 *
 * <pre>
 * java -cp $CP LogFileVectors <tmpdir> > log_file.txt
 * </pre>
 *
 * Output lines are escaped like RegExpDump; file contents are printed as their bytes, one
 * "bytes" line of two-digit hex per scenario. internal/diagnostic/log_file_test.go runs the same
 * scenarios.
 */
public final class LogFileVectors {
  private static StringBuilder out = new StringBuilder();

  interface Step {
    Object run() throws Exception;
  }

  private static void step(String name, Step s) {
    String result;
    try {
      Object r = s.run();
      result = r instanceof LogFile ? "ok" : String.valueOf(r);
    } catch (Throwable e) {
      result = e.getClass().getName() + ": " + e.getMessage();
      if (e.getCause() != null) {
        result += " / cause " + e.getCause().getClass().getName() + ": " + e.getCause().getMessage();
      }
    }
    out.append(RegExpDump.escape(name + " -> " + result)).append('\n');
  }

  private static void bytes(Path file) throws Exception {
    StringBuilder sb = new StringBuilder("bytes");
    if (Files.exists(file)) {
      sb.append(' ');
      for (byte b : Files.readAllBytes(file)) {
        sb.append(String.format("%02x", b & 0xff));
      }
    } else {
      sb.append(" <missing>");
    }
    out.append(sb).append('\n');
  }

  public static void main(String[] args) throws Exception {
    Path tmp = Path.of(args[0]);

    out.append("=== writing\n");
    Path f1 = tmp.resolve("a/b/one.log");
    LogFile[] log = new LogFile[1];
    step("createOrReopen", () -> log[0] = LogFile.createOrReopen(f1));
    step("isLogging", () -> log[0].isLogging());
    step("log(String)", () -> log[0].log("plain <&> é 😀"));
    step("log(Object)", () -> log[0].log((Object) List.of(1, "x")));
    step("log(Object null)", () -> log[0].log((Object) null));
    step("log(Supplier)", () -> log[0].log(() -> "supplied"));
    step("log(template)", () -> log[0].log("%s-%d-%s", "x", 3, null));
    step("logJson(Object)", () -> log[0].logJson((Object) Map.of("k", List.of("<v>"))));
    step("logJson(Supplier)", () -> log[0].logJson(() -> (Object) "s "));
    step("log high", () -> log[0].log("hi\ud83d"));
    step("log low", () -> log[0].log("\ude00lo"));
    step("log lone", () -> log[0].log("lone\ud800!"));
    step("log lone low", () -> log[0].log("\udc00"));
    step("close", () -> {
          log[0].close();
          return "void";
        });
    step("close again", () -> {
          log[0].close();
          return "void";
        });
    step("log after close", () -> log[0].log("x"));
    bytes(f1);

    out.append("=== reopen appends\n");
    step("createOrReopen", () -> log[0] = LogFile.createOrReopen(f1));
    step("log(String)", () -> log[0].log("second"));
    step("close", () -> {
          log[0].close();
          return "void";
        });
    bytes(f1);

    out.append("=== streamed json\n");
    Path f2 = tmp.resolve("streamed.log");
    step("createOrReopen", () -> log[0] = LogFile.createOrReopen(f2));
    step("log(String)", () -> log[0].log("before"));
    step(
        "logJson(producer)",
        () ->
            log[0].logJson(
                w -> {
                  w.beginObject();
                  w.name("a<").beginArray().value(1).value("<x>").nullValue().endArray();
                  w.name("n").nullValue();
                  w.endObject();
                }));
    step("log after producer", () -> log[0].log("after"));
    step("close", () -> {
          log[0].close();
          return "void";
        });
    bytes(f2);

    out.append("=== incomplete streamed json\n");
    Path f3 = tmp.resolve("incomplete.log");
    step("createOrReopen", () -> log[0] = LogFile.createOrReopen(f3));
    step("logJson(producer)", () -> log[0].logJson(w -> w.beginArray().value(true)));
    bytes(f3);

    out.append("=== throwing producer\n");
    Path f4 = tmp.resolve("throwing.log");
    step("createOrReopen", () -> log[0] = LogFile.createOrReopen(f4));
    step(
        "logJson(producer)",
        () ->
            log[0].logJson(
                w -> {
                  w.beginArray();
                  throw new IllegalStateException("boom");
                }));
    bytes(f4);

    out.append("=== bad paths\n");
    step("relative file without parent", () -> LogFile.createOrReopen(Path.of("rel-no-parent.log")));
    Files.writeString(tmp.resolve("afile"), "x");
    step("parent is a file", () -> LogFile.createOrReopen(tmp.resolve("afile/x.log")));
    Files.createDirectories(tmp.resolve("adir"));
    step("file is a directory", () -> LogFile.createOrReopen(tmp.resolve("adir")));

    out.append("=== no-op\n");
    step("createNoOp", () -> log[0] = LogFile.createNoOp());
    step("isLogging", () -> log[0].isLogging());
    step("log(Object null)", () -> log[0].log((Object) null));
    step("logJson(producer)", () -> log[0].logJson(w -> {
          throw new AssertionError("not called");
        }));
    step("close", () -> {
          log[0].close();
          return "void";
        });

    System.out.print(out);
  }
}

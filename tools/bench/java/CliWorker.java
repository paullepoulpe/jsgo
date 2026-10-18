// Persistent JSON worker for comparable Java/jsgo benchmarks. Each request uses
// a fresh CommandLineRunner/compiler. Optional immutable extern reuse is added
// by the explicitly selected benchmark overlay, never by modifying the jar.
import com.google.javascript.jscomp.CommandLineRunner;
import com.google.javascript.jscomp.parsing.BuiltinExternCache;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonElement;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonObject;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonParser;
import java.io.BufferedReader;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.InputStreamReader;
import java.io.IOException;
import java.io.UncheckedIOException;
import java.io.PrintStream;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.util.Base64;
import java.util.logging.Level;
import java.util.logging.Logger;

public final class CliWorker extends CommandLineRunner {
  private CliWorker(String[] args, byte[] input, PrintStream out, PrintStream err) {
    super(args, new ByteArrayInputStream(input), out, err);
  }
  private int runAsMain() {
    int[] status = {0};
    if (shouldRunCompiler()) {
      setExitCodeReceiver(code -> { status[0] = code; return null; });
      run();
    } else if (hasErrors()) { status[0] = -1; }
    return status[0] & 0xff;
  }
  private static Charset streamCharset(String property) {
    String name = System.getProperty(property, System.getProperty("native.encoding", "UTF-8"));
    try { return Charset.forName(name); }
    catch (RuntimeException e) { return Charset.defaultCharset(); }
  }
  public static void main(String[] args) {
    try { runWorker(); }
    catch (IOException error) { throw new UncheckedIOException(error); }
  }
  private static void runWorker() throws IOException {
    Logger.getLogger("com.google.javascript.jscomp.PhaseOptimizer").setLevel(Level.OFF);
    PrintStream protocol = System.out;
    Charset outCharset = streamCharset("stdout.encoding");
    Charset errCharset = streamCharset("stderr.encoding");
    var input = new BufferedReader(new InputStreamReader(System.in, StandardCharsets.UTF_8));
    for (String line; (line = input.readLine()) != null;) {
      JsonObject request;
      try { request = JsonParser.parseString(line).getAsJsonObject(); }
      catch (RuntimeException error) {
        System.err.println("java-worker: invalid request: " + error);
        System.exit(2);
        return;
      }
      var values = request.getAsJsonArray("args");
      String[] flags = new String[values.size()];
      for (int i = 0; i < flags.length; i++) { flags[i] = values.get(i).getAsString(); }
      byte[] stdin = request.has("stdin_base64")
          ? Base64.getDecoder().decode(request.get("stdin_base64").getAsString()) : new byte[0];
      var outBytes = new ByteArrayOutputStream();
      var errBytes = new ByteArrayOutputStream();
      var out = new PrintStream(outBytes, true, outCharset);
      var err = new PrintStream(errBytes, true, errCharset);
      PrintStream oldOut = System.out, oldErr = System.err;
      int exit;
      try {
        System.setOut(out);
        System.setErr(err);
        exit = new CliWorker(flags, stdin, out, err).runAsMain();
      } catch (Throwable error) {
        err.print("Exception in thread \"main\" ");
        error.printStackTrace(err);
        exit = 1;
      } finally {
        out.flush();
        err.flush();
        System.setOut(oldOut);
        System.setErr(oldErr);
      }
      var response = new JsonObject();
      JsonElement id = request.get("id");
      if (id != null) { response.add("id", id); }
      response.addProperty("exit", exit);
      response.addProperty("stdout_base64", Base64.getEncoder().encodeToString(outBytes.toByteArray()));
      response.addProperty("stderr_base64", Base64.getEncoder().encodeToString(errBytes.toByteArray()));
      protocol.println(response);
      protocol.flush();
    }
    if (Boolean.getBoolean("jsgo.bench.verifyCache")) {
      System.err.println("java-worker cache: " + BuiltinExternCache.statistics());
    }
  }
}

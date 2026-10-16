// Oracle tool: writes internal/jutil/testdata/throwable_message.tsv, the JDK's
// Throwable.getMessage() (or null) and toString() for the exceptions that
// jutil.GetMessageOrNull and jutil.InterruptedException model.
//   java -cp .cache/oracle-build JutilThrowableVectors internal/jutil/testdata
import java.io.PrintStream;
import java.nio.file.*;

public class JutilThrowableVectors {
  static void row(PrintStream out, String name, Throwable t) {
    String m = t.getMessage();
    out.println(name + "\t" + (m == null ? "<null>" : "=" + m) + "\t" + t);
  }

  public static void main(String[] args) throws Exception {
    Path dir = Paths.get(args[0]);
    try (PrintStream out = new PrintStream(Files.newOutputStream(dir.resolve("throwable_message.tsv")), true, "UTF-8")) {
      row(out, "InterruptedException()", new InterruptedException());
      row(out, "RuntimeException(InterruptedException())", new RuntimeException(new InterruptedException()));
      row(out, "RuntimeException(\"\")", new RuntimeException(""));
      row(out, "RuntimeException(\"msg\")", new RuntimeException("msg"));
      row(out, "NullPointerException()", new NullPointerException());
      row(out, "IllegalStateException()", new IllegalStateException());
      row(out, "IllegalStateException(\"x\")", new IllegalStateException("x"));
    }
  }
}

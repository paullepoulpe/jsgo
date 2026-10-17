import java.io.InputStream;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Map;

/**
 * Oracle for jscomp.resources.PropertiesParser (internal/resources).
 *
 * <pre>
 * java -cp $CP PropertiesParserDump <inputs.txt> > properties_parser.txt
 * </pre>
 *
 * Each line of inputs.txt is one input, escaped like RegExpDump (backslash doubled, chars outside
 * printable ASCII as backslash-u and four hex digits), or "@resource:<path>" for a classpath
 * resource of the pinned jar (a shaded com/google/javascript/jscomp/jarjar/ path also finds the
 * unshaded file). For each input it prints "=== " + the escaped input, then either
 * "exception <class>: <message>" or one "<key>=<value>" line per entry of
 * PropertiesParser.parse(input), in map order, escaped.
 */
public final class PropertiesParserDump {
  private static final String JARJAR = "com/google/javascript/jscomp/jarjar/";

  public static void main(String[] args) throws Exception {
    Class<?> parser = Class.forName("com.google.javascript.jscomp.resources.PropertiesParser");
    Method parse = parser.getDeclaredMethod("parse", String.class);
    parse.setAccessible(true);
    StringBuilder sb = new StringBuilder();
    for (String l : Files.readAllLines(Paths.get(args[0]), StandardCharsets.UTF_8)) {
      String input;
      if (l.startsWith("@resource:")) {
        String path = l.substring("@resource:".length());
        InputStream in = PropertiesParserDump.class.getResourceAsStream("/" + path);
        if (in == null && path.startsWith(JARJAR)) {
          // The Maven jar shades its dependencies under JARJAR; the npm one
          // (google-closure-compiler-java) bundles the same files unshaded.
          in = PropertiesParserDump.class.getResourceAsStream("/" + path.substring(JARJAR.length()));
        }
        input = new String(in.readAllBytes(), StandardCharsets.UTF_8);
      } else {
        input = RegExpDump.unescape(l);
      }
      sb.append("=== ").append(RegExpDump.escape(input)).append('\n');
      try {
        @SuppressWarnings("unchecked")
        Map<String, String> m = (Map<String, String>) parse.invoke(null, input);
        for (Map.Entry<String, String> e : m.entrySet()) {
          sb.append(RegExpDump.escape(e.getKey() + "=" + e.getValue())).append('\n');
        }
      } catch (InvocationTargetException ite) {
        Throwable e = ite.getCause();
        sb.append(RegExpDump.escape("exception " + e.getClass().getName() + ": " + e.getMessage()))
            .append('\n');
      }
    }
    System.out.print(sb);
  }
}

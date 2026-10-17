// Oracle tool: writes internal/args4j/testdata/args4j_vectors.txt, the
// behaviour of the args4j 2.33 that the pinned closure-compiler.jar shades
// (com.google.javascript.jscomp.jarjar.org.kohsuke.args4j) on option beans
// that use every feature Closure's command-line tools use: printUsage under
// several widths, filters and ParserProperties, printExample,
// printSingleLineUsage, and parseArgument on hand-written and seeded random
// argument vectors (with @-files), each followed by the bean's state and the
// exception thrown. The Go side, internal/args4j/args4j_vectors_test.go,
// declares the same beans, replays every "## " header and compares the whole
// text. Run it from internal/args4j (the @-files are relative paths there):
//   (cd internal/args4j && java -cp ../../.cache/$CLOSURE_JAR:../../.cache/oracle-build \
//      Args4jVectors > testdata/args4j_vectors.txt)
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.Argument;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.CmdLineException;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.CmdLineParser;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.NamedOptionDef;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.Option;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.OptionDef;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.OptionHandlerFilter;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.ParserProperties;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.spi.FieldSetter;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.spi.IntOptionHandler;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.spi.OptionHandler;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.spi.Parameters;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.spi.Setter;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.spi.StringArrayOptionHandler;
import com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.spi.StringOptionHandler;
import java.io.PrintStream;
import java.io.StringWriter;
import java.lang.reflect.AnnotatedElement;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.Set;
import java.util.function.Supplier;

public class Args4jVectors {
  enum Level { QUIET, DEFAULT, VERBOSE }

  enum Mode { BROWSER, NODE, WEBPACK, BROWSER_WITH_TRANSFORMED_PREFIXES }

  enum Fmt {
    PRETTY_PRINT,
    SINGLE_QUOTES;

    @Override
    public String toString() {
      return name().toLowerCase(java.util.Locale.ROOT);
    }
  }

  enum OutFmt { MD, JSON }

  /** CommandLineRunner.Flags.BooleanOptionHandler. */
  public static class ClosureBool extends OptionHandler<Boolean> {
    private static final Set<String> TRUES = Set.of("true", "on", "yes", "1");
    private static final Set<String> FALSES = Set.of("false", "off", "no", "0");

    public ClosureBool(CmdLineParser parser, OptionDef option, Setter<? super Boolean> setter) {
      super(parser, option, setter);
    }

    @Override
    public int parseArguments(Parameters params) throws CmdLineException {
      String param = null;
      try {
        param = params.getParameter(0);
      } catch (CmdLineException e) {
        param = null;
      }
      if (param == null) {
        setter.addValue(true);
        return 0;
      } else {
        String lowerParam = asciiLower(param);
        if (TRUES.contains(lowerParam)) {
          setter.addValue(true);
        } else if (FALSES.contains(lowerParam)) {
          setter.addValue(false);
        } else {
          setter.addValue(true);
          return 0;
        }
        return 1;
      }
    }

    @Override
    public String getDefaultMetaVariable() {
      return null;
    }
  }

  static String asciiLower(String s) {
    StringBuilder sb = new StringBuilder();
    for (char c : s.toCharArray()) {
      sb.append(c >= 'A' && c <= 'Z' ? (char) (c + 32) : c);
    }
    return sb.toString();
  }

  static final List<String> recorded = new ArrayList<>();

  /** CommandLineRunner.Flags.JsOptionHandler with its MultiFlagSetter. */
  public static class RecordingStringHandler extends StringOptionHandler {
    public RecordingStringHandler(CmdLineParser parser, OptionDef option, Setter<? super String> setter) {
      super(parser, option, new Recording(setter));
    }
  }

  @SuppressWarnings("unchecked")
  static class Recording implements Setter<String> {
    final Setter<? super String> proxy;

    Recording(Setter<? super String> proxy) {
      this.proxy = proxy;
    }

    @Override
    public boolean isMultiValued() {
      return proxy.isMultiValued();
    }

    @Override
    public Class<String> getType() {
      return (Class<String>) proxy.getType();
    }

    @Override
    public void addValue(String value) throws CmdLineException {
      String normalizedValue = value;
      if (value != null
          && value.length() > 0
          && (value.substring(0, 1).equals("'") || value.substring(0, 1).equals("\""))
          && value.substring(value.length() - 1).equals(value.substring(0, 1))) {
        normalizedValue = value.substring(1, value.length() - 1);
      }
      proxy.addValue(normalizedValue);
      recorded.add("JS:" + normalizedValue);
    }

    @Override
    public FieldSetter asFieldSetter() {
      return proxy.asFieldSetter();
    }

    @Override
    public AnnotatedElement asAnnotatedElement() {
      return proxy.asAnnotatedElement();
    }
  }

  static class BeanA {
    @Option(name = "--zeta", usage = "Last alphabetically but declared first")
    String zeta = "z";

    @Option(name = "--help", handler = ClosureBool.class, help = true,
        usage = "Displays this message on stdout and exit")
    boolean help = false;

    @Option(name = "--year",
        usage = "shortcut for defining goog.FEATURESET_YEAR=YYYY. The minimum valid value of the browser year is 2012")
    Integer year = 0;

    @Option(name = "--level", aliases = {"-W"},
        usage = "Specifies the warning level to use. Options: QUIET, DEFAULT, VERBOSE")
    Level level = Level.DEFAULT;

    @Option(name = "--mode", hidden = true, usage = "A hidden enum option")
    Mode mode = Mode.BROWSER;

    @Option(name = "--nullmode", usage = "A nullable enum option")
    Mode nullMode = null;

    @Option(name = "--fmt", usage = "Formatting options; may be repeated")
    List<Fmt> fmt = new ArrayList<>();

    @Option(name = "--js", handler = RecordingStringHandler.class,
        usage = "The JavaScript filename. You may specify multiple. The flag name is optional, "
            + "because args are interpreted as files by default. You may also use minimatch-style "
            + "glob patterns. For example, use --js='**.js' --js='!**_test.js' to recursively "
            + "include all js files that do not end in _test.js")
    List<String> js = new ArrayList<>();

    @Option(name = "--externs",
        usage = "The file containing JavaScript externs. You may specify multiple")
    List<String> externs = new ArrayList<>();

    @Option(name = "--define", aliases = {"--D", "-D"},
        usage = "Override the value of a variable annotated @define.\nThe format is <name>[=<val>], "
            + "where <name> is the name of a @define variable and <val> is a boolean, number, or a "
            + "single-quoted string that contains no single quotes. If [=<val>] is omitted, the "
            + "variable is marked true")
    List<String> define = new ArrayList<>();

    @Option(name = "--threads", handler = IntOptionHandler.class, usage = "Number of threads")
    int threads = 1;

    @Option(name = "--debug", handler = ClosureBool.class, usage = "Enable debugging options")
    boolean debug = false;

    @Option(name = "--plain", usage = "A boolean handled by args4j's own handler")
    boolean plain = false;

    @Option(name = "--boxed", usage = "A boxed Boolean")
    Boolean boxed = null;

    @Option(name = "--out", metaVar = "FILE",
        usage = "Primary output filename. If not specified, output is written to stdout")
    String out = "";

    @Option(name = "--prefix", usage = "A nullable string")
    String prefix = null;

    @Option(name = "--map", usage = "Prefixes to replace in ES6 import paths before resolving.")
    Map<String, String> map = new HashMap<>();

    @Option(name = "--nullmap", usage = "A map that starts null")
    Map<String, String> nullMap = null;

    @Option(name = "--words", handler = StringArrayOptionHandler.class, usage = "Greedy words")
    List<String> words = new ArrayList<>();

    @Option(name = "--nousage")
    String noUsage = "nu";

    @Option(name = "--dep", depends = {"--out"}, usage = "Needs --out")
    String dep = null;

    @Option(name = "--forb", forbids = {"--debug", "--plain"}, usage = "Cannot be used with --debug")
    String forb = null;

    @Option(name = "--nullint", usage = "A nullable Integer")
    Integer nullInt = null;

    @Option(name = "--ünï",
        usage = "Non-ASCII option name: ünïcödé usage text that is long "
            + "enough to wrap around the column limit at least once   here")
    String unicode = "é";

    @Argument List<String> arguments = new ArrayList<>();

    String dump() {
      return "zeta=" + zeta + " help=" + help + " year=" + year + " level=" + level + " mode=" + mode
          + " nullMode=" + nullMode + " fmt=" + fmt + " js=" + js + " externs=" + externs
          + " define=" + define + " threads=" + threads + " debug=" + debug + " plain=" + plain
          + " boxed=" + boxed + " out=" + out + " prefix=" + prefix + " map=" + map
          + " nullMap=" + nullMap + " words=" + words + " noUsage=" + noUsage + " dep=" + dep
          + " forb=" + forb + " nullInt=" + nullInt + " unicode=" + unicode
          + " arguments=" + arguments + " recorded=" + recorded;
    }
  }

  static class BeanB {
    @Option(name = "--mapping_file", required = true,
        usage = "The file name of the mapping generated by the production instrumentation pass.")
    String mapping = "";

    @Option(name = "--reports_directory", required = true,
        usage = "The folder/directory which contains all the reports created by the instrumented"
            + " production code.")
    String reports = "";

    @Option(name = "--result_output", required = true,
        usage = "Use this flag to provide the name of the final report that will be generated by"
            + "this reporter.")
    String result = "";

    @Option(name = "-h", help = true, usage = "Help")
    boolean h = false;

    String dump() {
      return "mapping=" + mapping + " reports=" + reports + " result=" + result + " h=" + h;
    }
  }

  static class BeanC {
    @Option(name = "--fix", usage = "Fix lint warnings automatically")
    boolean fix = false;

    @Option(name = "--resolution_mode",
        usage = "Expected path resolution mode for resolving module specifiers.")
    Mode resolutionMode = Mode.BROWSER;

    @Argument List<String> files = new ArrayList<>();

    String dump() {
      return "fix=" + fix + " resolutionMode=" + resolutionMode + " files=" + files;
    }
  }

  static class BeanD {
    @Option(name = "--output", usage = "Sets the desired output format.\n" + "Options: md, json" + "Default: json")
    OutFmt output = OutFmt.JSON;

    @Option(name = "--help", usage = "Usage")
    boolean help = false;

    @Argument(index = 1, metaVar = "SECOND", usage = "The second argument")
    String second = null;

    @Argument(index = 0, required = true, usage = "The first argument")
    String first = null;

    @Argument(index = 2, usage = "A boolean argument")
    boolean flag = false;

    String dump() {
      return "output=" + output + " help=" + help + " second=" + second + " first=" + first + " flag=" + flag;
    }
  }

  static class BeanDupName {
    @Option(name = "-x", usage = "x")
    String x;

    @Option(name = "-y", aliases = {"-x"}, usage = "y")
    String y;
  }

  static class BeanDupIndex {
    @Argument(index = 0)
    String a;

    @Argument(index = 0)
    String b;
  }

  static class BeanFinal {
    @Option(name = "-f", usage = "final")
    final String f = "";
  }

  static class BeanGap {
    @Argument(index = 1)
    String a;
  }

  static PrintStream out;

  static String esc(String s) {
    StringBuilder sb = new StringBuilder();
    for (char c : s.toCharArray()) {
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c == '\t') {
        sb.append("\\t");
      } else if (c == '\n') {
        sb.append("\\n");
      } else if (c == '\r') {
        sb.append("\\r");
      } else if (c < 0x20 || c > 0x7e) {
        sb.append(String.format("\\u%04x", (int) c));
      } else {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  static void block(String text) {
    for (String line : text.split("\n", -1)) {
      out.println("| " + esc(line));
    }
  }

  static Object newBean(String name) {
    switch (name) {
      case "A": return new BeanA();
      case "B": return new BeanB();
      case "C": return new BeanC();
      case "D": return new BeanD();
      case "DupName": return new BeanDupName();
      case "DupIndex": return new BeanDupIndex();
      case "Final": return new BeanFinal();
      case "Gap": return new BeanGap();
      default: throw new IllegalArgumentException(name);
    }
  }

  static String dump(Object bean) {
    if (bean instanceof BeanA) return ((BeanA) bean).dump();
    if (bean instanceof BeanB) return ((BeanB) bean).dump();
    if (bean instanceof BeanC) return ((BeanC) bean).dump();
    if (bean instanceof BeanD) return ((BeanD) bean).dump();
    return "-";
  }

  static final Set<String> CATEGORY = Set.of("js", "level", "externs", "out", "define", "help");

  static String usage(Object bean, String config) {
    ParserProperties props = ParserProperties.defaults();
    switch (config) {
      case "nodefaults": props.withShowDefaults(false); break;
      case "eqdelim": props.withOptionValueDelimiter("="); break;
      case "unsorted": props.withOptionSorter(null); break;
      default: break;
    }
    CmdLineParser parser = new CmdLineParser(bean, props);
    StringWriter w = new StringWriter();
    switch (config) {
      case "all":
        parser.printUsage(w, null, OptionHandlerFilter.ALL);
        break;
      case "w6": case "w20": case "w40": case "w60": case "w120": case "w5000":
        parser.setUsageWidth(Integer.parseInt(config.substring(1)));
        parser.printUsage(w, null);
        break;
      case "category":
        parser.printUsage(w, null, (optionHandler) -> {
          if (optionHandler.option instanceof NamedOptionDef namedOptionDef) {
            return !optionHandler.option.hidden()
                && CATEGORY.contains(namedOptionDef.name().replaceFirst("^--", ""));
          }
          return false;
        });
        break;
      case "example":
        w.write("[" + parser.printExample(OptionHandlerFilter.ALL) + "]\n");
        w.write("[" + parser.printExample(OptionHandlerFilter.REQUIRED) + "]\n");
        w.write("[" + parser.printExample(OptionHandlerFilter.PUBLIC) + "]\n");
        break;
      case "single":
        parser.printSingleLineUsage(w, null);
        break;
      default:
        parser.printUsage(w, null);
        break;
    }
    return w.toString();
  }

  static void runUsage(String bean, String config) {
    out.println("## usage " + bean + " " + config);
    try {
      block(usage(newBean(bean), config));
    } catch (Throwable t) {
      out.println("throw " + esc(t.toString()));
    }
  }

  static void runParse(String beanName, String... args) {
    StringBuilder header = new StringBuilder("## parse " + beanName);
    for (String a : args) {
      header.append('\t').append(esc(a));
    }
    out.println(header);
    recorded.clear();
    Object bean = newBean(beanName);
    try {
      CmdLineParser parser = new CmdLineParser(bean);
      try {
        parser.parseArgument(args);
        out.println("ok");
      } catch (CmdLineException e) {
        out.println("error " + esc(e.toString()));
        out.println("message " + esc(e.getMessage()));
        out.println("parser " + (e.getParser() == parser));
      }
      out.println("state " + esc(dump(bean)));
      if (Arrays.asList(args).contains("--help")) {
        StringWriter w = new StringWriter();
        parser.printUsage(w, null);
        block(w.toString());
      }
    } catch (Throwable t) {
      out.println("throw " + esc(t.toString()));
    }
  }

  static final String[] TOKENS_A = {
    "--zeta", "--help", "--year", "--level", "-W", "--mode", "--nullmode", "--fmt", "--js",
    "--externs", "--define", "--D", "-D", "--threads", "--debug", "--plain", "--boxed", "--out",
    "--prefix", "--map", "--nullmap", "--words", "--nousage", "--dep", "--forb", "--nullint",
    "--ünï", "--nope", "-", "--", "a.js", "'quoted.js'", "\"dq.js\"", "'", "x=y", "=v",
    "k=", "true", "off", "YES", "maybe", "0", "1", "verbose", "Verbose", "quiet", "node",
    "browser-with-transformed-prefixes", "pretty_print", "SINGLE-QUOTES", "12", "-3", "1e3",
    "2147483648", "", " ", "a b", "a  b c", "--js=foo.js", "--level=QUIET", "--define=x=1",
    "--threads=x", "--map=k=v", "--map=k", "--out=", "--boxed=false", "--help=yes",
    "@testdata/at1.txt", "@testdata/missing.txt", "@testdata", "@testdata/bom.txt", "ünï",
    "--js foo.js", "-W VERBOSE",
  };

  static final String[] TOKENS_B = {
    "--mapping_file", "--reports_directory", "--result_output", "-h", "m.json", "dir", "out.tsv",
    "--mapping_file=x", "extra", "--nope",
  };

  static final String[] TOKENS_C = {
    "--fix", "--resolution_mode", "NODE", "webpack", "bogus", "a.js", "b.js", "--fix=true", "-",
  };

  static final String[] TOKENS_D = {
    "--output", "md", "JSON", "xml", "--help", "first", "second", "true", "OFF", "maybe", "4th",
    "--output=md", "Yes",
  };

  static void random(Random rnd, String bean, String[] tokens, int n) {
    for (int i = 0; i < n; i++) {
      int len = rnd.nextInt(7);
      String[] args = new String[len];
      for (int j = 0; j < len; j++) {
        args[j] = tokens[rnd.nextInt(tokens.length)];
      }
      runParse(bean, args);
    }
  }

  public static void main(String[] argv) throws Exception {
    out = new PrintStream(System.out, true, "UTF-8");
    String[] configs = {
      "default", "all", "w6", "w20", "w40", "w60", "w120", "w5000", "nodefaults", "eqdelim",
      "unsorted", "category", "example", "single",
    };
    for (String bean : new String[] {"A", "B", "C", "D"}) {
      for (String c : configs) {
        runUsage(bean, c);
      }
    }
    for (String bean : new String[] {"DupName", "DupIndex", "Final", "Gap"}) {
      if (!bean.equals("Gap")) {
        // (Gap's usage throws the JVM's helpful NullPointerException message,
        // which describes bytecode.)
        runUsage(bean, "default");
      }
      runParse(bean, "-x", "v");
      runParse(bean, "v");
    }

    // Hand-written vectors.
    runParse("A");
    runParse("A", "--help");
    runParse("A", "--help", "--level", "VERBOSE", "--js", "a.js", "--js", "b.js", "--threads", "7");
    runParse("A", "--js", "a.js", "b.js", "--externs", "e.js", "c.js");
    runParse("A", "--js=a.js", "--js='q.js'", "--js", "\"d.js\"", "--js", "''", "--js", "'");
    runParse("A", "-W", "quiet", "--level", "Verbose", "--level=default");
    runParse("A", "--level", "LOUD");
    runParse("A", "--level");
    runParse("A", "--mode", "browser-with-transformed-prefixes", "--nullmode", "node");
    runParse("A", "--fmt", "pretty_print", "--fmt", "single-quotes", "--fmt", "PRETTY_PRINT");
    runParse("A", "--fmt", "bad");
    runParse("A", "--define", "a=1", "--D", "b", "-D", "c='x'", "--define=d=2");
    runParse("A", "--threads", "12");
    runParse("A", "--threads", "abc");
    runParse("A", "--threads", "2147483648");
    runParse("A", "--threads", "-0");
    runParse("A", "--threads", "+5");
    runParse("A", "--threads=abc");
    runParse("A", "--year", "2020", "--nullint", "-7");
    runParse("A", "--year", "x");
    runParse("A", "--debug", "--plain", "--boxed", "a.js");
    runParse("A", "--debug", "off", "--debug", "ON", "--debug", "maybe");
    runParse("A", "--debug=false");
    runParse("A", "--plain=true");
    runParse("A", "--out", "o.js", "--prefix", "p", "--dep", "d");
    runParse("A", "--dep", "d");
    runParse("A", "--forb", "f", "--debug");
    runParse("A", "--forb", "f", "--plain");
    runParse("A", "--forb", "f");
    runParse("A", "--map", "a=b", "--map", "c=", "--map", "d=e=f", "--map", "zz=1", "--map", "a=c");
    runParse("A", "--map", "nokey");
    runParse("A", "--map", "=v");
    runParse("A", "--nullmap", "k=v", "--nullmap", "k2=v2");
    runParse("A", "--words", "a", "b c", "d  e", "--debug", "x");
    runParse("A", "--words");
    runParse("A", "--words", "", " ", "-x");
    runParse("A", "--nousage", "n");
    runParse("A", "--nope");
    runParse("A", "-");
    runParse("A", "--");
    runParse("A", "--zeta");
    runParse("A", "--zeta", "--debug");
    runParse("A", "--js foo.js");
    runParse("A", "-W VERBOSE");
    runParse("A", "--ünï", "à😀");
    runParse("A", "@testdata/at1.txt");
    runParse("A", "@testdata/at2.txt", "--level", "QUIET");
    runParse("A", "@testdata/bom.txt");
    runParse("A", "@testdata/missing.txt");
    runParse("A", "@testdata");
    runParse("A", "@testdata//at1.txt/");
    runParse("A", "@");
    runParse("A", "--js", "@testdata/at3.txt");
    runParse("A", "a.js", "-", "b.js");
    runParse("A", "--help", "--year", "2018", "--fmt", "pretty_print", "--fmt", "single_quotes",
        "--map", "x=1", "--map", "y=", "--prefix", "", "--boxed", "--nullint", "3", "--words", "w1",
        "w2", "--level", "QUIET", "--nullmode", "WEBPACK", "--externs", "e1", "--externs", "e2");
    runParse("B");
    runParse("B", "-h");
    runParse("B", "--mapping_file", "m", "--reports_directory", "r", "--result_output", "o");
    runParse("B", "--mapping_file", "m", "--result_output", "o");
    runParse("B", "--mapping_file", "m", "x");
    runParse("C", "--fix", "a.js", "b.js", "--resolution_mode", "NODE", "c.js");
    runParse("C", "--resolution_mode", "bogus");
    runParse("D");
    runParse("D", "first");
    runParse("D", "first", "second", "TRUE");
    runParse("D", "first", "second", "maybe");
    runParse("D", "first", "second", "off", "extra");
    runParse("D", "--help");
    runParse("D", "--output", "md", "--help");
    runParse("D", "--output", "xml");
    runParse("D", "--output=json", "f");

    Random rnd = new Random(20261016L);
    random(rnd, "A", TOKENS_A, 1500);
    random(rnd, "B", TOKENS_B, 200);
    random(rnd, "C", TOKENS_C, 200);
    random(rnd, "D", TOKENS_D, 300);
  }
}

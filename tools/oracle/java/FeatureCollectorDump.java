// Oracle tool: for each file, parses it with ParserRunner.parse (the compiler's parsing.Config for
// --language_in=ECMASCRIPT_NEXT, as ParsingDump) and runs jscomp.parsing.FeatureCollector over the
// resulting SCRIPT. Prints one line per file: <rel-file> TAB <FeatureCollector.allFeatures()>
// (or "null" when parsing produced no AST). The Go side is TestFeatureCollectorGolden /
// TestFeatureCollectorOracle in internal/parsing/feature_collector_test.go.
//
//   javac -cp closure-compiler.jar -d build FeatureCollectorDump.java
//   java -Xss1g -cp closure-compiler.jar:build FeatureCollectorDump <root> <rel-file>... | @<listfile>
import com.google.javascript.jscomp.CompilerOptions;
import com.google.javascript.jscomp.SourceFile;
import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.FeatureCollector;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

public final class FeatureCollectorDump {
  public static void main(String[] args) throws Exception {
    Path root = Paths.get(args[0]);
    List<String> files = new ArrayList<>();
    for (int i = 1; i < args.length; i++) {
      String f = args[i];
      if (f.startsWith("@")) {
        for (String l : Files.readAllLines(Paths.get(f.substring(1)), StandardCharsets.UTF_8)) {
          if (!l.isEmpty()) {
            files.add(l);
          }
        }
      } else {
        files.add(f);
      }
    }
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
    Config config = (Config) getParserConfig.invoke(compiler, defaultCtx);
    PrintStream out = new PrintStream(System.out, true, StandardCharsets.UTF_8);
    for (String rel : files) {
      String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
      Node ast;
      try {
        ast = ParserRunner.parse(SourceFile.fromCode(rel, src), src, config, ErrorReporter.NULL_INSTANCE).ast;
      } catch (RuntimeException e) {
        ast = null;
      }
      if (ast == null) {
        out.println(rel + "\tnull");
        continue;
      }
      FeatureCollector fc = new FeatureCollector();
      fc.visitScript(ast);
      out.println(rel + "\t" + fc.allFeatures());
    }
  }
}

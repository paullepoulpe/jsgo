// Oracle tool: evaluates every JSDoc type expression in real-world files with the pinned Java
// JSTypeRegistry and dumps facts about the resulting types. It checks internal/rhino's type
// system on the types real code writes, not just hand-picked ones (TypeAlgebraDump).
//
//   javac -cp closure-compiler.jar -d build CorpusTypeDump.java
//   java -Xss1g -cp closure-compiler.jar:build CorpusTypeDump [--gzip] <outdir> <root> <rel-file>...|@<list>
//
// For each file: parse with ParserRunner.parse and the compiler's parsing.Config (as
// ParsingDump.java), create a fresh JSTypeRegistry, walk the AST in pre-order, and for every
// JSDocInfo, every JSTypeExpression in getTypeExpressions() order, evaluate it with
// registry.evaluateTypeExpressionInGlobalScope and print one line of facts. Registry warnings and
// errors are printed as they are reported. Every line is escaped to printable ASCII (as
// ParsingDump.java). The Go side is internal/parsing/corpus_types_oracle_test.go.
import com.google.javascript.jscomp.CompilerOptions;
import com.google.javascript.jscomp.SourceFile;
import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.JSTypeExpression;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.jstype.JSType;
import com.google.javascript.rhino.jstype.JSTypeNative;
import com.google.javascript.rhino.jstype.JSTypeRegistry;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class CorpusTypeDump {
  public static void main(String[] args) throws Exception {
    boolean gzip = false;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--gzip")) {
        gzip = true;
      } else {
        rest.add(a);
      }
    }
    Path outDir = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    List<String> files = new ArrayList<>();
    for (String f : rest.subList(2, rest.size())) {
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
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    Config config = compilerParserConfig();
    for (String rel : files) {
      StringBuilder sb = new StringBuilder();
      dumpOne(sb, config, rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8));
      Path out = outDir.resolve(rel + (gzip ? ".types.gz" : ".types"));
      Files.createDirectories(out.getParent());
      try (OutputStream os =
          gzip ? new GZIPOutputStream(Files.newOutputStream(out)) : Files.newOutputStream(out)) {
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  private static Config compilerParserConfig() throws Exception {
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
    Method m = com.google.javascript.jscomp.Compiler.class.getDeclaredMethod("getParserConfig", ctxClass);
    m.setAccessible(true);
    return (Config) m.invoke(compiler, defaultCtx);
  }

  private static void dumpOne(StringBuilder sb, Config config, String name, String src) {
    ErrorReporter parseReporter = ErrorReporter.NULL_INSTANCE;
    ParserRunner.ParseResult result =
        ParserRunner.parse(SourceFile.fromCode(name, src), src, config, parseReporter);
    if (result.ast == null) {
      line(sb, "# no ast");
      return;
    }
    ErrorReporter typeReporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int lineno, int lineOffset) {
            line(sb, "WARNING " + lineno + ":" + lineOffset + " " + message);
          }

          @Override
          public void error(String message, String sourceName, int lineno, int lineOffset) {
            line(sb, "ERROR " + lineno + ":" + lineOffset + " " + message);
          }
        };
    JSTypeRegistry registry = new JSTypeRegistry(typeReporter);
    JSType number = registry.getNativeType(JSTypeNative.NUMBER_TYPE);
    JSType string = registry.getNativeType(JSTypeNative.STRING_TYPE);
    JSType object = registry.getNativeType(JSTypeNative.OBJECT_TYPE);
    walk(sb, result.ast, registry, number, string, object);
  }

  private static void walk(
      StringBuilder sb, Node n, JSTypeRegistry registry, JSType number, JSType string, JSType object) {
    JSDocInfo info = n.getJSDocInfo();
    if (info != null) {
      for (JSTypeExpression expr : info.getTypeExpressions()) {
        Node root = expr.getRoot();
        String where = n.getLineno() + ":" + n.getCharno() + " @" + root.getLineno() + ":" + root.getCharno();
        JSType t;
        try {
          t = registry.evaluateTypeExpressionInGlobalScope(expr);
        } catch (RuntimeException e) {
          line(sb, where + " EXCEPTION " + e.getClass().getName() + ": " + e.getMessage());
          continue;
        }
        line(
            sb,
            where
                + " | "
                + t
                + " | ann="
                + t.toAnnotationString(JSType.Nullability.EXPLICIT)
                + " | unk="
                + t.isUnknownType()
                + " null="
                + t.isNullable()
                + " void="
                + t.isVoidable()
                + " union="
                + t.isUnionType()
                + " fn="
                + t.isFunctionType()
                + " rec="
                + t.isRecordType()
                + " tmpl="
                + t.isTemplatizedType()
                + " enumel="
                + t.isEnumElementType()
                + " noresolved="
                + t.isNoResolvedType()
                + " obj="
                + t.isObject()
                + " | <:num="
                + t.isSubtypeOf(number)
                + " <:str="
                + t.isSubtypeOf(string)
                + " <:obj="
                + t.isSubtypeOf(object)
                + " | nonnull="
                + t.restrictByNotNullOrUndefined()
                + " | autobox="
                + t.autobox()
                + " | lub(num)="
                + t.getLeastSupertype(number)
                + " | glb(obj)="
                + t.getGreatestSubtype(object));
      }
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      walk(sb, c, registry, number, string, object);
    }
  }

  private static void line(StringBuilder sb, String s) {
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c >= 0x20 && c < 0x7f) {
        sb.append(c);
      } else {
        sb.append(String.format("\\u%04x", (int) c));
      }
    }
    sb.append('\n');
  }
}

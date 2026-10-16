// Oracle tool: dumps the exact rhino.Node tree that the pinned Java Closure
// Compiler produces when parsing a file, plus every parse diagnostic.
//
// Build/run (see tools/oracle/README.md):
//   javac -cp closure-compiler.jar -d build AstDump.java
//   java -cp closure-compiler.jar:build AstDump [--lang=ECMASCRIPT_NEXT] [--jsdoc] <outdir> <root> <file>...
//
// For each <file> (a path relative to <root>), writes <outdir>/<file>.ast. The
// Go port's AST printer must reproduce this format byte for byte.
import com.google.javascript.jscomp.CompilerOptions;
import com.google.javascript.jscomp.JSError;
import com.google.javascript.jscomp.SourceFile;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.Node;
import java.io.IOException;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

public final class AstDump {
  public static void main(String[] args) throws IOException {
    CompilerOptions.LanguageMode lang = CompilerOptions.LanguageMode.ECMASCRIPT_NEXT;
    boolean jsdoc = false;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.startsWith("--lang=")) {
        lang = CompilerOptions.LanguageMode.valueOf(a.substring("--lang=".length()));
      } else if (a.equals("--jsdoc")) {
        jsdoc = true;
      } else {
        rest.add(a);
      }
    }
    Path outDir = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    for (String rel : rest.subList(2, rest.size())) {
      Path out = outDir.resolve(rel + ".ast");
      Files.createDirectories(out.getParent());
      try (PrintStream ps = new PrintStream(Files.newOutputStream(out), false, StandardCharsets.UTF_8)) {
        dumpOne(ps, root.resolve(rel), rel, lang, jsdoc);
      } catch (RuntimeException e) {
        try (PrintStream ps = new PrintStream(Files.newOutputStream(out), false, StandardCharsets.UTF_8)) {
          ps.println("# crash " + e.getClass().getName());
        }
      }
    }
  }

  private static void dumpOne(
      PrintStream ps, Path file, String name, CompilerOptions.LanguageMode lang, boolean jsdoc)
      throws IOException {
    com.google.javascript.jscomp.Compiler compiler =
        new com.google.javascript.jscomp.Compiler(new PrintStream(java.io.OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(lang);
    options.setParseJsDocDocumentation(
        jsdoc ? com.google.javascript.jscomp.parsing.Config.JsDocParsing.INCLUDE_DESCRIPTIONS_NO_WHITESPACE
              : com.google.javascript.jscomp.parsing.Config.JsDocParsing.TYPES_ONLY);
    compiler.initOptions(options);
    String src = Files.readString(file, StandardCharsets.UTF_8);
    Node rootNode = compiler.parse(SourceFile.fromCode(name, src));
    ps.println("# diagnostics");
    for (JSError e : compiler.getErrors()) {
      ps.println("ERROR " + e.getType().key + " " + e.getLineNumber() + ":" + e.getCharno());
    }
    for (JSError e : compiler.getWarnings()) {
      ps.println("WARNING " + e.getType().key + " " + e.getLineNumber() + ":" + e.getCharno());
    }
    ps.println("# tree");
    dumpTree(ps, rootNode, 0);
  }

  private static void dumpTree(PrintStream ps, Node n, int level) {
    String indent = "    ".repeat(level);
    ps.print(indent);
    ps.println(n.toString(true, true, false).replace("\n", "\\n"));
    JSDocInfo info = n.getJSDocInfo();
    if (info != null) {
      ps.print(indent);
      ps.println("  @jsdoc " + info.toStringVerbose().replace("\n", "\\n"));
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      dumpTree(ps, c, level + 1);
    }
  }
}

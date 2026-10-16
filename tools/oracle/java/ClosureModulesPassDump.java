// Differential oracle for the l11-closure-modules passes (internal/jscomp): ProcessClosurePrimitives,
// ProcessClosureProvidesAndRequires, ClosureCheckModule, ClosureOptimizePrimitives,
// ModuleImportResolver (isGoogModuleDependencyCall), modules/ClosureRequireProcessor,
// ConvertChunksToESModules, Es6CheckModule, Es6RenameReferences, Es6RenameTypeReferences,
// WhitespaceWrapGoogModules, ExtraRequireRemover, CodingConventions, deps/DepsFileRegexParser,
// ProcessTweaks, ExportTestFunctions, FindExportableNodes and GatherRawExports.
//
// For every file, each pass runs on a fresh Compiler set up the way CompilerTestCase sets one up:
// default CompilerOptions with --language_in=ECMASCRIPT_NEXT and --language_out=NO_TRANSPILE,
// Compiler.init(no externs, [SourceFile.fromCode(<rel>, <src>)], options) and Compiler.parse().
// Files with parse errors print their diagnostics and are skipped. The passes get
// (externsRoot, jsRoot), as compiler passes do; a pass that needs GatherModuleMetadata gets it run
// first (processCommonJsModules=false, ResolutionMode.BROWSER), as DefaultPassConfig does.
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside printable
// ASCII as \\uXXXX):
//   ## file <rel>
//   # parse ok|errors                then the parse diagnostics:
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   ## pass <name>
//   diag ...                         the diagnostics the pass added (the parse ones are removed,
//                                    one occurrence each), in ErrorManager order
//   exception <Throwable.toString()> if the pass threw (without the JVM's class-loader suffix of a
//                                    ClassCastException message, see CCE_LOADER_SUFFIX)
//   tree unchanged                   if jsRoot.toStringTree() did not change, else
//   tree changed <sha256>            the first 16 hex digits of the SHA-256 of the escaped lines
//                                    of jsRoot.toStringTree() (each followed by '\n'), and with
//                                    --full the escaped lines themselves
//   <pass-specific lines>
// The "conventions" section (no pass) prints, in preorder over the parsed jsRoot, the non-default
// answers of CodingConventions.getDefault(), ModuleImportResolver.isGoogModuleDependencyCall,
// ClosureRequireProcessor.getAllRequires and ConvertChunksToESModules.isDynamicImportCallback.
// The "deps" section prints what DepsFileRegexParser (shortcut mode off) finds in the file.
//
// Usage: java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.ClosureModulesPassDump \
//          [--full] [--sample=<N>] [--max-size=<B>] <out.gz> <root> <rel-file>... | --list=<listfile>
// --max-size=B prints "# skipped size <bytes>" instead of running anything on files larger than B
// bytes (the corpus run uses it to stay within memory and time).
// --full prints changed trees in full (for the checked-in golden set and for debugging a diff).
// --sample=N runs the passes that rewrite every file (convertChunksToESModules) on every Nth file
// only (in list order, starting with the first); the default is 1 (every file).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.deps.DependencyInfo;
import com.google.javascript.jscomp.deps.DepsFileRegexParser;
import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.HashBasedTable;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.Table;
import com.google.javascript.jscomp.modules.ClosureModulesPassDumpRequires;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;
import java.util.zip.GZIPOutputStream;

public final class ClosureModulesPassDump {

  /** A pass run: gets the compiler and (externsRoot, jsRoot); may append extra output lines. */
  private interface Run {
    void run(Compiler compiler, Node externs, Node js, List<String> extra);
  }

  private static int sample = 1;
  private static boolean full = false;
  private static long maxSize = Long.MAX_VALUE;

  public static void main(String[] args) throws Exception {
    List<String> rest = new ArrayList<>();
    List<String> files = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.startsWith("--max-size=")) {
        maxSize = Long.parseLong(a.substring("--max-size=".length()));
      } else if (a.startsWith("--sample=")) {
        sample = Integer.parseInt(a.substring("--sample=".length()));
      } else if (a.startsWith("--list=")) {
        for (String l :
            Files.readAllLines(Paths.get(a.substring("--list=".length())), StandardCharsets.UTF_8)) {
          if (!l.isEmpty()) {
            files.add(l);
          }
        }
      } else {
        rest.add(a);
      }
    }
    Path out = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    files.addAll(rest.subList(2, rest.size()));
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(out), true)) {
      int index = 0;
      for (String rel : files) {
        StringBuilder sb = new StringBuilder();
        long size = Files.size(root.resolve(rel));
        if (size > maxSize) {
          line(sb, "## file " + rel);
          line(sb, "# skipped size " + size);
        } else {
          String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
          dumpFile(sb, rel, src, index % sample == 0);
        }
        index++;
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
        os.flush();
        if (index % 500 == 0) {
          System.err.println("progress " + index + "/" + files.size());
        }
      }
    }
  }

  private static Compiler newCompiler(String rel, String src) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    compiler.init(
        ImmutableList.<SourceFile>of(), ImmutableList.of(SourceFile.fromCode(rel, src)), options);
    compiler.parse();
    return compiler;
  }

  private static List<String> diagnostics(Compiler compiler) {
    List<String> out = new ArrayList<>();
    for (JSError e : compiler.getErrors()) {
      out.add(diag("ERROR", e));
    }
    for (JSError e : compiler.getWarnings()) {
      out.add(diag("WARNING", e));
    }
    return out;
  }

  private static String diag(String level, JSError e) {
    return "diag "
        + level
        + " "
        + e.getType().key
        + " "
        + e.getLineno()
        + ":"
        + e.getCharno()
        + " \""
        + e.getDescription()
        + "\"";
  }

  private static void dumpFile(StringBuilder sb, String rel, String src, boolean sampled) {
    line(sb, "## file " + rel);
    Compiler parsed = newCompiler(rel, src);
    List<String> parseDiags = diagnostics(parsed);
    if (parsed.hasErrors()) {
      line(sb, "# parse errors");
      for (String d : parseDiags) {
        line(sb, d);
      }
      return;
    }
    line(sb, "# parse ok");
    for (String d : parseDiags) {
      line(sb, d);
    }

    pass(sb, rel, src, parseDiags, "closureCheckModule", (compiler, externs, js, extra) -> {
      gatherModuleMetadata(compiler, externs, js);
      new ClosureCheckModule(compiler, compiler.getModuleMetadataMap()).process(externs, js);
    });
    pass(sb, rel, src, parseDiags, "processClosurePrimitives", (compiler, externs, js, extra) -> {
      gatherModuleMetadata(compiler, externs, js);
      new ProcessClosurePrimitives(compiler).process(externs, js);
    });
    pass(sb, rel, src, parseDiags, "collectProvidedNames", (compiler, externs, js, extra) -> {
      ProcessClosureProvidesAndRequires p = new ProcessClosureProvidesAndRequires(compiler, true);
      Map<String, ProcessClosureProvidesAndRequires.ProvidedName> names =
          p.collectProvidedNames(externs, js);
      for (Map.Entry<String, ProcessClosureProvidesAndRequires.ProvidedName> e :
          names.entrySet()) {
        ProcessClosureProvidesAndRequires.ProvidedName pn = e.getValue();
        extra.add(
            "provided "
                + e.getKey()
                + " "
                + pn
                + " ns="
                + pn.getNamespace()
                + " implicitInit="
                + pn.hasImplicitInitialization()
                + " legacy="
                + pn.isFromLegacyModule()
                + " first="
                + nodeRef(pn.getFirstProvideCall())
                + " candidate="
                + nodeRef(pn.getCandidateDefinition()));
      }
      extra.add("exported " + p.getExportedVariableNames());
    });
    pass(
        sb, rel, src, parseDiags, "processClosureProvidesAndRequires",
        (compiler, externs, js, extra) -> {
          ProcessClosureProvidesAndRequires p =
              new ProcessClosureProvidesAndRequires(compiler, false);
          p.process(externs, js);
          extra.add("exported " + p.getExportedVariableNames());
        });
    pass(sb, rel, src, parseDiags, "closureOptimizePrimitives", (compiler, externs, js, extra) -> {
      new ClosureOptimizePrimitives(compiler, true).process(externs, js);
    });
    pass(
        sb, rel, src, parseDiags, "closureOptimizePrimitivesEs5",
        (compiler, externs, js, extra) -> {
          new ClosureOptimizePrimitives(compiler, false).process(externs, js);
        });
    pass(sb, rel, src, parseDiags, "extraRequireRemover", (compiler, externs, js, extra) -> {
      new ExtraRequireRemover(compiler).process(externs, js);
    });
    pass(sb, rel, src, parseDiags, "es6CheckModule", (compiler, externs, js, extra) -> {
      new Es6CheckModule(compiler).process(externs, js);
    });
    pass(sb, rel, src, parseDiags, "processTweaks", (compiler, externs, js, extra) -> {
      new ProcessTweaks(compiler, true).process(externs, js);
    });
    pass(sb, rel, src, parseDiags, "exportTestFunctions", (compiler, externs, js, extra) -> {
      new ExportTestFunctions(compiler, "goog.exportSymbol", "goog.exportProperty")
          .process(externs, js);
    });
    for (boolean allowLocalExports : new boolean[] {false, true}) {
      pass(
          sb, rel, src, parseDiags, "findExportableNodes local=" + allowLocalExports,
          (compiler, externs, js, extra) -> {
            FindExportableNodes f = new FindExportableNodes(compiler, allowLocalExports);
            NodeTraversal.traverse(compiler, js, f);
            for (Map.Entry<String, Node> e : f.getExports().entrySet()) {
              extra.add("export " + e.getKey() + " " + nodeRef(e.getValue()));
            }
            for (String s : f.getLocalExports()) {
              extra.add("local " + s);
            }
            for (Map.Entry<Node, String> e : f.getEs6ClassExports().entrySet()) {
              extra.add("es6class " + nodeRef(e.getKey()) + " " + e.getValue());
            }
          });
    }
    pass(sb, rel, src, parseDiags, "gatherRawExports", (compiler, externs, js, extra) -> {
      compiler.setLifeCycleStage(AbstractCompiler.LifeCycleStage.NORMALIZED);
      GatherRawExports g = new GatherRawExports(compiler);
      g.process(externs, js);
      extra.add("exported " + g.getExportedVariableNames());
    });
    pass(sb, rel, src, parseDiags, "whitespaceWrapGoogModules", (compiler, externs, js, extra) -> {
      new WhitespaceWrapGoogModules(compiler).process(externs, js);
    });
    if (sampled) {
      pass(
          sb, rel, src, parseDiags, "convertChunksToESModules",
          (compiler, externs, js, extra) -> {
            new ConvertChunksToESModules(compiler).process(externs, js);
          });
    }
    pass(sb, rel, src, parseDiags, "es6RenameReferences", (compiler, externs, js, extra) -> {
      NodeTraversal.traverse(
          compiler, js, new Es6RenameReferences(renameTable(compiler, js, false)));
    });
    pass(sb, rel, src, parseDiags, "es6RenameTypeReferences", (compiler, externs, js, extra) -> {
      NodeTraversal.traverse(
          compiler, js, new Es6RenameTypeReferences(renameTable(compiler, js, true)));
    });
    conventions(sb, parsed, parsed.getJsRoot());
    deps(sb, rel, src);
  }

  private static void gatherModuleMetadata(Compiler compiler, Node externs, Node js) {
    new GatherModuleMetadata(compiler, false, ResolutionMode.BROWSER).process(externs, js);
  }

  /**
   * Renames every name declared in every scope of js (in the scope where it is declared) to
   * name$r: (scope root, name) -> name + "$r". With moduleRootsOnly, only the scopes that
   * Es6RenameTypeReferences accepts (rooted at a MODULE_BODY, a SCRIPT or a function body) get
   * rows, as in Es6RewriteModules.
   */
  private static Table<Node, String, String> renameTable(
      Compiler compiler, Node js, boolean moduleRootsOnly) {
    Table<Node, String, String> table = HashBasedTable.create();
    NodeTraversal.traverse(
        compiler,
        js,
        new NodeTraversal.AbstractScopedCallback() {
          @Override
          public void enterScope(NodeTraversal t) {
            Scope s = t.getScope();
            Node r = s.getRootNode();
            if (moduleRootsOnly
                && !(r.isModuleBody() || r.isScript() || NodeUtil.isFunctionBlock(r))) {
              return;
            }
            for (Var v : s.getVarIterable()) {
              table.put(s.getRootNode(), v.getName(), v.getName() + "$r");
            }
          }

          @Override
          public void visit(NodeTraversal t, Node n, Node parent) {}
        });
    return table;
  }

  // HotSpot appends to a failed cast's message where each class was loaded from, e.g.
  // " (com.google.javascript.rhino.Node and com.google.javascript.rhino.Node$StringNode are in
  // unnamed module of loader 'app')". That is JVM environment detail (classpath vs module path,
  // the loader's name), not Closure behaviour, and the Go port has no class loaders, so the dump
  // drops it.
  private static final Pattern CCE_LOADER_SUFFIX =
      Pattern.compile("(cannot be cast to class \\S+) \\([^()]* module [^()]*\\)");

  private static void pass(
      StringBuilder sb, String rel, String src, List<String> parseDiags, String name, Run run) {
    line(sb, "## pass " + name);
    Compiler compiler = newCompiler(rel, src);
    Node externs = compiler.getExternsRoot();
    Node js = compiler.getJsRoot();
    String before = js.toStringTree();
    List<String> extra = new ArrayList<>();
    String exception = null;
    try {
      run.run(compiler, externs, js, extra);
    } catch (RuntimeException e) {
      exception = "exception " + CCE_LOADER_SUFFIX.matcher(e.toString()).replaceAll("$1");
    }
    List<String> remaining = new ArrayList<>(parseDiags);
    for (String d : diagnostics(compiler)) {
      if (!remaining.remove(d)) {
        line(sb, d);
      }
    }
    if (exception != null) {
      line(sb, exception);
    }
    String after = js.toStringTree();
    if (after.equals(before)) {
      line(sb, "tree unchanged");
    } else {
      StringBuilder tree = new StringBuilder();
      for (String l : after.split("\n", -1)) {
        line(tree, l);
      }
      line(sb, "tree changed " + sha256(tree.toString()));
      if (full) {
        sb.append(tree);
      }
    }
    for (String e : extra) {
      line(sb, e);
    }
  }

  private static String nodeRef(Node n) {
    if (n == null) {
      return "null";
    }
    return n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
  }

  private static void conventions(StringBuilder sb, Compiler compiler, Node js) {
    line(sb, "## conventions");
    CodingConvention c = CodingConventions.getDefault();
    for (Node script = js.getFirstChild(); script != null; script = script.getNext()) {
      line(sb, "package " + c.getPackageName(script.getStaticSourceFile()));
    }
    NodeUtil.visitPreOrder(js, n -> conventionsAt(sb, c, n));
  }

  private static void conventionsAt(StringBuilder sb, CodingConvention c, Node n) {
    String at = nodeRef(n);
    if (n.isCall()) {
      CodingConvention.SubclassRelationship rel = c.getClassesDefinedByCall(n);
      if (rel != null) {
        line(sb, at + " subclass " + rel.subclassName + " " + rel.superclassName + " " + rel.type);
      }
      for (boolean checkTypes : new boolean[] {false, true}) {
        CodingConvention.Bind bind = c.describeFunctionBind(n, checkTypes);
        if (bind != null) {
          line(
              sb,
              at
                  + " bind"
                  + (checkTypes ? "T" : "F")
                  + " "
                  + nodeRef(bind.target)
                  + " "
                  + nodeRef(bind.thisValue)
                  + " "
                  + nodeRef(bind.parameters)
                  + " "
                  + bind.getBoundParameterCount());
        }
      }
      if (c.isPropertyTestFunction(n)) {
        line(sb, at + " propertyTest");
      }
      if (c.isPropertyRenameFunction(n.getFirstChild())) {
        line(sb, at + " propertyRename");
      }
      if (ModuleImportResolver.isGoogModuleDependencyCall(n)) {
        line(sb, at + " googModuleDependencyCall");
      }
      if (ConvertChunksToESModules.isDynamicImportCallback(n)) {
        line(sb, at + " dynamicImportCallback");
      }
    }
    if (n.isExprResult() || n.isCall()) {
      if (c.isFunctionCallThatAlwaysThrows(n)) {
        line(sb, at + " alwaysThrows");
      }
      if (CodingConventions.defaultIsFunctionCallThatAlwaysThrows(n, "goog.asserts.fail")) {
        line(sb, at + " alwaysThrows goog.asserts.fail");
      }
    }
    if (n.isName() && !n.getString().isEmpty()) {
      String name = n.getString();
      if (c.isExported(name)) {
        line(sb, at + " exported " + name);
      }
      if (c.isExported(name, false)) {
        line(sb, at + " exportedGlobal " + name);
      }
    }
    if (n.isParamList()) {
      for (Node p = n.getFirstChild(); p != null; p = p.getNext()) {
        if (c.isVarArgsParameter(p)) {
          line(sb, nodeRef(p) + " varArgs");
        }
      }
    }
    if (NodeUtil.isNameDeclaration(n)) {
      for (String r : ClosureModulesPassDumpRequires.requires(n)) {
        line(sb, at + " " + r);
      }
    }
  }

  private static void deps(StringBuilder sb, String rel, String src) {
    line(sb, "## deps");
    PrintStreamErrorManager errorManager =
        new PrintStreamErrorManager(new PrintStream(OutputStream.nullOutputStream()));
    DepsFileRegexParser parser = new DepsFileRegexParser(errorManager);
    parser.setShortcutMode(false);
    List<DependencyInfo> infos = parser.parseFile(rel, src);
    for (DependencyInfo info : infos) {
      line(sb, "dep " + info);
    }
    for (JSError e : errorManager.getErrors()) {
      line(sb, diag("ERROR", e));
    }
    for (JSError e : errorManager.getWarnings()) {
      line(sb, diag("WARNING", e));
    }
    line(sb, "succeeded " + parser.didParseSucceed());
  }

  private static String sha256(String ascii) {
    try {
      byte[] d =
          java.security.MessageDigest.getInstance("SHA-256")
              .digest(ascii.getBytes(StandardCharsets.UTF_8));
      StringBuilder b = new StringBuilder();
      for (int i = 0; i < 8; i++) {
        b.append(String.format("%02x", d[i] & 0xff));
      }
      return b.toString();
    } catch (java.security.NoSuchAlgorithmException e) {
      throw new IllegalStateException(e);
    }
  }

  /** Backslash and every char outside printable ASCII become escapes. */
  static String escape(String s) {
    StringBuilder b = new StringBuilder(s.length());
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        b.append("\\\\");
      } else if (c >= 0x20 && c < 0x7f) {
        b.append(c);
      } else {
        b.append(String.format("\\u%04x", (int) c));
      }
    }
    return b.toString();
  }

  private static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }
}

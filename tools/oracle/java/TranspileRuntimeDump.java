// Oracle tool for the level-11 cluster l11-transpile-runtime: runs each of its passes (and
// drives each of its helper classes) on real files and prints the resulting AST, diagnostics and
// side results. The Go side is internal/jscomp/transpile_runtime_oracle_test.go; it must print
// the same bytes.
//
// It lives in package com.google.javascript.jscomp to reach package-private API:
//   javac -cp closure-compiler.jar -d <build> TranspileRuntimeDump.java
//   java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.TranspileRuntimeDump \
//       [--full] [--sections=a,b,...] <out.gz> <root> <rel-file>... | @<listfile>
//
// For every file and section, a fresh Compiler is initialized the way CompilerTestCase does it
// (Compiler.init with the file as the only source -- or, for the externs sections, as the only
// externs -- options from new CompilerOptions() with language_in ECMASCRIPT_NEXT plus the
// section's settings below, then Compiler.parse()). A file whose parse reports an error other
// than JSC_JS_MODULE_LOAD_WARNING (an import of a file outside the compilation) prints
// "# file <rel>" and "# parse-error" only.
//
// Output (gzipped; every line escaped like ParsingDump: '\' -> "\\", every char outside printable
// ASCII -> \\uXXXX):
//   # file <rel>
//   ## <section>
//   ERROR|WARNING <key> <line>:<charno> <description>   every diagnostic the section reported
//                                                      (parse diagnostics excluded), in order
//   <section-specific lines>
//   tree <same|changed> <lines> <crc32>   the tree the section works on (the sources root, or
//                                         the externs root for the externs sections) after the
//                                         section, as ParsingDump prints a tree (one line per
//                                         node: 4 spaces per depth + Node.toString(true, true,
//                                         false); "  @jsdoc " + JSDocInfo.toStringVerbose() after
//                                         a node with JSDoc; "  @shadow" and the shadow's tree
//                                         one level deeper after a closure-unaware shadow host):
//                                         whether it differs from the tree before the section,
//                                         its line count and the CRC-32 of its escaped lines
//                                         (each followed by '\n'); with --full, the lines
//                                         follow when the tree changed.
//   With --full, "# before" and the sources tree before any section (as above) follow the
//   "# file" line.
//   ! <exception>                          if the section threw: the root cause's toString(), a
//                                          NullPointerException's class name only, a
//                                          ClassCastException's without the JVM's module/loader
//                                          suffix; nothing else of the section follows.
//
// Sections:
//   substitute     SubstituteEs6Syntax.process.
//   shorthand      Es6NormalizeShorthandProperties.transpileSubtree on every node, in post-order
//                  (NodeUtil.visitPostOrder).
//   catch          RewriteCatchWithNoBinding.transpileSubtree on every node, in post-order.
//   trycatch       RewriteGlobalDeclarationsForTryCatchWrapping.process.
//   inject-<out>   InjectTranspilationRuntimeLibraries.process with language_out <out> (es5,
//                  es2015, es2017, es2019) in RuntimeLibraryMode.RECORD_ONLY, and
//                  inject-es2017-async with setInstrumentAsyncContext(true): "injected <lib>"
//                  for each RuntimeJsLibManager.getInjectedLibraries() entry. No tree.
//   inject-tree    the same with language_out ECMASCRIPT5 in RuntimeLibraryMode.INJECT, and the
//                  tree.
//   force          InjectRuntimeLibraries with the forced libraries es6/util/inherits and
//                  es6/symbol in RuntimeLibraryMode.INJECT: injected libraries and the tree.
//   unwrap         ManageClosureUnawareCode.unwrap(compiler).process.
//   caller         RewriteCallerCodeLocation.process.
//   strip          StripCode (the configuration of StripCodeTest, with tweak stripping) after
//                  Compiler.setLifeCycleStage(NORMALIZED) on the parsed (unnormalized) AST.
//   ids            ReplaceIdGenerators (the configuration of ReplaceIdGeneratorsTest, without
//                  pseudo names, no previous map): "map <line>" for each line of
//                  getSerializedIdMappings(), then the tree.
//   toggles        ReplaceToggles.process.
//   window         DeclaredGlobalExternsOnWindow.process on the file as externs.
//   externprops    GatherExternProperties (Mode.CHECK_AND_OPTIMIZE) on the file as externs:
//                  "prop <name>" for each of compiler.getExternProperties(). No tree.
//   accessors      GatherGetterAndSetterProperties.gather(compiler, ROOT): "accessor <name>
//                  <kind>" for each entry. No tree.
//   visibility     CollectFileOverviewVisibility.process on the sources root: "visibility
//                  <source name> <Visibility>" for each entry of getFileOverviewVisibilityMap().
//                  No tree.
//   access         a NodeTraversal over the sources root that, for each NAME with a Var in
//                  t.getScope(), computes AccessControlUtils.getEffectiveNameVisibility(name, var,
//                  the visibility map of the visibility section), and for each GETPROP,
//                  getEffectivePropertyVisibility(getprop, null, that map): "access <names>
//                  <getprops> <crc32>" (counts, and the CRC-32 of the result names, one per
//                  line, in traversal order). No tree.
//   template       a NodeTraversal over the sources root driving Es6TemplateLiterals as
//                  LateEs6ToEs3Converter does (insertion point: the SCRIPT's first child;
//                  visitTaggedTemplateLiteral for TAGGED_TEMPLATELIT, visitTemplateLiteral for
//                  a TEMPLATELIT not under one), after injectLibForField of
//                  $jscomp.createTemplateTagFirstArg and ...WithRaw in RECORD_ONLY mode.
//   logical        a NodeTraversal over the sources root calling
//                  RewriteLogicalAssignmentOperatorsHelper.visitLogicalAssignmentOperator for
//                  every ASSIGN_OR, ASSIGN_AND, ASSIGN_COALESCE (as
//                  RewriteLogicalAssignmentOperatorsPass does).
//   ctors          a NodeTraversal over the sources root calling
//                  SynthesizeExplicitConstructors.synthesizeClassConstructorIfMissing for every
//                  CLASS whose superclass is EMPTY or a qualified name.
//   arrows         Es6RewriteArrowFunction's use of ThisAndArgumentsReferenceUpdater and
//                  ThisAndArgumentsContext (contexts for SCRIPTs and non-arrow FUNCTIONs,
//                  lastSuperStatement for super() calls in constructors, the updater traversal of
//                  every arrow function's body, addVarDeclarations when leaving a context), for
//                  arrow functions with a BLOCK body only (the others are left alone).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableMap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSet;
import com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage;
import com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind;
import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.ThisAndArgumentsReferenceUpdater.ThisAndArgumentsContext;
import com.google.javascript.jscomp.js.RuntimeJsLibManager;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.JSDocInfo.Visibility;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.StaticSourceFile;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.io.PrintStream;
import java.io.Writer;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Deque;
import java.util.List;
import java.util.Map;
import java.util.function.Consumer;
import java.util.zip.CRC32;
import java.util.zip.GZIPOutputStream;

public final class TranspileRuntimeDump {
  static final String[] ALL_SECTIONS = {
    "substitute", "shorthand", "catch", "trycatch", "inject-es5", "inject-es2015",
    "inject-es2017", "inject-es2017-async", "inject-es2019", "inject-tree", "force", "unwrap",
    "caller", "strip", "ids", "toggles", "window", "externprops", "accessors", "visibility",
    "access", "template", "logical", "ctors", "arrows",
  };

  static Writer out;
  static boolean full = false;

  public static void main(String[] args) throws Exception {
    List<String> rest = new ArrayList<>();
    List<String> sections = List.of(ALL_SECTIONS);
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.startsWith("--sections=")) {
        sections = List.of(a.substring("--sections=".length()).split(","));
      } else {
        rest.add(a);
      }
    }
    Path outPath = Paths.get(rest.get(0));
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
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(outPath))) {
      out = new OutputStreamWriter(os, StandardCharsets.UTF_8);
      for (String rel : files) {
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        dumpFile(rel, src, sections);
        out.flush();
      }
      out.flush();
    }
  }

  // ---------------------------------------------------------------------------------------------
  // Shared helpers

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

  static void line(String s) throws java.io.IOException {
    out.write(escape(s));
    out.write('\n');
  }

  static String exc(Throwable e) {
    while (e.getCause() != null) {
      e = e.getCause();
    }
    if (e instanceof NullPointerException) {
      return e.getClass().getName();
    }
    if (e instanceof ClassCastException && e.getMessage() != null) {
      // Drop the JVM's " (<module/loader of each class>)" suffix.
      String m = e.getMessage();
      int paren = m.indexOf(" (");
      return e.getClass().getName() + ": " + (paren < 0 ? m : m.substring(0, paren));
    }
    return e.toString();
  }

  /** The tree as ParsingDump prints it: escaped lines. */
  static void treeLines(List<String> lines, Node n, int level) {
    String indent = "    ".repeat(level);
    lines.add(escape(indent + n.toString(true, true, false)));
    JSDocInfo info = n.getJSDocInfo();
    if (info != null) {
      lines.add(escape(indent + "  @jsdoc " + info.toStringVerbose()));
    }
    Node shadow = n.getClosureUnawareShadow();
    if (shadow != null) {
      lines.add(escape(indent + "  @shadow"));
      treeLines(lines, shadow, level + 1);
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      treeLines(lines, c, level + 1);
    }
  }

  static List<String> tree(Node n) {
    List<String> lines = new ArrayList<>();
    treeLines(lines, n, 0);
    return lines;
  }

  static long crc(List<String> lines) {
    CRC32 crc = new CRC32();
    for (String l : lines) {
      crc.update(l.getBytes(StandardCharsets.UTF_8));
      crc.update('\n');
    }
    return crc.getValue();
  }

  static void printTree(List<String> before, Node n) throws java.io.IOException {
    List<String> after = tree(n);
    out.write(
        "tree "
            + (after.equals(before) ? "same" : "changed")
            + " "
            + after.size()
            + " "
            + Long.toHexString(crc(after))
            + "\n");
    if (full && !after.equals(before)) {
      for (String l : after) {
        out.write(l);
        out.write('\n');
      }
    }
  }

  static Compiler newCompiler(
      String rel, String src, boolean asExterns, Consumer<CompilerOptions> settings) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    settings.accept(options);
    SourceFile file = SourceFile.fromCode(rel, src);
    if (asExterns) {
      compiler.init(ImmutableList.of(file), ImmutableList.of(), options);
    } else {
      compiler.init(ImmutableList.of(), ImmutableList.of(file), options);
    }
    compiler.parse();
    return compiler;
  }

  /**
   * Whether parsing reported an error other than JSC_JS_MODULE_LOAD_WARNING (an import of a file
   * that is not part of the compilation, as most corpus modules have).
   */
  static boolean hasParseError(Compiler compiler) {
    for (JSError e : compiler.getErrors()) {
      if (!e.getType().key.equals("JSC_JS_MODULE_LOAD_WARNING")) {
        return true;
      }
    }
    return false;
  }

  static void printDiagnostics(Compiler compiler, int skipErrors, int skipWarnings)
      throws java.io.IOException {
    ImmutableList<JSError> errors = compiler.getErrors();
    for (int i = skipErrors; i < errors.size(); i++) {
      JSError e = errors.get(i);
      line(
          "ERROR "
              + e.getType().key
              + " "
              + e.getLineNumber()
              + ":"
              + e.getCharno()
              + " "
              + e.getDescription());
    }
    ImmutableList<JSError> warnings = compiler.getWarnings();
    for (int i = skipWarnings; i < warnings.size(); i++) {
      JSError e = warnings.get(i);
      line(
          "WARNING "
              + e.getType().key
              + " "
              + e.getLineNumber()
              + ":"
              + e.getCharno()
              + " "
              + e.getDescription());
    }
  }

  interface Section {
    /** Runs the section; returns the tree to print, or null for none. */
    Node run(Compiler compiler, List<String> results) throws Exception;
  }

  static void dumpFile(String rel, String src, List<String> sections) throws Exception {
    line("# file " + rel);
    Compiler probe = newCompiler(rel, src, false, o -> {});
    if (hasParseError(probe)) {
      line("# parse-error");
      return;
    }
    // Every section starts from the same parse, so the "before" trees are computed once.
    List<String> jsBefore = tree(probe.getJsRoot());
    List<String> externsBefore = null;
    if (full) {
      line("# before");
      for (String l : jsBefore) {
        out.write(l);
        out.write('\n');
      }
    }
    for (String section : sections) {
      line("## " + section);
      boolean asExterns =
          section.equals("window") || section.equals("externprops");
      Compiler compiler = newCompiler(rel, src, asExterns, settings(section));
      if (hasParseError(compiler)) {
        line("# parse-error");
        continue;
      }
      int errors0 = compiler.getErrors().size();
      int warnings0 = compiler.getWarnings().size();
      if (asExterns && externsBefore == null) {
        externsBefore = tree(compiler.getExternsRoot());
      }
      List<String> before = asExterns ? externsBefore : jsBefore;
      List<String> results = new ArrayList<>();
      Node printed;
      try {
        printed = section(section).run(compiler, results);
      } catch (RuntimeException | AssertionError e) {
        printDiagnostics(compiler, errors0, warnings0);
        line("! " + exc(e));
        continue;
      }
      printDiagnostics(compiler, errors0, warnings0);
      for (String r : results) {
        line(r);
      }
      if (printed != null) {
        printTree(before, printed);
      }
    }
  }

  static Consumer<CompilerOptions> settings(String section) {
    switch (section) {
      case "inject-es5":
        return o -> injectSettings(o, LanguageMode.ECMASCRIPT5, false);
      case "inject-es2015":
        return o -> injectSettings(o, LanguageMode.ECMASCRIPT_2015, false);
      case "inject-es2017":
        return o -> injectSettings(o, LanguageMode.ECMASCRIPT_2017, false);
      case "inject-es2017-async":
        return o -> injectSettings(o, LanguageMode.ECMASCRIPT_2017, true);
      case "inject-es2019":
        return o -> injectSettings(o, LanguageMode.ECMASCRIPT_2019, false);
      case "inject-tree":
        return o -> {
          o.setLanguageOut(LanguageMode.ECMASCRIPT5);
          o.setRuntimeLibraryMode(RuntimeJsLibManager.RuntimeLibraryMode.INJECT);
        };
      case "force":
        return o -> o.setRuntimeLibraryMode(RuntimeJsLibManager.RuntimeLibraryMode.INJECT);
      case "template":
        return o -> o.setRuntimeLibraryMode(RuntimeJsLibManager.RuntimeLibraryMode.RECORD_ONLY);
      default:
        return o -> {};
    }
  }

  static void injectSettings(CompilerOptions o, LanguageMode out, boolean async) {
    o.setLanguageOut(out);
    o.setInstrumentAsyncContext(async);
    o.setRuntimeLibraryMode(RuntimeJsLibManager.RuntimeLibraryMode.RECORD_ONLY);
  }

  static void addInjected(Compiler compiler, List<String> results) {
    for (String lib : compiler.getRuntimeJsLibManager().getInjectedLibraries()) {
      results.add("injected " + lib);
    }
  }

  static Map<StaticSourceFile, Visibility> visibilityMap(Compiler compiler) {
    CollectFileOverviewVisibility pass = new CollectFileOverviewVisibility(compiler);
    pass.process(compiler.getExternsRoot(), compiler.getJsRoot());
    return pass.getFileOverviewVisibilityMap();
  }

  static Section section(String name) {
    switch (name) {
      case "substitute":
        return (c, r) -> {
          new SubstituteEs6Syntax(c).process(c.getExternsRoot(), c.getJsRoot());
          return c.getJsRoot();
        };
      case "shorthand":
        return (c, r) -> {
          Es6NormalizeShorthandProperties p = new Es6NormalizeShorthandProperties(c);
          NodeUtil.visitPostOrder(c.getJsRoot(), n -> p.transpileSubtree(n));
          return c.getJsRoot();
        };
      case "catch":
        return (c, r) -> {
          RewriteCatchWithNoBinding p = new RewriteCatchWithNoBinding(c);
          NodeUtil.visitPostOrder(c.getJsRoot(), n -> p.transpileSubtree(n));
          return c.getJsRoot();
        };
      case "trycatch":
        return (c, r) -> {
          new RewriteGlobalDeclarationsForTryCatchWrapping(c)
              .process(c.getExternsRoot(), c.getJsRoot());
          return c.getJsRoot();
        };
      case "inject-es5":
      case "inject-es2015":
      case "inject-es2017":
      case "inject-es2017-async":
      case "inject-es2019":
        return (c, r) -> {
          new InjectTranspilationRuntimeLibraries(c).process(c.getExternsRoot(), c.getJsRoot());
          addInjected(c, r);
          return null;
        };
      case "inject-tree":
        return (c, r) -> {
          new InjectTranspilationRuntimeLibraries(c).process(c.getExternsRoot(), c.getJsRoot());
          addInjected(c, r);
          return c.getJsRoot();
        };
      case "force":
        return (c, r) -> {
          new InjectRuntimeLibraries(c, ImmutableSet.of("es6/util/inherits", "es6/symbol"))
              .process(c.getExternsRoot(), c.getJsRoot());
          addInjected(c, r);
          return c.getJsRoot();
        };
      case "unwrap":
        return (c, r) -> {
          ManageClosureUnawareCode.unwrap(c).process(c.getExternsRoot(), c.getJsRoot());
          return c.getJsRoot();
        };
      case "caller":
        return (c, r) -> {
          new RewriteCallerCodeLocation(c).process(c.getExternsRoot(), c.getJsRoot());
          return c.getJsRoot();
        };
      case "strip":
        return (c, r) -> {
          c.setLifeCycleStage(LifeCycleStage.NORMALIZED);
          new StripCode(
                  c,
                  ImmutableSet.of(
                      "goog.debug.DebugWindow",
                      "goog.debug.FancyWindow",
                      "goog.debug.Formatter",
                      "goog.debug.HtmlFormatter",
                      "goog.debug.TextFormatter",
                      "goog.debug.Logger",
                      "goog.debug.LogManager",
                      "goog.debug.LogRecord",
                      "goog.net.BrowserChannel.LogSaver",
                      "goog.log",
                      "GA_GoogleDebugger"),
                  ImmutableSet.of(
                      "logger",
                      "logger_",
                      "debugWindow",
                      "debugWindow_",
                      "logFormatter_",
                      "logBuffer_"),
                  ImmutableSet.of("trace"),
                  true)
              .process(c.getExternsRoot(), c.getJsRoot());
          return c.getJsRoot();
        };
      case "ids":
        return (c, r) -> {
          RenamingMap idTestMap =
              new RenamingMap() {
                private final ImmutableMap<String, String> map =
                    ImmutableMap.of("foo", ":foo:", "bar", ":bar:");

                @Override
                public String get(String value) {
                  String replacement = map.get(value);
                  return replacement != null ? replacement : "unknown:" + value;
                }
              };
          ReplaceIdGenerators pass =
              new ReplaceIdGenerators(
                  c,
                  /* templateLiteralsAreTranspiled= */ false,
                  new ImmutableMap.Builder<String, RenamingMap>()
                      .put("goog.events.getUniqueId", RenamingToken.INCONSISTENT)
                      .put("goog.place.getUniqueId", RenamingToken.INCONSISTENT)
                      .put("id", idTestMap)
                      .put("get.id", idTestMap)
                      .buildOrThrow(),
                  /* generatePseudoNames= */ false,
                  /* previousMapSerialized= */ null,
                  /* xidHashFunction= */ null);
          pass.process(c.getExternsRoot(), c.getJsRoot());
          for (String l : pass.getSerializedIdMappings().split("\n", -1)) {
            r.add("map " + l);
          }
          return c.getJsRoot();
        };
      case "toggles":
        return (c, r) -> {
          new ReplaceToggles(c).process(c.getExternsRoot(), c.getJsRoot());
          return c.getJsRoot();
        };
      case "window":
        return (c, r) -> {
          new DeclaredGlobalExternsOnWindow(c).process(c.getExternsRoot(), c.getJsRoot());
          return c.getExternsRoot();
        };
      case "externprops":
        return (c, r) -> {
          new GatherExternProperties(c, GatherExternProperties.Mode.CHECK_AND_OPTIMIZE)
              .process(c.getExternsRoot(), c.getJsRoot());
          for (String p : c.getExternProperties()) {
            r.add("prop " + p);
          }
          return null;
        };
      case "accessors":
        return (c, r) -> {
          for (Map.Entry<String, PropertyAccessKind> e :
              GatherGetterAndSetterProperties.gather(c, c.getRoot()).entrySet()) {
            r.add("accessor " + e.getKey() + " " + e.getValue());
          }
          return null;
        };
      case "visibility":
        return (c, r) -> {
          for (Map.Entry<StaticSourceFile, Visibility> e : visibilityMap(c).entrySet()) {
            r.add("visibility " + e.getKey().getName() + " " + e.getValue());
          }
          return null;
        };
      case "access":
        return (c, r) -> {
          ImmutableMap<StaticSourceFile, Visibility> map =
              ImmutableMap.copyOf(visibilityMap(c));
          int[] counts = {0, 0};
          CRC32 crc = new CRC32();
          NodeTraversal.traverse(
              c,
              c.getJsRoot(),
              new NodeTraversal.AbstractPostOrderCallback() {
                @Override
                public void visit(NodeTraversal t, Node n, Node parent) {
                  Visibility v = null;
                  if (n.isName() && !n.getString().isEmpty()) {
                    Var var = t.getScope().getVar(n.getString());
                    if (var == null) {
                      return;
                    }
                    v = AccessControlUtils.getEffectiveNameVisibility(n, var, map);
                    counts[0]++;
                  } else if (n.isGetProp()) {
                    v = AccessControlUtils.getEffectivePropertyVisibility(n, null, map);
                    counts[1]++;
                  } else {
                    return;
                  }
                  crc.update((v.name() + "\n").getBytes(StandardCharsets.UTF_8));
                }
              });
          r.add(
              "access " + counts[0] + " " + counts[1] + " " + Long.toHexString(crc.getValue()));
          return null;
        };
      case "template":
        return (c, r) -> {
          RuntimeJsLibManager libs = c.getRuntimeJsLibManager();
          libs.injectLibForField("$jscomp.createTemplateTagFirstArg");
          libs.injectLibForField("$jscomp.createTemplateTagFirstArgWithRaw");
          Es6TemplateLiterals converter = new Es6TemplateLiterals(c);
          NodeTraversal.traverse(
              c,
              c.getJsRoot(),
              new NodeTraversal.Callback() {
                Node insertionPoint;

                @Override
                public boolean shouldTraverse(NodeTraversal t, Node n, Node parent) {
                  if (n.isScript()) {
                    insertionPoint = n.getFirstChild();
                  }
                  return true;
                }

                @Override
                public void visit(NodeTraversal t, Node n, Node parent) {
                  if (n.isTaggedTemplateLit()) {
                    converter.visitTaggedTemplateLiteral(t, n, insertionPoint);
                  } else if (n.isTemplateLit() && !parent.isTaggedTemplateLit()) {
                    converter.visitTemplateLiteral(t, n);
                  }
                }
              });
          return c.getJsRoot();
        };
      case "logical":
        return (c, r) -> {
          RewriteLogicalAssignmentOperatorsHelper helper =
              new RewriteLogicalAssignmentOperatorsHelper(
                  c, c.createAstFactory(), c.getUniqueIdSupplier());
          NodeTraversal.traverse(
              c,
              c.getJsRoot(),
              new NodeTraversal.AbstractPostOrderCallback() {
                @Override
                public void visit(NodeTraversal t, Node n, Node parent) {
                  switch (n.getToken()) {
                    case ASSIGN_OR, ASSIGN_AND, ASSIGN_COALESCE ->
                        helper.visitLogicalAssignmentOperator(t, n);
                    default -> {}
                  }
                }
              });
          return c.getJsRoot();
        };
      case "ctors":
        return (c, r) -> {
          SynthesizeExplicitConstructors synthesizer = new SynthesizeExplicitConstructors(c);
          NodeTraversal.traverse(
              c,
              c.getJsRoot(),
              new NodeTraversal.AbstractPostOrderCallback() {
                @Override
                public void visit(NodeTraversal t, Node n, Node parent) {
                  if (n.isClass()
                      && (n.getSecondChild().isEmpty() || n.getSecondChild().isQualifiedName())) {
                    synthesizer.synthesizeClassConstructorIfMissing(t, n);
                  }
                }
              });
          return c.getJsRoot();
        };
      case "arrows":
        return (c, r) -> {
          AstFactory astFactory = c.createAstFactory();
          Deque<ThisAndArgumentsContext> contextStack = new ArrayDeque<>();
          NodeTraversal.traverse(
              c,
              c.getJsRoot(),
              new NodeTraversal.Callback() {
                @Override
                public boolean shouldTraverse(NodeTraversal t, Node n, Node parent) {
                  switch (n.getToken()) {
                    case SCRIPT ->
                        contextStack.push(
                            new ThisAndArgumentsContext(
                                n, false, c.getUniqueIdSupplier().getUniqueId(t.getInput())));
                    case FUNCTION -> {
                      if (!n.isArrowFunction()) {
                        contextStack.push(
                            new ThisAndArgumentsContext(
                                n.getLastChild(),
                                NodeUtil.isEs6Constructor(n),
                                c.getUniqueIdSupplier().getUniqueId(t.getInput())));
                      }
                    }
                    case SUPER -> {
                      ThisAndArgumentsContext context = contextStack.peek();
                      if (context != null
                          && context.isConstructor
                          && parent.isCall()
                          && parent.getFirstChild() == n) {
                        Node s = parent;
                        while (s.getParent() != context.scopeBody) {
                          s = NodeUtil.getEnclosingStatement(s.getParent());
                        }
                        context.lastSuperStatement = s;
                      }
                    }
                    default -> {}
                  }
                  return true;
                }

                @Override
                public void visit(NodeTraversal t, Node n, Node parent) {
                  ThisAndArgumentsContext context = contextStack.peek();
                  if (n.isArrowFunction()) {
                    if (n.getLastChild().isBlock()) {
                      n.setIsArrowFunction(false);
                      NodeTraversal.traverse(
                          c,
                          n.getLastChild(),
                          new ThisAndArgumentsReferenceUpdater(c, context, astFactory));
                      t.reportCodeChange();
                    }
                  } else if (context != null && context.scopeBody == n) {
                    contextStack.pop();
                    context.addVarDeclarations(c, astFactory, t.getCurrentScript());
                  }
                }
              });
          return c.getJsRoot();
        };
      default:
        throw new IllegalArgumentException("unknown section " + name);
    }
  }
}

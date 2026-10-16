// Oracle tool for core group g2 (AstAnalyzer, ExpressionDecomposer, OptionalChainRewriter,
// MakeDeclaredNamesUnique, AstFactory): applies their API to real ASTs and prints the results.
// The Go side is internal/jscomp/core_g2_oracle_test.go; it must print the same bytes.
//
// It lives in package com.google.javascript.jscomp to reach package-private API:
//   javac -cp closure-compiler.jar -d <build> CoreG2Dump.java
//   java -Xss1g -cp closure-compiler.jar:<build> com.google.javascript.jscomp.CoreG2Dump \
//       <mode> <out.gz> [<root> <rel-file>... | @<listfile>]
//
// Every file is parsed with ParserRunner.parse and the parsing.Config the compiler uses for
// --language_in=ECMASCRIPT_NEXT (as ParsingDump does), from SourceFile.fromCode(rel, src).
// The output starts with that config ("# config ..." as in ParsingDump). Files whose parse
// reports an error, or throws, print "# file <rel>" and "# parse-error" only.
// Output is gzipped; every line is escaped like ParsingDump ('\' -> "\\", every char outside
// printable ASCII -> \\uXXXX). An exception prints as its root cause's toString(), except a
// NullPointerException, which prints its class name only (the JVM's helpful messages are not
// reproducible).
//
// Modes:
//   analysis  per node, in preorder: "<TOKEN> a0 a1 a2 a3 d0 d1"; each aK is five chars for the
//             AstAnalyzer configuration K (see ANALYZERS): mayHaveSideEffects,
//             mayEffectMutableState, nodeTypeMayHaveSideEffects, functionCallHasSideEffects
//             (CALL, OPTCHAIN_CALL, TAGGED_TEMPLATELIT only) and constructorCallHasSideEffects
//             (NEW only): '1' true, '0' false, '-' not applicable, '!' exception. dK is
//             ExpressionDecomposer K's canExposeExpression for nodes IR.mayBeExpression accepts:
//             'U', 'M', 'D', '-' or '!'. Each exception is followed by a line
//             "! <config> <query> <exception>".
//   mutate    on clones of the AST: ExpressionDecomposer.maybeExposeExpression on the first
//             8 DECOMPOSABLE nodes (decomposer d0), moveExpression on the first 4 MOVABLE CALLs,
//             OptionalChainRewriter.rewrite on the first 8 ends of full optional chains; after
//             each, the statements of the affected block from the injection point on (as many
//             as were added, plus one) as Node.toStringTree (at most 60 lines, then the line
//             count and the CRC-32 of all of them).
//   mdnu      on clones of the AST, MakeDeclaredNamesUnique with ContextualRenamer, then the
//             ContextualRenameInverter, then the assertOnChange check; InlineRenamer on detached
//             clones of the first 3 FUNCTIONs as FunctionToBlockMutator runs it (indices relative
//             to the function); BoilerplateRenamer;
//             TargettedRenamer. After each, every NAME or
//             IMPORT_STAR whose string differs from the original: "<preorder index> <old> <new>".
//   factory   (no files) every AstFactory create* method under each type mode and life-cycle
//             stage, on fixed inputs: Node.toStringTree of the result, or the exception.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableMap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSet;
import com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage;
import com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind;
import com.google.javascript.jscomp.colors.ColorRegistry;
import com.google.javascript.jscomp.colors.StandardColors;
import com.google.javascript.jscomp.js.RuntimeJsLibManager;
import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.StaticRef;
import com.google.javascript.rhino.StaticScope;
import com.google.javascript.rhino.StaticSlot;
import com.google.javascript.rhino.StaticSourceFile;
import com.google.javascript.rhino.Token;
import com.google.javascript.rhino.jstype.JSType;
import com.google.javascript.rhino.jstype.JSTypeNative;
import com.google.javascript.rhino.jstype.JSTypeRegistry;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.io.PrintStream;
import java.io.Writer;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.EnumSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Supplier;
import java.util.zip.GZIPOutputStream;

public final class CoreG2Dump {

  private static Writer out;

  public static void main(String[] args) throws Exception {
    String mode = args[0];
    Path outPath = Paths.get(args[1]);
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(outPath), 1 << 16)) {
      out = new OutputStreamWriter(os, StandardCharsets.UTF_8);
      if (mode.equals("factory")) {
        FactorySweep.run();
      } else {
        Path root = Paths.get(args[2]);
        List<String> files = new ArrayList<>();
        for (int i = 3; i < args.length; i++) {
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
        Config config = compilerParserConfig();
        line(
            "# config languageMode="
                + config.languageMode()
                + " jsDocParsingMode="
                + config.jsDocParsingMode()
                + " runMode="
                + config.runMode()
                + " strictMode="
                + config.strictMode());
        for (String rel : files) {
          String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
          line("# file " + rel);
          Node script = parse(config, rel, src);
          if (script == null) {
            line("# parse-error");
            continue;
          }
          switch (mode) {
            case "analysis" -> analysis(script);
            case "mutate" -> mutate(script);
            case "mdnu" -> mdnu(script);
            default -> throw new IllegalArgumentException(mode);
          }
        }
      }
      out.flush();
    }
  }

  // ---------------------------------------------------------------------------------------------
  // Shared helpers

  /** Backslash and every char outside printable ASCII become escapes (ParsingDump.escape). */
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
    // The root cause: NodeTraversal wraps what a callback throws in a RuntimeException whose
    // message embeds the cause's.
    while (e.getCause() != null) {
      e = e.getCause();
    }
    if (e instanceof NullPointerException) {
      return e.getClass().getName();
    }
    return e.toString();
  }

  /** The parsing.Config that CompilerInput passes to ParserRunner.parse (see ParsingDump). */
  private static Config compilerParserConfig() throws Exception {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    compiler.initOptions(options);
    Method getParserConfig =
        Compiler.class.getDeclaredMethod("getParserConfig", AbstractCompiler.ConfigContext.class);
    getParserConfig.setAccessible(true);
    return (Config) getParserConfig.invoke(compiler, AbstractCompiler.ConfigContext.DEFAULT);
  }

  private static Node parse(Config config, String name, String src) {
    boolean[] hadError = {false};
    ErrorReporter reporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int line, int lineOffset) {}

          @Override
          public void error(String message, String sourceName, int line, int lineOffset) {
            hadError[0] = true;
          }
        };
    try {
      ParserRunner.ParseResult result =
          ParserRunner.parse(SourceFile.fromCode(name, src), src, config, reporter);
      if (hadError[0] || result.ast == null) {
        return null;
      }
      return result.ast;
    } catch (RuntimeException e) {
      return null;
    }
  }

  /** The nodes of the tree in preorder (iterative: corpus trees are deep). */
  static List<Node> preorder(Node root) {
    List<Node> result = new ArrayList<>();
    ArrayList<Node> stack = new ArrayList<>();
    stack.add(root);
    while (!stack.isEmpty()) {
      Node n = stack.remove(stack.size() - 1);
      result.add(n);
      for (Node c = n.getLastChild(); c != null; c = c.getPrevious()) {
        stack.add(c);
      }
    }
    return result;
  }

  static Compiler newCompiler() {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    compiler.initOptions(new CompilerOptions());
    return compiler;
  }

  /** Node.toStringTree, one line per node (the tree string's lines, escaped). */
  static void printTree(Node n) throws java.io.IOException {
    for (String l : n.toStringTree().split("\n", -1)) {
      if (!l.isEmpty()) {
        line(l);
      }
    }
  }

  /** A supplier of "0", "1", ... as Compiler.getUniqueNameIdSupplier gives for a new compiler. */
  static Supplier<String> counter() {
    int[] next = {0};
    return () -> String.valueOf(next[0]++);
  }

  // ---------------------------------------------------------------------------------------------
  // analysis

  /** The AstAnalyzer configurations: a0 is what Compiler.getAstAnalyzer() builds by default. */
  static AstAnalyzer[] analyzers(Compiler compiler) {
    AccessorSummary summary =
        AccessorSummary.create(
            ImmutableMap.<String, PropertyAccessKind>builder()
                .put("getter", PropertyAccessKind.GETTER_ONLY)
                .put("setter", PropertyAccessKind.SETTER_ONLY)
                .put("length", PropertyAccessKind.GETTER_ONLY)
                .put("value", PropertyAccessKind.SETTER_ONLY)
                .put("x", PropertyAccessKind.GETTER_AND_SETTER)
                .put("prototype", PropertyAccessKind.GETTER_AND_SETTER)
                .put("a", PropertyAccessKind.GETTER_ONLY)
                .buildOrThrow());
    return new AstAnalyzer[] {
      compiler.getAstAnalyzer(),
      new AstAnalyzer(
          AstAnalyzer.Options.builder()
              .setUseTypesForLocalOptimization(false)
              .setAssumeGettersArePure(false)
              .setHasRegexpGlobalReferences(false)
              .setAssumeKnownBuiltinsArePure(true)
              .build(),
          compiler.getTypeRegistry(),
          summary),
      new AstAnalyzer(
          AstAnalyzer.Options.builder()
              .setUseTypesForLocalOptimization(false)
              .setAssumeGettersArePure(true)
              .setHasRegexpGlobalReferences(true)
              .setAssumeKnownBuiltinsArePure(false)
              .build(),
          compiler.getTypeRegistry(),
          null),
      new AstAnalyzer(
          AstAnalyzer.Options.builder()
              .setUseTypesForLocalOptimization(true)
              .setAssumeGettersArePure(true)
              .setHasRegexpGlobalReferences(false)
              .setAssumeKnownBuiltinsArePure(true)
              .build(),
          compiler.getTypeRegistry(),
          null),
    };
  }

  /** The ExpressionDecomposer configurations: d0 is createDefaultExpressionDecomposer(). */
  static ExpressionDecomposer[] decomposers(Compiler compiler) {
    return new ExpressionDecomposer[] {
      compiler.createDefaultExpressionDecomposer(),
      new ExpressionDecomposer(
          compiler,
          counter(),
          ImmutableSet.of("foo", "goog", "fn", "require", "Object"),
          Scope.createGlobalScope(new Node(Token.SCRIPT)),
          EnumSet.of(ExpressionDecomposer.Workaround.BROKEN_IE11_LOCATION_ASSIGN)),
    };
  }

  interface BoolQuery {
    boolean apply(AstAnalyzer a, Node n);
  }

  private static final String[] QUERY_NAMES = {"MHSE", "MEMS", "NTMHSE", "FCHSE", "CCHSE"};

  static void analysis(Node script) throws java.io.IOException {
    Compiler compiler = newCompiler();
    AstAnalyzer[] analyzers = analyzers(compiler);
    ExpressionDecomposer[] decomposers = decomposers(compiler);
    BoolQuery[] queries = {
      AstAnalyzer::mayHaveSideEffects,
      AstAnalyzer::mayEffectMutableState,
      AstAnalyzer::nodeTypeMayHaveSideEffects,
      AstAnalyzer::functionCallHasSideEffects,
      AstAnalyzer::constructorCallHasSideEffects,
    };
    for (Node n : preorder(script)) {
      StringBuilder sb = new StringBuilder(n.getToken().toString());
      List<String> exceptions = new ArrayList<>();
      boolean isCall = n.isCall() || n.isOptChainCall() || n.isTaggedTemplateLit();
      for (int a = 0; a < analyzers.length; a++) {
        sb.append(' ');
        for (int q = 0; q < queries.length; q++) {
          if ((q == 3 && !isCall) || (q == 4 && !n.isNew())) {
            sb.append('-');
            continue;
          }
          try {
            sb.append(queries[q].apply(analyzers[a], n) ? '1' : '0');
          } catch (RuntimeException e) {
            sb.append('!');
            exceptions.add("! a" + a + " " + QUERY_NAMES[q] + " " + exc(e));
          }
        }
      }
      for (int d = 0; d < decomposers.length; d++) {
        sb.append(' ');
        if (!IR.mayBeExpression(n)) {
          sb.append('-');
          continue;
        }
        try {
          sb.append(
              switch (decomposers[d].canExposeExpression(n)) {
                case UNDECOMPOSABLE -> 'U';
                case MOVABLE -> 'M';
                case DECOMPOSABLE -> 'D';
              });
        } catch (RuntimeException e) {
          sb.append('!');
          exceptions.add("! d" + d + " canExposeExpression " + exc(e));
        }
      }
      line(sb.toString());
      for (String e : exceptions) {
        line(e);
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // mutate

  /**
   * Prints block's children [from, from + count) as trees: the first REGION_LINES lines, then,
   * if there are more, "... <lines> lines crc32 <hex>", the CRC-32 of all the (escaped) lines.
   */
  static void printRegion(Node block, int from, int count) throws java.io.IOException {
    List<String> lines = new ArrayList<>();
    Node c = block.getChildAtIndex(from);
    for (int i = 0; i < count && c != null; i++, c = c.getNext()) {
      for (String l : c.toStringTree().split("\n", -1)) {
        if (!l.isEmpty()) {
          lines.add(escape(l));
        }
      }
    }
    java.util.zip.CRC32 crc = new java.util.zip.CRC32();
    for (int i = 0; i < lines.size(); i++) {
      if (i < REGION_LINES) {
        out.write(lines.get(i));
        out.write('\n');
      }
      crc.update((lines.get(i) + "\n").getBytes(StandardCharsets.UTF_8));
    }
    if (lines.size() > REGION_LINES) {
      line("... " + lines.size() + " lines crc32 " + Long.toHexString(crc.getValue()));
    }
  }

  private static final int REGION_LINES = 60;

  static void mutate(Node script) throws java.io.IOException {
    List<Node> nodes = preorder(script);
    // ExpressionDecomposer.maybeExposeExpression / moveExpression
    {
      Compiler compiler = newCompiler();
      ExpressionDecomposer probe = compiler.createDefaultExpressionDecomposer();
      int exposed = 0;
      int moved = 0;
      for (int i = 0; i < nodes.size() && (exposed < 8 || moved < 4); i++) {
        Node n = nodes.get(i);
        if (!IR.mayBeExpression(n)) {
          continue;
        }
        ExpressionDecomposer.DecompositionType type;
        try {
          type = probe.canExposeExpression(n);
        } catch (RuntimeException e) {
          continue;
        }
        boolean expose = type == ExpressionDecomposer.DecompositionType.DECOMPOSABLE && exposed < 8;
        boolean move =
            type == ExpressionDecomposer.DecompositionType.MOVABLE && n.isCall() && moved < 4;
        if (!expose && !move) {
          continue;
        }
        if (expose) {
          exposed++;
        } else {
          moved++;
        }
        line((expose ? "## expose " : "## move ") + i + " " + n.getToken());
        Node clone = script.cloneTree();
        Node target = preorder(clone).get(i);
        ExpressionDecomposer decomposer =
            compiler.createExpressionDecomposer(
                counter(), ImmutableSet.of(), Scope.createGlobalScope(new Node(Token.SCRIPT)));
        try {
          Node injectionPoint = ExpressionDecomposer.findInjectionPoint(target);
          Node block = injectionPoint.getParent();
          int from = block.getIndexOfChild(injectionPoint);
          int before = block.getChildCount();
          if (expose) {
            decomposer.maybeExposeExpression(target);
          } else {
            decomposer.moveExpression(target);
          }
          printRegion(block, from, block.getChildCount() - before + 1);
        } catch (RuntimeException e) {
          line("! " + exc(e));
        }
      }
    }
    // OptionalChainRewriter
    {
      Compiler compiler = newCompiler();
      int rewritten = 0;
      for (int i = 0; i < nodes.size() && rewritten < 8; i++) {
        Node n = nodes.get(i);
        if (!NodeUtil.isEndOfFullOptChain(n)) {
          continue;
        }
        rewritten++;
        line("## optchain " + i + " " + n.getToken());
        Node clone = script.cloneTree();
        Node target = preorder(clone).get(i);
        int[] next = {0};
        OptionalChainRewriter.Builder builder =
            OptionalChainRewriter.builder(compiler)
                .setTmpVarNameCreator(() -> "tmp$" + next[0]++);
        if (rewritten % 2 == 0) {
          builder.setScope(Scope.createGlobalScope(new Node(Token.SCRIPT)));
        }
        try {
          OptionalChainRewriter rewriter = builder.build(target);
          Node block = rewriter.insertionPoint.getParent();
          int from = block.getIndexOfChild(rewriter.insertionPoint);
          int before = block.getChildCount();
          rewriter.rewrite();
          printRegion(block, from, block.getChildCount() - before + 1);
        } catch (RuntimeException e) {
          line("! " + exc(e));
        }
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // mdnu

  /** Prints every NAME / IMPORT_STAR of tree whose string differs from orig's. */
  static void printRenames(List<Node> orig, Node tree) throws java.io.IOException {
    List<Node> now = preorder(tree);
    for (int i = 0; i < orig.size(); i++) {
      Node o = orig.get(i);
      if (o.isName() || o.isImportStar()) {
        String before = o.getString();
        String after = now.get(i).getString();
        if (!before.equals(after)) {
          line(i + " " + before + " " + after);
        }
      }
    }
  }

  static void mdnu(Node script) throws java.io.IOException {
    List<Node> orig = preorder(script);
    CodingConvention convention = CodingConventions.getDefault();
    // ContextualRenamer, then the inverter, then the assertOnChange check.
    {
      Compiler compiler = newCompiler();
      Node clone = script.cloneTree();
      line("## contextual");
      try {
        NodeTraversal.traverse(compiler, clone, MakeDeclaredNamesUnique.builder().build());
        printRenames(orig, clone);
        line("## inverted");
        MakeDeclaredNamesUnique.getContextualRenameInverter(compiler).process(null, clone);
        printRenames(orig, clone);
        line("## assert");
        Node again = clone.cloneTree();
        NodeTraversal.traverse(
            compiler,
            again,
            MakeDeclaredNamesUnique.builder().withAssertOnChange(true).build());
        line("ok");
      } catch (RuntimeException e) {
        line("! " + exc(e));
      }
    }
    // InlineRenamer on the first three functions.
    {
      Compiler compiler = newCompiler();
      int done = 0;
      for (int i = 0; i < orig.size() && done < 3; i++) {
        if (!orig.get(i).isFunction()) {
          continue;
        }
        done++;
        line("## inline " + i);
        // As FunctionToBlockMutator.makeLocalNamesUnique does, on a detached clone.
        Node fn = orig.get(i).cloneTree();
        try {
          NodeTraversal.traverseScopeRoots(
              compiler,
              ImmutableList.of(fn),
              MakeDeclaredNamesUnique.builder()
                  .withRenamer(
                      new MakeDeclaredNamesUnique.InlineRenamer(
                          convention, counter(), "inline_", done % 2 == 0, true, null))
                  .withMarkChanges(false)
                  .build(),
              true);
          printRenames(preorder(orig.get(i)), fn);
        } catch (RuntimeException e) {
          line("! " + exc(e));
        }
      }
    }
    // BoilerplateRenamer
    {
      Compiler compiler = newCompiler();
      Node clone = script.cloneTree();
      line("## boilerplate");
      try {
        NodeTraversal.traverse(
            compiler,
            clone,
            MakeDeclaredNamesUnique.builder()
                .withRenamer(
                    new MakeDeclaredNamesUnique.BoilerplateRenamer(
                        convention, counter(), "boiler_"))
                .withMarkChanges(false)
                .build());
        printRenames(orig, clone);
      } catch (RuntimeException e) {
        line("! " + exc(e));
      }
    }
    // TargettedRenamer
    {
      Compiler compiler = newCompiler();
      Node clone = script.cloneTree();
      line("## targetted");
      try {
        NodeTraversal.traverse(
            compiler,
            clone,
            MakeDeclaredNamesUnique.builder()
                .withRenamer(
                    new MakeDeclaredNamesUnique.TargettedRenamer(
                        new MakeDeclaredNamesUnique.ContextualRenamer(),
                        ImmutableSet.of("i", "e", "a", "x", "n", "t")))
                .build());
        printRenames(orig, clone);
      } catch (RuntimeException e) {
        line("! " + exc(e));
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // factory

  /** A StaticScope whose slots are declared by NAME nodes (cf. rhino.testing.MapBasedScope). */
  static final class TestScope implements StaticScope {
    private final Map<String, Node> declarations = new LinkedHashMap<>();

    TestScope declare(String name, Node declaration) {
      declarations.put(name, declaration);
      return this;
    }

    @Override
    public Node getRootNode() {
      return null;
    }

    @Override
    public StaticScope getParentScope() {
      return null;
    }

    @Override
    public StaticSlot getSlot(String name) {
      return getOwnSlot(name);
    }

    @Override
    public StaticSlot getOwnSlot(String name) {
      Node declaration = declarations.get(name);
      if (declaration == null) {
        return null;
      }
      return new TestSlot(name, declaration, this);
    }

    @Override
    public String toString() {
      return "TestScope" + declarations.keySet();
    }
  }

  static final class TestSlot implements StaticSlot, StaticRef {
    private final String name;
    private final Node node;
    private final StaticScope scope;

    TestSlot(String name, Node node, StaticScope scope) {
      this.name = name;
      this.node = node;
      this.scope = scope;
    }

    @Override
    public String getName() {
      return name;
    }

    @Override
    public StaticRef getDeclaration() {
      return this;
    }

    @Override
    public JSDocInfo getJSDocInfo() {
      return null;
    }

    @Override
    public StaticScope getScope() {
      return scope;
    }

    @Override
    public StaticSlot getSymbol() {
      return this;
    }

    @Override
    public Node getNode() {
      return node;
    }

    @Override
    public StaticSourceFile getSourceFile() {
      return null;
    }

    @Override
    public String toString() {
      return "TestSlot(" + name + ")";
    }
  }

  static final class FactorySweep {
    interface Case {
      Node build();
    }

    final AstFactory f;
    final String mode;
    final JSTypeRegistry registry;
    final ColorRegistry colorRegistry;
    final RuntimeJsLibManager libs;

    FactorySweep(
        AstFactory f,
        String mode,
        JSTypeRegistry registry,
        ColorRegistry colorRegistry,
        RuntimeJsLibManager libs) {
      this.f = f;
      this.mode = mode;
      this.registry = registry;
      this.colorRegistry = colorRegistry;
      this.libs = libs;
    }

    static void run() throws java.io.IOException {
      for (LifeCycleStage stage : new LifeCycleStage[] {LifeCycleStage.RAW, LifeCycleStage.NORMALIZED}) {
        for (String mode : new String[] {"NONE", "JSTYPE", "COLOR"}) {
          JSTypeRegistry registry = new JSTypeRegistry(null);
          ColorRegistry colorRegistry =
              ColorRegistry.builder().setDefaultNativeColorsForTesting().build();
          RuntimeJsLibManager libs =
              RuntimeJsLibManager.create(
                  RuntimeJsLibManager.RuntimeLibraryMode.RECORD_AND_VALIDATE_FIELDS,
                  (unused1, unused2) -> {
                    throw new UnsupportedOperationException();
                  },
                  new ChangeTracker(),
                  () -> {
                    throw new UnsupportedOperationException();
                  });
          AstFactory f =
              switch (mode) {
                case "NONE" -> AstFactory.createFactoryWithoutTypes(stage, libs);
                case "JSTYPE" -> AstFactory.createFactoryWithTypes(stage, registry, libs);
                default -> AstFactory.createFactoryWithColors(stage, colorRegistry, libs);
              };
          line("# factory " + stage + " " + mode);
          new FactorySweep(f, mode, registry, colorRegistry, libs).sweep();
        }
      }
    }

    void c(String name, Case body) throws java.io.IOException {
      line("## " + name);
      try {
        Node n = body.build();
        if (n == null) {
          line("null");
        } else {
          printTree(n);
        }
      } catch (RuntimeException e) {
        line("! " + exc(e));
      }
    }

    /** A NAME with the type of the given native (or the corresponding color). */
    Node typed(Node n, JSTypeNative nativeType, com.google.javascript.jscomp.colors.Color color) {
      switch (mode) {
        case "JSTYPE" -> n.setJSType(registry.getNativeType(nativeType));
        case "COLOR" -> n.setColor(color);
        default -> {}
      }
      return n;
    }

    Node num(double d) {
      return typed(IR.number(d), JSTypeNative.NUMBER_TYPE, StandardColors.NUMBER);
    }

    Node str(String s) {
      return typed(IR.string(s), JSTypeNative.STRING_TYPE, StandardColors.STRING);
    }

    Node name(String s) {
      return typed(IR.name(s), JSTypeNative.NUMBER_TYPE, StandardColors.NUMBER);
    }

    Node untypedName(String s) {
      return IR.name(s);
    }

    TestScope scope() {
      Node objectDecl =
          typed(IR.name("Object"), JSTypeNative.OBJECT_FUNCTION_TYPE, StandardColors.TOP_OBJECT);
      Node assignDecl =
          typed(IR.name("assign"), JSTypeNative.FUNCTION_TYPE, StandardColors.TOP_OBJECT);
      Node constDecl = name("k");
      constDecl.putBooleanProp(Node.IS_CONSTANT_NAME, true);
      return new TestScope()
          .declare("a", name("a"))
          .declare("a.b", str("b"))
          .declare("a.b.c", num(3))
          .declare("k", constDecl)
          .declare("untyped", untypedName("untyped"))
          .declare("Object", objectDecl)
          .declare("Object.assign", assignDecl)
          .declare("$jscomp", untypedName("$jscomp"))
          .declare("$jscomp.global", untypedName("global"))
          .declare("$jscomp.makeIterator", untypedName("makeIterator"))
          .declare("$jscomp.arrayFromIterator", untypedName("arrayFromIterator"))
          .declare("$jscomp.arrayFromIterable", untypedName("arrayFromIterable"))
          .declare("$jscomp.makeAsyncIterator", untypedName("makeAsyncIterator"))
          .declare(
              "$jscomp.asyncExecutePromiseGeneratorFunction",
              untypedName("asyncExecutePromiseGeneratorFunction"))
          .declare("$jscomp.global.Math", untypedName("Math"))
          .declare("$jscomp.global.Math.max", untypedName("max"))
          .declare("Math", untypedName("Math"))
          .declare("Math.max", untypedName("max"));
    }

    AstFactory.Type numberType() {
      return AstFactory.type(JSTypeNative.NUMBER_TYPE, StandardColors.NUMBER);
    }

    AstFactory.Type functionType() {
      return AstFactory.type(JSTypeNative.FUNCTION_TYPE, StandardColors.TOP_OBJECT);
    }

    void sweep() throws java.io.IOException {
      c("exprResult", () -> f.exprResult(num(1)));
      c("createEmpty", () -> f.createEmpty());
      c("createBlock", () -> f.createBlock(f.exprResult(num(1)), f.createEmpty()));
      c("createIf", () -> f.createIf(name("x"), f.createBlock()));
      c("createIfElse", () -> f.createIf(name("x"), f.createBlock(), f.createBlock()));
      c("createFor", () -> f.createFor(f.createEmpty(), name("x"), f.createEmpty(), f.createBlock()));
      c("createBreak", () -> f.createBreak());
      c("createLabel", () -> f.createLabel(IR.labelName("L"), f.createBlock()));
      c("createCatch", () -> f.createCatch(name("e"), f.createBlock()));
      c("createTryFinally", () -> f.createTryFinally(f.createBlock(), f.createBlock()));
      c(
          "createTryCatchFinally",
          () ->
              f.createTryCatchFinally(
                  f.createBlock(), f.createCatch(name("e"), f.createBlock()), f.createBlock()));
      c(
          "createTryCatchFinally/bad",
          () -> f.createTryCatchFinally(f.createBlock(), f.createBlock(), f.createBlock()));
      c("createThrow", () -> f.createThrow(str("x")));
      c("createReturn", () -> f.createReturn(num(2)));
      c("createYield", () -> f.createYield(numberType(), num(1)));
      c("createAwait", () -> f.createAwait(AstFactory.type(name("p")), name("p")));
      c("createString", () -> f.createString("hello"));
      c("createString/unicode", () -> f.createString("hé😀"));
      c("createNumber", () -> f.createNumber(0.5));
      c("createNumber/negative", () -> f.createNumber(-0.5));
      c("createBoolean/true", () -> f.createBoolean(true));
      c("createBoolean/false", () -> f.createBoolean(false));
      c("createNull", () -> f.createNull());
      c("createVoid", () -> f.createVoid(name("x")));
      c("createUndefinedValue", () -> f.createUndefinedValue());
      c("createNot", () -> f.createNot(name("x")));
      c("createThis", () -> f.createThis(numberType()));
      c("createSuper", () -> f.createSuper(AstFactory.type(name("s"))));
      c("createThisForEs6Class/notClass", () -> f.createThisForEs6Class(name("x")));
      c(
          "createSingleNameDeclaration/let",
          () -> f.createSingleNameDeclaration(Token.LET, "v", num(1)));
      c(
          "createSingleNameDeclaration/var",
          () -> f.createSingleNameDeclaration(Token.VAR, "v", num(1)));
      c(
          "createSingleNameDeclaration/const",
          () -> f.createSingleNameDeclaration(Token.CONST, "v", num(1)));
      c(
          "createSingleNameDeclaration/bad",
          () -> f.createSingleNameDeclaration(Token.NAME, "v", num(1)));
      c("createSingleLetNameDeclaration", () -> f.createSingleLetNameDeclaration("v"));
      c(
          "createSingleLetNameDeclaration/value",
          () -> f.createSingleLetNameDeclaration("v", str("s")));
      c("createSingleVarNameDeclaration", () -> f.createSingleVarNameDeclaration("v"));
      c(
          "createSingleVarNameDeclaration/value",
          () -> f.createSingleVarNameDeclaration("$jscomp$v", str("s")));
      c(
          "createSingleConstNameDeclaration",
          () -> f.createSingleConstNameDeclaration("v", str("s")));
      c(
          "createSingleConstObjectPatternDeclaration",
          () ->
              f.createSingleConstObjectPatternDeclaration(
                  IR.objectPattern(IR.stringKey("Foo")), name("x")));
      c(
          "createSingleConstObjectPatternDeclaration/bad",
          () -> f.createSingleConstObjectPatternDeclaration(name("y"), name("x")));
      c("createArgumentsReference", () -> f.createArgumentsReference());
      c("createArgumentsAliasDeclaration", () -> f.createArgumentsAliasDeclaration("args"));
      c("createName", () -> f.createName("n", numberType()));
      c("createName/jscompPrefix", () -> f.createName("$jscomp$x", numberType()));
      c("createName/jscomp", () -> f.createName("$jscomp", numberType()));
      c("createConstantName", () -> f.createConstantName("n", numberType()));
      c("createConstantName/jscomp", () -> f.createConstantName("$jscomp", numberType()));
      c("createName/scope", () -> f.createName(scope(), "a"));
      c("createName/scope/const", () -> f.createName(scope(), "k"));
      c("createName/scope/untyped", () -> f.createName(scope(), "untyped"));
      c("createName/scope/missing", () -> f.createName(scope(), "missing"));
      c("createName/nullScope", () -> f.createName((StaticScope) null, "a"));
      c("createName/scope/jscomp", () -> f.createName(scope(), "$jscomp"));
      c("createNameWithUnknownType", () -> f.createNameWithUnknownType("u"));
      c("createQName/dotted", () -> f.createQName(scope(), "a.b.c"));
      c("createQName/missing", () -> f.createQName(scope(), "a.b.missing"));
      c("createQName/iterable", () -> f.createQName(scope(), ImmutableList.of("a", "b")));
      c("createQName/varargs", () -> f.createQName(scope(), "a", "b", "c"));
      c(
          "createQName/baseAndIterable",
          () -> f.createQName(scope(), "a", ImmutableList.of("b", "c")));
      c("createQName/field/notInjected", () -> f.createQName(scope(), libs.getJsLibField("$jscomp.global")));
      c(
          "createQName/field",
          () -> {
            libs.injectLibForField("$jscomp.global");
            return f.createQName(scope(), libs.getJsLibField("$jscomp.global"));
          });
      c("createQNameWithUnknownType", () -> f.createQNameWithUnknownType("p.q.r"));
      c(
          "createQNameWithUnknownType/field",
          () -> f.createQNameWithUnknownType(libs.getJsLibField("$jscomp.global")));
      c(
          "createQNameWithUnknownType/baseAndIterable",
          () -> f.createQNameWithUnknownType("p", ImmutableList.of("q", "r")));
      c("createPrototypeAccess", () -> f.createPrototypeAccess(name("Foo")));
      c("createPrototypeAccess/untyped", () -> f.createPrototypeAccess(untypedName("Foo")));
      c("createJSCompDotGlobalAccess", () -> f.createJSCompDotGlobalAccess(scope(), "Math.max"));
      c("createGetPropWithoutColor", () -> f.createGetPropWithoutColor(name("o"), "p"));
      c(
          "createGetPropWithoutColor/jscompGlobal",
          () -> f.createGetPropWithoutColor(untypedName("$jscomp"), "global"));
      c(
          "createGetPropWithoutColor/string",
          () -> f.createGetPropWithoutColor(str("s"), "length"));
      c("createGetProp", () -> f.createGetProp(name("o"), "p", numberType()));
      c(
          "createGetPropsWithoutColors",
          () -> f.createGetPropsWithoutColors(name("o"), ImmutableList.of("p", "q")));
      c(
          "createGetPropsWithUnknownType",
          () -> f.createGetPropsWithUnknownType(name("o"), ImmutableList.of("p", "q")));
      c("createGetPropWithUnknownType", () -> f.createGetPropWithUnknownType(name("o"), "p"));
      c(
          "createStartOptChainGetprop",
          () -> f.createStartOptChainGetprop(name("o"), "p", numberType()));
      c(
          "createContinueOptChainGetprop",
          () ->
              f.createContinueOptChainGetprop(
                  f.createStartOptChainGetprop(name("o"), "p", numberType()), "q", numberType()));
      c(
          "createStartOptChainGetelem",
          () -> f.createStartOptChainGetelem(name("o"), num(0), numberType()));
      c(
          "createContinueOptChainGetelem",
          () ->
              f.createContinueOptChainGetelem(
                  f.createStartOptChainGetelem(name("o"), num(0), numberType()),
                  num(1),
                  numberType()));
      c(
          "createStartOptChainCall",
          () -> f.createStartOptChainCall(name("o"), numberType(), num(1), str("x")));
      c(
          "createContinueOptChainCall",
          () ->
              f.createContinueOptChainCall(
                  f.createStartOptChainGetprop(name("o"), "p", numberType()), numberType()));
      c("createGetElem", () -> f.createGetElem(name("o"), num(0)));
      c("createDelProp", () -> f.createDelProp(f.createGetPropWithUnknownType(name("o"), "p")));
      c("createStringKey", () -> f.createStringKey("k", str("v")));
      c("createComputedProperty", () -> f.createComputedProperty(str("k"), num(1)));
      c("createGetterDef", () -> f.createGetterDef("g", num(1)));
      c("createIn", () -> f.createIn(str("k"), name("o")));
      c("createComma", () -> f.createComma(name("a"), str("b")));
      c("createCommas", () -> f.createCommas(name("a"), str("b"), num(1), name("c")));
      c("createAnd", () -> f.createAnd(name("a"), str("b")));
      c("createAnd/untyped", () -> f.createAnd(untypedName("a"), str("b")));
      c("createOr", () -> f.createOr(str("a"), num(1)));
      c("createAdd", () -> f.createAdd(str("a"), num(1)));
      c("createSub", () -> f.createSub(name("a"), num(1)));
      c("createInc/pre", () -> f.createInc(name("a"), false));
      c("createInc/post", () -> f.createInc(name("a"), true));
      c("createLessThan", () -> f.createLessThan(name("a"), num(1)));
      c("createBitwiseAnd", () -> f.createBitwiseAnd(name("a"), num(1)));
      c("createRightShift", () -> f.createRightShift(name("a"), num(1)));
      c("createCall", () -> f.createCall(name("fn"), numberType(), num(1), str("x")));
      c("createCall/getprop", () -> f.createCall(f.createGetPropWithUnknownType(name("o"), "m"), numberType()));
      c("createCallWithUnknownType", () -> f.createCallWithUnknownType(name("fn"), num(1)));
      c(
          "createObjectDotAssignCall",
          () -> f.createObjectDotAssignCall(scope(), numberType(), f.createObjectLit(), name("a")));
      c("createNewNode/unknownCtor", () -> f.createNewNode(untypedName("C"), num(1)));
      c(
          "createNewNode/objectCtor",
          () ->
              f.createNewNode(
                  typed(
                      IR.name("Object"),
                      JSTypeNative.OBJECT_FUNCTION_TYPE,
                      colorRegistry.get(StandardColors.ARRAY_ID)),
                  num(1)));
      c(
          "createConstructorCall",
          () ->
              f.createConstructorCall(
                  AstFactory.type(JSTypeNative.OBJECT_FUNCTION_TYPE, StandardColors.TOP_OBJECT),
                  IR.superNode(),
                  num(1)));
      c("createAssignStatement", () -> f.createAssignStatement(name("a"), str("b")));
      c("createAssign", () -> f.createAssign(name("a"), str("b")));
      c("createAssign/name", () -> f.createAssign("a", str("b")));
      c("createObjectLit", () -> f.createObjectLit(f.createStringKey("k", num(1))));
      c(
          "createObjectLit/type",
          () -> f.createObjectLit(numberType(), f.createStringKey("k", num(1))));
      c("createQuotedStringKey", () -> f.createQuotedStringKey("k", num(1)));
      c("createEmptyFunction", () -> f.createEmptyFunction(functionType()));
      c("createEmptyGeneratorFunction", () -> f.createEmptyGeneratorFunction(functionType()));
      c(
          "createFunction",
          () ->
              f.createFunction(
                  "fn", f.createParamList("p", "q"), f.createBlock(), functionType()));
      c("createParamList", () -> f.createParamList("p", "$jscomp$q"));
      c(
          "createZeroArgFunction",
          () ->
              f.createZeroArgFunction(
                  "z", f.createBlock(f.createReturn(num(1))), registry.getNativeType(JSTypeNative.NUMBER_TYPE)));
      c(
          "createZeroArgFunction/nullReturn",
          () -> f.createZeroArgFunction("z", f.createBlock(), null));
      c(
          "createZeroArgGeneratorFunction",
          () -> f.createZeroArgGeneratorFunction("g", f.createBlock(), null));
      c(
          "createZeroArgArrowFunctionForExpression",
          () -> f.createZeroArgArrowFunctionForExpression(num(1)));
      c(
          "createZeroArgArrowFunctionForExpression/untyped",
          () -> f.createZeroArgArrowFunctionForExpression(untypedName("u")));
      c(
          "createMemberFunctionDef",
          () -> f.createMemberFunctionDef("m", f.createEmptyFunction(functionType())));
      c(
          "createMemberFunctionDef/named",
          () ->
              f.createMemberFunctionDef(
                  "m", f.createFunction("n", IR.paramList(), f.createBlock(), functionType())));
      c("createSheq", () -> f.createSheq(name("a"), num(1)));
      c("createEq", () -> f.createEq(name("a"), f.createNull()));
      c("createNe", () -> f.createNe(name("a"), f.createNull()));
      c("createHook", () -> f.createHook(name("a"), str("b"), num(1)));
      c("createHook/untyped", () -> f.createHook(name("a"), untypedName("b"), num(1)));
      c("createArraylit", () -> f.createArraylit(num(1), str("b")));
      c("createArraylit/iterable", () -> f.createArraylit(ImmutableList.of(num(1), name("c"))));
      c(
          "createJSCompMakeIteratorCall/notInjected",
          () -> f.createJSCompMakeIteratorCall(name("arr"), scope()));
      c(
          "createJSCompMakeIteratorCall",
          () -> {
            libs.injectLibForField("$jscomp.makeIterator");
            return f.createJSCompMakeIteratorCall(name("arr"), scope());
          });
      c(
          "createJscompArrayFromIteratorCall",
          () -> {
            libs.injectLibForField("$jscomp.arrayFromIterator");
            return f.createJscompArrayFromIteratorCall(name("it"), scope());
          });
      c(
          "createJscompArrayFromIterableCall",
          () -> {
            libs.injectLibForField("$jscomp.arrayFromIterable");
            return f.createJscompArrayFromIterableCall(name("it"), scope());
          });
      c(
          "createJSCompMakeAsyncIteratorCall",
          () -> {
            libs.injectLibForField("$jscomp.makeAsyncIterator");
            return f.createJSCompMakeAsyncIteratorCall(name("it"), scope());
          });
      c(
          "createEmptyAsyncGeneratorWrapperArgument/unknown",
          () ->
              f.createEmptyAsyncGeneratorWrapperArgument(
                  registry.getNativeType(JSTypeNative.UNKNOWN_TYPE)));
      c(
          "createEmptyAsyncGeneratorWrapperArgument/function",
          () ->
              f.createEmptyAsyncGeneratorWrapperArgument(
                  registry.createFunctionType(
                      registry.getNativeType(JSTypeNative.VOID_TYPE),
                      registry.getNativeType(JSTypeNative.NUMBER_TYPE))));
      c(
          "createJscompAsyncExecutePromiseGeneratorFunctionCall",
          () -> {
            libs.injectLibForField("$jscomp.asyncExecutePromiseGeneratorFunction");
            return f.createJscompAsyncExecutePromiseGeneratorFunctionCall(
                scope(), f.createEmptyGeneratorFunction(functionType()));
          });
      c("type/colorId", () -> f.createThis(AstFactory.type(StandardColors.ARRAY_ID)));
      c(
          "type/jstypeNative",
          () -> f.createThis(AstFactory.type(JSTypeNative.STRING_TYPE)));
      c("type/color", () -> f.createThis(AstFactory.type(StandardColors.STRING)));
      c(
          "type/jstype",
          () -> f.createThis(AstFactory.type(registry.getNativeType(JSTypeNative.STRING_TYPE))));
      line("## isAddingTypes " + f.isAddingTypes() + " isAddingColors " + f.isAddingColors());
    }
  }
}

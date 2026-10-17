package com.google.javascript.jscomp;

// TypedScopeDump: the oracle for TypedScopeCreator (and FunctionTypeBuilder).
//
//   java -Xss1g -cp closure-compiler.jar:build com.google.javascript.jscomp.TypedScopeDump \
//        <out.txt.gz> <externs.js> <root> <rel-file>... | @<listfile>
//
// For each input file (or directory: all its files, in name order, as one compilation), compiled
// as the only source next to <externs.js> (the DEFAULT_EXTERNS of
// CompilerTestCase, checked in as internal/jscomp/testdata/typed_scope/default_externs.js) with
// the options CompilerTestCase.getOptions() sets for TypedScopeCreatorTest, this runs the
// prerequisites TypedScopeCreatorTest.getProcessor() runs (GatherModuleMetadata,
// ModuleMapCreator, InferConsts), then the scope-building half of
// TypeInferencePass.inferAllScopes: TypedScopeCreator.createScope(ROOT, null) and a scoped
// traversal that creates every scope, inside JSTypeResolver.openForDefinition(), then
// resolveWeakImportsPreResolution() and finishAndFreeze(). No type inference runs. Output, one
// section per file ("=== <rel-file>"), every line escaped to printable ASCII:
//
//   P <count>                       parse errors, followed by the D lines of parsing (nothing
//                                   else is printed for the file)
//   X <class>: <message>            an exception escaped (the innermost cause, first message line)
//   D <level> <key> <line>:<col> <description>
//                                   each diagnostic reported, in report order
//   S <idx> <parentIdx> <token> <source>:<line>:<col>
//                                   each scope, in creation order (externs-rooted scopes skipped)
//   V <name> | <type> | <annotation> | inferred=<b> esc=<b> once=<b> input=<name> node=<pos>
//                                   each TypedVar of the scope, in declaration order; in the
//                                   global scope, lines identical to the baseline are skipped
//   T <row> <name> | <type>         every entry of JSTypeRegistry's scoped name table, sorted, except
//                                   lines identical to the baseline; <row> is GLOBAL or the row's
//                                   root node <source>:<token>@<line>:<col> (with a trailing "!"
//                                   if no scope was created for it)
//   N <token> <source>:<line>:<col> | <type>
//                                   the JSType TypedScopeCreator attached to each non-literal node
//                                   of the source (pre-order)
//
// The baseline is the same dump of an empty source file, printed first as "=== <baseline>".
// The Go side is internal/jscomp/typed_scope_creator_oracle_test.go.

import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.NodeTraversal.AbstractScopedCallback;
import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.js.RuntimeJsLibManager;
import com.google.javascript.jscomp.modules.ModuleMapCreator;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.Table;
import com.google.javascript.jscomp.parsing.Config.JsDocParsing;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.Token;
import com.google.javascript.rhino.jstype.JSType;
import com.google.javascript.rhino.jstype.JSTypeRegistry;
import com.google.javascript.rhino.jstype.JSTypeResolver;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Field;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.IdentityHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class TypedScopeDump {
  static final String BASELINE = "<baseline>";

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    String externs = Files.readString(Paths.get(args[1]), StandardCharsets.UTF_8);
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
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(out))) {
      Set<String> baseline = new HashSet<>();
      List<String> lines =
          dump(externs, ImmutableList.of(SourceFile.fromCode(BASELINE, "")), null);
      baseline.addAll(lines);
      write(os, BASELINE, lines);
      for (String rel : files) {
        Path p = root.resolve(rel);
        List<SourceFile> inputs = new ArrayList<>();
        if (Files.isDirectory(p)) {
          // A directory is one compilation of all its files, in name order.
          List<String> names = new ArrayList<>();
          try (var stream = Files.list(p)) {
            stream.forEach(f -> names.add(f.getFileName().toString()));
          }
          Collections.sort(names);
          for (String name : names) {
            inputs.add(
                SourceFile.fromCode(
                    rel + "/" + name, Files.readString(p.resolve(name), StandardCharsets.UTF_8)));
          }
        } else {
          inputs.add(SourceFile.fromCode(rel, Files.readString(p, StandardCharsets.UTF_8)));
        }
        write(os, rel, dump(externs, inputs, baseline));
      }
    }
  }

  private static void write(OutputStream os, String name, List<String> lines) throws Exception {
    StringBuilder sb = new StringBuilder();
    sb.append("=== ").append(escape(name)).append('\n');
    for (String l : lines) {
      sb.append(escape(l)).append('\n');
    }
    os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
  }

  /** The options of CompilerTestCase.getOptions() as TypedScopeCreatorTest sets them up. */
  static CompilerOptions options() {
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.UNSUPPORTED);
    options.setEmitUseStrict(false);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    options.setModuleResolutionMode(ResolutionMode.BROWSER);
    options.setParseJsDocDocumentation(JsDocParsing.TYPES_ONLY);
    options.setAssumeStaticInheritanceIsNotUsed(true);
    options.setPreserveTypeAnnotations(true);
    options.setAssumeGettersArePure(false);
    options.setCheckSymbols(true);
    options.setWarningLevel(DiagnosticGroups.INVALID_CASTS, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.MISPLACED_MSG_ANNOTATION, CheckLevel.WARNING);
    options.setWarningLevel(DiagnosticGroups.MISSING_PROPERTIES, CheckLevel.WARNING);
    options.setCodingConvention(new GoogleCodingConvention());
    options.setPolymerVersion(1);
    options.setRuntimeLibraryMode(RuntimeJsLibManager.RuntimeLibraryMode.RECORD_AND_VALIDATE_FIELDS);
    // TypedScopeCreatorTest.getOptions()
    options.setClosurePass(false);
    return options;
  }

  /** Records every report in order. */
  static final class RecordingErrorManager implements ErrorManager {
    final List<String> lines = new ArrayList<>();
    final List<JSError> errors = new ArrayList<>();
    final List<JSError> warnings = new ArrayList<>();
    double typedPercent;

    @Override
    public void report(CheckLevel level, JSError error) {
      lines.add(
          "D "
              + level
              + " "
              + error.type().key
              + " "
              + error.lineno()
              + ":"
              + error.charno()
              + " "
              + error.description());
      if (level == CheckLevel.ERROR) {
        errors.add(error);
      } else if (level == CheckLevel.WARNING) {
        warnings.add(error);
      }
    }

    @Override
    public void generateReport() {}

    @Override
    public int getErrorCount() {
      return errors.size();
    }

    @Override
    public int getWarningCount() {
      return warnings.size();
    }

    @Override
    public ImmutableList<JSError> getErrors() {
      return ImmutableList.copyOf(errors);
    }

    @Override
    public ImmutableList<JSError> getWarnings() {
      return ImmutableList.copyOf(warnings);
    }

    @Override
    public void setTypedPercent(double typedPercent) {
      this.typedPercent = typedPercent;
    }

    @Override
    public double getTypedPercent() {
      return typedPercent;
    }
  }

  static List<String> dump(String externs, List<SourceFile> inputs, Set<String> baseline) {
    List<String> out = new ArrayList<>();
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    compiler.setAllowableFeatures(LanguageMode.UNSUPPORTED.toFeatureSet());
    compiler.setPreferRegexParser(false);
    RecordingErrorManager em = new RecordingErrorManager();
    compiler.init(
        ImmutableList.of(SourceFile.fromCode("externs", externs)),
        inputs,
        options());
    compiler.setErrorManager(em);
    Node root = compiler.parseInputs();
    if (root == null || compiler.getErrorCount() > 0) {
      out.add("P " + compiler.getErrorCount());
      out.addAll(em.lines);
      return out;
    }
    em.lines.clear();
    Node externsRoot = root.getFirstChild();
    Node jsRoot = root.getLastChild();
    TypedScopeCreator scopeCreator = null;
    try {
      new GatherModuleMetadata(compiler, false, ResolutionMode.BROWSER)
          .process(externsRoot, jsRoot);
      new ModuleMapCreator(compiler, compiler.getModuleMetadataMap()).process(externsRoot, jsRoot);
      new InferConsts(compiler).process(externsRoot, jsRoot);

      JSTypeRegistry registry = compiler.getTypeRegistry();
      scopeCreator = new TypedScopeCreator(compiler);
      try (JSTypeResolver.Closer closer = registry.getResolver().openForDefinition()) {
        TypedScope topScope = scopeCreator.createScope(root, null);
        NodeTraversal.builder()
            .setCompiler(compiler)
            .setCallback(
                new AbstractScopedCallback() {
                  @Override
                  public void enterScope(NodeTraversal t) {
                    t.getTypedScope();
                  }

                  @Override
                  public void visit(NodeTraversal t, Node n, Node parent) {}
                })
            .setScopeCreator(scopeCreator)
            .traverseWithScope(root, topScope);
        scopeCreator.resolveWeakImportsPreResolution();
      }
      scopeCreator.finishAndFreeze();
    } catch (Throwable e) {
      out.addAll(em.lines);
      // Compiler wraps exceptions of passes in RuntimeException("INTERNAL COMPILER ERROR...");
      // report the original one.
      while (e.getCause() != null && e.getCause() != e) {
        e = e.getCause();
      }
      String msg = String.valueOf(e.getMessage());
      int nl = msg.indexOf('\n');
      out.add("X " + e.getClass().getName() + ": " + (nl >= 0 ? msg.substring(0, nl) : msg));
      return out;
    }
    out.addAll(em.lines);

    List<TypedScope> scopes = new ArrayList<>(scopeCreator.getAllMemoizedScopes());
    Collections.reverse(scopes);
    Map<Node, Integer> scopeIndex = new IdentityHashMap<>();
    for (int i = 0; i < scopes.size(); i++) {
      scopeIndex.put(scopes.get(i).getRootNode(), i);
    }
    for (int i = 0; i < scopes.size(); i++) {
      TypedScope s = scopes.get(i);
      Node r = s.getRootNode();
      if (i > 0 && r.isFromExterns()) {
        continue;
      }
      Integer parent = s.getParent() == null ? -1 : scopeIndex.get(s.getParent().getRootNode());
      out.add("S " + i + " " + parent + " " + r.getToken() + " " + r.getSourceFileName() + ":"
          + r.getLineno() + ":" + r.getCharno());
      for (TypedVar v : s.getVarIterable()) {
        String line = varLine(v);
        if (i == 0 && baseline != null && baseline.contains(line)) {
          continue;
        }
        out.add(line);
      }
    }

    List<String> table = new ArrayList<>();
    try {
      JSTypeRegistry registry = compiler.getTypeRegistry();
      Field tf = JSTypeRegistry.class.getDeclaredField("scopedNameTable");
      tf.setAccessible(true);
      Field gf = JSTypeRegistry.class.getDeclaredField("nameTableGlobalRoot");
      gf.setAccessible(true);
      @SuppressWarnings("unchecked")
      Table<Node, String, JSType> t = (Table<Node, String, JSType>) tf.get(registry);
      Node globalRoot = (Node) gf.get(registry);
      for (Table.Cell<Node, String, JSType> c : t.cellSet()) {
        String row;
        if (c.getRowKey() == globalRoot) {
          row = "GLOBAL";
        } else {
          Node r = c.getRowKey();
          row = r.getSourceFileName() + ":" + pos(r) + (scopeIndex.containsKey(r) ? "" : "!");
        }
        String line = "T " + row + " " + c.getColumnKey() + " | " + typeString(c.getValue());
        if (baseline != null && baseline.contains(line)) {
          continue;
        }
        table.add(line);
      }
    } catch (ReflectiveOperationException e) {
      throw new LinkageError("reflection", e);
    }
    Collections.sort(table);
    out.addAll(table);

    nodeTypes(out, jsRoot);
    return out;
  }

  static String varLine(TypedVar v) {
    JSType type = v.getType();
    String ann = "-";
    if (type != null && !v.isTypeInferred()) {
      try {
        ann = type.toAnnotationString(JSType.Nullability.EXPLICIT);
      } catch (RuntimeException e) {
        ann = "!" + e.getClass().getSimpleName();
      }
    }
    return "V "
        + v.getName()
        + " | "
        + typeString(type)
        + " | "
        + ann
        + " | inferred="
        + v.isTypeInferred()
        + " esc="
        + v.isMarkedEscaped()
        + " once="
        + v.isMarkedAssignedExactlyOnce()
        + " input="
        + v.getInputName()
        + " node="
        + (v.getNameNode() == null ? "null" : pos(v.getNameNode()));
  }

  static String typeString(JSType type) {
    if (type == null) {
      return "null";
    }
    try {
      return type.toString();
    } catch (RuntimeException e) {
      return "!" + e.getClass().getSimpleName();
    }
  }

  static String pos(Node n) {
    return n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
  }

  private static void nodeTypes(List<String> out, Node n) {
    JSType type = n.getJSType();
    if (type != null && !isLiteral(n.getToken())) {
      out.add("N " + n.getToken() + " " + n.getSourceFileName() + ":" + n.getLineno() + ":"
          + n.getCharno() + " | "
          + typeString(type));
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      nodeTypes(out, c);
    }
  }

  private static boolean isLiteral(Token t) {
    switch (t) {
      case NULL:
      case VOID:
      case STRINGLIT:
      case TEMPLATELIT_STRING:
      case NUMBER:
      case BIGINT:
      case TRUE:
      case FALSE:
      case REGEXP:
      case ARRAYLIT:
        return true;
      default:
        return false;
    }
  }

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

  private TypedScopeDump() {}
}

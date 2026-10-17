// Differential oracle for the TypedAST serialization passes (internal/jscomp/
// serialization_typed_ast_serializer.go, serialization_serialize_typed_ast_pass.go,
// serialization_convert_types_to_colors.go) and, through the round trip, the deserializer
// (serialization_typed_ast_deserializer.go, serialization_script_node_deserializer.go).
//
// For every file, each variant starts from a fresh Compiler with TypeCheckDump.options()
// (CompilerTypeTestCase's options), the externs <file> as "externs.js" and the input as one
// SourceFile named <rel>: Compiler.init, Compiler.parse. Files with parse errors print their
// diagnostics and are skipped. The variants:
//   typed     type checking exactly as TypeCheckDump does it (GatherModuleMetadata,
//             ModuleMapCreator, InferConsts, then DefaultPassConfig's inferTypes and checkTypes),
//             then SerializeTypedAstPass.createFromOutputStream(compiler, out,
//             SerializationOptions.builder().setIncludeDebugInfo(false).setRuntimeLibraries([]))
//             .process(externsRoot, jsRoot): the JSType branch of TypedAstSerializer (as the
//             checks-phase serializeTypedAst pass of DefaultPassConfig).
//   colors    the same type checking, then RemoveCastNodes and ConvertTypesToColors(includeDebugInfo
//             false) + setLifeCycleStage(COLORS_AND_SIMPLIFIED_JSDOC) (DefaultPassConfig's
//             removeCastNodes and typesToColors), then the same serialization: the Color branch
//             (as Compiler.saveState after the checks).
//   untyped   no type checking: the serialization only (TypedAstSerializer's NoOpTypeSerializer).
// After each serialization the bytes are deserialized in another fresh Compiler (same options,
// init with the same SourceFile objects) with TypedAstDeserializer.deserializeFullAst(compiler,
// SourceFile.fromCode("syntheticExterns", "", EXTERN), {externs.js, <rel>}, bytes,
// includeTypeInformation=true, resolveSourceMapAnnotations, parseInlineSourceMaps from the
// options) and the scripts of externs.js and <rel> are taken from its filesystem, as
// SerializeAndDeserializeAstTest does.
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside printable
// ASCII as \\uXXXX):
//   ## file <rel>
//   # parse ok|errors                 then the parse diagnostics:
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   ## variant typed|colors|untyped
//   checks <n> <sha16>                the diagnostics reported before serialization (type checking
//                                     and ConvertTypesToColors), count and hash of their diag lines
//   exception <stage> <class>: <message>   if a step threw (the variant ends); stage is checks,
//                                     serialize or deserialize
//   typedast <bytes> <sha16>          the serialized TypedAst.List: size and SHA-256 of its bytes
//   text <n> <sha16>                  the TypedAst in text form: TypedAst.toString() without its
//                                     LazyAsts, then per LazyAst "extern_ast" or "code_ast", the
//                                     LazyAst.toString() without its script and
//                                     AstNode.parseFrom(script).toString(); count and hash of the
//                                     lines, and with --full the lines themselves
//   deserialized <n> <sha16>          per script (externs.js, then <rel>): "script <name>" and one
//                                     line per node ("script <name> absent" if the filesystem
//                                     has no such file, e.g. for @nocompile), 4 spaces per depth + Node.toString(true,
//                                     true, false) + " : <ColorId>" when it has a color; then one
//                                     line per distinct color in first-seen order:
//                                     "color <id> inv=.. ctor=.. keep=.. assert=.. box=<id|null>
//                                     own=[..] proto=[ids] inst=[ids] union=[ids]";
//                                     count and hash, and with --full the lines
// sha16 is the first 16 hex digits of the SHA-256 of the escaped lines (each followed by '\n').
//
// Usage: java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//          com.google.javascript.jscomp.SerializationDump --externs=<file> [--full]
//          [--sample=<N>] [--full-every=<M>] <out.gz> <root> <rel-file>... | --list=<listfile>
// (compile together with TypeCheckDump.java). --sample=N dumps every Nth file of the list only;
// --full-every=M prints the full lines for every Mth dumped file (--full: all of them).
// The Go side is internal/jscomp/serialization_oracle_test.go.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage;
import com.google.javascript.jscomp.colors.Color;
import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.modules.ModuleMapCreator;
import com.google.javascript.jscomp.serialization.AstNode;
import com.google.javascript.jscomp.serialization.ConvertTypesToColors;
import com.google.javascript.jscomp.serialization.LazyAst;
import com.google.javascript.jscomp.serialization.SerializationOptions;
import com.google.javascript.jscomp.serialization.SerializeTypedAstPass;
import com.google.javascript.jscomp.serialization.TypedAst;
import com.google.javascript.jscomp.serialization.TypedAstDeserializer;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.StaticSourceFile.SourceKind;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSet;
import com.google.javascript.jscomp.jarjar.com.google.protobuf.CodedInputStream;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class SerializationDump {

  private static boolean full = false;
  private static int sample = 1;
  private static int fullEvery = 0;

  static final String[] VARIANTS = {"typed", "colors", "untyped"};

  public static void main(String[] args) throws Exception {
    String externsFile = null;
    List<String> rest = new ArrayList<>();
    List<String> files = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
      } else if (a.startsWith("--sample=")) {
        sample = Integer.parseInt(a.substring("--sample=".length()));
      } else if (a.startsWith("--full-every=")) {
        fullEvery = Integer.parseInt(a.substring("--full-every=".length()));
      } else if (a.startsWith("--externs=")) {
        externsFile = a.substring("--externs=".length());
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
    if (externsFile == null) {
      throw new IllegalArgumentException("--externs=<file> is required");
    }
    String externs = Files.readString(Paths.get(externsFile), StandardCharsets.UTF_8);
    Path out = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    files.addAll(rest.subList(2, rest.size()));
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(out))) {
      int index = 0;
      int dumped = 0;
      for (String rel : files) {
        if (index++ % sample != 0) {
          continue;
        }
        boolean fullThis = full || (fullEvery > 0 && dumped % fullEvery == 0);
        dumped++;
        StringBuilder sb = new StringBuilder();
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        dumpFile(sb, externs, rel, src, fullThis);
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  private static final class Setup {
    final TypeCheckDump.RecordingErrorManager errors = new TypeCheckDump.RecordingErrorManager();
    final Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    final CompilerOptions options = TypeCheckDump.options();
    final SourceFile externsFile;
    final SourceFile inputFile;

    Setup(SourceFile externsFile, SourceFile inputFile) {
      this.externsFile = externsFile;
      this.inputFile = inputFile;
      compiler.setErrorManager(errors);
      compiler.init(ImmutableList.of(externsFile), ImmutableList.of(inputFile), options);
    }
  }

  static void dumpFile(StringBuilder sb, String externs, String rel, String src, boolean fullThis)
      throws Exception {
    line(sb, "## file " + rel);
    SourceFile externsFile = SourceFile.fromCode("externs.js", externs);
    SourceFile inputFile = SourceFile.fromCode(rel, src);
    {
      Setup s = new Setup(externsFile, inputFile);
      try {
        s.compiler.parse();
      } catch (RuntimeException | StackOverflowError e) {
        line(sb, "# parse exception " + e.getClass().getName() + ": " + e.getMessage());
        return;
      }
      boolean ok = s.errors.getErrorCount() == 0;
      line(sb, "# parse " + (ok ? "ok" : "errors"));
      if (!ok) {
        for (String d : s.errors.lines) {
          line(sb, "diag " + d);
        }
        return;
      }
    }
    for (String variant : VARIANTS) {
      line(sb, "## variant " + variant);
      dumpVariant(sb, variant, externsFile, inputFile, fullThis);
    }
  }

  private static void dumpVariant(
      StringBuilder sb, String variant, SourceFile externsFile, SourceFile inputFile,
      boolean fullThis)
      throws Exception {
    Setup s = new Setup(externsFile, inputFile);
    String stage = "checks";
    try {
      s.compiler.parse();
      Node externsRoot = s.compiler.getExternsRoot();
      Node jsRoot = s.compiler.getJsRoot();
      if (!variant.equals("untyped")) {
        typeCheck(s.compiler, s.options);
      }
      if (variant.equals("colors")) {
        new RemoveCastNodes(s.compiler).process(externsRoot, jsRoot);
        new ConvertTypesToColors(
                s.compiler, SerializationOptions.builder().setIncludeDebugInfo(false).build())
            .process(externsRoot, jsRoot);
        s.compiler.setLifeCycleStage(LifeCycleStage.COLORS_AND_SIMPLIFIED_JSDOC);
      }
      line(sb, "checks " + s.errors.lines.size() + " " + sha16(prefixed(s.errors.lines)));

      stage = "serialize";
      ByteArrayOutputStream bytes = new ByteArrayOutputStream();
      SerializeTypedAstPass.createFromOutputStream(
              s.compiler,
              bytes,
              SerializationOptions.builder()
                  .setIncludeDebugInfo(false)
                  .setRuntimeLibraries(ImmutableList.of())
                  .build())
          .process(externsRoot, jsRoot);
      byte[] data = bytes.toByteArray();
      line(sb, "typedast " + data.length + " " + sha16Bytes(data));
      List<String> text = new ArrayList<>();
      TypedAst ast = TypedAst.List.parseFrom(data).getTypedAsts(0);
      addLines(text, ast.toBuilder().clearExternAst().clearCodeAst().build().toString());
      for (LazyAst lazy : ast.getExternAstList()) {
        text.add("extern_ast");
        addLazyAst(text, lazy);
      }
      for (LazyAst lazy : ast.getCodeAstList()) {
        text.add("code_ast");
        addLazyAst(text, lazy);
      }
      section(sb, "text", text, fullThis);

      stage = "deserialize";
      Setup d = new Setup(externsFile, inputFile);
      TypedAstDeserializer.DeserializedAst deserialized =
          TypedAstDeserializer.deserializeFullAst(
              d.compiler,
              SourceFile.fromCode("syntheticExterns", "", SourceKind.EXTERN),
              ImmutableSet.of(externsFile, inputFile),
              new ByteArrayInputStream(data),
              true,
              d.options.getResolveSourceMapAnnotations(),
              d.options.getParseInlineSourceMaps());
      List<String> tree = new ArrayList<>();
      Map<Object, Color> colors = new LinkedHashMap<>();
      for (SourceFile f : ImmutableList.of(externsFile, inputFile)) {
        var supplier = deserialized.getFilesystem().get(f);
        if (supplier == null) {
          tree.add("script " + f.getName() + " absent");
          continue;
        }
        tree.add("script " + f.getName());
        dumpTree(tree, colors, supplier.get(), 0);
      }
      for (Color c : colors.values()) {
        tree.add(colorLine(c));
      }
      section(sb, "deserialized", tree, fullThis);
    } catch (RuntimeException | StackOverflowError e) {
      if (System.getenv("SERIALIZATION_DUMP_TRACE") != null) {
        e.printStackTrace(); // debugging aid, on stderr only
      }
      line(sb, "exception " + stage + " " + e.getClass().getName() + ": " + e.getMessage());
    }
  }

  /** TypeCheckDump.dumpOne's type checking steps. */
  static void typeCheck(Compiler compiler, CompilerOptions options) {
    Node externsRoot = compiler.getExternsRoot();
    Node jsRoot = compiler.getJsRoot();
    new GatherModuleMetadata(compiler, false, ResolutionMode.BROWSER).process(externsRoot, jsRoot);
    new ModuleMapCreator(compiler, compiler.getModuleMetadataMap()).process(externsRoot, jsRoot);
    new InferConsts(compiler).process(externsRoot, jsRoot);

    // DefaultPassConfig.inferTypes
    Node globalRoot = jsRoot.getParent();
    TypeInferencePass inferencePass =
        new TypeInferencePass(
            compiler,
            compiler.getReverseAbstractInterpreter(),
            (TypedScopeCreator) compiler.getTypedScopeCreator());
    compiler.setTypeCheckingHasRun(true);
    compiler.setTopScope(inferencePass.inferAllScopes(globalRoot));

    // DefaultPassConfig.checkTypes
    TypeCheck check =
        new TypeCheck(
                compiler,
                compiler.getReverseAbstractInterpreter(),
                compiler.getTypeRegistry(),
                compiler.getTopScope(),
                (TypedScopeCreator) compiler.getTypedScopeCreator())
            .reportUnknownTypes(
                options.enables(DiagnosticGroup.forType(TypeCheck.UNKNOWN_EXPR_TYPE)))
            .reportMissingProperties(
                !options.disables(DiagnosticGroup.forType(TypeCheck.INEXISTENT_PROPERTY)));
    check.process(externsRoot, jsRoot);
    compiler.getErrorManager().setTypedPercent(check.getTypedPercent());
  }

  private static void addLazyAst(List<String> text, LazyAst lazy) throws Exception {
    addLines(text, lazy.toBuilder().clearScript().build().toString());
    // As ScriptNodeDeserializer reads the script: without protobuf's default nesting limit (100),
    // which deep ASTs exceed.
    CodedInputStream script = lazy.getScript().newCodedInput();
    script.setRecursionLimit(Integer.MAX_VALUE);
    addLines(text, AstNode.parseFrom(script).toString());
  }

  private static void addLines(List<String> out, String s) {
    for (String l : s.split("\n", -1)) {
      out.add(l);
    }
    // String.split keeps a trailing empty string for a text ending in "\n": drop it.
    if (!out.isEmpty() && out.get(out.size() - 1).isEmpty() && s.endsWith("\n")) {
      out.remove(out.size() - 1);
    }
  }

  private static void dumpTree(List<String> out, Map<Object, Color> colors, Node n, int level) {
    StringBuilder b = new StringBuilder();
    b.append("    ".repeat(level)).append(n.toString(true, true, false));
    Color c = n.getColor();
    if (c != null) {
      b.append(" : ").append(c.getId());
      collectColor(colors, c);
    }
    out.add(b.toString());
    for (Node child = n.getFirstChild(); child != null; child = child.getNext()) {
      dumpTree(out, colors, child, level + 1);
    }
  }

  private static void collectColor(Map<Object, Color> colors, Color c) {
    if (colors.containsKey(c.getId())) {
      return;
    }
    colors.put(c.getId(), c);
    for (Color x : c.getPrototypes()) {
      collectColor(colors, x);
    }
    for (Color x : c.getInstanceColors()) {
      collectColor(colors, x);
    }
    for (Color x : c.getUnionElements()) {
      collectColor(colors, x);
    }
  }

  private static String colorLine(Color c) {
    return "color "
        + c.getId()
        + " inv="
        + c.isInvalidating()
        + " ctor="
        + c.isConstructor()
        + " keep="
        + c.getPropertiesKeepOriginalName()
        + " assert="
        + c.isClosureAssert()
        + " box="
        + c.getBoxId()
        + " own="
        + c.getOwnProperties()
        + " proto="
        + ids(c.getPrototypes())
        + " inst="
        + ids(c.getInstanceColors())
        + " union="
        + ids(c.getUnionElements());
  }

  private static String ids(ImmutableSet<Color> cs) {
    List<String> out = new ArrayList<>();
    for (Color c : cs) {
      out.add(c.getId().toString());
    }
    return out.toString();
  }

  private static List<String> prefixed(List<String> diags) {
    List<String> out = new ArrayList<>();
    for (String d : diags) {
      out.add("diag " + d);
    }
    return out;
  }

  private static void section(StringBuilder sb, String name, List<String> lines, boolean fullThis)
      throws Exception {
    line(sb, name + " " + lines.size() + " " + sha16(lines));
    if (fullThis) {
      for (String l : lines) {
        line(sb, l);
      }
    }
  }

  private static String sha16(List<String> lines) throws Exception {
    MessageDigest md = MessageDigest.getInstance("SHA-256");
    for (String l : lines) {
      md.update((TypeCheckDump.escape(l) + "\n").getBytes(StandardCharsets.UTF_8));
    }
    return hex16(md.digest());
  }

  private static String sha16Bytes(byte[] data) throws Exception {
    return hex16(MessageDigest.getInstance("SHA-256").digest(data));
  }

  private static String hex16(byte[] digest) {
    StringBuilder hex = new StringBuilder();
    for (byte b : digest) {
      hex.append(String.format("%02x", b));
    }
    return hex.substring(0, 16);
  }

  private static void line(StringBuilder sb, String s) {
    sb.append(TypeCheckDump.escape(s)).append('\n');
  }
}

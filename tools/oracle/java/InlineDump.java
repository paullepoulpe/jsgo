// Differential oracle for the inlining passes and GlobalNamespace
// (internal/jscomp/inline_variables.go, flow_sensitive_inline_variables.go and the reaching-def/use
// analyses, inline_functions.go with function_injector.go / function_to_block_mutator.go /
// inline_cost_estimator.go, inline_and_collapse_properties.go with global_namespace.go and
// destructuring_global_name_extractor.go, inline_simple_methods.go).
//
// For every file, each pass starts from a fresh Compiler set up the way CompilerTestCase sets one
// up: default CompilerOptions with --language_in=ECMASCRIPT_NEXT and --language_out=NO_TRANSPILE,
// the moduleLoad group off (as NormalizeTest.getOptions sets it), the GoogleCodingConvention (CompilerTestCase's default), externs EXTERNS,
// Compiler.init(externs, [SourceFile.fromCode(<rel>, <src>)], options) and Compiler.parse(), then
// Normalize.createNormalizeForOptimizations(compiler).process(externs, js) (as
// CompilerTestCase.enableNormalize() and DefaultPassConfig do). Files with parse errors print their
// diagnostics and are skipped. Then compiler.beforePass(<pass>) (as PhaseOptimizer does; the
// indexed decisions log of InlineAndCollapseProperties needs it) and the pass runs once:
//   inline-variables-locals  new InlineVariables(compiler, Mode.LOCALS_ONLY)
//   inline-variables-all     new InlineVariables(compiler, Mode.ALL)
//   flow-sensitive           new FlowSensitiveInlineVariables(compiler)
//   inline-functions         new InlineFunctions(compiler, compiler.getUniqueNameIdSupplier(),
//                              Reach.ALL, false, false, UNLIMITED_FUN_SIZE_AFTER_INLINING)
//                              (InlineFunctionsTest's defaults)
//   collapse-all             InlineAndCollapseProperties.builder(compiler)
//                              .setPropertyCollapseLevel(ALL).setChunkOutputType(GLOBAL_NAMESPACE)
//                              .setHaveModulesBeenRewritten(false)
//                              .setModuleResolutionMode(BROWSER).build()
//                              (InlineAndCollapsePropertiesTest's processor, without
//                              Es6NormalizeClasses)
//   collapse-none            the same with PropertyCollapseLevel.NONE (InlineAliases)
//   inline-simple-methods    new InlineSimpleMethods(compiler), after
//                              compiler.setExternProperties(ImmutableSet.of())
// and, from another fresh normalized compiler, a GlobalNamespace dump:
//   new GlobalNamespace(compiler, externsRoot, jsRoot); for every Name of getNameIndex() (in
//   order): "name <Name.toString()>" then "  decisions <canCollapseOrInline()>
//   <canCollapseOrInlineChildNames()> <calculateInlinability()> needsStub=<needsToBeStubbed()>
//   nsObjLit=<isNamespaceObjectLit()> simpleStub=<isSimpleStubDeclaration()>" (each value
//   "!<exception class>" when it throws) and one "  ref ..." line per Ref (getRefs() order) with
//   its type, node, parent token, scope root and whether it is the declaration / initialization.
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside printable
// ASCII as \\uXXXX):
//   ## file <rel>
//   # parse ok|errors                then the parse diagnostics:
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   # normalize exception <class>: <message>   if Normalize threw (the file ends)
//   ## pass <name>
//   diag ...                         the diagnostics added by the pass, in ErrorManager order
//   exception <class>: <message>     if the pass threw (then the pass section ends)
//   tree unchanged                   if root.toStringTree() is the same as after normalization,
//   tree changed <sha256>            else the first 16 hex digits of the SHA-256 of the escaped
//                                    lines of root.toStringTree() (each followed by '\n'), and
//                                    with --full the escaped lines themselves
//   changes <n> <sha256>             the change-scope roots (SCRIPT and FUNCTION nodes) whose
//                                    change time is non-zero, as "<token>@<line>:<charno> <time>"
//                                    lines: count and hash, and with --full the lines
//   ## namespace
//   exception <class>: <message>     if building the namespace threw
//   namespace <n> <sha256>           the names' lines (count of names, hash of all lines), and
//                                    with --full the lines themselves
//
// Usage: java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//          com.google.javascript.jscomp.InlineDump [--full] [--sample=<N>] <out.gz> <root> \
//          <rel-file>... | --list=<listfile>
// -XX:-ShowCodeDetailsInExceptionMessages turns off the JVM's "helpful" NullPointerException
// messages (which describe Java bytecode), so that an NPE has a null message, as a Go nil
// dereference does.
// --sample=N dumps every Nth file of the list only (in list order, starting with the first).
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.CompilerOptions.ChunkOutputType;
import com.google.javascript.jscomp.CompilerOptions.LanguageMode;
import com.google.javascript.jscomp.CompilerOptions.PropertyCollapseLevel;
import com.google.javascript.jscomp.CompilerOptions.Reach;
import com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSet;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Callable;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class InlineDump {

  private static boolean full = false;
  private static int sample = 1;

  static final String EXTERNS =
      "var window; var Arguments; var arguments; /** @constructor */ function String() {}\n"
          + "/** @nosideeffects */ function nochg(){}\nfunction chg(){}\n";

  static final String[] PASSES = {
    "inline-variables-locals",
    "inline-variables-all",
    "flow-sensitive",
    "inline-functions",
    "collapse-all",
    "collapse-none",
    "inline-simple-methods",
  };

  public static void main(String[] args) throws Exception {
    List<String> rest = new ArrayList<>();
    List<String> files = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--full")) {
        full = true;
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
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(out))) {
      int index = 0;
      for (String rel : files) {
        if (index++ % sample != 0) {
          continue;
        }
        StringBuilder sb = new StringBuilder();
        String src = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
        dumpFile(sb, rel, src);
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  private static Compiler newCompiler(String rel, String src) {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(LanguageMode.ECMASCRIPT_NEXT);
    options.setLanguageOut(LanguageMode.NO_TRANSPILE);
    // As NormalizeTest.getOptions does: imports of files outside the input are not errors.
    options.setWarningLevel(DiagnosticGroups.MODULE_LOAD, CheckLevel.OFF);
    options.setCodingConvention(new GoogleCodingConvention());
    compiler.init(
        ImmutableList.of(SourceFile.fromCode("externs.js", EXTERNS)),
        ImmutableList.of(SourceFile.fromCode(rel, src)),
        options);
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

  private static void collectChanges(Node n, List<String> out) {
    if ((n.isScript() || n.isFunction()) && n.getChangeTime() != 0) {
      out.add(n.getToken() + "@" + n.getLineno() + ":" + n.getCharno() + " " + n.getChangeTime());
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      collectChanges(c, out);
    }
  }

  private static String sha16(List<String> lines) throws Exception {
    MessageDigest md = MessageDigest.getInstance("SHA-256");
    for (String l : lines) {
      md.update((l + "\n").getBytes(StandardCharsets.UTF_8));
    }
    StringBuilder hex = new StringBuilder();
    for (byte b : md.digest()) {
      hex.append(String.format("%02x", b));
    }
    return hex.substring(0, 16);
  }

  private static String exception(Throwable t) {
    return t.getClass().getName() + ": " + t.getMessage();
  }

  private static CompilerPass newPass(String name, Compiler c) {
    switch (name) {
      case "inline-variables-locals":
        return new InlineVariables(c, InlineVariables.Mode.LOCALS_ONLY);
      case "inline-variables-all":
        return new InlineVariables(c, InlineVariables.Mode.ALL);
      case "flow-sensitive":
        return new FlowSensitiveInlineVariables(c);
      case "inline-functions":
        return new InlineFunctions(
            c,
            c.getUniqueNameIdSupplier(),
            Reach.ALL,
            /* assumeStrictThis= */ false,
            CompilerOptions.UNLIMITED_FUN_SIZE_AFTER_INLINING);
      case "collapse-all":
        return InlineAndCollapseProperties.builder(c)
            .setPropertyCollapseLevel(PropertyCollapseLevel.ALL)
            .setChunkOutputType(ChunkOutputType.GLOBAL_NAMESPACE)
            .setHaveModulesBeenRewritten(false)
            .setModuleResolutionMode(ResolutionMode.BROWSER)
            .build();
      case "collapse-none":
        return InlineAndCollapseProperties.builder(c)
            .setPropertyCollapseLevel(PropertyCollapseLevel.NONE)
            .setChunkOutputType(ChunkOutputType.GLOBAL_NAMESPACE)
            .setHaveModulesBeenRewritten(false)
            .setModuleResolutionMode(ResolutionMode.BROWSER)
            .build();
      case "inline-simple-methods":
        c.setExternProperties(ImmutableSet.of());
        return new InlineSimpleMethods(c);
      default:
        throw new IllegalArgumentException(name);
    }
  }

  /** A fresh compiler, parsed and normalized; null (after printing why) if normalize threw. */
  private static Compiler normalized(StringBuilder sb, String rel, String src) {
    Compiler c = newCompiler(rel, src);
    try {
      Normalize.createNormalizeForOptimizations(c).process(c.getExternsRoot(), c.getJsRoot());
    } catch (Throwable t) {
      line(sb, "# normalize exception " + exception(t));
      return null;
    }
    return c;
  }

  private static void dumpFile(StringBuilder sb, String rel, String src) throws Exception {
    line(sb, "## file " + rel);
    Compiler compiler = newCompiler(rel, src);
    List<String> parseDiags = diagnostics(compiler);
    line(sb, "# parse " + (compiler.hasErrors() ? "errors" : "ok"));
    for (String d : parseDiags) {
      line(sb, d);
    }
    if (compiler.hasErrors()) {
      return;
    }
    for (String pass : PASSES) {
      StringBuilder head = new StringBuilder();
      Compiler c = normalized(head, rel, src);
      if (c == null) {
        sb.append(head);
        return;
      }
      line(sb, "## pass " + pass);
      List<String> baseDiags = diagnostics(c);
      String before = c.getRoot().toStringTree();
      String exception = null;
      try {
        c.beforePass(pass);
        newPass(pass, c).process(c.getExternsRoot(), c.getJsRoot());
      } catch (Throwable t) {
        exception = "exception " + exception(t);
      }
      List<String> diags = diagnostics(c);
      for (String d : baseDiags) {
        diags.remove(d);
      }
      for (String d : diags) {
        line(sb, d);
      }
      if (exception != null) {
        line(sb, exception);
        continue;
      }
      String after = c.getRoot().toStringTree();
      if (after.equals(before)) {
        line(sb, "tree unchanged");
      } else {
        List<String> lines = new ArrayList<>();
        for (String l : after.split("\n", -1)) {
          lines.add(escape(l));
        }
        if (lines.get(lines.size() - 1).isEmpty()) {
          lines.remove(lines.size() - 1);
        }
        line(sb, "tree changed " + sha16(lines));
        if (full) {
          for (String l : lines) {
            sb.append(l).append('\n');
          }
        }
      }
      List<String> changes = new ArrayList<>();
      collectChanges(c.getRoot(), changes);
      List<String> escapedChanges = new ArrayList<>();
      for (String l : changes) {
        escapedChanges.add(escape(l));
      }
      line(sb, "changes " + changes.size() + " " + sha16(escapedChanges));
      if (full) {
        for (String l : escapedChanges) {
          sb.append(l).append('\n');
        }
      }
    }
    Compiler c = normalized(sb, rel, src);
    if (c == null) {
      return;
    }
    line(sb, "## namespace");
    List<String> lines = new ArrayList<>();
    int count = 0;
    try {
      GlobalNamespace namespace = new GlobalNamespace(c, c.getExternsRoot(), c.getJsRoot());
      for (GlobalNamespace.Name name : namespace.getNameIndex().values()) {
        count++;
        dumpName(lines, name);
      }
    } catch (Throwable t) {
      line(sb, "exception " + exception(t));
      return;
    }
    line(sb, "namespace " + count + " " + sha16(lines));
    if (full) {
      for (String l : lines) {
        sb.append(l).append('\n');
      }
    }
  }

  private static String attempt(Callable<Object> f) {
    try {
      return String.valueOf(f.call());
    } catch (Throwable t) {
      return "!" + t.getClass().getName();
    }
  }

  private static void dumpName(List<String> lines, GlobalNamespace.Name name) {
    lines.add(escape("name " + name));
    lines.add(
        escape(
            "  decisions "
                + attempt(name::canCollapseOrInline)
                + " "
                + attempt(name::canCollapseOrInlineChildNames)
                + " "
                + attempt(name::calculateInlinability)
                + " needsStub="
                + attempt(name::needsToBeStubbed)
                + " nsObjLit="
                + attempt(name::isNamespaceObjectLit)
                + " simpleStub="
                + attempt(name::isSimpleStubDeclaration)));
    for (GlobalNamespace.Ref ref : name.getRefs()) {
      Node n = ref.getNode();
      String parent = n == null || n.getParent() == null ? "null" : n.getParent().getToken().toString();
      String scope =
          ref.scope == null ? "null" : node(ref.scope.getRootNode());
      lines.add(
          escape(
              "  ref "
                  + refType(ref)
                  + " "
                  + node(n)
                  + " parent="
                  + parent
                  + " scope="
                  + scope
                  + (ref == name.getDeclaration() ? " decl" : "")
                  + (ref == name.getInitialization() ? " init" : "")));
    }
  }

  private static String refType(GlobalNamespace.Ref ref) {
    // Ref.type is private; Ref.toString() starts with "Ref{type=<TYPE>".
    String s = ref.toString();
    int start = s.indexOf("type=") + 5;
    int end = start;
    while (end < s.length() && (Character.isLetter(s.charAt(end)) || s.charAt(end) == '_')) {
      end++;
    }
    return s.substring(start, end);
  }

  private static String node(Node n) {
    if (n == null) {
      return "null";
    }
    String s = n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
    if (n.isQualifiedName()) {
      s += " " + n.getQualifiedName();
    } else if (n.isStringKey()
        || n.isMemberFunctionDef()
        || n.isMemberFieldDef()
        || n.isGetterDef()
        || n.isSetterDef()) {
      s += " " + n.getString();
    }
    return s;
  }

  private static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }

  static String escape(String s) {
    StringBuilder sb = new StringBuilder(s.length());
    for (int i = 0; i < s.length(); i++) {
      char ch = s.charAt(i);
      if (ch == '\\') {
        sb.append("\\\\");
      } else if (ch >= 0x20 && ch < 0x7f) {
        sb.append(ch);
      } else {
        sb.append(String.format("\\u%04x", (int) ch));
      }
    }
    return sb.toString();
  }
}

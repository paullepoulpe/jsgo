// Differential oracle for ReplaceMessages (internal/jscomp/replace_messages.go) with a
// deterministic message bundle built from each input, in the style of ReplaceMessagesTest's
// SimpleMessageBundle (messages keyed by id, no IdGenerator, so a message's id is its key).
//
// For every file: a fresh Compiler with TypeCheckDump.options() (CompilerTypeTestCase's options),
// the externs <file> as "externs.js" and the input as one SourceFile named <rel>: Compiler.init,
// Compiler.parse. Files with parse errors print their diagnostics and are skipped. A collecting
// JsMessageVisitor (no IdGenerator) then records every goog.getMsg() and declareIcuTemplate()
// message in visiting order, and the bundle holds, for the i-th distinct message id (i = 0, 1,
// ...) except every 4th one (i % 4 == 3, left out to exercise missing translations), a
// "translation": a JsMessage with the same key and id whose parts are the original parts in
// reverse order, every string part prefixed with "~". Each variant then starts from a fresh
// Compiler (same setup), type checks exactly as SerializationDump.typeCheck, runs RemoveCastNodes
// and ConvertTypesToColors(includeDebugInfo false) + setLifeCycleStage(COLORS_AND_SIMPLIFIED_JSDOC)
// (ReplaceMessagesTest enables type checking with colors), and runs
//   full         new ReplaceMessages(compiler, bundle, false).getFullReplacementPass()
//   full-strict  the same with strictReplacement = true
//   protect      getMsgProtectionPass()
//   complete     getMsgProtectionPass(), then, on the same AST, a new ReplaceMessages'
//                getReplacementCompletionPass()
//
// Output (gzipped; every line escaped like ParsingDump: backslash, and every char outside printable
// ASCII as \\uXXXX):
//   ## file <rel>
//   # parse ok|errors                 then the parse diagnostics:
//   diag ERROR|WARNING <key> <line>:<charno> "<description>"
//   bundle <n> <sha16>                the bundle: per message "<id>" then one line per part
//                                     ("  s <string>" or "  p <canonical placeholder name>")
//   ## variant <name>
//   checks <n> <sha16>                the diagnostics reported by the prerequisites
//   exception <stage> <class>: <message>   if a step threw (the variant ends); stage is checks or pass
//   diags <n> <sha16>                 the diagnostics the pass(es) reported
//   tree <n> <sha16>                  jsRoot afterwards: one line per node, 2 spaces per depth +
//                                     Node.toString(true, true, false)
//   externs <n> <sha16>               the same for externsRoot
//   ... with --full (or for every --full-every'th file) the lines of each section follow it.
// sha16 is the first 16 hex digits of the SHA-256 of the escaped lines (each followed by '\n').
//
// Usage: java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp closure-compiler.jar:<build> \
//          com.google.javascript.jscomp.ReplaceMessagesDump --externs=<file> [--full]
//          [--sample=<N>] [--full-every=<M>] <out.gz> <root> <rel-file>... | --list=<listfile>
// (compile together with TypeCheckDump.java and SerializationDump.java). The Go side is
// internal/jscomp/replace_messages_oracle_test.go.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.AbstractCompiler.LifeCycleStage;
import com.google.javascript.jscomp.serialization.ConvertTypesToColors;
import com.google.javascript.jscomp.serialization.SerializationOptions;
import com.google.javascript.rhino.Node;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
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

public final class ReplaceMessagesDump {

  private static boolean full = false;
  private static int sample = 1;
  private static int fullEvery = 0;

  static final String[] VARIANTS = {"full", "full-strict", "protect", "complete"};

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

    Setup(SourceFile externsFile, SourceFile inputFile) {
      compiler.setErrorManager(errors);
      compiler.init(ImmutableList.of(externsFile), ImmutableList.of(inputFile), options);
    }
  }

  /** Records every message the visitor finds, in visiting order. */
  private static final class Collector extends JsMessageVisitor {
    final List<JsMessage> found = new ArrayList<>();

    Collector(AbstractCompiler compiler) {
      super(compiler, null);
    }

    @Override
    protected void processJsMessageDefinition(JsMessageDefinition definition) {
      found.add(definition.getMessage());
    }

    @Override
    protected void processIcuTemplateDefinition(IcuTemplateDefinition definition) {
      found.add(definition.getMessage());
    }
  }

  private static final class Bundle implements MessageBundle {
    final Map<String, JsMessage> messages = new LinkedHashMap<>();

    @Override
    public JsMessage getMessage(String id) {
      return messages.get(id);
    }

    @Override
    public Iterable<JsMessage> getAllMessages() {
      return messages.values();
    }

    @Override
    public JsMessage.IdGenerator idGenerator() {
      return null;
    }
  }

  static Bundle buildBundle(List<JsMessage> found) {
    Map<String, JsMessage> byId = new LinkedHashMap<>();
    for (JsMessage m : found) {
      byId.putIfAbsent(m.getId(), m);
    }
    Bundle bundle = new Bundle();
    int i = 0;
    for (JsMessage m : byId.values()) {
      if (i++ % 4 == 3) {
        continue;
      }
      JsMessage.Builder b = new JsMessage.Builder().setKey(m.getKey()).setId(m.getId());
      List<JsMessage.Part> parts = new ArrayList<>(m.getParts());
      for (int j = parts.size() - 1; j >= 0; j--) {
        JsMessage.Part p = parts.get(j);
        if (p.isPlaceholder()) {
          b.appendPart(p);
        } else {
          b.appendStringPart("~" + p.getString());
        }
      }
      bundle.messages.put(m.getId(), b.build());
    }
    return bundle;
  }

  static void dumpFile(StringBuilder sb, String externs, String rel, String src, boolean fullThis)
      throws Exception {
    line(sb, "## file " + rel);
    SourceFile externsFile = SourceFile.fromCode("externs.js", externs);
    SourceFile inputFile = SourceFile.fromCode(rel, src);
    Bundle bundle;
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
      try {
        Collector collector = new Collector(s.compiler);
        collector.process(s.compiler.getExternsRoot(), s.compiler.getJsRoot());
        bundle = buildBundle(collector.found);
      } catch (RuntimeException | StackOverflowError e) {
        line(sb, "exception bundle " + e.getClass().getName() + ": " + e.getMessage());
        return;
      }
    }
    List<String> bundleLines = new ArrayList<>();
    for (JsMessage m : bundle.messages.values()) {
      bundleLines.add(m.getId());
      for (JsMessage.Part p : m.getParts()) {
        bundleLines.add(
            p.isPlaceholder() ? "  p " + p.getCanonicalPlaceholderName() : "  s " + p.getString());
      }
    }
    section(sb, "bundle", bundleLines, fullThis);
    for (String variant : VARIANTS) {
      line(sb, "## variant " + variant);
      dumpVariant(sb, variant, externsFile, inputFile, bundle, fullThis);
    }
  }

  private static void dumpVariant(
      StringBuilder sb,
      String variant,
      SourceFile externsFile,
      SourceFile inputFile,
      Bundle bundle,
      boolean fullThis)
      throws Exception {
    Setup s = new Setup(externsFile, inputFile);
    String stage = "checks";
    try {
      s.compiler.parse();
      Node externsRoot = s.compiler.getExternsRoot();
      Node jsRoot = s.compiler.getJsRoot();
      SerializationDump.typeCheck(s.compiler, s.options);
      new RemoveCastNodes(s.compiler).process(externsRoot, jsRoot);
      new ConvertTypesToColors(
              s.compiler, SerializationOptions.builder().setIncludeDebugInfo(false).build())
          .process(externsRoot, jsRoot);
      s.compiler.setLifeCycleStage(LifeCycleStage.COLORS_AND_SIMPLIFIED_JSDOC);
      int before = s.errors.lines.size();
      line(sb, "checks " + before + " " + sha16(prefixed(s.errors.lines)));

      stage = "pass";
      switch (variant) {
        case "full" ->
            new ReplaceMessages(s.compiler, bundle, false)
                .getFullReplacementPass()
                .process(externsRoot, jsRoot);
        case "full-strict" ->
            new ReplaceMessages(s.compiler, bundle, true)
                .getFullReplacementPass()
                .process(externsRoot, jsRoot);
        case "protect" ->
            new ReplaceMessages(s.compiler, bundle, false)
                .getMsgProtectionPass()
                .process(externsRoot, jsRoot);
        case "complete" -> {
          new ReplaceMessages(s.compiler, bundle, false)
              .getMsgProtectionPass()
              .process(externsRoot, jsRoot);
          new ReplaceMessages(s.compiler, bundle, false)
              .getReplacementCompletionPass()
              .process(externsRoot, jsRoot);
        }
        default -> throw new IllegalArgumentException(variant);
      }
      List<String> diags = prefixed(s.errors.lines.subList(before, s.errors.lines.size()));
      section(sb, "diags", diags, fullThis);
      List<String> tree = new ArrayList<>();
      dumpTree(tree, jsRoot, 0);
      section(sb, "tree", tree, fullThis);
      List<String> ext = new ArrayList<>();
      dumpTree(ext, externsRoot, 0);
      section(sb, "externs", ext, fullThis);
    } catch (RuntimeException | StackOverflowError e) {
      if (System.getenv("REPLACE_MESSAGES_DUMP_TRACE") != null) {
        e.printStackTrace(); // debugging aid, on stderr only
      }
      line(sb, "exception " + stage + " " + e.getClass().getName() + ": " + e.getMessage());
    }
  }

  private static void dumpTree(List<String> out, Node n, int level) {
    out.add("  ".repeat(level) + n.toString(true, true, false));
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      dumpTree(out, c, level + 1);
    }
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
    StringBuilder hex = new StringBuilder();
    for (byte b : md.digest()) {
      hex.append(String.format("%02x", b));
    }
    return hex.substring(0, 16);
  }

  private static void line(StringBuilder sb, String s) {
    sb.append(TypeCheckDump.escape(s)).append('\n');
  }
}

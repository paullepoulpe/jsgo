// Oracle tool for the l11-messages-types-misc ports (JsMessageVisitor, XtbMessageBundle,
// TypeTransformation, AngularPass, ChromePass). It lives in package com.google.javascript.jscomp to
// reach package-private API. Output lines are escaped to printable ASCII as in ParsingDump.java.
//
//   javac -cp closure-compiler.jar -d build MessagesTypesMiscDump.java
//   java -cp closure-compiler.jar:build com.google.javascript.jscomp.MessagesTypesMiscDump \
//       <mode> <out.gz|-> <args...>
//
// Modes:
//   msgs <root> <rel-file>...|@<list>   For each file: new Compiler(), parseTestCode(code), run a
//       JsMessageVisitor (GoogleJsMessageIdGenerator with a null project id) that records every
//       goog.getMsg()/declareIcuTemplate() definition and getMsgWithFallback() call, then print the
//       messages and the errors and warnings reported after parsing. Files with nothing to report
//       print nothing.
//   xtb <dir>   For each *.xtb file (sorted): new XtbMessageBundle(stream, "TestProject"), print
//       every message (getAllMessages() order) or the exception.
//   ttl <expressions-file>   Each non-empty line is a TTL expression; evaluated exactly as
//       TypeTransformationTest.testTTL does (same type variables, name variables and extra type
//       definitions); prints the parser errors, the warnings and the resulting type.
//   angular|chrome <root> <rel-file>...|@<list>   For each file: new Compiler(), parseTestCode(code),
//       run AngularPass / ChromePass with process(null, script); prints the pass's diagnostics and,
//       the AST changes (see diffTree). Files with nothing to report print nothing.
// The Go side is internal/jscomp/messages_types_misc_oracle_test.go.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableMap;
import com.google.javascript.jscomp.parsing.TypeTransformationParser;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.jstype.JSType;
import com.google.javascript.rhino.jstype.JSTypeNative;
import com.google.javascript.rhino.jstype.JSTypeRegistry;
import com.google.javascript.rhino.jstype.JSTypeResolver;
import com.google.javascript.rhino.jstype.ObjectType;
import com.google.javascript.rhino.jstype.RecordTypeBuilder;
import java.io.ByteArrayInputStream;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Stream;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class MessagesTypesMiscDump {
  static String escape(String s) {
    if (s == null) {
      return "<null>";
    }
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

  public static void main(String[] args) throws Exception {
    String mode = args[0];
    OutputStream os =
        args[1].equals("-")
            ? System.out
            : new GZIPOutputStream(Files.newOutputStream(Paths.get(args[1])));
    StringBuilder sb = new StringBuilder();
    switch (mode) {
      case "msgs", "angular", "chrome" -> {
        Path root = Paths.get(args[2]);
        for (String rel : files(args, 3)) {
          String code = Files.readString(root.resolve(rel), StandardCharsets.UTF_8);
          StringBuilder one = new StringBuilder();
          if (mode.equals("msgs")) {
            dumpMessages(one, code);
          } else {
            dumpPass(one, mode, code);
          }
          if (one.length() > 0) {
            sb.append("== ").append(escape(rel)).append('\n').append(one);
          }
          if (sb.length() > 1 << 20) {
            os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
            sb.setLength(0);
          }
        }
      }
      case "xtb" -> {
        List<Path> xtbs = new ArrayList<>();
        try (Stream<Path> s = Files.list(Paths.get(args[2]))) {
          s.filter(p -> p.toString().endsWith(".xtb")).sorted().forEach(xtbs::add);
        }
        for (Path p : xtbs) {
          sb.append("== ").append(p.getFileName()).append('\n');
          dumpXtb(sb, Files.readAllBytes(p));
        }
      }
      case "ttl" -> {
        for (String line : Files.readAllLines(Paths.get(args[2]), StandardCharsets.UTF_8)) {
          if (!line.isEmpty()) {
            sb.append("== ").append(escape(line)).append('\n');
            dumpTtl(sb, line);
          }
        }
      }
      default -> throw new IllegalArgumentException(mode);
    }
    os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
    os.close();
  }

  static List<String> files(String[] args, int from) throws Exception {
    List<String> files = new ArrayList<>();
    for (int i = from; i < args.length; i++) {
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
    return files;
  }

  static Compiler newCompiler() {
    return new Compiler(new PrintStream(OutputStream.nullOutputStream()));
  }

  static void dumpDiagnostics(StringBuilder sb, Compiler compiler) {
    dumpDiagnostics(sb, compiler, 0, 0);
  }

  // Prints the errors and warnings reported after the first fromErrors / fromWarnings ones (the
  // corpus modes skip the parser's diagnostics, which belong to the parser oracles).
  static void dumpDiagnostics(StringBuilder sb, Compiler compiler, int fromErrors, int fromWarnings) {
    List<JSError> errors = compiler.getErrors();
    for (JSError e : errors.subList(fromErrors, errors.size())) {
      sb.append("error ").append(escape(e.toString())).append('\n');
    }
    List<JSError> warnings = compiler.getWarnings();
    for (JSError e : warnings.subList(fromWarnings, warnings.size())) {
      sb.append("warning ").append(escape(e.toString())).append('\n');
    }
  }

  // ---------------------------------------------------------------------------------------------
  // msgs

  static String parts(List<JsMessage.Part> parts) {
    StringBuilder b = new StringBuilder("[");
    for (JsMessage.Part p : parts) {
      if (b.length() > 1) {
        b.append(", ");
      }
      if (p.isPlaceholder()) {
        b.append("PH(").append(p.getJsPlaceholderName()).append('|')
            .append(p.getCanonicalPlaceholderName()).append(')');
      } else {
        b.append('"').append(p.getString()).append('"');
      }
    }
    return b.append(']').toString();
  }

  static String map(Map<String, String> m) {
    StringBuilder b = new StringBuilder("{");
    for (Map.Entry<String, String> e : m.entrySet()) {
      if (b.length() > 1) {
        b.append(", ");
      }
      b.append(e.getKey()).append('=').append(e.getValue());
    }
    return b.append('}').toString();
  }

  static void message(StringBuilder sb, JsMessage m) {
    sb.append("  key=").append(escape(m.getKey())).append('\n');
    sb.append("  id=").append(escape(m.getId())).append('\n');
    sb.append("  anonymous=").append(m.isAnonymous()).append(" external=").append(m.isExternal())
        .append(" empty=").append(m.isEmpty()).append('\n');
    sb.append("  source=").append(escape(m.getSourceName())).append('\n');
    sb.append("  desc=").append(escape(m.getDesc())).append('\n');
    sb.append("  meaning=").append(escape(m.getMeaning())).append('\n');
    sb.append("  alternateId=").append(escape(m.getAlternateId())).append('\n');
    sb.append("  parts=").append(escape(parts(m.getParts()))).append('\n');
    sb.append("  jsPlaceholders=").append(escape(m.jsPlaceholderNames().toString())).append('\n');
    sb.append("  examples=").append(escape(map(m.getPlaceholderNameToExampleMap()))).append('\n');
    sb.append("  originalCode=").append(escape(map(m.getPlaceholderNameToOriginalCodeMap())))
        .append('\n');
  }

  static String node(Node n) {
    if (n == null) {
      return "<null>";
    }
    return n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
  }

  static void dumpMessages(StringBuilder sb, String code) {
    Compiler compiler = newCompiler();
    Node root = compiler.parseTestCode(code);
    int parseErrors = compiler.getErrors().size();
    int parseWarnings = compiler.getWarnings().size();
    StringBuilder found = new StringBuilder();
    JsMessageVisitor visitor =
        new JsMessageVisitor(compiler, new GoogleJsMessageIdGenerator(null)) {
          @Override
          protected void processJsMessageDefinition(JsMessageDefinition definition) {
            found.append("getMsg ").append(node(definition.getMessageNode())).append(" text=")
                .append(node(definition.getTemplateTextNode())).append(" values=")
                .append(node(definition.getPlaceholderValuesNode())).append(" valueMap=")
                .append(escape(definition.getPlaceholderValueMap().keySet().toString()))
                .append(" html=").append(definition.shouldEscapeLessThan())
                .append(" unescape=").append(definition.shouldUnescapeHtmlEntities()).append('\n');
            message(found, definition.getMessage());
          }

          @Override
          protected void processIcuTemplateDefinition(IcuTemplateDefinition definition) {
            found.append("icu ").append(node(definition.getMessageNode())).append(" text=")
                .append(node(definition.getTemplateTextNode())).append('\n');
            message(found, definition.getMessage());
          }

          @Override
          void processMessageFallback(Node callNode, JsMessage message1, JsMessage message2) {
            found.append("fallback ").append(node(callNode)).append(' ')
                .append(escape(message1.getKey())).append(' ').append(escape(message2.getKey()))
                .append('\n');
          }
        };
    visitor.process(null, root);
    sb.append(found);
    dumpDiagnostics(sb, compiler, parseErrors, parseWarnings);
  }

  // ---------------------------------------------------------------------------------------------
  // xtb

  static void dumpXtb(StringBuilder sb, byte[] xtb) {
    XtbMessageBundle bundle;
    try {
      bundle = new XtbMessageBundle(new ByteArrayInputStream(xtb), "TestProject");
    } catch (RuntimeException e) {
      sb.append("exception ").append(e.getClass().getName());
      if (e.getClass() == RuntimeException.class && e.getCause() != null) {
        sb.append(" cause ").append(e.getCause().getClass().getName());
      } else {
        sb.append(' ').append(escape(e.getMessage()));
      }
      sb.append('\n');
      return;
    }
    for (JsMessage m : bundle.getAllMessages()) {
      sb.append("message ").append(escape(m.getId())).append('\n');
      message(sb, m);
      for (JsMessage.GrammaticalGenderCase g : m.getGenderedMessageVariants()) {
        sb.append("  gender ").append(g).append(' ')
            .append(escape(parts(m.getGenderedMessageParts(g)))).append('\n');
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // ttl: the environment of TypeTransformationTest

  static final String EXTRA_TYPE_DEFS =
      """
      /** @type {number} */ Array.prototype.length;
      /** @typedef {!Array<?>} */ var ArrayAlias;

      /** @constructor */
      function Bar() {}

      /** @type {number} */
      var n = 10;
      """;

  static JSType record(JSTypeRegistry registry, Object... propsAndTypes) {
    RecordTypeBuilder builder = new RecordTypeBuilder(registry);
    for (int i = 0; i < propsAndTypes.length; i += 2) {
      builder.addProperty((String) propsAndTypes[i], (JSType) propsAndTypes[i + 1], null);
    }
    return builder.build();
  }

  static void dumpTtl(StringBuilder sb, String ttlExp) {
    Compiler compiler = newCompiler();
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    options.setCheckTypes(true);
    compiler.initOptions(options);
    JSTypeRegistry registry = compiler.getTypeRegistry();
    StringBuilder reported = new StringBuilder();
    ErrorReporter errorReporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int line, int lineOffset) {
            reported.append("parser warning ").append(escape(message)).append('\n');
          }

          @Override
          public void error(String message, String sourceName, int line, int lineOffset) {
            reported.append("parser error ").append(escape(message)).append('\n');
          }
        };
    try (JSTypeResolver.Closer closer = registry.getResolver().openForDefinition()) {
      JSType str = registry.getNativeType(JSTypeNative.STRING_TYPE);
      JSType num = registry.getNativeType(JSTypeNative.NUMBER_TYPE);
      JSType bool = registry.getNativeType(JSTypeNative.BOOLEAN_TYPE);
      ObjectType arr = registry.getNativeObjectType(JSTypeNative.ARRAY_TYPE);
      ImmutableMap<String, JSType> typeVars =
          new ImmutableMap.Builder<String, JSType>()
              .put("S", str)
              .put("N", num)
              .put("B", bool)
              .put("BOT", registry.getNativeType(JSTypeNative.NO_TYPE))
              .put("TOP", registry.getNativeType(JSTypeNative.ALL_TYPE))
              .put("UNK", registry.getNativeType(JSTypeNative.UNKNOWN_TYPE))
              .put("CHKUNK", registry.getNativeType(JSTypeNative.CHECKED_UNKNOWN_TYPE))
              .put("SO", registry.getNativeType(JSTypeNative.STRING_OBJECT_TYPE))
              .put("NO", registry.getNativeType(JSTypeNative.NUMBER_OBJECT_TYPE))
              .put("BO", registry.getNativeType(JSTypeNative.BOOLEAN_OBJECT_TYPE))
              .put("NULL", registry.getNativeType(JSTypeNative.NULL_TYPE))
              .put("OBJ", registry.getNativeType(JSTypeNative.OBJECT_TYPE))
              .put("UNDEF", registry.getNativeType(JSTypeNative.VOID_TYPE))
              .put("ARR", arr)
              .put("ARRNUM", registry.createTemplatizedType(arr, num))
              .put("REC", record(registry, "n", num, "s", str, "b", bool))
              .put("NESTEDREC", record(registry, "s", str, "r", record(registry, "s", str, "b", bool)))
              .put("ASYNCH", record(registry, "service", registry.createTemplatizedType(arr, num)))
              .buildOrThrow();
      ImmutableMap<String, String> nameVars =
          new ImmutableMap.Builder<String, String>()
              .put("s", "string")
              .put("n", "number")
              .put("b", "boolean")
              .buildOrThrow();
      TypeTransformationParser ttlParser =
          new TypeTransformationParser(
              ttlExp, SourceFile.fromCode("[testcode]", ttlExp), errorReporter, 0, 0);
      if (ttlParser.parseTypeTransformation()) {
        Node ast = ttlParser.getTypeTransformationAst();
        Node extraTypeDefs = compiler.parseTestCode(EXTRA_TYPE_DEFS);
        TypedScope scope =
            new TypedScopeCreator(compiler)
                .createScope(IR.root(IR.root(), IR.root(extraTypeDefs)), null);
        TypeTransformation typeTransformation = new TypeTransformation(compiler, scope);
        // printType() writes to System.out: capture it.
        PrintStream savedOut = System.out;
        java.io.ByteArrayOutputStream captured = new java.io.ByteArrayOutputStream();
        System.setOut(new PrintStream(captured, true, StandardCharsets.UTF_8));
        try {
          JSType resultType = typeTransformation.eval(ast, typeVars, nameVars);
          reported.append("result ").append(escape(resultType.toString())).append('\n');
        } catch (RuntimeException e) {
          reported.append("exception ").append(e.getClass().getName()).append(' ')
              .append(escape(e.getMessage())).append('\n');
        } finally {
          System.setOut(savedOut);
        }
        if (captured.size() > 0) {
          reported.append("stdout ").append(escape(captured.toString(StandardCharsets.UTF_8)))
              .append('\n');
        }
      } else {
        reported.append("unparsed\n");
      }
    }
    sb.append(reported);
    dumpDiagnostics(sb, compiler);
  }

  // ---------------------------------------------------------------------------------------------
  // angular / chrome

  // Prints the pass's diagnostics and the AST changes, found by comparing a clone of the parsed
  // script with the script after the pass (see diffTree).
  static void dumpPass(StringBuilder sb, String mode, String code) {
    Compiler compiler = newCompiler();
    Node root = compiler.parseTestCode(code);
    int parseErrors = compiler.getErrors().size();
    int parseWarnings = compiler.getWarnings().size();
    Node before = root.cloneTree();
    CompilerPass pass = mode.equals("angular") ? new AngularPass(compiler) : new ChromePass(compiler);
    try {
      pass.process(null, root);
    } catch (RuntimeException e) {
      // Only the first message line: the rest of an internal error's message is the JVM's text for
      // its cause.
      String message = e.getMessage();
      sb.append("exception ").append(e.getClass().getName()).append(' ')
          .append(escape(message == null ? null : message.split("\n", -1)[0]));
      if (e.getCause() != null) {
        sb.append(" cause ").append(e.getCause().getClass().getName());
      }
      sb.append('\n');
    }
    dumpDiagnostics(sb, compiler, parseErrors, parseWarnings);
    diffTree(sb, before, root, "0");
  }

  // Compares the node before the pass with the node after it. A changed toString() prints
  // "changed <path> <after.toString()>". Children are compared pairwise when their number is
  // unchanged; otherwise the after children whose toStringTree() is not among the before children
  // (matched as a multiset, in order) print "inserted <path>/<index>" and their tree, followed by
  // "unmatched-before <path> <count of before children left unmatched>".
  static void diffTree(StringBuilder sb, Node before, Node after, String path) {
    String afterLine = after.toString();
    if (!before.toString().equals(afterLine)) {
      sb.append("changed ").append(path).append(' ').append(escape(afterLine)).append('\n');
    }
    if (before.getChildCount() == after.getChildCount()) {
      int i = 0;
      for (Node b = before.getFirstChild(), a = after.getFirstChild();
          b != null;
          b = b.getNext(), a = a.getNext(), i++) {
        diffTree(sb, b, a, path + "/" + i);
      }
      return;
    }
    Map<String, Integer> unmatched = new java.util.HashMap<>();
    for (Node b = before.getFirstChild(); b != null; b = b.getNext()) {
      unmatched.merge(b.toStringTree(), 1, Integer::sum);
    }
    int left = before.getChildCount();
    int i = 0;
    for (Node a = after.getFirstChild(); a != null; a = a.getNext(), i++) {
      String tree = a.toStringTree();
      Integer count = unmatched.get(tree);
      if (count != null && count > 0) {
        unmatched.put(tree, count - 1);
        left--;
        continue;
      }
      sb.append("inserted ").append(path).append('/').append(i).append('\n');
      for (String line : tree.split("\n", -1)) {
        sb.append("| ").append(escape(line)).append('\n');
      }
    }
    sb.append("unmatched-before ").append(path).append(' ').append(left).append('\n');
  }
}

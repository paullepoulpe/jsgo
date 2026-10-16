// Oracle tool: writes the Java reference vectors for the jscomp-a port (error
// formatting with LightweightMessageFormatter / VerboseMessageFormatter over
// real source files, the JSON and print-stream error report generators,
// SortingErrorManager ordering, warnings guards, RhinoErrorReporter's message
// mapping, ModuleLoader and its resolvers, ModuleNames, SortedDependencies,
// DependencyInfo.Util.writeAddDependency, SourceMapResolver and
// JSDocSerializer) into internal/jscomp/testdata/jscomp_a.tsv. Every case
// records its inputs, so internal/jscomp/jscomp_a_differential_test.go
// replays it without a seeded generator of its own.
//
// It lives in package com.google.javascript.jscomp to reach package-private
// API (VerboseMessageFormatter, RhinoErrorReporter, ErrorWithLevel,
// JSDocSerializer); deps.ModuleNames is reached by reflection. Run with a
// UTF-8 locale, from the repository root:
//   LC_ALL=C.UTF-8 java -cp .cache/closure-compiler.jar:.cache/oracle-build \
//       com.google.javascript.jscomp.JscompAVectors .cache/corpus .cache/corpus-files.txt \
//       internal/jscomp/testdata/jscomp_a.tsv
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableList;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableMap;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSet;
import com.google.debugging.sourcemap.proto.Mapping.OriginalMapping;
import com.google.javascript.jscomp.deps.BrowserModuleResolver;
import com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver;
import com.google.javascript.jscomp.deps.DependencyInfo;
import com.google.javascript.jscomp.deps.ModuleLoader;
import com.google.javascript.jscomp.deps.NodeModuleResolver;
import com.google.javascript.jscomp.deps.SimpleDependencyInfo;
import com.google.javascript.jscomp.deps.SortedDependencies;
import com.google.javascript.jscomp.deps.WebpackModuleResolver;
import com.google.javascript.jscomp.ijs.CheckTypeSummaryWarningsGuard;
import com.google.javascript.jscomp.serialization.JsdocTag;
import com.google.javascript.jscomp.serialization.OptimizationJsdoc;
import com.google.javascript.rhino.JSDocInfo;
import com.google.javascript.rhino.Msg;
import com.google.javascript.rhino.Node;
import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Random;
import jsgo.oracle.CorpusExclusions;

public class JscompAVectors {
  static final Random r = new Random(20260415);
  static PrintStream out;

  static String esc(String s) {
    if (s == null) {
      return "\\N";
    }
    StringBuilder b = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      switch (c) {
        case '\\' -> b.append("\\\\");
        case '\n' -> b.append("\\n");
        case '\t' -> b.append("\\t");
        case '\r' -> b.append("\\r");
        default -> b.append(c);
      }
    }
    return b.toString();
  }

  static void row(String... cols) {
    StringBuilder b = new StringBuilder();
    for (int i = 0; i < cols.length; i++) {
      if (i > 0) {
        b.append('\t');
      }
      b.append(esc(cols[i]));
    }
    out.println(b);
  }

  static String exc(Throwable t) {
    return "!" + t.getClass().getSimpleName() + ": " + t.getMessage();
  }

  // ---------------------------------------------------------------------------
  // Diagnostic types and groups shared with the Go test.

  static final DiagnosticType[] TYPES = {
    DiagnosticType.error("JSA_E0", "error zero {0}"),
    DiagnosticType.warning("JSA_W1", "warning one <{0}>"),
    DiagnosticType.disabled("JSA_D2", "disabled two"),
    ClosurePrimitiveErrors.DUPLICATE_NAMESPACE_AND_MODULE,
    RhinoErrorReporter.TYPE_PARSE_ERROR,
    StrictWarningsGuard.UNRAISABLE_WARNING,
    DiagnosticType.warning("JSA_W6", "w6 '{0}' \"q\" & é"),
  };

  static final DiagnosticGroup[] GROUPS = {
    new DiagnosticGroup("g0", TYPES[0], TYPES[1]),
    new DiagnosticGroup("g1", TYPES[1]),
    new DiagnosticGroup("g2", TYPES[2], TYPES[6]),
    DiagnosticGroup.forType(TYPES[6]),
    new DiagnosticGroup("g4", new DiagnosticGroup[] {
      new DiagnosticGroup("g0", TYPES[0], TYPES[1]), new DiagnosticGroup("g2", TYPES[2], TYPES[6])}),
    new DiagnosticGroup("g5", TYPES[4], TYPES[3]),
  };

  static final CheckLevel[] LEVELS = {CheckLevel.ERROR, CheckLevel.WARNING, CheckLevel.OFF};

  // ---------------------------------------------------------------------------
  // Source excerpts.

  static final Map<String, SourceFile> files = new LinkedHashMap<>();

  /**
   * A SourceExcerptProvider over the "src" files. Names starting with "orig/" are the reverse
   * mapped originals (the same file); lines whose number is 3 mod 5 have such a mapping.
   */
  static final SourceExcerptProvider PROVIDER =
      new SourceExcerptProvider() {
        SourceFile file(String name) {
          if (name == null) {
            return null;
          }
          if (name.startsWith("orig/")) {
            name = name.substring(5);
          }
          return files.get(name);
        }

        @Override
        public String getSourceLine(String sourceName, int lineNumber) {
          SourceFile f = file(sourceName);
          return f == null ? null : f.getLine(lineNumber);
        }

        @Override
        public Region getSourceLines(String sourceName, int lineNumber, int length) {
          SourceFile f = file(sourceName);
          return f == null ? null : f.getLines(lineNumber, length);
        }

        @Override
        public Region getSourceRegion(String sourceName, int lineNumber) {
          SourceFile f = file(sourceName);
          return f == null ? null : f.getRegion(lineNumber);
        }

        @Override
        public OriginalMapping getSourceMapping(
            String sourceName, int lineNumber, int columnNumber) {
          if (sourceName == null || sourceName.startsWith("orig/") || lineNumber % 5 != 3) {
            return null;
          }
          return OriginalMapping.newBuilder()
              .setOriginalFile("orig/" + sourceName)
              .setLineNumber(lineNumber + 1)
              .setColumnPosition(columnNumber + 1)
              .build();
        }
      };

  static String srcExcerpt(Path p) throws Exception {
    String code = Files.readString(p, StandardCharsets.UTF_8);
    // Keep the first 40 lines, at most 3000 chars.
    int end = 0;
    int lines = 0;
    while (end < code.length() && lines < 40 && end < 3000) {
      int nl = code.indexOf('\n', end);
      if (nl < 0) {
        end = code.length();
        break;
      }
      end = nl + 1;
      lines++;
    }
    String s = code.substring(0, Math.min(end, code.length()));
    if (s.length() > 3000) {
      s = s.substring(0, 3000);
      // Do not split a surrogate pair.
      if (Character.isHighSurrogate(s.charAt(s.length() - 1))) {
        s = s.substring(0, s.length() - 1);
      }
    }
    return s;
  }

  static JSError makeError(String name, int line, int charno, int len, int nodeMode, int typeIdx) {
    DiagnosticType type = TYPES[typeIdx];
    if (nodeMode == 0) {
      return JSError.make(name, line, charno, type, "arg" + line);
    }
    Node n = Node.newString("x").setLinenoCharno(line, charno);
    n.setLength(len);
    if (name != null) {
      n.setStaticSourceFile(files.get(name));
    }
    return JSError.make(n, type, "arg" + line);
  }

  static void formatting(Path corpus, List<String> names) throws Exception {
    for (String name : names) {
      String code = srcExcerpt(corpus.resolve(name));
      SourceFile sf = SourceFile.fromCode(name, code);
      files.put(name, sf);
      row("src", name, code);
      int numLines = sf.getNumLines();
      SortingErrorManager manager = new SortingErrorManager(ImmutableSet.of());
      List<String[]> reported = new ArrayList<>();
      for (int k = 0; k < 6; k++) {
        int line = r.nextInt(8) == 0 ? r.nextInt(3) - 1 : 1 + r.nextInt(numLines + 1);
        String text = sf.getLine(Math.max(line, 1));
        int lineLen = text == null ? 0 : text.length();
        int charno =
            switch (r.nextInt(6)) {
              case 0 -> -1;
              case 1 -> lineLen;
              case 2 -> lineLen + 1 + r.nextInt(3);
              default -> lineLen == 0 ? 0 : r.nextInt(lineLen);
            };
        int len = r.nextInt(4) == 0 ? 0 : 1 + r.nextInt(r.nextBoolean() ? 8 : 200);
        int nodeMode = r.nextInt(3) == 0 ? 0 : 1;
        int typeIdx = r.nextInt(TYPES.length);
        int levelIdx = r.nextInt(3);
        String srcName = r.nextInt(15) == 0 ? null : name;
        JSError error = makeError(srcName, line, charno, len, nodeMode, typeIdx);
        List<String> cols = new ArrayList<>();
        cols.add("fmt");
        cols.add(srcName);
        cols.add(Integer.toString(line));
        cols.add(Integer.toString(charno));
        cols.add(Integer.toString(len));
        cols.add(Integer.toString(nodeMode));
        cols.add(Integer.toString(typeIdx));
        cols.add(Integer.toString(levelIdx));
        for (ErrorFormat f : ErrorFormat.values()) {
          for (boolean colorize : new boolean[] {false, true}) {
            if (colorize && f != ErrorFormat.SINGLELINE) {
              continue;
            }
            MessageFormatter formatter = f.toFormatter(PROVIDER, colorize);
            cols.add(safe(() -> formatter.formatError(error)));
            cols.add(safe(() -> formatter.formatWarning(error)));
          }
        }
        VerboseMessageFormatter verbose = new VerboseMessageFormatter(PROVIDER);
        cols.add(safe(() -> verbose.formatError(error)));
        LightweightMessageFormatter noLevel =
            new LightweightMessageFormatter(PROVIDER).setIncludeLevel(false).setIncludeLocation(false);
        cols.add(safe(() -> noLevel.formatWarning(error)));
        cols.add(safe(() -> new LightweightMessageFormatter(PROVIDER).getExcerptWithPosition(error)));
        row(cols.toArray(new String[0]));
        manager.report(LEVELS[levelIdx], error);
      }
      int typed = r.nextInt(3);
      manager.setTypedPercent(typed == 0 ? 0.0 : typed == 1 ? 12.25 : 99.95);
      ByteArrayOutputStream json = new ByteArrayOutputStream();
      new JsonErrorReportGenerator(new PrintStream(json, true, StandardCharsets.UTF_8), PROVIDER)
          .generateReport(manager);
      ByteArrayOutputStream ps1 = new ByteArrayOutputStream();
      new PrintStreamErrorReportGenerator(
              ErrorFormat.SINGLELINE.toFormatter(PROVIDER, false),
              new PrintStream(ps1, true, StandardCharsets.UTF_8),
              1)
          .generateReport(manager);
      ByteArrayOutputStream ps3 = new ByteArrayOutputStream();
      new PrintStreamErrorReportGenerator(
              ErrorFormat.FULL.toFormatter(PROVIDER, false),
              new PrintStream(ps3, true, StandardCharsets.UTF_8),
              3)
          .generateReport(manager);
      ByteArrayOutputStream ps2 = new ByteArrayOutputStream();
      PrintStreamErrorManager pem =
          new PrintStreamErrorManager(new PrintStream(ps2, true, StandardCharsets.UTF_8));
      pem.setSummaryDetailLevel(2);
      for (JSError e : manager.getErrors()) {
        pem.report(CheckLevel.ERROR, e);
      }
      for (JSError e : manager.getWarnings()) {
        pem.report(CheckLevel.WARNING, e);
      }
      pem.setTypedPercent(manager.getTypedPercent());
      pem.generateReport();
      row(
          "report",
          name,
          Integer.toString(typed),
          json.toString(StandardCharsets.UTF_8),
          ps1.toString(StandardCharsets.UTF_8),
          ps3.toString(StandardCharsets.UTF_8),
          ps2.toString(StandardCharsets.UTF_8),
          manager.getErrorCount() + "/" + manager.getWarningCount() + "/" + manager.hasHaltingErrors());
    }
  }

  interface ThrowingSupplier {
    String get() throws Exception;
  }

  static String safe(ThrowingSupplier s) {
    try {
      return s.get();
    } catch (Exception e) {
      return exc(e);
    }
  }

  // ---------------------------------------------------------------------------
  // SortingErrorManager ordering.

  static final String[] SOURCES = {null, "a.js", "b.js", "dir/a.js", "foo/x.i.js", "gen/z.js", "B.js"};

  static void sorting() {
    for (int c = 0; c < 60; c++) {
      SortingErrorManager m = new SortingErrorManager(ImmutableSet.of());
      StringBuilder spec = new StringBuilder();
      int n = 1 + r.nextInt(10);
      for (int i = 0; i < n; i++) {
        String src = SOURCES[r.nextInt(SOURCES.length)];
        int line = r.nextInt(4) - 1;
        int charno = r.nextInt(4) - 1;
        int typeIdx = r.nextInt(TYPES.length);
        int levelIdx = r.nextInt(3);
        if (i > 0) {
          spec.append(';');
        }
        spec.append(src == null ? "~" : src)
            .append('|')
            .append(line)
            .append('|')
            .append(charno)
            .append('|')
            .append(typeIdx)
            .append('|')
            .append(levelIdx);
        m.report(LEVELS[levelIdx], JSError.make(src, line, charno, TYPES[typeIdx], "a"));
      }
      StringBuilder sorted = new StringBuilder();
      for (SortingErrorManager.ErrorWithLevel e : m.getSortedDiagnostics()) {
        sorted
            .append(e.level)
            .append(' ')
            .append(e.error.sourceName())
            .append(':')
            .append(e.error.getLineNumber())
            .append(':')
            .append(e.error.charno())
            .append(' ')
            .append(e.error.description())
            .append('\n');
      }
      row(
          "sorting",
          spec.toString(),
          sorted.toString(),
          m.getErrorCount() + "/" + m.getWarningCount() + "/" + m.hasHaltingErrors() + "/"
              + m.getErrors().size() + "/" + m.getWarnings().size());
    }
  }

  // ---------------------------------------------------------------------------
  // Warnings guards. Specs: see the Go test's jsaParseGuard.

  static String randomPaths() {
    String[] parts = {"/foo/", "bar", "a.js", "gen/", "x.i.js", ""};
    int n = 1 + r.nextInt(2);
    StringBuilder b = new StringBuilder();
    for (int i = 0; i < n; i++) {
      if (i > 0) {
        b.append(',');
      }
      b.append(parts[r.nextInt(parts.length)]);
    }
    return b.toString();
  }

  static String randomGuardSpec(int depth) {
    switch (r.nextInt(depth > 1 ? 8 : 9)) {
      case 0:
        return "bypath:" + (r.nextBoolean() ? "F" : "X") + ":" + (r.nextBoolean() ? "ERROR" : "OFF")
            + ":" + randomPaths();
      case 1:
        return "show:" + (r.nextBoolean() ? "I" : "E") + ":" + randomPaths();
      case 2:
        return "strict";
      case 3:
      case 4:
        return "dg:" + r.nextInt(GROUPS.length) + ":" + LEVELS[r.nextInt(3)];
      case 5:
        return "dgps:" + r.nextInt(GROUPS.length) + ":" + new String[] {"gen/", "a", "/"}[r.nextInt(3)];
      case 6:
        return "summary:" + LEVELS[r.nextInt(3)];
      case 7:
        return "cu";
      default:
        {
          int n = r.nextInt(5);
          StringBuilder b = new StringBuilder("compose(");
          for (int i = 0; i < n; i++) {
            if (i > 0) {
              b.append(';');
            }
            if (i > 0 && r.nextInt(5) == 0) {
              b.append("=").append(r.nextInt(i));
            } else {
              b.append(randomGuardSpec(depth + 1));
            }
          }
          return b.append(')').toString();
        }
    }
  }

  static int findClose(String s, int open) {
    int depth = 0;
    for (int i = open; i < s.length(); i++) {
      if (s.charAt(i) == '(') {
        depth++;
      } else if (s.charAt(i) == ')') {
        depth--;
        if (depth == 0) {
          return i;
        }
      }
    }
    throw new IllegalStateException(s);
  }

  static List<String> splitTop(String s) {
    List<String> parts = new ArrayList<>();
    int depth = 0;
    int start = 0;
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '(') {
        depth++;
      } else if (c == ')') {
        depth--;
      } else if (c == ';' && depth == 0) {
        parts.add(s.substring(start, i));
        start = i + 1;
      }
    }
    if (start < s.length()) {
      parts.add(s.substring(start));
    }
    return parts;
  }

  static WarningsGuard parseGuard(String spec) {
    String[] f = spec.split(":", -1);
    switch (f[0]) {
      case "bypath":
        {
          List<String> paths = List.of(f[3].split(",", -1));
          CheckLevel level = CheckLevel.valueOf(f[2]);
          return f[1].equals("F")
              ? ByPathWarningsGuard.forPath(paths, level)
              : ByPathWarningsGuard.exceptPath(paths, level);
        }
      case "show":
        return new ShowByPathWarningsGuard(
            f[2].split(",", -1),
            f[1].equals("I")
                ? ShowByPathWarningsGuard.ShowType.INCLUDE
                : ShowByPathWarningsGuard.ShowType.EXCLUDE);
      case "strict":
        return new StrictWarningsGuard();
      case "dg":
        return new DiagnosticGroupWarningsGuard(
            GROUPS[Integer.parseInt(f[1])], CheckLevel.valueOf(f[2]));
      case "dgps":
        return new DiagnosticGroupPathSuppressingWarningsGuard(GROUPS[Integer.parseInt(f[1])], f[2]);
      case "summary":
        return new CheckTypeSummaryWarningsGuard(CheckLevel.valueOf(f[1]));
      case "cu":
        return new ClosureUnawareCodeWarningsGuard();
      default:
        {
          String inner = spec.substring("compose(".length(), spec.length() - 1);
          List<WarningsGuard> guards = new ArrayList<>();
          for (String part : splitTop(inner)) {
            if (part.startsWith("=")) {
              guards.add(guards.get(Integer.parseInt(part.substring(1))));
            } else {
              guards.add(parseGuard(part));
            }
          }
          return new ComposeWarningsGuard(guards);
        }
    }
  }

  static void guards() {
    for (int c = 0; c < 120; c++) {
      String spec = c < 10 ? randomGuardSpec(9) : "compose(" + randomGuardSpec(1) + ";" + randomGuardSpec(1) + ")";
      if (c >= 10 && c % 2 == 0) {
        spec = randomGuardSpec(0);
      }
      WarningsGuard guard;
      try {
        guard = parseGuard(spec);
      } catch (RuntimeException e) {
        row("guard", spec, "", exc(e));
        continue;
      }
      StringBuilder levels = new StringBuilder();
      StringBuilder errors = new StringBuilder();
      for (int i = 0; i < 12; i++) {
        String src = SOURCES[r.nextInt(SOURCES.length)];
        int typeIdx = r.nextInt(TYPES.length);
        int line = r.nextInt(14);
        if (i > 0) {
          errors.append(';');
          levels.append(';');
        }
        errors.append(src == null ? "~" : src).append('|').append(typeIdx).append('|').append(line);
        try {
          levels.append(guard.level(JSError.make(src, line, -1, TYPES[typeIdx])));
        } catch (RuntimeException e) {
          levels.append(exc(e));
        }
      }
      StringBuilder must = new StringBuilder();
      for (int g = 0; g < GROUPS.length; g++) {
        if (g > 0) {
          must.append(';');
        }
        try {
          must.append(guard.mustRunChecks(GROUPS[g]));
        } catch (RuntimeException e) {
          must.append(exc(e));
        }
      }
      row(
          "guard",
          spec,
          errors.toString(),
          levels.toString(),
          must.toString(),
          Integer.toString(guard.getPriority()),
          guard instanceof ComposeWarningsGuard cg ? Integer.toString(cg.getGuards().size()) : "-");
    }
  }

  // ---------------------------------------------------------------------------
  // RhinoErrorReporter.mapError and forOldRhino.

  static void mapError() throws Exception {
    Method mapError = RhinoErrorReporter.class.getDeclaredMethod("mapError", String.class);
    mapError.setAccessible(true);
    List<String> messages = new ArrayList<>();
    for (Msg m : Msg.values()) {
      messages.add(m.format());
      messages.add(m.format("x", "y", "z"));
    }
    messages.add("Trailing comma is not legal in an ECMA-262 object initializer");
    messages.add("Duplicate parameter name \"a\"");
    messages.add("Duplicate parameter name \"a\" extra");
    messages.add("Unnecessary escape: '\\a'");
    messages.add("invalid param name \"x\"");
    messages.add("Too many template parameters\nFound 2");
    messages.add("Too many template parameters");
    messages.add("Missing type declaration.");
    messages.add("x Unknown type y");
    messages.add("x Unknown type y\nz");
    messages.add("Bad type annotation. Import in typedef");
    messages.add("Bad type annotation.\nmore");
    messages.add("Octal literal");
    messages.add("Octal integer literal in strict mode");
    messages.add("This language feature is only supported for ECMASCRIPT_2016 mode or better: x");
    messages.add("This language feature is not currently supported by the compiler: y");
    messages.add("something else");
    messages.add("");
    for (String message : messages) {
      DiagnosticType t = (DiagnosticType) mapError.invoke(null, message);
      List<String> reported = new ArrayList<>();
      ErrorHandler h = (level, error) -> reported.add(
          level + " " + error.type().key + " " + error.defaultLevel() + " " + error.sourceName()
              + ":" + error.getLineNumber() + ":" + error.charno() + " " + error.description());
      com.google.javascript.rhino.ErrorReporter reporter = RhinoErrorReporter.forOldRhino(h);
      reporter.warning(message, "w.js", 3, 4);
      reporter.error(message, "e.js", 5, 6);
      row("maperr", message, t == null ? null : t.key, String.join("\n", reported));
    }
  }

  // ---------------------------------------------------------------------------
  // ModuleNames, ModuleLoader and the resolvers.

  static String invokeModuleNames(String fn, String arg) throws Exception {
    Method m = Class.forName("com.google.javascript.jscomp.deps.ModuleNames").getDeclaredMethod(fn, String.class);
    m.setAccessible(true);
    try {
      return (String) m.invoke(null, arg);
    } catch (java.lang.reflect.InvocationTargetException e) {
      return exc(e.getCause());
    }
  }

  static final String[] PATH_PARTS = {
    "a", "b", "..", ".", "node_modules", "lib", "x:y", "with space", "@scope", "pkg-name",
    "index.js", "f.js", "f", "dir", "", "é", "a.b.js", "c\\d.js",
  };

  static String randomPath() {
    StringBuilder b = new StringBuilder();
    int kind = r.nextInt(4);
    if (kind == 0) {
      b.append('/');
    } else if (kind == 1) {
      b.append("./");
    } else if (kind == 2 && r.nextBoolean()) {
      b.append("../");
    }
    int n = 1 + r.nextInt(4);
    for (int i = 0; i < n; i++) {
      if (i > 0) {
        b.append('/');
      }
      b.append(PATH_PARTS[r.nextInt(PATH_PARTS.length)]);
    }
    if (r.nextInt(3) == 0) {
      b.append(".js");
    }
    return b.toString();
  }

  static void moduleNames(List<String> corpusNames) throws Exception {
    String[] fns = {
      "escapePath", "canonicalizePath", "toModuleName", "toJSIdentifier", "fileToModuleName",
      "fileToJsIdentifier"
    };
    List<String> inputs = new ArrayList<>();
    for (int i = 0; i < 80; i++) {
      inputs.add(randomPath());
    }
    for (int i = 0; i < 40; i++) {
      inputs.add(corpusNames.get(r.nextInt(corpusNames.size())));
    }
    inputs.add("file://a/b.jar!com/example/test");
    inputs.add("/a/../../../b");
    for (String in : inputs) {
      List<String> cols = new ArrayList<>(List.of("modnames", in));
      for (String fn : fns) {
        cols.add(invokeModuleNames(fn, in));
      }
      row(cols.toArray(new String[0]));
    }
    for (int i = 0; i < 80; i++) {
      String from = randomPath();
      String to = r.nextInt(4) == 0 ? new String[] {"m0.js", "./x.js", "/abs.js", "\\w.js"}[r.nextInt(4)] : randomPath();
      if (r.nextInt(5) == 0) {
        from = "m.js";
      }
      String result;
      try {
        result = ModuleLoader.relativePathFrom(from, to);
      } catch (RuntimeException e) {
        result = exc(e);
      }
      row("relpath", from, to, result);
    }
  }

  static String mapSpec(Map<String, String> m) {
    StringBuilder b = new StringBuilder();
    for (Map.Entry<String, String> e : m.entrySet()) {
      if (b.length() > 0) {
        b.append(',');
      }
      b.append(e.getKey()).append('=').append(e.getValue());
    }
    return b.toString();
  }

  static void modules(List<String> corpusNames) {
    for (int c = 0; c < 70; c++) {
      List<String> inputs = new ArrayList<>();
      int n = 1 + r.nextInt(8);
      for (int i = 0; i < n; i++) {
        String p;
        switch (r.nextInt(4)) {
          case 0 -> p = corpusNames.get(r.nextInt(corpusNames.size()));
          case 1 -> p = "/node_modules/" + PATH_PARTS[r.nextInt(PATH_PARTS.length)] + "/index.js";
          case 2 -> p = (r.nextBoolean() ? "/" : "") + "src/" + PATH_PARTS[r.nextInt(PATH_PARTS.length)] + ".js";
          default -> p = randomPath();
        }
        inputs.add(p);
      }
      if (r.nextInt(3) == 0) {
        inputs.add("/node_modules/pkg/package.json");
        inputs.add("/node_modules/pkg/lib/main.js");
      }
      List<String> roots = new ArrayList<>();
      for (int i = r.nextInt(3); i > 0; i--) {
        roots.add(new String[] {".", "src", "/node_modules/", "@angular", "gen/", "/"}[r.nextInt(6)]);
      }
      Map<String, String> map = new LinkedHashMap<>();
      String factorySpec;
      ModuleLoader.ModuleResolverFactory factory;
      switch (r.nextInt(4)) {
        case 0 -> {
          factorySpec = "browser";
          factory = BrowserModuleResolver.FACTORY;
        }
        case 1 -> {
          if (r.nextBoolean()) {
            map.put("/node_modules/pkg/package.json", "/node_modules/pkg/lib/main.js");
            map.put("node_modules/pkg/lib/x.js", ModuleLoader.JSC_BROWSER_SKIPLISTED_MARKER);
          }
          // "nodenull:" passes a null map, "node:" the (possibly empty) map.
          boolean nullMap = map.isEmpty() && r.nextBoolean();
          factorySpec = (nullMap ? "nodenull:" : "node:") + mapSpec(map);
          factory = new NodeModuleResolver.Factory(nullMap ? null : map);
        }
        case 2 -> {
          map.put("@p/", "/src/");
          if (r.nextBoolean()) {
            map.put("@p/x/", "/node_modules/");
          }
          map.put("~/", "/");
          factorySpec = "prefix:" + mapSpec(map);
          factory = new BrowserWithTransformedPrefixesModuleResolver.Factory(ImmutableMap.copyOf(map));
        }
        default -> {
          for (int i = 0; i < inputs.size(); i++) {
            map.put(Integer.toString(i), inputs.get(i));
          }
          factorySpec = "webpack:" + mapSpec(map);
          factory = new WebpackModuleResolver.Factory(map);
        }
      }
      boolean escape = r.nextInt(4) != 0;
      List<String> errors = new ArrayList<>();
      ModuleLoader loader;
      String id = Integer.toString(c);
      ImmutableList.Builder<DependencyInfo> deps = ImmutableList.builder();
      for (String in : inputs) {
        deps.add(SimpleDependencyInfo.builder(in, in).build());
      }
      try {
        loader =
            ModuleLoader.builder()
                .setModuleRoots(roots)
                .setInputs(deps.build())
                .setFactory(factory)
                .setErrorHandler((level, error) -> errors.add(level + " " + error.type().key + " " + error.description()))
                .setPathEscaper(escape ? ModuleLoader.PathEscaper.ESCAPE : ModuleLoader.PathEscaper.CANONICALIZE_ONLY)
                .build();
      } catch (RuntimeException e) {
        row("modsetup", id, factorySpec, String.join(",", roots), String.join("\u0001", inputs), Boolean.toString(escape), exc(e));
        continue;
      }
      row("modsetup", id, factorySpec, String.join(",", roots), String.join("\u0001", inputs), Boolean.toString(escape), mapSpec(loader.getPackageJsonMainEntries()));
      for (int q = 0; q < 10; q++) {
        String from = r.nextBoolean() ? inputs.get(r.nextInt(inputs.size())) : randomPath();
        String addr;
        switch (r.nextInt(5)) {
          case 0 -> addr = inputs.get(r.nextInt(inputs.size()));
          case 1 -> addr = map.isEmpty() ? "pkg" : map.keySet().iterator().next() + "f.js";
          case 2 -> addr = r.nextBoolean() ? "pkg" : "./" + PATH_PARTS[r.nextInt(PATH_PARTS.length)];
          case 3 -> addr = Integer.toString(r.nextInt(4));
          default -> addr = randomPath();
        }
        errors.clear();
        String resolved;
        String js;
        String silent;
        String asPath;
        try {
          ModuleLoader.ModulePath p = loader.resolve(from);
          resolved = p.toString() + " " + p.toModuleName();
          ModuleLoader.ModulePath j = p.resolveJsModule(addr, "s.js", 1, 2);
          js = j == null ? null : j.toString();
          ModuleLoader.ModulePath s = p.resolveJsModuleSilently(addr);
          silent = s == null ? null : s.toString();
          asPath = p.resolveModuleAsPath(addr).toString();
        } catch (RuntimeException e) {
          resolved = exc(e);
          js = silent = asPath = "";
        }
        row("modres", id, from, addr, resolved, js, silent, asPath, String.join("\n", errors));
      }
    }
  }

  // ---------------------------------------------------------------------------
  // SortedDependencies and DependencyInfo.Util.writeAddDependency.

  static void deps() throws Exception {
    for (int c = 0; c < 60; c++) {
      int n = 1 + r.nextInt(7);
      List<SimpleDependencyInfo> infos = new ArrayList<>();
      StringBuilder spec = new StringBuilder();
      for (int i = 0; i < n; i++) {
        String name = "f" + i + ".js";
        List<String> provides = new ArrayList<>();
        for (int k = r.nextInt(3); k > 0; k--) {
          provides.add("n" + r.nextInt(8));
        }
        List<String> requires = new ArrayList<>();
        for (int k = r.nextInt(3); k > 0; k--) {
          requires.add("n" + r.nextInt(8));
        }
        List<String> typeRequires = new ArrayList<>();
        for (int k = r.nextInt(2); k > 0; k--) {
          typeRequires.add("n" + r.nextInt(8));
        }
        SimpleDependencyInfo.Builder b =
            SimpleDependencyInfo.builder("rel/" + name, name).setProvides(provides).setTypeRequires(typeRequires);
        List<DependencyInfo.Require> reqs = new ArrayList<>();
        for (String q : requires) {
          reqs.add(DependencyInfo.Require.googRequireSymbol(q));
        }
        b.setRequires(reqs);
        Map<String, String> flags = new LinkedHashMap<>();
        if (r.nextInt(3) == 0) {
          flags.put("module", r.nextBoolean() ? "goog" : "es6");
        }
        if (r.nextInt(4) == 0) {
          flags.put("lang", "es'8");
        }
        b.setLoadFlags(flags);
        infos.add(b.build());
        if (i > 0) {
          spec.append(';');
        }
        spec.append(name)
            .append('>')
            .append(String.join(",", provides))
            .append('|')
            .append(String.join(",", requires))
            .append('|')
            .append(String.join(",", typeRequires))
            .append('|')
            .append(mapSpec(flags));
        StringBuilder sb = new StringBuilder();
        DependencyInfo.Util.writeAddDependency(sb, infos.get(i));
        row("adddep", spec.substring(spec.lastIndexOf(name)), sb.toString(), infos.get(i).toString(),
            Boolean.toString(infos.get(i).isEs6Module()) + "/" + infos.get(i).isGoogModule());
      }
      List<SimpleDependencyInfo> roots = new ArrayList<>();
      StringBuilder rootSpec = new StringBuilder();
      for (int k = 1 + r.nextInt(2); k > 0; k--) {
        int idx = r.nextInt(n);
        roots.add(infos.get(idx));
        if (rootSpec.length() > 0) {
          rootSpec.append(',');
        }
        rootSpec.append(idx);
      }
      String sorted;
      String strong;
      String weak;
      String providing;
      try {
        SortedDependencies<SimpleDependencyInfo> sd = new SortedDependencies<>(infos);
        sorted = names(sd.getSortedList());
        strong = names(sd.getSortedStrongDependenciesOf(roots));
        weak = names(sd.getSortedWeakDependenciesOf(roots));
        StringBuilder pb = new StringBuilder();
        for (int k = 0; k < 8; k++) {
          SimpleDependencyInfo p = sd.maybeGetInputProviding("n" + k);
          pb.append(p == null ? "-" : p.getName()).append(',');
        }
        pb.append(names(sd.getInputsWithoutProvides()));
        providing = pb.toString();
      } catch (Exception e) {
        sorted = exc(e);
        strong = weak = providing = "";
      }
      row("deps", spec.toString(), rootSpec.toString(), sorted, strong, weak, providing);
    }
  }

  static String names(List<? extends DependencyInfo> l) {
    List<String> out = new ArrayList<>();
    for (DependencyInfo d : l) {
      out.add(d.getName());
    }
    return String.join(",", out);
  }

  // ---------------------------------------------------------------------------
  // SourceMapResolver.

  static void sourceMapResolver() {
    for (int i = 0; i < 40; i++) {
      String base = randomPath();
      String rel = r.nextInt(4) == 0 ? "data:application/json;base64,eyJ2IjozfQ==" : randomPath();
      String name;
      String inline;
      try {
        name = SourceMapResolver.getRelativePath(base, rel).getName();
      } catch (RuntimeException e) {
        name = exc(e);
      }
      try {
        SourceFile sf = SourceMapResolver.extractSourceMap(SourceFile.fromCode(base + "x", ""), rel, true);
        // Only an inline map has code; the others name files that do not exist.
        inline = sf == null ? null : sf.getName() + (sf.getName().endsWith(".inline.map") ? " " + sf.getCode() : "");
      } catch (Exception e) {
        inline = exc(e);
      }
      row("smr", base, rel, name, inline);
    }
  }

  // ---------------------------------------------------------------------------
  // JSDocSerializer.

  static void jsdoc() {
    for (int c = 0; c < 60; c++) {
      JSDocInfo.Builder b = JSDocInfo.builder();
      StringBuilder ops = new StringBuilder();
      for (int k = r.nextInt(6); k >= 0; k--) {
        int op = r.nextInt(20);
        ops.append(op).append(',');
        switch (op) {
          case 0 -> b.recordConstancy();
          case 1 -> b.recordNoSideEffects();
          case 2 -> b.recordModifies(ImmutableSet.of("this", "arguments"));
          case 3 -> b.recordDescription("desc " + c);
          case 4 -> b.recordMeaning("meaning");
          case 5 -> b.recordAlternateMessageId("alt");
          case 6 -> b.recordSuppressions(ImmutableSet.of("partialAlias", "messageConventions"));
          case 7 -> b.recordSuppressions(ImmutableSet.of("untranspilableFeatures", "other"));
          case 8 -> b.recordStableIdGenerator();
          case 9 -> b.recordConsistentIdGenerator();
          case 10 -> b.recordIdGenerator();
          case 11 -> b.recordConstructor();
          case 12 -> b.recordInterface();
          case 13 -> b.recordAbstract();
          case 14 -> b.recordNoInline();
          case 15 -> b.recordTypeSummary();
          case 16 -> b.addLicense("lic");
          case 17 -> b.recordFileOverview("fo");
          case 18 -> b.recordThrowsAnnotation("{Error}");
          default -> b.recordClosureUnawareCode(JSDocInfo.PerFileClosureUnawareMode.SIMPLE);
        }
      }
      JSDocInfo info = b.build();
      com.google.javascript.jscomp.serialization.StringPool.Builder pool =
          com.google.javascript.jscomp.serialization.StringPool.builder();
      String ser;
      String de;
      try {
        OptimizationJsdoc s = serializeJsdoc(info, pool);
        ser = s == null ? null : s.toString();
        JSDocInfo d = com.google.javascript.jscomp.serialization.JSDocSerializer.convertJSDocInfoForOptimizations(info);
        de = d == null ? null : jsdocSummary(d);
      } catch (Exception e) {
        ser = exc(e);
        de = "";
      }
      row("jsdoc", ops.toString(), ser, de);
    }
  }

  static OptimizationJsdoc serializeJsdoc(
      JSDocInfo info, com.google.javascript.jscomp.serialization.StringPool.Builder pool) throws Exception {
    Method m =
        com.google.javascript.jscomp.serialization.JSDocSerializer.class.getDeclaredMethod(
            "serializeJsdoc", JSDocInfo.class, com.google.javascript.jscomp.serialization.StringPool.Builder.class);
    m.setAccessible(true);
    return (OptimizationJsdoc) m.invoke(null, info, pool);
  }

  static String jsdocSummary(JSDocInfo d) {
    return d.isConstant() + " " + d.isNoSideEffects() + " " + d.getModifies() + " " + d.getDescription()
        + " " + d.getMeaning() + " " + d.getAlternateMessageId() + " " + d.getSuppressions() + " "
        + d.isStableIdGenerator() + d.isConsistentIdGenerator() + d.isIdGenerator() + " "
        + d.isConstructor() + d.isInterface() + d.isAbstract() + d.isNoInline() + d.isTypeSummary()
        + " " + d.getLicense() + " " + d.hasFileOverview() + " " + d.getThrowsAnnotations() + " "
        + d.getPerFileClosureUnawareMode();
  }

  public static void main(String[] args) throws Exception {
    Path corpus = Paths.get(args[0]);
    List<String> corpusNames = Files.readAllLines(Paths.get(args[1]));
    CorpusExclusions.removeFrom(corpusNames); // tools/oracle/gpl-excluded-corpus.txt
    out = new PrintStream(Files.newOutputStream(Paths.get(args[2])), true, StandardCharsets.UTF_8);
    List<String> picked = new ArrayList<>();
    for (int i = 0; i < 32; i++) {
      String name = corpusNames.get(r.nextInt(corpusNames.size()));
      if (!picked.contains(name)) {
        picked.add(name);
      }
    }
    formatting(corpus, picked);
    sorting();
    guards();
    mapError();
    moduleNames(corpusNames);
    modules(corpusNames);
    deps();
    sourceMapResolver();
    jsdoc();
    out.close();
  }
}

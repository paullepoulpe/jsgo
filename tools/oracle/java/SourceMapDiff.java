package com.google.debugging.sourcemap;

import com.google.debugging.sourcemap.proto.Mapping.OriginalMapping;
import com.google.javascript.jscomp.jarjar.com.google.gson.Gson;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonElement;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.OutputStreamWriter;
import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Collection;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.zip.GZIPOutputStream;

/**
 * Differential oracle for internal/sourcemap: a seeded script of calls on the public API of
 * com.google.debugging.sourcemap (SourceMapGeneratorV3 with mappings, prefixes, starting
 * positions, sources content, source roots, extensions, resets, merges and index maps;
 * SourceMapConsumerV3 / SourceMapConsumerFactory / SourceMapObjectParser on the generated maps,
 * on mutated maps and on hand-written edge cases; Base64, Base64VLQ and Util.escapeString), each
 * call followed by what it returned or printed.
 *
 * <p>Output (gzipped): internal/sourcemap/testdata/differential.txt.gz, one line per call
 * ("&gt;\top\targ...") followed by its results ("&lt;\ttext"). Every field is escaped to printable
 * ASCII like GsonVectors ("\\0" is null). internal/sourcemap/differential_test.go replays the
 * calls on the Go port and compares.
 *
 * <p>The class is in the package so that it can call the package-private Util.escapeString. JVM
 * exceptions without a Java-defined message are reported by kind only: "!NPE" for
 * NullPointerException, "!INDEX" for (Array|String)IndexOutOfBoundsException.
 *
 * <p>Usage: java -cp closure-compiler.jar:. com.google.debugging.sourcemap.SourceMapDiff out.gz
 */
public final class SourceMapDiff {
  private static final Gson GSON = new Gson();

  private final PrintWriter out;
  private final List<String> buffers = new ArrayList<>();
  private SourceMapGeneratorV3 g = new SourceMapGeneratorV3();
  private SourceMapConsumerV3 c = new SourceMapConsumerV3();
  private int ops = 0;

  private SourceMapDiff(PrintWriter out) {
    this.out = out;
  }

  static String esc(String s) {
    if (s == null) {
      return "\\0";
    }
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char ch = s.charAt(i);
      if (ch == '\\') {
        sb.append("\\\\");
      } else if (ch == '\t') {
        sb.append("\\t");
      } else if (ch == '\n') {
        sb.append("\\n");
      } else if (ch >= 0x20 && ch <= 0x7e) {
        sb.append(ch);
      } else {
        sb.append(String.format("\\u%04x", (int) ch));
      }
    }
    return sb.toString();
  }

  static String exc(Throwable t) {
    if (t instanceof SourceMapParseException) {
      return "!SMPE " + t.getMessage();
    }
    if (t instanceof NullPointerException) {
      return "!NPE";
    }
    if (t instanceof IndexOutOfBoundsException) {
      return "!INDEX";
    }
    return "!" + t;
  }

  static String describe(Object o) {
    if (o == null) {
      return "null";
    }
    return o.getClass().getSimpleName() + ":" + o;
  }

  static String str(String s) {
    return s == null ? "<null>" : s;
  }

  static String strs(String[] a) {
    if (a == null) {
      return "<nullarray>";
    }
    StringBuilder sb = new StringBuilder("[");
    for (int i = 0; i < a.length; i++) {
      sb.append(i == 0 ? "" : "|").append(str(a[i]));
    }
    return sb.append("]").toString();
  }

  static String strs(Collection<String> a) {
    if (a == null) {
      return "<nullarray>";
    }
    return strs(a.toArray(new String[0]));
  }

  static String om(OriginalMapping m) {
    if (m == null) {
      return "-";
    }
    return m.getOriginalFile()
        + ":"
        + m.getLineNumber()
        + ":"
        + m.getColumnPosition()
        + (m.hasIdentifier() ? "=" + m.getIdentifier() : "")
        + ":"
        + m.getPrecision().getNumber();
  }

  private String resolve(String arg) {
    if (arg != null && arg.startsWith("@")) {
      return buffers.get(Integer.parseInt(arg.substring(1)));
    }
    return arg;
  }

  private void result(String s) {
    out.println("<\t" + esc(s));
  }

  /** Writes the call, runs it, writes its results. */
  private void emit(String... op) {
    ops++;
    StringBuilder sb = new StringBuilder(">");
    for (String a : op) {
      sb.append('\t').append(esc(a));
    }
    out.println(sb);
    try {
      run(op);
    } catch (Throwable t) {
      result(exc(t));
    }
  }

  private void run(String[] op) throws Exception {
    switch (op[0]) {
      case "case" -> buffers.clear();
      case "gnew" -> g = new SourceMapGeneratorV3();
      case "gfactory" ->
          g =
              (SourceMapGeneratorV3)
                  SourceMapGeneratorFactory.getInstance(SourceMapFormat.valueOf(op[1]));
      case "greset" -> g.reset();
      case "gvalidate" -> g.validate(Boolean.parseBoolean(op[1]));
      case "gsourceroot" -> g.setSourceRoot(op[1]);
      case "gprefix" -> g.setWrapperPrefix(op[1]);
      case "gstart" -> g.setStartingPosition(Integer.parseInt(op[1]), Integer.parseInt(op[2]));
      case "gmap" ->
          g.addMapping(
              op[1],
              op[2],
              new FilePosition(Integer.parseInt(op[3]), Integer.parseInt(op[4])),
              new FilePosition(Integer.parseInt(op[5]), Integer.parseInt(op[6])),
              new FilePosition(Integer.parseInt(op[7]), Integer.parseInt(op[8])));
      case "gcontent" -> g.addSourcesContent(op[1], op[2]);
      case "gext" -> {
        Object value =
            switch (op[2]) {
              case "str" -> op[3];
              case "int" -> Integer.parseInt(op[3]);
              case "bool" -> Boolean.parseBoolean(op[3]);
              case "json" -> GSON.fromJson(op[3], JsonElement.class);
              default -> null;
            };
        g.addExtension(op[1], value);
      }
      case "gextrm" -> g.removeExtension(op[1]);
      case "gexthas" -> result(String.valueOf(g.hasExtension(op[1])));
      case "gextget" -> result(describe(g.getExtension(op[1])));
      case "gappend" -> {
        StringBuilder sb = new StringBuilder();
        try {
          g.appendTo(sb, op[1]);
        } finally {
          buffers.add(sb.toString());
          result(sb.toString());
        }
      }
      case "gindex" -> {
        List<SourceMapSection> sections = new ArrayList<>();
        int n = Integer.parseInt(op[2]);
        for (int i = 0; i < n; i++) {
          String type = op[3 + 4 * i];
          String value = resolve(op[4 + 4 * i]);
          int line = Integer.parseInt(op[5 + 4 * i]);
          int column = Integer.parseInt(op[6 + 4 * i]);
          sections.add(
              type.equals("URL")
                  ? SourceMapSection.forURL(value, line, column)
                  : SourceMapSection.forMap(value, line, column));
        }
        StringBuilder sb = new StringBuilder();
        try {
          g.appendIndexMapTo(sb, op[1], sections);
        } finally {
          buffers.add(sb.toString());
          result(sb.toString());
        }
      }
      case "gmerge" -> {
        g.mergeMapSection(Integer.parseInt(op[1]), Integer.parseInt(op[2]), resolve(op[3]));
        result("ok");
      }
      case "gmergex" -> {
        g.mergeMapSection(
            Integer.parseInt(op[1]),
            Integer.parseInt(op[2]),
            resolve(op[3]),
            (key, cur, nw) -> "merged(" + key + "," + describe(cur) + "," + describe(nw) + ")");
        result("ok");
      }
      case "cparse" -> {
        c = new SourceMapConsumerV3();
        c.parse(resolve(op[1]));
        result("ok");
      }
      case "cfactory" -> {
        Map<String, String> maps = new LinkedHashMap<>();
        for (int i = 2; i < op.length; i++) {
          int eq = op[i].indexOf('=');
          maps.put(op[i].substring(0, eq), resolve(op[i].substring(eq + 1)));
        }
        SourceMapping m =
            SourceMapConsumerFactory.parse(
                resolve(op[1]),
                url -> {
                  String v = maps.get(url);
                  if ("!io".equals(v)) {
                    throw new IOException("cannot read " + url);
                  }
                  return v;
                });
        c = (SourceMapConsumerV3) m;
        result("ok");
      }
      case "proto" -> protoScenario();
      case "cfactory0" -> {
        c = (SourceMapConsumerV3) SourceMapConsumerFactory.parse(resolve(op[1]));
        result("ok");
      }
      case "gbig" -> {
        // A map with many sources/names (the consumer's non-compact layout).
        StringBuilder sb = new StringBuilder("{\"version\":3,\"mappings\":\"");
        sb.append(op[3]).append("\",\"sources\":[");
        for (int i = 0; i < Integer.parseInt(op[1]); i++) {
          sb.append(i == 0 ? "" : ",").append("\"s").append(i).append('"');
        }
        sb.append("],\"names\":[");
        for (int i = 0; i < Integer.parseInt(op[2]); i++) {
          sb.append(i == 0 ? "" : ",").append("\"n").append(i).append('"');
        }
        sb.append("]}");
        buffers.add(sb.toString());
        result(String.valueOf(sb.length()));
      }
      case "clookup" -> {
        for (int i = 1; i + 1 < op.length; i += 2) {
          result(om(c.getMappingForLine(Integer.parseInt(op[i]), Integer.parseInt(op[i + 1]))));
        }
      }
      case "creverse" -> {
        StringBuilder sb = new StringBuilder();
        for (OriginalMapping m : c.getReverseMapping(op[1], Integer.parseInt(op[2]), 0)) {
          sb.append(' ').append(m.getLineNumber()).append(',').append(m.getColumnPosition());
        }
        result(sb.toString());
      }
      case "cvisitcount" -> {
        int[] n = {0};
        StringBuilder sb = new StringBuilder();
        c.visitMappings(
            (sourceName, symbolName, src, start, end) -> {
              if (n[0]++ < 6) {
                sb.append(str(sourceName)).append('|').append(str(symbolName)).append(' ');
              }
            });
        result(n[0] + " " + sb);
      }
      case "cfast" -> {
        c = SourceMapConsumerFactory.parseFast(resolve(op[1]));
        result("ok");
      }
      case "cdump" -> dumpConsumer();
      case "oparse" -> dumpObject(SourceMapObjectParser.parse(resolve(op[1])));
      case "ofast" -> dumpObject(SourceMapObjectParser.parseFast(resolve(op[1])));
      case "venc" -> {
        StringBuilder sb = new StringBuilder();
        Base64VLQ.encode(sb, Integer.parseInt(op[1]));
        result(sb.toString());
      }
      case "vdec" -> {
        String s = op[1];
        int[] pos = {0};
        Base64VLQ.CharIterator it =
            new Base64VLQ.CharIterator() {
              @Override
              public boolean hasNext() {
                return pos[0] < s.length();
              }

              @Override
              public char next() {
                return s.charAt(pos[0]++);
              }
            };
        StringBuilder sb = new StringBuilder();
        try {
          while (it.hasNext()) {
            sb.append(Base64VLQ.decode(it)).append(' ');
          }
        } finally {
          result(sb.toString());
        }
      }
      case "b64" -> result(Base64.base64EncodeInt(Integer.parseInt(op[1])));
      case "b64from" -> result(String.valueOf(Base64.fromBase64(op[1].charAt(0))));
      case "esc" -> result(Util.escapeString(op[1]));
      default -> throw new AssertionError(op[0]);
    }
  }

  static String proto(OriginalMapping m) {
    return m.hasOriginalFile() + "," + m.getOriginalFile() + "," + m.hasLineNumber() + ","
        + m.getLineNumber() + "," + m.hasColumnPosition() + "," + m.getColumnPosition() + ","
        + m.hasIdentifier() + "," + m.getIdentifier() + "," + m.hasPrecision() + ","
        + m.getPrecision();
  }

  /** The reconstructed OriginalMapping (proto_mapping.go) against the generated class. */
  private void protoScenario() {
    OriginalMapping d = OriginalMapping.getDefaultInstance();
    result("default " + proto(d) + " " + d.isInitialized());
    OriginalMapping a =
        OriginalMapping.newBuilder()
            .setOriginalFile("a.js")
            .setLineNumber(3)
            .setColumnPosition(-4)
            .setIdentifier("id")
            .setPrecision(OriginalMapping.Precision.APPROXIMATE_LINE)
            .build();
    result("a " + proto(a));
    OriginalMapping b = a.toBuilder().clearIdentifier().clearPrecision().build();
    result("b " + proto(b));
    OriginalMapping c = OriginalMapping.newBuilder(b).setLineNumber(0).clearOriginalFile().build();
    result("c " + proto(c));
    OriginalMapping e =
        OriginalMapping.newBuilder().setPrecision(OriginalMapping.Precision.UNKNOWN_PRECISION).build();
    result("e " + proto(e));
    OriginalMapping f =
        OriginalMapping.newBuilder().setOriginalFile("").setLineNumber(0).setIdentifier("").build();
    result("f " + proto(f));
    OriginalMapping g = OriginalMapping.newBuilder().mergeFrom(b).mergeFrom(f).build();
    result("g " + proto(g));
    OriginalMapping h = a.toBuilder().clear().setColumnPosition(9).build();
    result("h " + proto(h));
    OriginalMapping.Builder i = OriginalMapping.newBuilder(a);
    result("i " + i.hasIdentifier() + i.getIdentifier() + i.getPrecision() + i.getLineNumber());
    OriginalMapping[] all = {d, a, b, c, e, f, g, h, a.toBuilder().build(), d.toBuilder().build()};
    StringBuilder sb = new StringBuilder("equals");
    for (OriginalMapping x : all) {
      sb.append(' ');
      for (OriginalMapping y : all) {
        sb.append(x.equals(y) ? '1' : '0');
      }
    }
    result(sb.toString());
    result("forNumber " + OriginalMapping.Precision.forNumber(0) + " "
        + OriginalMapping.Precision.forNumber(2) + " " + OriginalMapping.Precision.forNumber(3) + " "
        + OriginalMapping.Precision.EXACT.getNumber() + " "
        + java.util.Arrays.toString(OriginalMapping.Precision.values()));
  }

  private interface Body {
    void run() throws Exception;
  }

  private void guarded(Body b) {
    try {
      b.run();
    } catch (Throwable t) {
      result(exc(t));
    }
  }

  private void dumpObject(SourceMapObject o) {
    result("version=" + o.getVersion() + " lineCount=" + o.getLineCount());
    result("file=" + str(o.getFile()));
    result("mappings=" + str(o.getMappings()));
    result("sourceRoot=" + str(o.getSourceRoot()));
    result("sources=" + strs(o.getSources()));
    result("sourcesContent=" + strs(o.getSourcesContent()));
    result("names=" + strs(o.getNames()));
    if (o.getSections() == null) {
      result("sections=null");
    } else {
      for (SourceMapSection s : o.getSections()) {
        result("section " + s.getSectionType() + " " + s.getLine() + " " + s.getColumn() + " "
            + str(s.getSectionValue()));
      }
    }
    for (Map.Entry<String, Object> e : o.getExtensions().entrySet()) {
      result("ext " + e.getKey() + "=" + describe(e.getValue()));
    }
  }

  private void dumpConsumer() {
    SourceMapConsumerV3 c = this.c;
    guarded(() -> result("file=" + str(c.getFile()) + " lineCount=" + c.getLineCount()
        + " sourceRoot=" + str(c.getSourceRoot())));
    guarded(() -> result("sources=" + strs(c.getOriginalSources())));
    guarded(() -> result("sourcesContent=" + strs(c.getOriginalSourcesContent())));
    guarded(() -> result("names=" + strs(c.getOriginalNames())));
    guarded(() -> {
      for (Map.Entry<String, Object> e : c.getExtensions().entrySet()) {
        result("ext " + e.getKey() + "=" + describe(e.getValue()));
      }
    });
    int[] max = {0, 0, 0}; // line, column, source line
    guarded(() ->
        c.visitMappings(
            (sourceName, symbolName, src, start, end) -> {
              result("visit " + str(sourceName) + "|" + str(symbolName) + "|" + src.getLine()
                  + "," + src.getColumn() + "|" + start.getLine() + "," + start.getColumn() + "|"
                  + end.getLine() + "," + end.getColumn());
              max[0] = Math.max(max[0], Math.max(start.getLine(), end.getLine()));
              max[1] = Math.max(max[1], Math.max(start.getColumn(), end.getColumn()));
              max[2] = Math.max(max[2], src.getLine());
            }));
    int maxLine = Math.min(Math.max(c.getLineCount(), max[0]) + 2, 300);
    int maxCol = Math.min(max[1] + 3, 300);
    for (int line = 1; line <= maxLine; line++) {
      int l = line;
      guarded(() -> {
        StringBuilder sb = new StringBuilder("L" + l + ":");
        String prev = null;
        int runStart = 0;
        for (int col = 0; col <= maxCol + 1; col++) {
          String cur = col <= maxCol ? om(c.getMappingForLine(l, col)) : null;
          if (prev != null && !prev.equals(cur)) {
            sb.append(' ').append(runStart).append('-').append(col - 1).append('=').append(prev);
            runStart = col;
          }
          if (prev == null) {
            runStart = col;
          }
          prev = cur;
        }
        result(sb.toString());
      });
    }
    int[][] probes = {
      {0, 1}, {-1, 1}, {1, -5}, {1, Integer.MAX_VALUE}, {1, Integer.MIN_VALUE},
      {Integer.MAX_VALUE, 1}, {Integer.MIN_VALUE, 1}, {2, Integer.MIN_VALUE}
    };
    for (int[] p : probes) {
      guarded(() -> result("P" + p[0] + "," + p[1] + "=" + om(c.getMappingForLine(p[0], p[1]))));
    }
    guarded(() -> {
      LinkedHashSet<String> files = new LinkedHashSet<>();
      for (String s : c.getOriginalSources()) {
        if (s != null) {
          files.add(s);
        }
      }
      files.add("not-a-source.js");
      for (String f : files) {
        for (int line = -1; line <= Math.min(max[2] + 1, 200); line++) {
          Collection<OriginalMapping> rs = c.getReverseMapping(f, line, 0);
          if (!rs.isEmpty()) {
            StringBuilder sb = new StringBuilder("R " + f + " " + line + ":");
            for (OriginalMapping m : rs) {
              sb.append(' ').append(m.getLineNumber()).append(',').append(m.getColumnPosition());
            }
            result(sb.toString());
          }
        }
      }
    });
  }

  // ---------------------------------------------------------------------------------------------
  // Script generation

  private static final String[] POOL = {
    "a.js", "b.js", "dir/c.js", "", "c:\\myfile.js", "x\"y", "it's", "tab\there", "nl\nx", "cr\rx",
    "</script>", "</SCRIPT", "<!-- x", "-->", "]]>", "<", ">", "a-->b", "a]]>b", "--", "]]", "\\",
    "/", "\u00e9.js", "\u65e5\u672c.js", "\ud83d\ude00.js", "\ud800", "\udc00x", "\u2028",
    "\u0000", "\u001f", "\u007f", "\u0080", "</scr\u0131pt", "</\u017fcript", "</scr\u0130pt",
    "<\u0000", "foo", "bar", "__BASIC__", "$jscomp", "x_y",
  };

  private static final char[] ALPHABET = {
    'a', 'b', 'z', 'S', 'C', 'R', 'I', 'P', 'T', 's', 'c', 'r', 'i', 'p', 't', ' ', '"', '\'',
    '\\', '\n', '\t', '\r', '<', '>', '/', '-', ']', '!', '\u00e9', '\u65e5', '\ud800', '\u2028',
    '\u0000', '\u001f', '\u007f', '\u0080', '\u0131', '\u017f', '\u0130', '\ufeff', '{', '}', ',',
    ';', ':',
  };

  static String randomString(Random r) {
    if (r.nextInt(3) > 0) {
      return POOL[r.nextInt(POOL.length)];
    }
    StringBuilder sb = new StringBuilder();
    int n = r.nextInt(12);
    for (int i = 0; i < n; i++) {
      int k = r.nextInt(ALPHABET.length + 3);
      if (k == ALPHABET.length) {
        sb.append("\ud83d\ude00");
      } else if (k == ALPHABET.length + 1) {
        sb.append("</script");
      } else if (k == ALPHABET.length + 2) {
        sb.append("<!--");
      } else {
        sb.append(ALPHABET[k]);
      }
    }
    return sb.toString();
  }

  static final String[] JSON_VALUES = {
    "{}", "[]", "{\"a\":[1,2,{\"b\":null}],\"c\":\"\\u00e9<>&='\"}", "\"str\\u00e9\\n\"", "1.50",
    "-0", "1e400", "true", "false", "null", "[\"\\u2028\",{}]", "{'single':unquoted}", "12345678901234567890",
  };

  private int last() {
    return buffers.size() - 1;
  }

  private void randomGeneratorCase(Random r, int caseNo) {
    emit("case", String.valueOf(caseNo));
    if (r.nextInt(5) == 0) {
      emit("gfactory", r.nextBoolean() ? "DEFAULT" : "V3");
    } else {
      emit("gnew");
    }
    if (r.nextInt(4) == 0) {
      emit("gsourceroot", randomString(r));
    }
    if (r.nextInt(4) == 0) {
      StringBuilder p = new StringBuilder();
      int n = r.nextInt(4);
      for (int i = 0; i < n; i++) {
        p.append(randomString(r)).append(r.nextBoolean() ? "\n" : "");
      }
      emit("gprefix", p.toString());
    }
    if (r.nextInt(6) == 0) {
      emit("gvalidate", String.valueOf(r.nextBoolean()));
    }
    int rounds = r.nextInt(5) == 0 ? 2 : 1;
    for (int round = 0; round < rounds; round++) {
      if (round > 0) {
        emit("greset");
      }
      if (r.nextInt(4) == 0) {
        emit("gstart", String.valueOf(r.nextInt(30) - (r.nextInt(30) == 0 ? 31 : 0)),
            String.valueOf(r.nextInt(50)));
      }
      int count = r.nextInt(5) == 0 ? r.nextInt(150) : r.nextInt(40);
      randomMappings(r, count);
      int extensions = r.nextInt(3) == 0 ? r.nextInt(5) : 0;
      for (int i = 0; i < extensions; i++) {
        randomExtensionOp(r);
      }
      int contents = r.nextInt(3) == 0 ? r.nextInt(4) : 0;
      for (int i = 0; i < contents; i++) {
        emit("gcontent", POOL[r.nextInt(8)], randomString(r));
      }
      emit("gappend", r.nextInt(6) == 0 ? null : randomString(r));
      int map = last();
      emit("cparse", "@" + map);
      emit("cdump");
      if (r.nextInt(3) == 0) {
        emit("cfast", "@" + map);
        emit("cdump");
      }
      if (r.nextInt(3) == 0) {
        emit("ofast", "@" + map);
      }
      if (r.nextInt(4) == 0) {
        emit("oparse", "@" + map);
      }
      if (r.nextInt(3) == 0) {
        String mutated = mutate(r, buffers.get(map));
        emit("cparse", mutated);
        emit("cdump");
        emit("ofast", mutated);
        emit("oparse", mutated);
      }
    }
    if (r.nextInt(3) == 0) {
      int sources = last();
      emit("gnew");
      int n = 1 + r.nextInt(3);
      int line = 0;
      for (int i = 0; i < n; i++) {
        line += r.nextInt(5);
        String col = String.valueOf(r.nextInt(20));
        if (r.nextBoolean()) {
          emit("gmerge", String.valueOf(line), col, "@" + (sources - r.nextInt(sources + 1)));
        } else {
          emit("gmergex", String.valueOf(line), col, "@" + (sources - r.nextInt(sources + 1)));
        }
        if (r.nextInt(2) == 0) {
          randomExtensionOp(r);
        }
      }
      emit("gappend", "merged.js");
      emit("cparse", "@" + last());
      emit("cdump");
    }
    if (r.nextInt(3) == 0) {
      int maps = last();
      List<String> op = new ArrayList<>();
      op.add("gindex");
      op.add(randomString(r));
      int n = 1 + r.nextInt(3);
      op.add(String.valueOf(n));
      List<String> supplier = new ArrayList<>();
      int line = 0;
      for (int i = 0; i < n; i++) {
        line += r.nextInt(6);
        String ref = "@" + (maps - r.nextInt(maps + 1));
        if (r.nextBoolean()) {
          op.add("MAP");
          op.add(ref);
        } else {
          op.add("URL");
          op.add("u" + i);
          int k = r.nextInt(10);
          supplier.add("u" + i + "=" + (k == 0 ? "!io" : ref));
        }
        op.add(String.valueOf(line));
        op.add(String.valueOf(r.nextInt(30)));
      }
      emit("gnew");
      emit(op.toArray(new String[0]));
      List<String> f = new ArrayList<>();
      f.add("cfactory");
      f.add("@" + last());
      f.addAll(supplier);
      emit(f.toArray(new String[0]));
      emit("cdump");
    }
  }

  private void randomExtensionOp(Random r) {
    String name = r.nextInt(8) == 0 ? "bad_" + r.nextInt(3) : "x_ext" + r.nextInt(4);
    switch (r.nextInt(8)) {
      case 0, 1 -> emit("gext", name, "str", randomString(r));
      case 2 -> emit("gext", name, "int", String.valueOf(r.nextInt() >> r.nextInt(31)));
      case 3 -> emit("gext", name, "bool", String.valueOf(r.nextBoolean()));
      case 4, 5 -> emit("gext", name, "json", JSON_VALUES[r.nextInt(JSON_VALUES.length)]);
      case 6 -> emit(r.nextBoolean() ? "gexthas" : "gextget", name);
      default -> emit(r.nextInt(6) == 0 ? "gext" : "gextrm", name, "null", "");
    }
  }

  private void randomMappings(Random r, int count) {
    // A pre-order walk of nested output ranges, with occasional overlapping, inverted, out of
    // order and far-away ranges.
    List<int[]> stack = new ArrayList<>(); // {endLine, endCol}
    int line = 0;
    int col = 0;
    for (int i = 0; i < count; i++) {
      while (!stack.isEmpty() && r.nextInt(3) == 0) {
        int[] top = stack.remove(stack.size() - 1);
        if (top[0] > line || (top[0] == line && top[1] > col)) {
          line = top[0];
          col = top[1];
        }
      }
      if (r.nextInt(6) == 0) {
        line += 1 + r.nextInt(3);
        col = r.nextInt(10);
      } else {
        col += r.nextInt(8);
      }
      if (r.nextInt(60) == 0) {
        col += 1000 + r.nextInt(100000);
      }
      int startLine = line;
      int startCol = col;
      if (r.nextInt(80) == 0) {
        startCol = Math.max(0, col - 1 - r.nextInt(5)); // out of order
      }
      int endLine = startLine;
      int endCol = startCol + 1 + r.nextInt(20);
      if (r.nextInt(10) == 0) {
        endLine = startLine + 1 + r.nextInt(2);
        endCol = r.nextInt(30);
      }
      if (!stack.isEmpty() && r.nextInt(20) > 0) {
        int[] parent = stack.get(stack.size() - 1);
        if (endLine > parent[0] || (endLine == parent[0] && endCol > parent[1])) {
          endLine = parent[0];
          endCol = parent[1];
        }
      }
      if (r.nextInt(40) == 0) {
        endLine = startLine;
        endCol = startCol - r.nextInt(3); // empty or inverted
      }
      String src = r.nextInt(25) == 0 ? null : POOL[r.nextInt(r.nextInt(3) == 0 ? POOL.length : 4)];
      String name = r.nextInt(5) < 2 ? null : randomString(r);
      int sl = r.nextInt(30) == 0 ? -1 : r.nextInt(r.nextInt(5) == 0 ? 5000 : 50);
      int sc = r.nextInt(30) == 0 ? -r.nextInt(5) : r.nextInt(80);
      if (r.nextInt(200) == 0) {
        sc = Integer.MAX_VALUE - r.nextInt(3);
      }
      emit("gmap", src, name, String.valueOf(sl), String.valueOf(sc), String.valueOf(startLine),
          String.valueOf(startCol), String.valueOf(endLine), String.valueOf(endCol));
      if (r.nextInt(2) == 0) {
        stack.add(new int[] {endLine, endCol});
      }
      if (r.nextInt(15) == 0) {
        emit("gstart", String.valueOf(r.nextInt(5)), String.valueOf(r.nextInt(30)));
      }
    }
  }

  static String mutate(Random r, String s) {
    if (s.isEmpty()) {
      return "{";
    }
    String[] junk = {",", ";", "}", "]", "\"", "\\", "null", "-1", "1e9", "/*x*/", "\u00e9", "@", "A", "!"};
    int at = r.nextInt(s.length());
    switch (r.nextInt(4)) {
      case 0:
        return s.substring(0, at);
      case 1:
        return s.substring(0, at) + junk[r.nextInt(junk.length)] + s.substring(at + 1);
      case 2:
        return s.substring(0, at) + junk[r.nextInt(junk.length)] + s.substring(at);
      default:
        return s + junk[r.nextInt(junk.length)];
    }
  }

  // ---------------------------------------------------------------------------------------------
  // Hand-written cases

  private static final String BASIC =
      "{\"version\":3,\"file\":\"testcode\",\"lineCount\":1,\"mappings\":\"AAAAA,QAASA,UAAS,EAAG;\","
          + "\"sources\":[\"testcode\"],\"names\":[\"__BASIC__\"]}";

  static String map(String fields) {
    return "{\"version\":3," + fields + "}";
  }

  private static final String[] CONSUMER_CASES = {
    BASIC,
    map("\"file\":\"\",\"mappings\":\"\",\"sources\":[],\"names\":[]"),
    "{\"version\":2,\"mappings\":\"\",\"sources\":[],\"names\":[]}",
    "{\"mappings\":\"\",\"sources\":[],\"names\":[]}",
    "{\"version\":\"3\",\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    "{\"version\":3.0,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    "{\"version\":3e0,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    "{\"version\":[3],\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    "{\"version\":true,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    "{\"version\":{},\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    "{\"version\":null,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    "{\"version\":4294967299,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    map("\"lineCount\":1,\"mappings\":\"AAAA;AAAA;AAAA\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"lineCount\":0,\"mappings\":\"AAAA;\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"lineCount\":0,\"mappings\":\"\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"lineCount\":-7,\"mappings\":\"AAAA;;CAAC\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"lineCount\":1e1,\"mappings\":\"AAAA;;CAAC\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"lineCount\":\"2\",\"mappings\":\"AAAA;;CAAC\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"AA;\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"AAA;\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"AAAAAA;\",\"sources\":[\"a\"],\"names\":[\"n\"]"),
    map("\"mappings\":\"AAAA,,AAAA;,;\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"AA!A;\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"A\u00e9AA;\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"g\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"gggggggggggA,AAAA\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"//////D,AAAA\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"+/////D,AAAA;D,B,F\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"ADAA,CAAA\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"AFAA,CAAA\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"ACAA,CAAA\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"AAAAD,CAAAD\",\"sources\":[\"a\"],\"names\":[\"n\"]"),
    map("\"mappings\":\"AAAAC\",\"sources\":[\"a\"],\"names\":[\"n\"]"),
    map("\"mappings\":\"AAAA,CAAAA\",\"sources\":[\"a\"]"),
    map("\"mappings\":\"AAAA,CAAA\",\"sources\":[\"a\"]"),
    map("\"mappings\":\"AAAA,CAAA\",\"names\":[]"),
    map("\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":null,\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"AAAA\",\"sources\":[null,\"b\"],\"names\":[null]"),
    map("\"mappings\":\"ACAA,CDAAA\",\"sources\":[null,\"b\"],\"names\":[null]"),
    map("\"mappings\":\"ACAA\",\"sources\":\"a.js\",\"names\":[]"),
    map("\"mappings\":\"ACAA\",\"sources\":[{}],\"names\":[]"),
    map("\"mappings\":\"ACAA\",\"sources\":[[\"x\"]],\"names\":[]"),
    map("\"mappings\":\"ACAA\",\"sources\":[1,true,2.5],\"names\":[]"),
    map("\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[],\"sourceRoot\":null,\"file\":null"),
    map("\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[],\"sourceRoot\":5,\"file\":5"),
    map("\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[],\"sourcesContent\":[null,\"x\"]"),
    map("\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[],\"x_a\":1,\"x_b\":{\"k\":[1,null]},"
        + "\"x_c\":[\"\\u2028<\"],\"x_\":null,\"x\":1,\"X_y\":2,\"x_d\":1.50,\"x_e\":'single'"),
    map("\"mappings\":\"AAAA;AACA;;AAEA,CAAC;\",\"sources\":[\"a\",\"b\"],\"names\":[\"p\",\"q\"]"),
    map("\"mappings\":\"AAAA,EAAC,DAAE,DAAE,IAAI;;;;\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"KAAK,GAAG,AAAA\",\"sources\":[\"a\"],\"names\":[]"),
    map("\"mappings\":\"\\u0041AAA,C\\u0041AA\",\"sources\":[\"\\ud83d\\ude00\",\"\\ud800\"],"
        + "\"names\":[]"),
    map("\"mappings\":\"AAAA\",\"sources\":[\"a\\u12G4\"],\"names\":[]"),
    map("\"mappings\":\"AAAA\",\"sources\":[\"a\\u12\"],\"names\":[]"),
    map("\"mappings\":\"AAAA\",\"sources\":[\"a\\u+12f\"],\"names\":[]"),
    map("\"mappings\":\"AAAA\",\"sources\":[\"a\\u-12f\"],\"names\":[]"),
    map("\"mappings\":\"AAAA\",\"sources\":[\"a\\x\\/\\b\\f\\n\\r\\t\\\"\\\\\"],\"names\":[]"),
    "{\"version\":3,\"version\":2,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    "{\"version\":2,\"version\":3,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    "{/* c */\"version\":3,'mappings':'AAAA',sources:[a],names:[]}",
    "{\"version\":3,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[],}",
    "{\"version\":3,\"mappings\":\"AAAA\",\"sources\":[\"a\",],\"names\":[]}",
    "{\"version\":3,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]} x",
    "{\"version\":3,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}}",
    "{\"version\":3,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]",
    "{\"version\":3,\"mappings\":\"AAAA",
    "{\"version\":3,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[],\"other\":{\"a\":[1,"
        + "\"]\",{}],\"b\":-1.5e+3},\"x_raw\":{\"n\":[true,false,null]}}",
    "{\"version\":3,\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[],\"other\":tru}",
    "{\"version\":-,\"mappings\":\"AAAA\"}",
    "{\"version\":99999999999,\"mappings\":\"AAAA\"}",
    "{\"version\":3 ,\"lineCount\": 2 , \"mappings\" : \"AAAA;AACA\" , \"sources\":[ \"a\" ] , "
        + "\"names\":[ ]\n}\n",
    "",
    "   ",
    "null",
    "[]",
    "{}",
    "\"str\"",
    "{\"version\":3}",
    // Index maps.
    "{\"version\":3,\"file\":\"idx\",\"sections\":[{\"offset\":{\"line\":0,\"column\":0},\"map\":"
        + BASIC + "},{\"offset\":{\"line\":2,\"column\":5},\"map\":" + BASIC + "}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":1,\"column\":2},\"map\":"
        + "\"{\\\"version\\\":3,\\\"mappings\\\":\\\"AAAA\\\",\\\"sources\\\":[\\\"s\\\"],"
        + "\\\"names\\\":[]}\"}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0,\"column\":0},\"url\":\"u\"}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0,\"column\":0},\"url\":\"u\",\"map\":{}}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0,\"column\":0}}]}",
    "{\"version\":3,\"sections\":[{\"map\":" + BASIC + "}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"column\":0},\"map\":" + BASIC + "}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0},\"map\":" + BASIC + "}]}",
    "{\"version\":3,\"sections\":[{\"offset\":[],\"map\":" + BASIC + "}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":\"x\",\"column\":0},\"map\":" + BASIC + "}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":-1,\"column\":0},\"map\":" + BASIC + "}]}",
    "{\"version\":3,\"sections\":[3]}",
    "{\"version\":3,\"sections\":{}}",
    "{\"version\":3,\"sections\":[]}",
    "{\"version\":3,\"sections\":null}",
    "{\"version\":3,\"lineCount\":1,\"sections\":[]}",
    "{\"version\":3,\"mappings\":\"\",\"sections\":[]}",
    "{\"version\":3,\"sources\":[],\"sections\":[]}",
    "{\"version\":3,\"names\":[],\"sections\":[]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0,\"column\":0},\"map\":{\"version\":3,"
        + "\"sections\":[{\"offset\":{\"line\":1,\"column\":1},\"map\":" + BASIC + "}]}}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0,\"column\":0},\"map\":null}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0,\"column\":0},\"map\":5}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0,\"column\":0},\"map\":" + BASIC
        + "},{\"offset\":{\"line\":0,\"column\":0},\"map\":" + BASIC + "}]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":3,\"column\":0},\"map\":" + BASIC
        + "},{\"offset\":{\"line\":1,\"column\":0},\"map\":" + BASIC + "}]}",
    "{\"version\":3,\"file\":\"\",\"sections\":[]}",
    "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0,\"column\":0},\"map\":{\"version\":3,"
        + "\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[],\"x_e\":1}}]}",
  };

  private static final String[] FACTORY_CASES = {
    "/** Begin line maps. **/{}", "garbage", "", " {}", "{\"version\":4}", "{\"version\":4,"
        + "\"mappings\":\"\",\"sources\":[],\"names\":[]}",
  };

  private static final String[] ESCAPES = {
    "", "abc", "a\"b", "a'b", "a\\b", "\n\r\t", "\b\f\u000b", "</script>", "</SCRIPT", "</scr",
    "<!--", "<!-", "<!-x", "-->", "--->", "->", "]]>", "]>", "a]]>", "<", ">", "<<", ">>",
    "</scr\u0131pt", "</\u017fcript", "</scr\u0130pt", "</scripT", "\u00e9", "\u0080", "\u007f",
    "\u001f", " ", "~", "\ud83d\ude00", "\ud800", "\udc00", "\uffff", "\u2028\u2029", "\ufeff",
    "x</script", "<</script", "--->>", "]]]>",
  };

  private static final int[] INTS = {
    0, 1, -1, 2, -2, 15, 16, 31, 32, -31, -32, 1023, 1024, -1024, 65535, 1 << 20, (1 << 30) - 1,
    1 << 30, Integer.MAX_VALUE, Integer.MAX_VALUE - 1, Integer.MIN_VALUE, Integer.MIN_VALUE + 1,
    0x70000000, -100, 42,
  };

  private static final String[] VLQS = {
    "", "A", "B", "C", "D", "gB", "hB", "+", "/", "AAAA", "//////", "gggggggA", "gggggggggggA",
    "+/////D", "//////D", "!", "\u00e9", "g", "AAg", ";", "a;b", "\u0100",
  };

  static String vlqs(int... values) {
    StringBuilder sb = new StringBuilder();
    try {
      for (int v : values) {
        Base64VLQ.encode(sb, v);
      }
    } catch (IOException e) {
      throw new AssertionError(e);
    }
    return sb.toString();
  }

  private void handWritten() {
    emit("case", "fixed");
    for (int v : INTS) {
      emit("venc", String.valueOf(v));
      emit("b64", String.valueOf(v));
    }
    for (String v : VLQS) {
      emit("vdec", v);
    }
    for (char ch : new char[] {'A', 'z', '0', '+', '/', '=', '!', '\u00ff'}) {
      emit("b64from", String.valueOf(ch));
    }
    emit("proto");
    for (String s : ESCAPES) {
      emit("esc", s);
    }
    for (String s : POOL) {
      emit("esc", s);
    }
    for (String m : CONSUMER_CASES) {
      emit("case", "consumer");
      emit("cparse", m);
      emit("cdump");
      emit("cfast", m);
      emit("cdump");
      emit("oparse", m);
      emit("ofast", m);
      emit("cfactory", m, "u=" + BASIC);
      emit("cdump");
      emit("gnew");
      emit("gmerge", "1", "3", m);
      emit("gmergex", "0", "0", m);
      emit("gappend", "m.js");
    }
    for (String m : FACTORY_CASES) {
      emit("cfactory", m);
      emit("cfactory0", m);
      emit("cfast", m);
    }
    // Large maps: the non-compact entry layout (65535 or more sources or names) and the
    // compact layout's 0xFFFF sentinels.
    emit("case", "big");
    String[][] big = {
      {"3", "65536", vlqs(0, 0, 0, 0, 65534) + "," + vlqs(2, 1, 1, 1, 1) + ";" + vlqs(1, 1, 1, 1, -3)
          + "," + vlqs(4) + "," + vlqs(1, -2, 0, 0)},
      {"65535", "2", vlqs(0, 65534, 3, 4) + "," + vlqs(1, -65534, 0, 0, 1) + ";" + vlqs(0, 7, 1, 1)},
      {"65534", "65534", vlqs(0, 65533, 0, 0, 65533) + "," + vlqs(1, -1, 1, 1, -1) + ";"
          + vlqs(0, 0, 0, 0, 2)},
      {"65534", "3", vlqs(0, 65535, 0, 0)},
      {"2", "65534", vlqs(0, 0, 0, 0, 65535) + "," + vlqs(1, 1, 0, 0)},
      {"70000", "70000", vlqs(0, 69999, 5, 5, 69999) + ";" + vlqs(3, -69999, -5, -5, -69999)},
    };
    for (String[] b : big) {
      emit("gbig", b[0], b[1], b[2]);
      int m = last();
      for (String kind : new String[] {"cparse", "cfast", "cfactory0"}) {
        emit(kind, "@" + m);
        emit("clookup", "1", "1", "1", "2", "1", "3", "1", "5", "2", "1", "2", "4", "2", "9", "3", "1");
        emit("creverse", "s0", "0");
        emit("creverse", "s1", "1");
        emit("creverse", "s65534", "3");
        emit("creverse", "s69999", "5");
        emit("cvisitcount");
      }
      emit("gnew");
      emit("gmerge", "2", "1", "@" + m);
      emit("gappend", "big.js");
      emit("cparse", "@" + last());
      emit("cvisitcount");
    }
    // Index maps with supplier failures.
    emit("case", "suppliers");
    String idx =
        "{\"version\":3,\"sections\":[{\"offset\":{\"line\":0,\"column\":0},\"url\":\"u\"},"
            + "{\"offset\":{\"line\":4,\"column\":2},\"url\":\"v\"}]}";
    emit("cfactory", idx, "u=" + BASIC, "v=" + BASIC);
    emit("cdump");
    emit("cfactory", idx, "u=" + BASIC);
    emit("cfactory", idx, "u=!io", "v=" + BASIC);
    emit("cfactory", idx, "u=" + BASIC, "v={");
    // Generator edge cases.
    emit("case", "generator");
    emit("gnew");
    emit("gvalidate", "true");
    emit("gvalidate", "false");
    emit("gappend", null);
    emit("gappend", "");
    emit("gstart", "-1", "0");
    emit("gstart", "0", "-1");
    emit("gext", "bad", "int", "1");
    emit("gext", "x_n", "null", "");
    emit("gappend", "a");
    emit("gextrm", "x_n");
    emit("gmap", "a", null, "0", "0", "0", "5", "0", "3");
    emit("gmap", "a", null, "0", "0", "0", "4", "0", "8");
    emit("gappend", "a");
    emit("gnew");
    emit("gprefix", "ab\ncd\n\u00e9\ud83d\ude00");
    emit("gmap", "a", "n", "1", "2", "0", "0", "0", "3");
    emit("gmap", "a", "n", "1", "2", "0", "3", "2", "3");
    emit("gmap", "b", null, "4", "4", "1", "0", "1", "1");
    emit("gcontent", "a", "content of a\n</script>");
    emit("gcontent", "zz", "unused");
    emit("gappend", "w");
    emit("cparse", "@" + last());
    emit("cdump");
    List<String> sections = new ArrayList<>(List.of("gindex", "idx", "3"));
    sections.addAll(List.of("URL", "u\"1", "0", "0", "MAP", "@" + last(), "3", "4"));
    sections.addAll(List.of("URL", "</script>", "2147483647", "-1"));
    emit(sections.toArray(new String[0]));
  }

  public static void main(String[] args) throws Exception {
    try (PrintWriter out =
        new PrintWriter(
            new OutputStreamWriter(
                new GZIPOutputStream(new FileOutputStream(args[0])), StandardCharsets.UTF_8))) {
      SourceMapDiff d = new SourceMapDiff(out);
      d.handWritten();
      Random r = new Random(20260415L);
      int cases = args.length > 1 ? Integer.parseInt(args[1]) : 300;
      for (int i = 0; i < cases; i++) {
        d.randomGeneratorCase(r, i);
      }
      System.err.println(d.ops + " calls");
    }
  }
}

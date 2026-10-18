// Benchmark-only equivalent of jsgo's immutable built-in extern parse cache.
// The pinned Apache-2.0 Closure sources are overlaid by build_java_worker.py;
// the stock reference jar is never modified.
package com.google.javascript.jscomp.parsing;

import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.StaticSourceFile;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.atomic.AtomicLong;

public final class BuiltinExternCache {
  private static final boolean ENABLED = Boolean.getBoolean("jsgo.bench.externCache");
  private static final boolean VERIFY = Boolean.getBoolean("jsgo.bench.verifyCache");
  private static final AtomicLong HITS = new AtomicLong();
  private static final AtomicLong MISSES = new AtomicLong();
  private record Key(String name, String source) {}
  private static final class Group {
    final Config config;
    final Map<Key, ParserRunner.ParseResult> files = new HashMap<>();
    Group(Config config) { this.config = config; }
  }
  private static final List<Group> GROUPS = new ArrayList<>();

  public static boolean eligible(StaticSourceFile file, Config config) {
    return ENABLED && file.isExtern() && !file.isClosureUnawareCode()
        && !file.isTypeScriptSource() && file.getName().startsWith("externs.zip//")
        && config.jsDocParsingMode() == Config.JsDocParsing.TYPES_ONLY;
  }

  static ParserRunner.ParseResult parse(
      StaticSourceFile file, String source, Config config, ErrorReporter reporter) {
    Key key = new Key(file.getName(), source);
    ParserRunner.ParseResult template = null;
    synchronized (GROUPS) {
      for (Group group : GROUPS) {
        if (group.config.equals(config)) {
          template = group.files.get(key);
          break;
        }
      }
    }
    if (template != null) {
      if (VERIFY) { HITS.incrementAndGet(); }
      return copy(template, file);
    }
    if (VERIFY) { MISSES.incrementAndGet(); }
    RecordingReporter record = new RecordingReporter(reporter);
    ParserRunner.ParseResult parsed = ParserRunner.parseUncached(file, source, config, record);
    if (record.reported || parsed.ast == null || !parsed.comments.isEmpty()
        || !parsed.topLevelStatementRanges.isEmpty() || parsed.sourceMapURL != null) {
      return parsed;
    }
    template = copy(parsed, file);
    synchronized (GROUPS) {
      for (Group group : GROUPS) {
        if (group.config.equals(config)) {
          if (group.files.size() < 512) { group.files.put(key, template); }
          return parsed;
        }
      }
      if (GROUPS.size() == 4) { GROUPS.remove(0); }
      Group group = new Group(config);
      group.files.put(key, template);
      GROUPS.add(group);
    }
    return parsed;
  }

  private static ParserRunner.ParseResult copy(ParserRunner.ParseResult template, StaticSourceFile file) {
    Node ast = template.ast.cloneTree(false);
    bind(ast, file);
    return new ParserRunner.ParseResult(ast, template.comments, template.topLevelStatementRanges,
        template.features, template.sourceMapURL);
  }

  private static void bind(Node node, StaticSourceFile file) {
    node.setStaticSourceFile(file);
    var doc = node.getJSDocInfo();
    if (doc != null) {
      doc = doc.cloneForParseTemplate();
      node.setJSDocInfo(doc);
      for (Node type : doc.getTypeNodes()) { bind(type, file); }
      for (Node transform : doc.getTypeTransformations().values()) { bind(transform, file); }
    }
    for (Node child = node.getFirstChild(); child != null; child = child.getNext()) { bind(child, file); }
  }

  public static String statistics() { return "hits=" + HITS.get() + " misses=" + MISSES.get(); }

  private static final class RecordingReporter implements ErrorReporter {
    final ErrorReporter target;
    boolean reported;
    RecordingReporter(ErrorReporter target) { this.target = target; }
    public void warning(String message, String source, int line, int offset) {
      reported = true;
      target.warning(message, source, line, offset);
    }
    public void error(String message, String source, int line, int offset) {
      reported = true;
      target.error(message, source, line, offset);
    }
  }
}

package com.google.javascript.jscomp.parsing;

import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.SimpleSourceFile;
import com.google.javascript.rhino.StaticSourceFile;

/** Isolation/metadata checks for the benchmark-only Java extern cache. */
public final class BuiltinExternCacheTest {
  private static final ErrorReporter REPORTER = new ErrorReporter() {
    public void warning(String message, String source, int line, int offset) {
      throw new AssertionError(message);
    }
    public void error(String message, String source, int line, int offset) {
      throw new AssertionError(message);
    }
  };
  private static String fingerprint(Node root) {
    StringBuilder out = new StringBuilder();
    fingerprint(root, out);
    return out.toString();
  }
  private static void fingerprint(Node node, StringBuilder out) {
    out.append(node.toString()).append('\n');
    var doc = node.getJSDocInfo();
    if (doc != null) {
      out.append(doc.toStringVerbose()).append(" inline=").append(doc.isInlineType());
      for (Node type : doc.getTypeNodes()) { fingerprint(type, out); }
      for (Node transform : doc.getTypeTransformations().values()) { fingerprint(transform, out); }
    }
    for (Node child = node.getFirstChild(); child != null; child = child.getNext()) { fingerprint(child, out); }
  }
  private static void checkBinding(Node node, StaticSourceFile file) {
    if (node.getStaticSourceFile() != file) { throw new AssertionError("stale source metadata"); }
    var doc = node.getJSDocInfo();
    if (doc != null) {
      for (Node type : doc.getTypeNodes()) { checkBinding(type, file); }
      for (Node transform : doc.getTypeTransformations().values()) { checkBinding(transform, file); }
    }
    for (Node child = node.getFirstChild(); child != null; child = child.getNext()) { checkBinding(child, file); }
  }
  private static void mutate(Node tree) {
    var doc = tree.getJSDocInfo();
    if (doc != null) {
      for (Node type : doc.getTypeNodes()) {
        type.setLinenoCharno(999, 3); type.detachChildren();
      }
      for (Node transform : doc.getTypeTransformations().values()) {
        transform.setLinenoCharno(777, 2); transform.detachChildren();
      }
    }
    for (Node child = tree.getFirstChild(); child != null; child = child.getNext()) { mutate(child); }
  }
  public static void main(String[] args) {
    if (!Boolean.getBoolean("jsgo.bench.externCache")
        || !Boolean.getBoolean("jsgo.bench.verifyCache")) {
      throw new AssertionError("run with externCache and verifyCache enabled");
    }
    Config config = ParserRunner.createConfig(Config.LanguageMode.ES_NEXT, null, Config.StrictMode.SLOPPY);
    String source = "function f(/** number */ x) {}\n/** @type {string} */ var value;\n"
        + "/** @template T := 'number' =: */ function transform() {}";
    var file = new SimpleSourceFile("externs.zip//cache-test.js", StaticSourceFile.SourceKind.EXTERN);
    Node original = ParserRunner.parseUncached(file, source, config, REPORTER).ast;
    Node first = ParserRunner.parse(file, source, config, REPORTER).ast;
    // A miss returns the raw parser tree; compare cached clones after rebinding.
    first.detachChildren();
    String cachedFingerprint = null;
    for (int i = 0; i < 4; i++) {
      var current = new SimpleSourceFile(file.getName(), StaticSourceFile.SourceKind.EXTERN);
      Node tree = ParserRunner.parse(current, source, config, REPORTER).ast;
      checkBinding(tree, current);
      if (!original.isEquivalentTo(tree)) { throw new AssertionError("clone differs from uncached parse"); }
      String got = fingerprint(tree);
      // Mutating one returned tree must leave subsequent clones unchanged.
      if (cachedFingerprint == null) { cachedFingerprint = got; }
      else if (!cachedFingerprint.equals(got)) { throw new AssertionError("mutable AST/JSDoc leaked"); }
      mutate(tree);
      tree.detachChildren();
    }
    if (!fingerprint(original).equals(cachedFingerprint)) {
      throw new AssertionError("cached clone changed parser/JSDoc facts");
    }
    if (!cachedFingerprint.contains("inline=true")) { throw new AssertionError("INLINE_TYPE was lost"); }
    if (fingerprint(ParserRunner.parse(file, "var different;", config, REPORTER).ast).equals(cachedFingerprint)) {
      throw new AssertionError("source change reused an old AST");
    }
    long hits = BuiltinExternCache.hitsForTest();
    Config strict = ParserRunner.createConfig(Config.LanguageMode.ES_NEXT, null, Config.StrictMode.STRICT);
    if (config.equals(strict)) { throw new AssertionError("test configurations unexpectedly equal"); }
    ParserRunner.parse(file, "var strictOptionFixture;", strict, REPORTER);
    if (BuiltinExternCache.hitsForTest() != hits) { throw new AssertionError("config change hit old group"); }
    ParserRunner.parse(file, "var strictOptionFixture;", strict, REPORTER);
    if (BuiltinExternCache.hitsForTest() != hits + 1) { throw new AssertionError("matching config did not hit"); }
    Config noInlineDocs = ParserRunner.createConfig(Config.LanguageMode.ES_NEXT,
        Config.JsDocParsing.INCLUDE_DESCRIPTIONS_WITH_WHITESPACE, Config.RunMode.STOP_AFTER_ERROR,
        null, false, Config.StrictMode.SLOPPY);
    if (BuiltinExternCache.eligible(file, noInlineDocs)) { throw new AssertionError("wrong JSDoc mode eligible"); }
    var custom = new SimpleSourceFile("custom.js", StaticSourceFile.SourceKind.EXTERN);
    if (BuiltinExternCache.eligible(custom, config)) { throw new AssertionError("custom extern eligible"); }
    int[] warnings = {0};
    ErrorReporter diagnosticReporter = new ErrorReporter() {
      public void warning(String message, String sourceName, int line, int offset) { warnings[0]++; }
      public void error(String message, String sourceName, int line, int offset) { warnings[0]++; }
    };
    long before = BuiltinExternCache.hitsForTest();
    for (int i = 0; i < 2; i++) {
      ParserRunner.parse(file, "/** @type { */ var invalid;", config, diagnosticReporter);
    }
    if (warnings[0] < 2 || BuiltinExternCache.hitsForTest() != before) {
      throw new AssertionError("diagnostic parse was cached or failed to replay diagnostics");
    }
    System.out.println("cache isolation, inline types, TTL trees and source rebinding PASS; "
        + BuiltinExternCache.statistics());
  }
}

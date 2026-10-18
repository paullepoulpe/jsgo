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
  public static void main(String[] args) {
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
      // Source bindings in type nodes are deliberately refreshed, so compare
      // the first bound clone against subsequent independently mutated clones.
      if (cachedFingerprint == null) { cachedFingerprint = got; }
      else if (!cachedFingerprint.equals(got)) { throw new AssertionError("mutable AST/JSDoc leaked"); }
      for (Node child = tree.getFirstChild(); child != null; child = child.getNext()) {
        var doc = child.getJSDocInfo();
        if (doc != null) {
          for (Node type : doc.getTypeNodes()) { type.setLineno(999); type.detachChildren(); }
          for (Node transform : doc.getTypeTransformations().values()) { transform.setLineno(777); transform.detachChildren(); }
        }
      }
      tree.detachChildren();
    }
    if (!cachedFingerprint.contains("inline=true")) { throw new AssertionError("INLINE_TYPE was lost"); }
    if (fingerprint(ParserRunner.parse(file, "var different;", config, REPORTER).ast).equals(cachedFingerprint)) {
      throw new AssertionError("source change reused an old AST");
    }
    System.out.println("cache isolation, inline types, TTL trees and source rebinding PASS; "
        + BuiltinExternCache.statistics());
  }
}

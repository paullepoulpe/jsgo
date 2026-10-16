/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
(/** @closureUnaware */ function() {
// This nested usage of the @closureUnaware annotation is not a parse error, but it
// could be confusing if a human author expected there to be multiple shadowed ASTs /
// sections of closure-aware and closure-unaware code.
// This test-case is just codifying that it is always the top-most AST node that is
// converted into a shadow AST node, and that there won't be nested shadow ASTs.
  (/** @closureUnaware */ function() {
    window['foo'] = 10;
  }).call(globalThis);
}).call(globalThis);

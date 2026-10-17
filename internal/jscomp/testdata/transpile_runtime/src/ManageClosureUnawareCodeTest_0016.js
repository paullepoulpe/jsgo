/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (goog.DEBUG) {
  /** @closureUnaware */
  (function() {
    function bar() {
      window['foo'] = 5;
    }
    bar();
  }).call(globalThis);
}

/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (goog.DEBUG) {
  /** @closureUnaware */
  (function() {
    (function() {
      window['foo'] = 10;
    }).call(globalThis);
  }).call(globalThis);
}

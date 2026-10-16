/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (!window['foo']) {
  /** @closureUnaware */
  (function() {
    window['foo'] = 5;
  }).call(globalThis);
}

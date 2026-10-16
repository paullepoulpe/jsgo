/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
if (!window['foo']) {
  if (goog.DEBUG) {
    /** @closureUnaware */
    (function() {
      window['foo'] = 5;
    }).call(globalThis);
  } else {
    /** @closureUnaware */
    (function() {
       window['foo'] = 10;
    }).call(globalThis);
  }
}

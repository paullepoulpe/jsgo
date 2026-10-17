/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
/** @closureUnaware */
(function() {
  /** @dependency */
  window['foo'] = 5;
}).call(globalThis);

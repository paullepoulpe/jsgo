/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
/** @closureUnaware */
(function() {
    window['foo'] = 5;
  }).call(globalThis);

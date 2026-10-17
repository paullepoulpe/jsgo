/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('foo.bar.baz_raw');
goog.require('foo.bar');
const {a} = goog.require('foo.baz');
/** @closureUnaware */
(function() {
  window['foo'] = 5;
}).call(globalThis);
exports = globalThis['foo'];

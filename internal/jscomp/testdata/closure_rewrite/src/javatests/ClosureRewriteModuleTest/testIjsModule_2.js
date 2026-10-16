goog.module('ns.a');
var b = goog.require('external1');
var c = goog.requireType('external2');
/** @type {b} */ new b;
/** @typedef {c} */ var d;

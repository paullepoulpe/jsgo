goog.module('p.B');
var A = goog.require('p.A');
/** @constructor @extends {A} */
function B() {}
B.prototype = new A;
exports = B;

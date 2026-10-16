goog.module('my.mod');
const {b, a} = goog.require('x.y');
const z = goog.require('a.b');
class C {}
/** @enum {number} */
C.E = {A: 1};
/** @typedef {number} */
C.T;
exports = C;

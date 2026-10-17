goog.module('nested');
class C {}
let obj = {};
/** @enum {number} */
C.E = {A: 1};
/** @interface */
obj.I = class {};
C.D = class {};
/** @typedef {number} */
C.T;
C.F = function() {};
obj.S = '';
C.prototype.m = function() {};
exports.SomeEnum = {};
function f() { /** @enum {number} */ C.Inner = {}; }

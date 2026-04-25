/** @param {!Array} x */ function g(x) {}
/** @this {goog.MyTypedef} */ function f() { g(this); }
var goog = {};
/** @typedef {(Array|null|undefined)} */ goog.MyTypedef;

/** @param {!Array} x */ function g(x) {}
/** @this {goog.MyTypedef} */ function f() { g(this); }
var goog = {};
/** @typedef {(RegExp|null|undefined)} */ goog.MyTypedef;

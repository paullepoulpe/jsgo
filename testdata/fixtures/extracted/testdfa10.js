/** @param {null} x */ function g(x) {}
/** @param {string?} x */function f(x) {
if (!x) { x = ''; }
if (g(x)) { return 0; } else { return 1; } }

/** @param {number} x */ function f(x) {}
/** @param {...?} x */ function g(x) {}
g(function(y) { f(y); }, true);

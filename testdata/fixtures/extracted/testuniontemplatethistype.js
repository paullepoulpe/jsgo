/** @constructor */ function F() {}
/** @return {F|Array} */ function g() { return []; }
/** @param {F} x */ function h(x) { }
/**
* @param {T} x
* @param {function(this:T, ...)} y
* @template T
*/
function f(x, y) {}
f(g(), function() { h(this); });

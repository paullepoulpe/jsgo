/**
 * @constructor
 * @param {T} x
 * @template T
 */
function C(x) {}
/**
 * @template {number|string} U
 * @param {C<U>} x
 * @param {U} y
 */
function f(x,y) {}
var /** C<number|string> */ c;
f(c,0);

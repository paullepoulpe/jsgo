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
var /** C<string> */ c;
f(c,false);

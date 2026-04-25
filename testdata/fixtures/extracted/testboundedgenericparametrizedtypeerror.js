/**
 * @constructor
 * @param {T} x
 * @template T
 */
function C(x) {}
/**
 * @template {number|string} U
 * @param {C<U>} x
 */
function f(x) {}
var /** C<boolean> */ c;
f(c);

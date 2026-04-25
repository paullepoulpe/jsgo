/**
 * @constructor
 * @template T
 */
function C(x) {}
/**
 * @template {number|string} U
 * @param {C<boolean>} x
 * @return {C<U>}
 */
function f(x) { return x; }

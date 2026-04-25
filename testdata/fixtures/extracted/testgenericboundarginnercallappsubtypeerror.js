/**
 * @param {number} x
 */
function numID(x) { return x; }
/**
 * @param {T} x
 * @template {number|boolean|string} T
 */
function foo(x) { return numID(x); }

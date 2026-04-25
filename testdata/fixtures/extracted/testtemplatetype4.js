/**
 * @param {...T} p
 * @return {T}
 * @template T
 */
function fn(p) { return p; }
/** @type {!Object} */ var x;
x = fn(3, null);

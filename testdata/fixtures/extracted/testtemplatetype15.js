/**
 * @param {{foo:T}} p
 * @return {T}
 * @template T
 */
function fn(p) { return p.foo; }
/** @type {!Object} */ var x;
x = fn({foo:3});

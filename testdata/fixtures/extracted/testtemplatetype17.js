/** @constructor */ function C() {}
C.prototype.foo = 1;
/**
 * @param {{foo:T}} p
 * @return {T}
 * @template T
 */
function fn(p) { return p.foo; }
/** @type {!Object} */ var x;
x = fn(new C());

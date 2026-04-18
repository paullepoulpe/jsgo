/** @constructor */ function C() {
  /** @type {number} */ this.foo = 1
}
/**
 * @param {{foo:T}} p
 * @return {T}
 * @template T
 */
function fn(p) { return p.foo; }
/** @type {!Object} */ var x;
x = fn(new C());

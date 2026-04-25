/** @constructor @implements {Iterable<VALUE>} @template VALUE */
function Foo() {
  /** @type {!Map<VALUE, VALUE>} */ this.map = new Map;
}
/** @override @return {!Iterator<VALUE>} */
Foo.prototype[Symbol.iterator] = function() {};

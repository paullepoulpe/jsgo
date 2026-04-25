/** @constructor
 * @extends {Base} */ function Sub() {}
/** @return {number} */ function f() { return (new Sub()).foo; }
/** @constructor */ function Base() {}
/** @type {boolean} */ Base.prototype.foo = true;

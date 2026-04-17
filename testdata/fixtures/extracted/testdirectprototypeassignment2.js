/** @constructor */ function Base() {}
/** @constructor
 * @extends {Base} */ function A() {}
A.prototype = new Base();
A.prototype.foo = 3;
/** @return {string} */ function foo() { return (new Base).foo; }

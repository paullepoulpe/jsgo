/** @constructor */ function Base() {}
Base.prototype.foo = 3;
/** @constructor
 * @extends {Base} */ function A() {}
A.prototype = new Base();
/** @return {string} */ function foo() { return (new A).foo; }

function f(/** !IArrayLike<string> */ x){};
/** @constructor */
function Foo() {}
/** @type {number} */ Foo.prototype.length
f(new Foo)

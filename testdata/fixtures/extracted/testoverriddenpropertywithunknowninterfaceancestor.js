/** @interface @extends {Unknown}  */ function Foo() {}
/** @type {string} */ Foo.prototype.data;
/** @constructor @implements {Foo} */ function SubFoo() {}
/** @type {string|Object} */
SubFoo.prototype.data = null;

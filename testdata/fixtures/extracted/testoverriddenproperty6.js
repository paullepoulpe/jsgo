/** @constructor */ function Foo() {}
/** @type {?number} */ Foo.prototype.bar = null;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
SubFoo.prototype.bar = 3;

/** @constructor */ function Foo() {}
/** @type {?number} */ Foo.prototype.bar = null;

/**
 * @constructor
 * @extends {Foo}
 */
function SubFoo() {}
/** @type {?} */
SubFoo.prototype.bar = 'not a number';

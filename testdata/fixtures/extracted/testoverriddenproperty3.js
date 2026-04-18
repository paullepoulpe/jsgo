/** @constructor */ function Foo() {
}
/** @type {string} */ Foo.prototype.data;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/** @type {string|Object}
 * @override */
SubFoo.prototype.data = null;

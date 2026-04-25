/** @constructor */ function Foo() {}
/** @type {Object} */
Foo.prototype.bar = {};
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @type {Array}
 * @override
 */
SubFoo.prototype.bar = [];

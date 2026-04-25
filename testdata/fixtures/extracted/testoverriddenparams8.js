/** @constructor
 * @template T */ function Foo() {}
/** @param {T} x */
Foo.prototype.bar = function(x) { };
/**
 * @constructor
 * @extends {Foo<string>}
 */ function SubFoo() {}
/**
 * @param {number} x
 * @override
 */
SubFoo.prototype.bar = function(x) {};

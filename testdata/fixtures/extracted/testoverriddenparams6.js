/** @constructor */ function Foo() {}
/** @param {number} x */
Foo.prototype.bar = function(x) { };
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @override
 */
SubFoo.prototype.bar = function() {};
(new SubFoo()).bar(true);

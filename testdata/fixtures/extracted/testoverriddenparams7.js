/** @interface */ function Foo() {}
/** @param {number} x */
Foo.prototype.bar = function(x) { };
/**
 * @interface
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @override
 */
SubFoo.prototype.bar = function() {};
var subFoo = /** @type {SubFoo} */ (null);
subFoo.bar(true);

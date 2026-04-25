/** @constructor */ function Foo() {}
/** @param {...number} var_args */
Foo.prototype.bar = function(var_args) { };
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @param {number} x
 * @override
 */
SubFoo.prototype.bar = function(x) {};

/** @constructor */ function Foo() {}
/** @param {...?} var_args */
Foo.prototype.bar = function(var_args) {};
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @param {number} x
 * @override
 */
SubFoo.prototype.bar = function(x) {};

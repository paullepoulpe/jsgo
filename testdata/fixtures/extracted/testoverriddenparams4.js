/** @constructor */ function Foo() {}
/** @type {function(...number)} */
Foo.prototype.bar = function(var_args) {};
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @type {function(number)}
 * @override
 */
SubFoo.prototype.bar = function(x) {};

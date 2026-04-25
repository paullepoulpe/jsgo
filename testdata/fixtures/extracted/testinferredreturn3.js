/** @constructor */ function Foo() {}
Foo.prototype.bar = function() {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @return {number}
 * @override  */
SubFoo.prototype.bar = function() { return 3; };

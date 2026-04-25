/** @constructor */
function Foo() {}
/** @return {number} */
Foo.prototype.foo = function() {};

/** @constructor */ var FooAlias = Foo;

var x = /** @type {!FooAlias} */ ({});
var /** null */ n = x.foo();

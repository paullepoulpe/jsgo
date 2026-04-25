/** @constructor */
function Foo() {}
/** @return {number} */
Foo.prototype.foo = function() {};

/** @const */ var FooAlias = Foo;

var x = /** @type {!FooAlias} */ ({});
var /** null */ n = x.foo();

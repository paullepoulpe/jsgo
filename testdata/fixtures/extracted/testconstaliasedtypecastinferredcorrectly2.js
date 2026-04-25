/** @constructor */
function Foo() {}
/** @return {number} */
Foo.prototype.foo = function() {};

var ns = {};
/** @const */ ns.FooAlias = Foo;

var x = /** @type {!ns.FooAlias} */ ({});
var /** null */ n = x.foo();

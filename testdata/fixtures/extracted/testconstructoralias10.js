/**
 * @param {number} x
 * @constructor */
var Foo = function(x) {};
/** @constructor */ var FooAlias = Foo;
/** @return {number} */ function foo() {
  return new FooAlias(1); }

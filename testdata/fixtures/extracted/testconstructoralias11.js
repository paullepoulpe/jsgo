/**
 * @param {number} x
 * @constructor */
var Foo = function(x) {};
/** @const */ var FooAlias = Foo;
/** @const */ var FooAlias2 = FooAlias;
/** @return {FooAlias2} */ function foo() {
  return 1; }

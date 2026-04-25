/** @constructor */ var Foo = function() {};
/** @constructor */ var FooAlias = Foo;
/** @return {Foo} */ function foo() {
  return new FooAlias(); }

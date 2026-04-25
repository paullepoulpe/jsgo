/** @constructor */ var Foo = function() {};
/** @constructor */ var FooAlias = Foo;
/** @return {FooAlias} */ function foo() {
  return new Foo(); }

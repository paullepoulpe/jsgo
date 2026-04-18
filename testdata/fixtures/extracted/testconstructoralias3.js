/** @constructor */ var Foo = function() {};
/** @type {number} */ Foo.prototype.bar = 3;
/** @constructor */ var FooAlias = Foo;
/** @return {string} */ function foo() {
  return (new FooAlias()).bar; }

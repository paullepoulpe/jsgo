/** @constructor */ var Foo = function() {};
var FooAlias = Foo;
/** @type {number} */ FooAlias.prototype.bar = 3;
/** @return {string} */ function foo() {
  return (new Foo()).bar; }

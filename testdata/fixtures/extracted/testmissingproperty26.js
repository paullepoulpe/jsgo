/** @constructor */ var Foo = function() {};
/** @constructor */ var FooAlias = Foo;
FooAlias.prototype.bar = function() {};
(new Foo()).bar();

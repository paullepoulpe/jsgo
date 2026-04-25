/** @constructor */ var Foo = function() {};
Foo.prototype.bar = function() {};
/** @constructor */ var FooAlias = Foo;
(new FooAlias()).bar();

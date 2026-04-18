/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
(new Foo).b = 3;

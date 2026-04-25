/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
(new Foo).a = 3;

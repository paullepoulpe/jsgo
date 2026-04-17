/** @constructor */ function Foo() {}
/** @param {number} x */ Foo.prototype.bar = function(x) {};
Foo.prototype.bar = function(x) { return 3; };

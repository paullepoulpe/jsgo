/** @constructor */ function Foo() {}
/** @param {number} x */ Foo.prototype.bar = function(x) {};
/** @param {string} x */ function f(x) {}
Foo.prototype.bar = function(y) { f(y); };

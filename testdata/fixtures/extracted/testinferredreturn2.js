/** @constructor */ function Foo() {}
Foo.prototype.bar = function() {};
/** @param {number} x */ function g(x) {}
g((new Foo()).bar());

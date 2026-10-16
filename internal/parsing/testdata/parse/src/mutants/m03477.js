/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
/** @param {Bar|Baz} x */ function g(x) { return x.specialProp; }
/** @param {Foo|Bar} x */ function f(x) { x.specialProp = 1; }
/** @constructor */ function Baz() {}

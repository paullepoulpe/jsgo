/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
/** @constructor */ function Baz() {}
/** @param {Foo|Bar} x */ function f(x) { x.specialProp = 1; }
/** @param {Bar|Baz} x */ function g(x) { return x.specialProp; }

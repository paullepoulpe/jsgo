function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends */ ({bar: 1}));

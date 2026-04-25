function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends {Foo.prototype|Foo} */ ({bar: 1}));

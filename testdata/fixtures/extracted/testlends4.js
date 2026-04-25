function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends {Foo} */ ({bar: 1}));
alert(Foo.bar);

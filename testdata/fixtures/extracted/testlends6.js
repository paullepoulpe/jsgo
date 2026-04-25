function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends {Foo.prototype} */ ({bar: 1}));
alert((new Foo()).bar);

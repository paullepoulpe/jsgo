function extend(x, y) {}
/** @type {number} */ var Foo = 3;
extend(Foo, /** @lends {Foo} */ ({bar: 1}));

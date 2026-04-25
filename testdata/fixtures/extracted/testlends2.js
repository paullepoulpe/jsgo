function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends {Foob} */ ({bar: 1}));

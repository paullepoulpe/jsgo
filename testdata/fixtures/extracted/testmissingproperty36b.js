/** @constructor */ function Foo() {}
Foo.prototype.baz = 0;
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
SubFoo.prototype.bar = 0;
/** @param {{bar: number}} x */ function f(x) { return x.baz; }

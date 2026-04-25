/** @record */
function Foo() {}
/** @type {number} */
Foo.prototype.a;
function f(/** !Foo */ x) {
  var y = x;
  /** @type {number} */
  y.b = 123;
}
function g(/** !Foo */ x) {
  var /** null */ n = x.b;
}

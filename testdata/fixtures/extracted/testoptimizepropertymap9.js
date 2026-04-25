/** @constructor */
function Foo() {
  this.a = 123;
}
function f(/** {a: number, c: string} */ x) {
  x.b = 234;
}
function g(/** !Foo */ x) {
  return x.b + 5;
}

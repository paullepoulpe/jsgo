/**
 * @param {T} x
 * @template {boolean} T
 */
function foo(x) { return x; }
class Foo {
  /**
   * @param {T} x
   * @template {string} T
   */
  foo(x) { return x; }
}
var F = new Foo();
F.foo('Hi');
foo(true);

/**
 * @template {string} T
 */
class Foo {
  /**
   * @param {T} x
   * @template {number} T
   */
  foo(x) { return x; }
}
const /** !Foo<string> */ f = new Foo();
f.foo(3);

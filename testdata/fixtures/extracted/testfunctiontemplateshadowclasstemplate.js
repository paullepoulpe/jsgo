/**
 * @template T
 */
class Foo {
  /**
   * @param {T} x
   * @template T
   */
  foo(x) { return x; }
}
const /** !Foo<string> */ f = new Foo();
f.foo(3);

/**
 * @template {number|string} T
 * @template {number|null} S
 */
class Foo {
  /**
   * @param {T} x
   */
  bar(x) { }
  /**
   * @param {S} x
   */
  baz(x) { }
}
var /** Foo<number> */ foo;
foo.baz(3);
foo.baz(null);
foo.baz('3');

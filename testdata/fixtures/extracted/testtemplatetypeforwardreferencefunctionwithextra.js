/** @param {!Foo<string, boolean>} x */
function f(x) {}

/**
 * @constructor
 * @template T
 */
function Foo() {}

/** @param {!Foo<number>} x */
function g(x) {
  f(x);
}

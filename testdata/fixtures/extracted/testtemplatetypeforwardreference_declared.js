/** @param {!Foo<string>} x */
function f(x) {}

/**
 * @template T
 * @constructor
 */
function Foo() {}

/** @param {!Foo<number>} x */
function g(x) {
  f(x);
}

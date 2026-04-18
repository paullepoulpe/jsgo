/** @param {!Foo<string>} x */
function f(x) {}

/**
 * @template T
 * @constructor
 */
var Foo = function() {}

/** @param {!Foo<number>} x */
function g(x) {
  f(x);
}

class A {
  foo() {}
}
/** @param {?A} a */
function f(a) {
// TypeCheck should not complain that foo() is getting called without a correctly typed
// `this` value.
  a?.foo();
}

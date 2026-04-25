/** @constructor */ function Foo() {}
/** @type {number?} */ Foo.prototype.bar = null;
/** @type {number?} */ Foo.prototype.baz = null;
/** @param {Foo} foo */
function f(foo) {
  while (true) {
    if (!foo.baz) break;
    foo.bar = null;
  }
  // Tests a bug where this condition always evaluated to true.
  return foo.bar == null;
}

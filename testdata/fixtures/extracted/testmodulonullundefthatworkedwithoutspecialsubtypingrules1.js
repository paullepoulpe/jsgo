/** @constructor */
function Foo() {}
function f(/** function(?Foo, !Foo) */ x) {
  return /** @type {function(!Foo, ?Foo)} */ (x);
}

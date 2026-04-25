/** @constructor @template A, B */
function Foo(/** A */ a, /** B */ b) {}
/** @template T */
function f(/** T */ t) {
  /** @template S */
  function g(/** S */ s) {
    var /** !Foo<T, S> */ foo = new Foo(t, s);
  }
}

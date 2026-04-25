/** @constructor */
function Foo() { }

/**
 * @template {!Foo} X
 * @param {X} x
 * @return {X}
 */
function identity(x) {
  return x;
}

var /** function(!Foo): !Foo */ y = identity;

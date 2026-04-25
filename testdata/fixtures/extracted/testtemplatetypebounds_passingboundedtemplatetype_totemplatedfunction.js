/** @constructor */
function Foo() { }

/**
 * @template {!Foo} X
 * @param {X} x
 * @return {X}
 */
function clone(x) {
// The focus of this test is that `X` (already a template variable) is bound to `Y` at
// this callsite. We confirm that by ensuring an `X` is returned from `cloneInternal`.
  return cloneInternal(x);
}

/**
 * @template {!Foo} Y
 * @param {Y} x
 * @return {Y}
 */
function cloneInternal(x) {
  return x;
}

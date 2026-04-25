class Foo {
  takesString(/** string */ stringParameter) {}
}
/** @param {*} all */
function f(all) {
  /** some jsdoc */
  (/** @type {!Foo} */ (all)).takesString(0);
}

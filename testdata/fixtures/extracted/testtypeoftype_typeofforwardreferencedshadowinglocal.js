/** @const {string} */
var x = 'x';
function f() {
  function g(/** typeof x */ a) {}
  /** @const {number} */
  var x = 1;
  g(null);
}

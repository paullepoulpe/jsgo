/** @return {number} */
function f() {
  var x = 3;
  /** @return {number} */
  function g() { var x = 3;x = true; return x; }
  return x;
}

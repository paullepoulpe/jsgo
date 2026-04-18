var x = null;
/** @return {number} */
function f() {
  x = '3';
  /** @return {number} */
  function g() { x = 3; return x; }
  return x;
}

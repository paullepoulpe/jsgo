var x = null;
/** @return {number} */
function f() {
  x = 3;
  /** @return {number} */
  function g() { x = true; return x; }
  return x;
}

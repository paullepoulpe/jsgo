/** @return {boolean} */
function f() {
  /** @type {?string} */ var x = null;
  function g() { x = 'y'; } g();
  return x == null;
}

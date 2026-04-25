function f() {
  /** @type {number} */ var x = 3;
  function g() { x = null; }
  return x;
}

function f() {
  /** @type {?number} */ var x = null;
  /** @return {string} */
  function g() {
    if (!x) {
      x = 1;
    }
    return x;
  }
}

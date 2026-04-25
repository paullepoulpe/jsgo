/** @param {?number} x */
function f(x) {
  if (x) {
    /** @const */ var y = x;
    /** @return {boolean} */ function g() { return y; }
  }
}

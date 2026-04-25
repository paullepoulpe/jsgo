function f() {
  /** @type {?number} */ var x = null;
  /** @return {number} */
  function g() {
    x = 1;
    h();
    return x;
  }
  function h() {
    x = null;
  }
}

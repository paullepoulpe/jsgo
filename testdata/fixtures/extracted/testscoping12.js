/** @constructor */ function F() {}
/** @type {number} */ F.prototype.bar = 3;
/** @param {!F} f */ function g(f) {
  /** @return {string} */
  function h() {
    return f.bar;
  }
}

function fn() {
  /** @type {number} */
  var out = 0;
  try {
    foo();
  } catch (/** @type {string} */ e) {
    out = e;
  }
}

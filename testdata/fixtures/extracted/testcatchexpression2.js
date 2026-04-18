function fn() {
  /** @type {number} */
  var out = 0;
  /** @type {string} */
  var e;
  try {
    foo();
  } catch (e) {
    out = e;
  }
}

/** @param {number} x */ function foo(x) {}
/** @param {{length:*}} y */ function bar(y) {
  /** @type {string} */ y.length;
  foo(y.length);
}

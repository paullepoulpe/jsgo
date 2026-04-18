var obj = {};
/**
 * @param {string} x
 * @return {Object} */ function f(x) { return obj; }
/** @param {?number} x */ function g(x) {
  if (x) { delete f(x)['a']; }
}

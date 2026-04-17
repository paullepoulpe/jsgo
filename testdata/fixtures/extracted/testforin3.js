/** @param {boolean} x */ function f(x) {}
/** @type {Object<number>} */ var obj = {};
for (var k in obj) {
  f(obj[k]);
}

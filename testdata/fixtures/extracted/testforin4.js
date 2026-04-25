/** @param {boolean} x */ function f(x) {}
/** @enum {string} */ var E = {FOO: 'bar'};
/** @type {Object<E, Array>} */ var obj = {};
for (var k in obj) {
  f(obj[k]);
}

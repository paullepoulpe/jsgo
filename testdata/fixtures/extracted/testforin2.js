/** @param {boolean} x */ function f(x) {}
/** @enum {string} */ var E = {FOO: 'bar'};
/** @type {Object<E, string>} */ var obj = {};
var k = null;
for (k in obj) {
  f(k);
}

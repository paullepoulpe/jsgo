/**
 * @param {string} x
 * @param {number} y
 * @param {number} z
 */
function g(x, y, z) {}
function f() {
  var x = 'a'; g(x, x = 3, x);
}

/**
 * @param {number} x
 * @return {number}
 */
function g(x) { return 5; }
function f() {
  var x = {};
  x.foo = '3';
  return g(x.foo) in {};
}

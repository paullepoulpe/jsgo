/**
 * @param {number} x
 * @return {number}
 */
function g(x) { return 5; }
function f() {
  var x = {};
  x.foo = '3';
  switch (3) { case g(x.foo): return 3; }
}

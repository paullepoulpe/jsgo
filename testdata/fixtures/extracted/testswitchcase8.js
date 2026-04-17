/**
 * @param {number} x
 * @return {number}
 */
function g(x) { return 5; }
function f() {
  var x = {};
  x.foo = '3';
  switch (g(x.foo)) { case 3: return 3; }
}

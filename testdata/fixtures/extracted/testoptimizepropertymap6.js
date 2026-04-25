function f(/** {a:number} */ x) {
  var y = x;
  /** @type {number} */
  y.b = 123;
}
function g(/** {c:number} */ x) {
  var /** null */ n = x.b;
}

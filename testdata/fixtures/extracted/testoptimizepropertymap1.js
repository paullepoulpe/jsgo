/** @return {!Function} */
function f() {
  var x = function() {};
  /** @type {number} */
  x.prop = 123;
  return x;
}
function g(/** !Function */ x) {
  var /** null */ n = x.prop;
}

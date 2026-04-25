/** @return {string} */
function f() {
  var x = function() {};
  x = /** @type {function(): number} */ (function() { return 3; });
  return x();
}

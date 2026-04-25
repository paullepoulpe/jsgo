/** @return {string} */
function f() {
  var x = function() {};
  if (f())
    x = /** @type {function(): number} */
        (function() { return 3; });
  return x();
}

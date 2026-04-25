/** @return {number} */
function f() {
  var x = null;
  function g() { x = 3; }
  g();
  return x;
}

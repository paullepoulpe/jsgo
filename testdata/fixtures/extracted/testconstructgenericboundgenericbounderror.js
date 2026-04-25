/**
 * @param {T} x
 * @param {U} y
 * @template {boolean|string} T
 * @template {T} U
 */
function foo(x,y) {
  var /** U */ z = x;
  x = y;
}

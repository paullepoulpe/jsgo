/**
 * @param {!Array<T>} x
 * @param {!Array<number|boolean>} y
 * @template {number} T
 */
function foo(x,y) {
  var /** !Array<T> */ z = y;
  y = x;
}

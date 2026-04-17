/**
 * @param {number|null} x
 * @return {number}
 */
function f(x) {
  return typeof x == 'object' ? 1 : x;
}

/**
 * @param {!Object} o
 * @param {number?} x
 * @return {string}
 */
function f(o, x) {
  if (x) {
    o.x = x;
  } else {
    o.x = 'a';
  }
  return o.x;
}

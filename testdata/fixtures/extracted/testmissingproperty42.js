/** @param {Object} x */
function f(x) {
  if (typeof x.impossible == 'undefined') throw Error();
  return x.impossible;
}

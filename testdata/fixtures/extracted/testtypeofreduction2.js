/**
 * @param {string|number} x
 * @return {string}
 */
function f(x) {
  return typeof x != 'string' ? String(x) : x;
}

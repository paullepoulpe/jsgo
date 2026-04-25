/** @param {Object} a */
function f(a) {
  a.prototype = '__proto'
}
/** @param {Object} b
 *  @return {!Object}
 */
function g(b) {
  return b.prototype
}

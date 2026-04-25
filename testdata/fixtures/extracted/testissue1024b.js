/** @param {Object} a */
function f(a) {
  a.prototype = {foo:3};
}
/** @param {Object} b
 */
function g(b) {
  b.prototype = function(){};
}

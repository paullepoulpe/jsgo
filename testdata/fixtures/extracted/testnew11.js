/** @param {Function} c1 */
function f(c1) {
  var c2 = function(){};
  c1.prototype = new c2;
}

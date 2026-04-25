/** @param {string} x
 * @constructor
 */
var Bar = function(x) {};
/** @param {string} x */ function g(x) { return true; }
/** @param {string|number} opt_x */
function f(opt_x) {
  if (opt_x) { new Bar(g(opt_x) && 'x'); }
}

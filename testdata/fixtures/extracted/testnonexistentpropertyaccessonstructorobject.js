/**
 * @constructor
 * @struct
 */
var A = function() {};
/** @param {A|Object} a */
function foo(a) {
  if (a.bar) { a.bar(); }
}

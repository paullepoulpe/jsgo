goog.module('a');
/** @type {z} */ var c;
var z = goog.forwardDeclare('x.y.z');
function f() {
  z = goog.module.get('x.y.z');
  c = new z;
}

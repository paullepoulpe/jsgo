goog.module('a');
var x = goog.forwardDeclare('x.y.z');
function f() {
  x = goog.module.get('x.y.z');
  new x;
}

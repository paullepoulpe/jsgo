goog.module('a');
var z = goog.forwardDeclare('x.y.z');
function f() {
  z = goog.module.get('x.y.z');
  new z;
}

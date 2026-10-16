goog.module('a');
var x = goog.forwardDeclare('b');
function f() {
  x = goog.module.get('b');
  new x;
}

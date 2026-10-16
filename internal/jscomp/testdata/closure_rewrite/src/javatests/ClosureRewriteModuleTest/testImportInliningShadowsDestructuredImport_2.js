goog.module('a.b.d');
const {d} = goog.requireType('a.b.c');
function foo() {
  var a = 10;
  var b = d;
}

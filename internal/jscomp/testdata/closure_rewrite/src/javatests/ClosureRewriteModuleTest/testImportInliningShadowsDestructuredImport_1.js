goog.module('a.b.d');
const {d} = goog.require('a.b.c');
function foo() {
  var a = 10;
  var b = d;
}

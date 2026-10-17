goog.module('x.y.z');

function f() {
  return new (goog.module.get('a.b.c'))();
}

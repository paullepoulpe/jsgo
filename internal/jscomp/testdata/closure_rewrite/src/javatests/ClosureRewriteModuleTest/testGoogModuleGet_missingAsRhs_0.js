goog.module('x.y.z');

function f() {
  const c = goog.module.get('a.b.c');
  return c;
}
f();

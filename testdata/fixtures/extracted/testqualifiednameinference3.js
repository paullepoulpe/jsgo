var x = {};
x.y = c;
function f(a, b) {
  if (a) {
    if (b)
      x.y = 2;
    else
      x.y = 1;
  }
  return x.y == null;
} function g() { x.y = null; }

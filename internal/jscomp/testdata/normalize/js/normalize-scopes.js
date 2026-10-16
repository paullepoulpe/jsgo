var x = 1;
function f(x) {
  var y = x;
  { let x = 2; const y2 = x; }
  try { throw x; } catch (x) { var y = x; }
  for (let x = 0; x < 2; x++) { let y = x; }
  switch (x) { case 1: let z = 1; break; default: let z2 = 2; }
  return function g(x) { return x + y; };
}
class D { constructor(x) { this.x = x; } static m(x) { return x; } }
label: { var inLabel = 1; }
var q = function q() { return q; };
(function iife(a) { var a; })(1);
var { [computed]: cv } = obj;
for (const [k, v] of Object.entries(obj)) { var acc = k + v; }

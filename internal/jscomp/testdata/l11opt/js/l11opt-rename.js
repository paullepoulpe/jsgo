var globalCounter = 0;
function outer(alpha, beta) {
  var local1 = alpha + beta;
  var f = function inner(x) { return x > 0 ? inner(x - 1) : local1; };
  label1: for (var i = 0; i < 3; i++) {
    label2: for (var j = 0; j < 3; j++) {
      if (j == 1) continue label1;
      if (i == 2) break label2;
      globalCounter++;
    }
  }
  unused: { globalCounter += f(alpha); }
  return local1;
}
function second(a, b, c) { let q = a; const r = b; { let q = c; r + q; } return q; }
try { outer(1, 2); } catch (err) { second(err, err, err); }
var arguments_like = function() { return arguments[0]; };
outer(second(1, 2, 3), arguments_like(4));

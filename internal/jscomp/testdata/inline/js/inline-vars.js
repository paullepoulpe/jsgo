var a = 1; var b = a; use(b);
function f(x) { var y = x; var z = y + 1; return z; }
function g() { const k = 'str'; use(k, k); let m = {}; use(m); }
function h(arg) { const argAlias = arg; const argAliasAlias = argAlias; use(argAliasAlias); }
function args() { var p = arguments[0]; arguments[1] = 2; return p; }
var neverAssigned; use(neverAssigned, neverAssigned);
function fdecl() { function inner() { return 1; } var q = inner; return q(); }
/** @noinline */ var noinl = 3; use(noinl);
var THIS_IS_CONST = 5; use(THIS_IS_CONST + THIS_IS_CONST);
function loops() { for (var i = 0; i < 3; i++) { var t = i; use(t); } }
function gp() { var o = window.location; o(); var p2 = window.foo; use(p2); }

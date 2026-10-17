function pure(a) { return a + 1; }
function impure() { alert('x'); }
function mutatesArg(o) { o.x = 1; }
function throws() { throw new Error('e'); }
function callsPure() { return pure(1); }
function callsImpure() { impure(); }
var arrow = (x) => x * 2;
class C { m() { return 1; } n() { this.y = 2; } }
pure(1);
impure();
mutatesArg({});
callsPure();
callsImpure();
arrow(3);
new C().m();
new C().n();
externPure();
alert(throws);
var tag = function(s) { return s; };
tag`hello`;
window.f?.();
var x = 0;
function incX() { x++; }
incX();

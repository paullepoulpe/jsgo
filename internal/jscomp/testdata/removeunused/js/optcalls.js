function unusedParams(a, b, c) { return a; }
unusedParams(1, 2, 3);
unusedParams(4, 5, 6);
function constArg(x, y) { return x + y; }
constArg(1, 2);
constArg(1, 3);
function unusedReturn() { alert(1); return 5; }
unusedReturn();
function usedReturn() { return 6; }
alert(usedReturn());
class K { constructor(a, b) { this.a = a; } }
new K(1, 2);
new K(3, 4);
/** @constructor */ function A() {}
A.prototype.foo = function() { return this.bar; };
A.prototype.baz = function(z) { alert(z); };
var aa = new A();
aa.foo();
aa.baz(1);
function emptyFn() {}
function identity(v) { return v; }
function returnsConst() { return 42; }
function getter() { return this.prop; }
function setter(v) { this.prop = v; }
var o = {emptyFn: emptyFn, identity: identity, returnsConst: returnsConst};
emptyFn(); identity(1); returnsConst();

/** @constructor */
function Foo() { this.a = 1; this.a = 2; this.b = 3; }
Foo.prototype.m1 = function() { return this.a; };
Foo.prototype.m2 = function() { return this.b; };
Foo.prototype.m3 = function() { return this.m1(); };
Foo.prototype.m4 = function() { return this.m2(); };
Foo.prototype.m5 = 5;
Foo.prototype.m6 = function() { return Foo.prototype.m5; };
Foo.prototype.m7 = function(x) { x.a = 1; x.a = 2; return x; };
Foo.prototype.m8 = function() {};
Foo.staticMethod = function() { return 'static'; };
Foo.staticProp = 42;
/** @constructor @extends {Foo} */
function Bar() { Foo.call(this); }
$jscomp.inherits(Bar, Foo);
Bar.prototype = { n1: function() { return 1; }, 'quoted': 2, n2() { return super.toString; } };
Object.defineProperties(Foo, { getterProp: { get: function() { return 1; } } });
class Baz extends Foo { constructor() { super(); this.c = 1; } method() { return () => super.m1(); } static sm() {} get g() { return 1; } }
var alias = Foo;
alias.aliasProp = 1;
new Bar().m3();

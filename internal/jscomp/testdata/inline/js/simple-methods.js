function A() {} A.prototype.getB = function() { return this.b; };
A.prototype.getC = function() { return this.c.d; };
A.prototype.empty = function() {};
A.prototype.lit = function() { return 3; };
function B() {} B.prototype.getB = function() { return this.b; };
var o = new A(); use(o.getB(), o.getC(), o.lit()); o.empty(); var e = o.empty();
var objlit = { meth: function() { return this.z; } }; use(objlit.meth());
class D { m() { return this.v; } get gg() { return 1; } } use(new D().m());
A.prototype.diff = function() { return this.x; }; B.prototype.diff = function() { return this.y; }; use(o.diff());

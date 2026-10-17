this.x = 1;
var a = this.y;
function f() { this.z = 2; }
/** @constructor */ function C() { this.w = 3; }
C.prototype.m = function() { this.v = 4; };
var o = {m: function() { this.u = 5; }};
(a = this).b = 1;
var arrow = () => this.q;
class K { static { this.s = 1; } f = this.g; }
/** @this {Object} */ function h() { this.t = 1; }
/** @type {function(this:Object)} */ var tf = function() { this.r = 1; };
/** @lends {X.prototype} */ ({ m: function() { this.p = 1; } });
foo(function() { this.n = 1; });

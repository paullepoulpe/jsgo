goog.provide('my.ns.sub');
my.ns.sub.fn = function() { return 1; };
my.ns.sub.val = 2;
var top = {a: {b: {c: 1}}};
top.a.b.d = 2;
function setter() { top.e = 3; }
var x = top.a;
var z = top.a.b.c;
var w = (top.f = 4);
use(top.prototype);
Object.defineProperties(top, {gp: {get: function() { return 1; }}});
for (top.loop in {}) {}
var {pa, pb} = top;
class Cls { static sm() {} }
Cls.prototype.inst = function() {};
var inherits = {}; goog.inherits(Sub, Cls);

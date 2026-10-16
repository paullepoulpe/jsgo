var x = a ** -b;
var y = (-a) ** b;
var z = (a ** b) ** c, w = a ** (b ** c);
var q = a ?? (b || c);
var r = (a && b) ?? c;
var s = (a ?? b) || c;
var t = new (a())();
var u = new (a.b().c)();
var v = new (a?.b)();
var aa = (a?.b).c;
var bb = a?.b.c;
var cc = new a.b.c();
var dd = new a;
var ee = new (void 0)();
var ff = typeof void !a;
var gg = a < !--b;
var hh = a-- > b;
var ii = (a, b);
var jj = (a = b) ? c : (d, e);
var kk = a ? b : c ? d : e;
var ll = (a ? b : c) ? d : e;
for (var i = ("x" in y); i < 1; i++) {}
for (let j = (a in b) ? 1 : 2;;) {}
for (const k of (a, b)) {}
for (x in y) {}
for ((x in y) ? a : b;;) {}
var mm = () => ({});
var nn = () => ({a} = b);
var oo = async () => {};
var pp = async (a, b) => a + b;
var qq = a => b => c;
var rr = (function () {})();
var ss = (class {}).name;
({}).x;
({a} = b);
[a] = b;
(function () {}).call(this);
(async function () {})();
var tt = (a => a)(1);
var uu = (() => {}) || 1;
var vv = typeof (() => 1);
var ww = !function () {}();
var xx = delete a.b, yy = void a;
var zz = a instanceof (b, c);
var ab = (yield_ = 1);
var ac = a && b || c && d;
var ad = a || b && c;
var ae = (a || b) && c;
var af = a + (b + c), ag = (a + b) + c, ah = a - (b - c), ai = a * (b / c);
var aj = a = b = c;
var ak = (a += 1) + (b -= 2);
var al = a.b.c.d(e).f[g](h)?.i?.[j]?.(k);
var am = (a, b) => (c, d);
var an = x => x ? 1 : 2;

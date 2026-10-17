// Constant folding: numbers, bigints, strings, comparisons, shifts.
var a = 1 + 2 * 3 - 4 / 2 % 3;
var b = 2 ** 10, c = 2 ** -14, d = 1.5 ** 2, e = 10 ** 21, f = 2.7182818 ** 2.2;
var g = 7 % -3, h = -7 % 3, i0 = 5.5 % 2, j = 0 / 0, k = 1 / 0, l = -1 / 0;
var m = 1 << 31, n = -1 >>> 0, o = -8 >> 1, p = 1 << 33, q = 3.5 | 0, r = ~5, s = ~-1.5;
var t = 10n + 5n, u = 10n - 20n, v = 3n * 4n, w = 7n / 2n, x = 7n % 3n, y = 2n ** 64n;
var z = 5n & 3n, aa = 5n | 3n, bb = 5n ^ 3n, cc = ~5n, dd = -(-5n), ee = 1n ** -1n;
var ff = "a" + "b" + 1 + 2, gg = 1 + 2 + "c", hh = "x" + null + undefined + true;
var ii = "abc" < "abd", jj = "a" == "a", kk = 1 == "1", ll = null == undefined;
var mm = NaN === NaN, nn = 1n < 2, oo = 3 > 2n, pp = 0n == 0, qq = 9007199254740993n > 9007199254740992;
var rr = typeof 1, ss = typeof "s", tt = typeof function () {}, uu = typeof void 0;
var vv = [1, 2, 3][1], ww = "hello"[1], xx = [1, 2, 3].length, yy = "abc".length;
var zz = {a: 1, b: 2}.a, ab = {get a() { return 1; }}.a, ac = [...[1, 2], 3];
var ad = `a${1}b${"c"}d${true}`, ae = `x` + `y${z}`, af = "p" + `q${z}`;
var ag = x || true, ah = x && false, ai = null ?? 1, aj = 0 ?? 1, ak = (x(), null) ?? 2;
var al = !0, am = !1, an = !!x, ao = -NaN, ap = +"3", aq = -(-4);
var ar = 1 instanceof Object, as = {} instanceof Object;
x = x + 1; x = x * y; x = y + x; x = x - 1 - 2; x += 1; x -= 2;
var at = 0.1 + 0.2, au = 1e300 * 1e10, av = 123456789 * 987654321, aw = 4.35 * 100;
var ax = 1 / 3, ay = 2 / 3 * 3, az = 1e21 + 1, ba = -0 + 0, bb2 = 0 * -1;
Object.defineProperties(o, {});
if ("a" < 1) { foo(); }

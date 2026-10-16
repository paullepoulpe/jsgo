// Known methods.
var a = "abcdef".indexOf("bc"), b = "abcdefbc".lastIndexOf("bc", 4), c = "abc".substr(1, 1);
var d = "abc".substring(1), e = "abc".slice(0, 2), f = "abc".charAt(1), g = "abc".charCodeAt(2);
var h = "a,b,c".split(","), i0 = "abc".split(""), j = "abc".split(), k = "abc".split(",", 0);
var l = "ABC".toLowerCase(), m = "abcß".toUpperCase(), n = " \t abc \n".trim();
var o = "abcabc".replace("b", "x"), p = "abcabc".replaceAll("b", "x");
var q = parseInt("123"), r = parseInt("0x1f"), s = parseInt("11", 2), t = parseFloat("1.50");
var u = parseInt("-0"), v = parseFloat("1e3"), w = Number.parseInt("7"), x = Number.isNaN(NaN);
var y = Number.isFinite(1), z = Number.isSafeInteger(9007199254740991);
var aa = Math.abs(-3), bb = Math.ceil(1.2), cc = Math.floor(-1.2), dd = Math.round(2.5);
var ee = Math.round(-2.5), ff = Math.sign(-3), gg = Math.trunc(-3.7), hh = Math.clz32(1);
var ii = Math.max(1, 2, NaN), jj = Math.min(), kk = Math.imul(3, 4), ll = Math.pow(2, 0.5);
var mm = Math.fround(1.5), nn = Math.fround(1.1), oo = Math.pow(3, -2), pp = Math.max(Infinity, NaN);
var qq = Array.of(1, 2), rr = [1, 2].concat([3]), ss = "𝄞".charAt(0);
var tt = "x".repeat(3), uu = "ab".indexOf("b", -5), vv = "ab".lastIndexOf("", 10);

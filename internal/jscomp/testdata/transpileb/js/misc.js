async function af() { await 1; }
function* gen() { yield 1; }
var {a, b: [c, d = 2], ...rest} = obj;
var arrow = (x) => x * 2;
var sp = [...a, ...b];
var t = `x${a}y`;
var e = a ** 2;
var q = a ?? b;
var oc = a?.b?.c;
a ||= b;
var objs = {...a, b};
try { f(); } catch { g(); }
var o = { get x() { return 1; }, set x(v) {} };
for (const x of [1, 2]) { setTimeout(() => x); }
label: { break label; }

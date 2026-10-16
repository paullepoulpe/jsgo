// Conditions, dead code, exits, statement fusion, property collection.
function f1(x) { if (x) { return 1; } else { return 2; } }
function f2(x) { if (!x) foo(); else bar(); }
function f3(x) { if (x) { a = 1; } else { a = 2; } }
function f4(x) { if (x) return; foo(); return; }
function f5() { while (true) { if (x) break; foo(); } }
function f6() { for (;;) { if (x) { continue; } bar(); } }
function f7() { switch (x) { case 1: foo(); break; default: bar(); break; } }
function f8() { switch (1) { case 1: foo(); case 2: bar(); break; case 3: baz(); } }
function f9() { try { } catch (e) { foo(); } finally { bar(); } }
function f10() { do { foo(); } while (false); }
function f11() { if (false) { var v = 1; } else { bar(); } return v; }
function f12() { throw x; foo(); }
function f13(a) { a ? b() : c(); a ? true : false; a ? x : x; }
function f14() { label: { break label; } }
function f15() { var o = {}; o.a = 1; o.b = 2; o["c"] = 3; var arr = []; arr[0] = 1; arr[2] = 3; }
function f16() { a(); b(); c(); if (d) e(); }
function f17() { x = 1; y = 1; z = x; }
function f18() { var x = /a/; if (x) foo(x); }
function f19() { if (x) { if (y) { foo(); } } }
function f20() { return x ? y : z ? w : v; }
function f21() { a = b ? c : d; if (!(a == b)) foo(); if (!(a === b)) bar(); }
function f22() { for (var i = 0; false; i++) foo(); }
function f23() { var [] = foo(); var {a: {}} = bar(); [] = baz(); }
function f24(a = undefined, b = void 0) { return a; }
function f25() { (void 0)?.x; null?.(); (void foo())?.(); }
function f26() { x["a"] = 1; x["0"] = 2; x["a-b"] = 3; var o = {"a": 1, "b c": 2, get "g"() {}}; }
class C1 { ["m"]() {} static ["s"] = 1; get ["g"]() {} }
function f27() { new Object(); new Array(); new Array(1, 2); Array(3); new RegExp("a"); Boolean(x); String(1); }
function f28() { return undefined; }
function f29() { a, b, c; return x; }
function f30() { if (x) { return f(); } return f(); }
function f31() { while (a) { return f(); } return f(); }
function f32(y) { var x = y.bind(this, 1)(2); window.Object(); window.foo(); }
function f33() { if (x) return 1; if (y) return 1; }
function f34() { x = [1, 2, 3].join(""); y = ["a", "b"].join(","); z = [a, "b", "c"].join("-"); }

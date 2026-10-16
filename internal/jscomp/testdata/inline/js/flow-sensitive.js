function a1() { var x; x = 1; print(x); }
function a2() { var x = 1; x = x + 1; print(x); }
function a3(p) { var x = p.foo; p.bar = 2; print(x); }
function a4() { var x = 1; if (cond()) { x = 2; } print(x); }
function a5() { var x = noSFX(); var y = x; print(y); hasSFX(); }
function a6() { var x = 1, y = 2; var z = x + y; print(z); }
function a7() { var x; for (x in obj) { print(x); } }
function a8() { var x = foo(); try { print(x); } catch (e) { print(e); } }
function a9() { const c = 5; print(c); }
function a10() { var a = 1; var b = a ? 2 : 3; print(b); }
function a11() { var [x, y] = arr(); var z = x; print(z); }
function a12() { var x = 1; label: { print(x); } }

function f() { return; var x = 1; }
function g() { throw 1; foo(); }
function h() { while (true) {} bar(); }
function i() { if (false) { baz(); } }
function j() { for (;;) { break; } qux(); }
function k() { do { continue; } while (false); quux(); }
function l() { return 1; ; }
function m() { switch (x) { case 1: return; break; } }
function n() { try { return; } finally { cleanup(); } after(); }
function o() { if (1) { return; } else { other(); } }
var arrow = () => { return; dead(); };
class C { static { if (x) { throw 1; } after(); } m() { return; dead(); } }

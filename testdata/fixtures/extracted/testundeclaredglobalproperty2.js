/** @const */ var x = {}; x.y = null;
function f() { x.y = 3; }
/** @param {string} a */ function g(a) { }
function h() { g(x.y); }

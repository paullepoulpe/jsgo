/** @const */ var x = {}; x.y = null;
function f(a) { x.y = a; }
/** @param {string} a */ function g(a) { }
function h() { g(x.y); }

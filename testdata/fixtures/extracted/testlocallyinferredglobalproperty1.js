/** @constructor */ function F() {}
/** @type {number} */ F.prototype.z;
/** @const */ var x = {}; /** @type {F} */ x.y;
function f() { x.y.z = 'abc'; }
/** @param {number} x */ function g(x) {}
function h() { g(x.y.z); }

/** @return {{prop: (number|undefined)}} x */
function f() { var x = {}; g(x); return x; }
/** @param {number} x */
function g(x) {}

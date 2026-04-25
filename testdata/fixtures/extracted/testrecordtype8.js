/** @return {{prop: (number|string)}} x */
function f() { var x = {prop: 3}; g(x.prop); return x; }
/** @param {string} x */
function g(x) {}

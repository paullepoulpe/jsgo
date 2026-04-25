/** @param {{prop: (number|undefined)}} x */
function f(x) {}
/** @param {{otherProp: (string|undefined)}} x */
function g(x) {}
var x = {}; f(x); g(x);

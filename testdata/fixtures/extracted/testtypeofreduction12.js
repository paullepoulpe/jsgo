/** @enum {string} */ var E = {A: 'a', B: 'b'};
/** @param {E|Array} x
 @return {Array} */
function f(x) { return typeof x == 'object' ? x : []; }

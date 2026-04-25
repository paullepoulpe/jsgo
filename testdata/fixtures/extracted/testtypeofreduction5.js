/** @enum {string} */ var E = {A: 'a', B: 'b'};
/** @param {!E|number} x
 @return {string} */
function f(x) { return typeof x != 'number' ? x : 'a'; }

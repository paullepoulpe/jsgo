/** @return {(string|null)} */
function f() { return null;}
/** @type {(string|null)} */ var a = f();
/** @type {string} */
var b = a ? a : 'default';

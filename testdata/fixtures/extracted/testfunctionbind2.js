/** @type {function(number): boolean} */
function f(x) { return true; }
f(f.bind(null, 3)());

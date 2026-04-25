/** @type {function(number, string): boolean} */
function f(x, y) { return true; }
f.bind(null, 3)(true);

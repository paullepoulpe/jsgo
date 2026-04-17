/** @type {string|number} */
var x = 1;
let g = function(/** typeof x */ a) {}
x = 'str';
g(null);

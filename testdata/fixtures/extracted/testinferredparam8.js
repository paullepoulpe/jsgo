/** @param {function(string)} callback */
function foo(callback) {}
foo(/** random JSDoc */ function(x) { var /** number */ n = x; });

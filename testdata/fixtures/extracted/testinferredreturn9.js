/** @param {function():string} x */
function f(x) {}
f(/** asdf */ function() { return 123; });

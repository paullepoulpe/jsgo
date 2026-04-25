/**
* @param {T} x
* @param {T} y
* @param {function(this:T, ...)} z
* @template T
*/
function f(x, y, z) {}
f(this, this, function() { this });

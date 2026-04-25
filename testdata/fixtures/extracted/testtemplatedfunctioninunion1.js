/**
* @param {T} x
* @param {function(this:T, ...)|{fn:Function}} z
* @template T
*/
function f(x, z) {}
f([], function() { /** @type {string} */ var x = this });

/**
* @param {T} x
* @param {function(this:T, ...)} y
* @template T
*/
function f(x, y) {}
f(0, function() {});

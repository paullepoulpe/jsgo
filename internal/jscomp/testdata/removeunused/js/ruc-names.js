var unusedGlobal = 1;
var usedGlobal = 2;
function unusedFn(a, b) { return a + b; }
function usedFn(x, unusedParam) { var local = x * 2; var dead = 3; return local; }
alert(usedFn(usedGlobal));
var sideEffect = alert('kept');
let [d1, d2] = [1, 2];
const {p, q: renamed, ...others} = window;
alert(p);
for (var i in window) {}
try { throw 1; } catch (e) {}
var f = function named() { return named; };
(function iife(unused) { var x = 1; })();

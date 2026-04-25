/** @param {null} x */ function g(x) {}
/** @param {number?} x */
function f(x) {
while (x) {
if (g(x)) { x = 1; }
x = x-1;
}
}

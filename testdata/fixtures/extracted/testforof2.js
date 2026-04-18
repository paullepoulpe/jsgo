/** @param {boolean} x */ function f(x) {}
/** @type {!Iterable<number>} */ var it = [1, 2, 3];
for (let x of it) { f(x); }

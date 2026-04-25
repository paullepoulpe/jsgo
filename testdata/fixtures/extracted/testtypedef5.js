/** @param {AB} x */ function f(x) {}
/** @constructor */ function A() {}
/** @constructor */ function B() {}
/** @typedef {(A|B)} */ var AB;
f(new A()); f(new B()); f(1);

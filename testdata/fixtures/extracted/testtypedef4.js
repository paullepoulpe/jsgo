/** @constructor */ function A() {}
/** @constructor */ function B() {}
/** @typedef {(A|B)} */ var AB;
/** @param {AB} x */ function f(x) {}
f(new A()); f(new B()); f(1);

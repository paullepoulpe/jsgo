/** @constructor
 * @extends {B}  */ function C() {}
/** @constructor
 * @extends {D}  */ function E() {}
/** @constructor
 * @extends {C}  */ function D() {}
/** @constructor
 * @extends {A} */ function B() {}
/** @constructor */ function A() {}
/** @param {number} x */ function f(x) {} f(new E());

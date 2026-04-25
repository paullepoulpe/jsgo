/** @constructor
 * @extends {Super} */ function Sub() {}
/** @constructor
 * @extends {Sub} */ function Sub2() {}
/** @constructor */ function Super() {}
/** @param {Super} x */ function foo(x) {}
foo(new Sub2());

/** @constructor
 * @extends {bad} */ function Sub() {}
/** @constructor
 * @extends {Sub} */ function Sub2() {}
/** @param {Sub} x */ function foo(x) {}
foo(new Sub2());

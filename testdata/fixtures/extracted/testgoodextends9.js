/** @constructor */ function Super() {}
/** @constructor
 * @extends {Super} */ function Sub() {}
Sub.prototype = new Super();
/** @return {Super} */ function foo() { return new Sub(); }

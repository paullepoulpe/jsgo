/** @constructor */ function Super() {}
/** @param {boolean} x */ Super.prototype.foo = function(x) {};
/** @constructor
 * @extends {Super} */ function Sub() {}
Sub.prototype = new Super();
(new Sub()).foo(0);

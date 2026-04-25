/** @constructor */ function Super() {}
Super.prototype.foo = function() {};
/** @constructor
 * @extends {Super} */ function Sub() {}
Sub.prototype = new Super();
/** @override */ Sub.prototype.foo = function() {};

var goog = {};
/** @constructor */ goog.F = function() {};
/** @return {number} */ goog.F.prototype.foo = function() { return 3; };
/** @constructor
 * @extends {goog.F} */
goog.G = function() {};
/** @override */ goog.G.prototype.foo = function() { return true; };

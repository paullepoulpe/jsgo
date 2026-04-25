var goog = {};
/** @constructor */goog.Super = function() {};
/** @type {number} */ goog.Super.prototype.foo = 3;
/** @constructor
 @extends {goog.Super} */goog.Sub = function() {};
/** @type {number} */ goog.Sub.prototype.foo = 5;

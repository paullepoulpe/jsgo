var goog = {};
/** @constructor */goog.Super = function() {};
/** @param {number} x */ goog.Super.prototype.foo = function(x) {};
/** @constructor
 @extends {goog.Super} */goog.Sub = function() {};
/** @override @param {string} x */ goog.Sub.prototype.foo = function(x) {};

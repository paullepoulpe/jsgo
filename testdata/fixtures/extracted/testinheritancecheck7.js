var goog = {};
/** @constructor */goog.Super = function() {};
goog.Super.prototype.foo = 3;
/** @constructor
 @extends {goog.Super} */goog.Sub = function() {};
goog.Sub.prototype.foo = 5;

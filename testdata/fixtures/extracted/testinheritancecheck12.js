var goog = {};
/** @constructor */goog.Super = function() {};
goog.Super.prototype.foo = 3;
/** @constructor
 @extends {goog.Super} */goog.Sub = function() {};
/** @override */goog.Sub.prototype.foo = "some string";

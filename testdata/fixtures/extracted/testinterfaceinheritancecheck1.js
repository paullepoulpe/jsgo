/** @interface */function Super() {};
Super.prototype.foo = function() {};
/** @constructor
 @implements {Super} */function Sub() {};
Sub.prototype.foo = function() {};

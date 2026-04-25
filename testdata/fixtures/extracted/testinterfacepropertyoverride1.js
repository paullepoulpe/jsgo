/** @interface */function Super() {};
Super.prototype.foo = function() {};
/** @interface
 @extends {Super} */function Sub() {};
Sub.prototype.foo = function() {};

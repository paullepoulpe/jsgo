/** @interface */function Super() {};
Super.prototype.foo = function() {};
/** @constructor
 @implements {Super} */function Sub() {};
/** @override */Sub.prototype.foo = function() {};

/** @constructor */function Super() {};
Super.prototype.foo = function() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override */Sub.prototype.foo = function() {};

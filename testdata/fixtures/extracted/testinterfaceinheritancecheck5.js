/** @interface */function Super() {};/** @return {string} */Super.prototype.foo = function() {};/** @constructor
 @implements {Super} */function Sub() {};/** @override
 @return {number} */Sub.prototype.foo = function() { return 1; };

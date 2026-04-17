/** @interface */function Root() {};
/** @return {number} */Root.prototype.foo = function() {};
/** @interface
 @extends {Root} */function Super() {};
/** @constructor
 @implements {Super} */function Sub() {};
/** @return {number} */Sub.prototype.foo = function() { return 1;};

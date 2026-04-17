/** @constructor */function Root() {};
Root.prototype.foo = function() {};
/** @constructor
 @extends {Root} */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
Sub.prototype.foo = function() {};

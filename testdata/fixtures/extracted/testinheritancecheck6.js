/** @constructor */function Root() {};
Root.prototype.foo = function() {};
/** @constructor
 @extends {Root} */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override */Sub.prototype.foo = function() {};

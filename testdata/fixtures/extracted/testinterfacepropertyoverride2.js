/** @interface */function Root() {};
Root.prototype.foo = function() {};
/** @interface
 @extends {Root} */function Super() {};
/** @interface
 @extends {Super} */function Sub() {};
Sub.prototype.foo = function() {};

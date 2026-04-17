/** @interface */function Root() {};
/** @return {number} */Root.prototype.foo = function() {};
/** @interface
 @extends {Root} */function Super() {};
/** @constructor
 @implements {Super} */function Sub() {};
/** @override
 * @return {number} */Sub.prototype.foo =
function() { return 1;};

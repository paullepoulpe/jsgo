/** @constructor */function Root() {};
/** @return {number} */
Root.prototype.foo = function() { return 1; };
/** @constructor
 @extends {Root} */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
 @return {string} */Sub.prototype.foo =
function() { return "some string" };

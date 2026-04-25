/** @constructor */function Super() {};
/** @return {number} */
Super.prototype.foo = function() { return 1; };
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
 @return {string} */Sub.prototype.foo =
function() { return "some string" };

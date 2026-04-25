/** @constructor */function Super() {};
/** @return {number} */
Super.prototype.foo = function() { return 1; };
/** @constructor
 @extends {Super} */function Sub() {};
/** @override */Sub.prototype.foo =
function() {};

/** @constructor */function Super() {};
Super.prototype.foo = function() { return 3; };
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
 @return {number} */Sub.prototype.foo =
function() { return 1; };

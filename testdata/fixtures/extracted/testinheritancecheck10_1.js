/** @constructor */function Root() {};
Root.prototype.foo = function() { return 3; };
/** @constructor
 @extends {Root} */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
 @return {number} */Sub.prototype.foo =
function() { return 1; };

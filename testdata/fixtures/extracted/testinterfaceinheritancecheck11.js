/** @constructor */function Super() {};
/** @param {number} bar */Super.prototype.foo = function(bar) {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
  @param {string} bar */Sub.prototype.foo =
function(bar) {};

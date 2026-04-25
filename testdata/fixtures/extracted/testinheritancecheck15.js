/** @constructor */function Super() {};
/** @param {number} bar */Super.prototype.foo;
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
  @param {number} bar */Sub.prototype.foo =
function(bar) {};

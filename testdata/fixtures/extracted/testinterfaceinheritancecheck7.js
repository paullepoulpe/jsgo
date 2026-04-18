/** @interface */function Super() {};
/** @param {number} bar */Super.prototype.foo = function(bar) {};
/** @constructor
 @implements {Super} */function Sub() {};
/** @override
  @param {string} bar */Sub.prototype.foo =
function(bar) {};

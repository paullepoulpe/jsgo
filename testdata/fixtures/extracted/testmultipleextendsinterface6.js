/** @interface */function Super1() {};
/** @interface */function Super2() {};
/** @param {number} bar */Super2.prototype.foo = function(bar) {};
/** @interface @extends {Super1} @extends {Super2} */function Sub() {};
/** @override @param {string} bar */Sub.prototype.foo =
function(bar) {};

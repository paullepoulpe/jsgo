/** @constructor */ function F() {}
/** @return {number} */ F.prototype.foo = function() { return 3; };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */ G.prototype.foo = function() { return true; };

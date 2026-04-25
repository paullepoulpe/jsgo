/** @constructor */ function F() {}
/** @param {boolean} x
 * @return {number} */
F.prototype.foo = function(x) { return 3; };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */
G.prototype.foo = function(y) { return y; };

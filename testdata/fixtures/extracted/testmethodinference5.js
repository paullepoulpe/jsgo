/** @constructor */ function F() {}
/** @param {number} x
 * @return {string} */
F.prototype.foo = function(x) { return 'x'; };
/** @constructor
 * @extends {F} */
function G() {}
/** @type {number} */ G.prototype.num = 3;
/** @override */
G.prototype.foo = function(y) { return this.num + y; };

/** @constructor */ function F() {}
F.prototype.foo = function() { };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */ G.prototype.foo = function(x, y) { };

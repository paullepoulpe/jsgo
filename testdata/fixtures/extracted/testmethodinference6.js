/** @constructor */ function F() {}
/** @param {number} x */ F.prototype.foo = function(x) { };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */ G.prototype.foo = function() { };
(new G()).foo(1);

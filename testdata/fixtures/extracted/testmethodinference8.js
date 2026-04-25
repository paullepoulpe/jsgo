/** @constructor */ function F() {}
F.prototype.foo = function() { };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */
G.prototype.foo = function(opt_b, var_args) { };
(new G()).foo(1, 2, 3);

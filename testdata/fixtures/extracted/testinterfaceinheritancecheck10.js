/** @interface */ function I() {}
/** @return {number} */ I.prototype.bar = function() {};
/** @constructor */ function F() {}
/** @return {number} */ F.prototype.foo = function() {return 3; };
/** @constructor
 * @extends {F}
 * @implements {I} */
function G() {}
/** @return {number}
 * @override */
G.prototype.bar = G.prototype.foo;
/** @return {string} */ function f() { return new G().bar(); }

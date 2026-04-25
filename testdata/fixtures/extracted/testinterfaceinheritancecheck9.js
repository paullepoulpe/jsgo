/** @interface */ function I() {}
/** @return {number} */ I.prototype.bar = function() {};
/** @constructor */ function F() {}
/** @return {number} */ F.prototype.bar = function() {return 3; };
/** @return {number} */ F.prototype.foo = function() {return 3; };
/** @constructor
 * @extends {F}
 * @implements {I} */
function G() {}
/** @return {string} */ function f() { return new G().bar(); }

/** @interface */ function I() {};
/** @type {string} */ I.prototype.foobar;
/**\s
 * @constructor\s
 * @implements {I} */
function C() {
/**\s
 * @type {number} */ this.foobar = 2;};
/** @type {I} */\s
 var test = new C(); alert(test.foobar);

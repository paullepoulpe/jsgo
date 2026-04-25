/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @interface
 * @extends {Int1} */ function Int2() {};
/** @interface
 * @extends {Int0}
 * @extends {Int2} */
function Int3() {};

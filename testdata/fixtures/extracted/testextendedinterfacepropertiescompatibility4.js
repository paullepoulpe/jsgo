/** @interface */function Int0() {};
/** @interface
 * @extends {Int0} */ function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @interface */function Int2() {};
/** @interface
 * @extends {Int2} */ function Int3() {};
/** @type {string} */
Int2.prototype.foo;
/** @interface
 * @extends {Int1}
 * @extends {Int3} */
function Int4() {};

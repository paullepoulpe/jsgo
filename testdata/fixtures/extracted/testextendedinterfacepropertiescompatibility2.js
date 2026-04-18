/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @interface */function Int2() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @type {Object} */
Int2.prototype.foo;
/** @interface
 * @extends {Int0}
 * @extends {Int1}
 * @extends {Int2}*/
function Int3() {};

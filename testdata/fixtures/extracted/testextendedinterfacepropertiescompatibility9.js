/** @interface
 * @template T */function Int0() {};
/** @interface
 * @template T */function Int1() {};
/** @type {T} */
Int0.prototype.foo;
/** @type {T} */
Int1.prototype.foo;
/** @interface
 * @extends {Int0<number>}
 * @extends {Int1<string>} */
function Int2() {};

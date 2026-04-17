/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @interface
 * @extends {Int0}
 * @extends {Int1}
 * @extends {number} */
function Int2() {};
/** @constructor
 * @implements {Int2} */function Foo() {};

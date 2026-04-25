/** @interface */function Int0() {};
/** @interface */function Int1() {};
Int0.prototype.foo = function() {};
/** @interface
 * @extends {Int0}
 * @extends {Int1} */
function Int2() {};
/** @constructor
 * @implements {Int2} */function Foo() {};

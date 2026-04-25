/** @interface */function Int() {};
Int.prototype.foo = function() {};
/** @interface
 @extends {Int} */function Int2() {};
/** @constructor
 @implements {Int2} */function Foo() {};

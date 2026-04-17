/** @interface @template T */ function A() {}
/** @constructor @implements {A<Foo>} */
function Foo() {}
/** @type {!Foo} */
var x = new Foo();

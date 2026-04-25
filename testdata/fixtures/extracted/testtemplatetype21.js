/** @interface @template T */ function A() {}
/** @constructor @implements {A<Foo>} */
function Foo() {}
/** @constructor @implements {A<Bar>} */
function Bar() {}
/** @type {!Foo} */
var x = new Bar();

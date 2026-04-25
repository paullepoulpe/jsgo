/** @interface @template T */ function A() {}
/** @interface @template T */ function B() {}
/** @constructor @implements {A<Foo>} */
function Foo() {}
/** @constructor @implements {B<Foo>} */
function Bar() {}
/** @constructor @implements {B<Foo>} */
function Qux() {}
/** @type {!Qux} */
var x = new Bar();

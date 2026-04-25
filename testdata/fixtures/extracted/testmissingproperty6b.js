/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
/** @constructor
 * @extends {Foo} */
function Bar() { this.a = 3; };

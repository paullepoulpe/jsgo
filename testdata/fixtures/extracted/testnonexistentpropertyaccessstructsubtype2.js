/**
 * @constructor
 * @struct
 */
function Foo() {
  this.x = 123;
}
var objlit = /** @struct */ { y: 234 };
Foo.prototype = objlit;
var n = objlit.x;

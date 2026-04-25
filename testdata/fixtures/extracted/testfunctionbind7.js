/** @constructor */
function MyType() {
  /** @type {number} */
  this.x = 0;
}
var m = new MyType;
(function f() {this.x = 'str';}).bind(m);

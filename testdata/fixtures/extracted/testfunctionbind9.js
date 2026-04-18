/** @constructor */
function MyType() {}

/** @constructor */
function AnotherType() {}
AnotherType.prototype.foo = function() {};

var m = new MyType;
(function f() {this.foo();}).bind(m);

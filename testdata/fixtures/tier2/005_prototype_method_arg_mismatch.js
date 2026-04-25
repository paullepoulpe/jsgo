/** @constructor */
function Foo() {}
/** @param {number} x */
Foo.prototype.bar = function(x) {};
var f = new Foo();
f.bar('hi');

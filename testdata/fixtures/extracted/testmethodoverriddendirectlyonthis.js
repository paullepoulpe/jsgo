/** @constructor */
function Foo() {
  this.bar = function() { return 'str'; };
}
/** @return {number} */
Foo.prototype.bar = function() {};

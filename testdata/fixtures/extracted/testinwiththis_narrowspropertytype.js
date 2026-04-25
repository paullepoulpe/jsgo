/** @constructor */
function Foo() {}
Foo.prototype.method = function() {
  if ('x' in this) {
    return this.x; // this access is allowed
  }
  return this.y; // this access causes a warning
}

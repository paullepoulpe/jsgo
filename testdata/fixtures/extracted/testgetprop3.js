/** @constructor */
function Foo() { /** @type {?Object} */ this.x = null; }
Foo.prototype.initX = function() { this.x = {foo: 1}; };
Foo.prototype.bar = function() {
  if (this.x == null) { this.initX(); alert(this.x.foo); }
};

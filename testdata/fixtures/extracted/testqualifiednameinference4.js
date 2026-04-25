/** @param {string} x */ function f(x) {}
/**
 * @param {?string} x
 * @constructor
 */
function Foo(x) { this.x_ = x; }
Foo.prototype.bar = function() {
  if (this.x_) { f(this.x_); }
};

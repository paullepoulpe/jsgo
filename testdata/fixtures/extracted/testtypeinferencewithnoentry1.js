/** @param {number} x */ function f(x) {}
/** @constructor */ function Foo() {}
Foo.prototype.init = function() {
  /** @type {?{baz: number}} */ this.bar = {baz: 3};
};
/**
 * @extends {Foo}
 * @constructor
 */
function SubFoo() {}
/** Method */
SubFoo.prototype.method = function() {
  for (var i = 0; i < 10; i++) {
    f(this.bar);
    f(this.bar.baz);
  }
};

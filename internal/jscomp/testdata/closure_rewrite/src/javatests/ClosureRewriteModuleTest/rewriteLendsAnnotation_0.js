goog.module('mod');
class Foo {
  constructor() {
    this.x = /** @lends {Foo.prototype} */ {};
  }
}

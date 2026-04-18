function f(/** !IArrayLike */ x){};
/** @constructor */
function Foo() {
  /** @type {number} */ this.length = 5;
}
f(new Foo)

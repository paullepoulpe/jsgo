/** @constructor */ function Foo() {
  /** @type {!Bar} */
  this.baz = new Bar();
}
Foo.prototype.foo = function() {
  this.baz.bar();
};
/** @constructor */ function Bar() {
  /** @type {!Foo} */
  this.baz = new Foo();
}
Bar.prototype.bar = function() {
  this.baz.foo();
};

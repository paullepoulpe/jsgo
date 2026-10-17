var Foo = class {
  /**
   * @return {number}
   */
  m() {return 6;};
};
goog.exportSymbol('Foo', Foo);
goog.exportProperty(Foo.prototype, 'm', Foo.prototype.m);

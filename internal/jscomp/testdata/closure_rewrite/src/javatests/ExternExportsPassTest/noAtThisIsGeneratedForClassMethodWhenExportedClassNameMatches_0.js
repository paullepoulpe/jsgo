class Foo {
  /**
   * @param {string} a
   * @return {number}
   */
  method(a) {}
};
goog.exportSymbol('Foo', Foo)
goog.exportProperty(Foo.prototype, 'method', Foo.prototype.method)

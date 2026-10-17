/** @const */ var a = {};
/** @const */ a.b = {};
a.b.c = class {
  constructor(d, e, f) {
    /** @export {number} */ this.thisProp = 1;
  }
  /**
   * @param {string} a
   * @return {number}
   */
  method(a) {}

  static staticMethod() {}
};
goog.exportSymbol('foobar', a.b.c)
goog.exportProperty(a.b.c.prototype, 'exportedMethod', a.b.c.prototype.method)
goog.exportProperty(a.b.c, 'exportedStaticMethod', a.b.c.staticMethod)

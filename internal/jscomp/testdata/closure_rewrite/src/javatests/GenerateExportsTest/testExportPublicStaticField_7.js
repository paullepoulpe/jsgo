let ns = {};
ns = class Bar {
  /** @export */ static x;
};
goog.exportSymbol('ns.x', ns.x);

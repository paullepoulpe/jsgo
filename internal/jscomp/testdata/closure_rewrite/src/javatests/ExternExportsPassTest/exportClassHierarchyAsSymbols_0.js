/** @const */ var a = {};
/** @const */ a.b = {};
a.b.c = class {
  constructor(d, e, f) {}
};
goog.exportSymbol('foobar', a.b.c)
a.b.d = class extends a.b.c {
  constructor(d, e, f) {}
};
goog.exportSymbol('bazboff', a.b.d)

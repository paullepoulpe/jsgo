const a = {};
/** @const */ a.b = {};
a.b.c = class {
  constructor(d, e, f) {}
};
goog.exportProperty(a.b, 'cprop', a.b.c)

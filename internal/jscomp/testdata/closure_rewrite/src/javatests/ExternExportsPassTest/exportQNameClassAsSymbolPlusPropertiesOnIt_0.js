const a = {};
a.b = class {
  constructor(p1) {
  }
};
a.b.c = function(d, e, f) {};
a.b.prototype.c = function(g, h, i) {};
goog.exportSymbol('a.b', a.b);
goog.exportProperty(a.b, 'c', a.b.c);
goog.exportProperty(a.b.prototype, 'c', a.b.prototype.c);

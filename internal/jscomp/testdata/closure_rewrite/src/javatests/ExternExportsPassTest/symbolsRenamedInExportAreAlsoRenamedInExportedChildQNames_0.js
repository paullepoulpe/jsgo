/** @const */ var a = {};
a.b = function(p1) {};
a.b.c = function(d, e, f) {};
a.b.prototype.c = function(g, h, i) {};
goog.exportSymbol('hello', a);
goog.exportProperty(a.b, 'c', a.b.c);
goog.exportProperty(a.b.prototype, 'c', a.b.prototype.c);

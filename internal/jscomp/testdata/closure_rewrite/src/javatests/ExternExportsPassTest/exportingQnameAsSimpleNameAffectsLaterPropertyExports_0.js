/** @const */ var a = {};
a.b = function(p1) {};
a.b.c = function(d, e, f) {};
a.b.prototype.c = function(g, h, i) {};
goog.exportSymbol('prefix', a.b);
goog.exportProperty(a.b, 'c', a.b.c);

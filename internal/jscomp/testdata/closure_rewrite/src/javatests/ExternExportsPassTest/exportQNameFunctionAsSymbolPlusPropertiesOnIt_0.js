/** @const */ var a = {};
a.b = function(p1) {};
/** @constructor */
a.b.c = function(d, e, f) {};
/** @return {number} */
a.b.c.prototype.method = function(g, h, i) {};
goog.exportSymbol('a.b', a.b);
goog.exportProperty(a.b, 'c', a.b.c);
goog.exportProperty(a.b.c.prototype, 'exportedMethod', a.b.c.prototype.method);

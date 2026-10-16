/** @const */ var a = {};
/** @const */ a.b = {};
a.b.c = function(d, e, f) {};
goog.exportProperty(a.b, 'cprop', a.b.c)

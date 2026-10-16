/** @const */ var a = {};
/** @const */ a.b = {};
/** @const */ var d = {};
a.b.c = d;
goog.exportSymbol('foobar', a.b.c)

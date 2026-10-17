/** @const */ var a = {};
/** @const */ a.b = {};
var d = function() {};
a.b.c = d;
goog.exportSymbol('foobar', a.b.c())

/** @const */ var a = {};
/** @const */ a.b = {};
a.b.c = function(d, e, f) {};
var x;
// Ensure we don't have a recurrence of a bug that caused us to consider the last
// assignment containing 'a.b.c' to be its definition, even if it was on the rhs.
x = a.b.c;
goog.exportSymbol('foobar', a.b.c)

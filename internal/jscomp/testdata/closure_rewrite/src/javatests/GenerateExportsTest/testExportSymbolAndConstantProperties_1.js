/** @export */function foo(){}
goog.exportSymbol("foo",foo);
/** @export */foo.BAR=5;
goog.exportProperty(foo,"BAR",foo.BAR)

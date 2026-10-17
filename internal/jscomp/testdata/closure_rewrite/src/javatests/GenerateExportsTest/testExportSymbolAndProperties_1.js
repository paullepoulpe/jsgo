/** @export */function foo(){}
goog.exportSymbol("foo",foo);
/** @export */foo.prototype.bar=function(){};
goog.exportProperty(foo.prototype,"bar",foo.prototype.bar)

var foo = function() { /** @export */ this.foo = 1; };
/** @export */ foo.method = function(){};
goog.exportSymbol('foo.method', foo.method);

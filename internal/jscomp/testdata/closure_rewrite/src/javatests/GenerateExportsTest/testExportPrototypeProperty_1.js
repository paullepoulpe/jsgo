function Foo() {}
/** @export */ Foo.prototype.bar = function(){};
goog.exportProperty(Foo.prototype, 'bar', Foo.prototype.bar);

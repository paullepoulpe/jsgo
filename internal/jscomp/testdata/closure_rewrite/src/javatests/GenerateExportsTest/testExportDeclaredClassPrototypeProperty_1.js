class Foo {}
/** @export */
Foo.prototype.x = 3;
goog.exportProperty(Foo.prototype, 'x', Foo.prototype.x);

/**
 * @constructor
 */
function Foo() {}
goog.exportSymbol('Foo', Foo);
/**
 * @param {!Foo} x
 * @return {!Foo}
 */
Foo.f = function(x) { return x; };
goog.exportProperty(Foo, 'f', Foo.f);

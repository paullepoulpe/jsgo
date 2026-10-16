goog.provide('foo.baz');
goog.provide('other.thing');

goog.scope(function() {

const a = goog.module.get('other.thing');
const b = a.b;
foo.baz = b

}); // goog.scope

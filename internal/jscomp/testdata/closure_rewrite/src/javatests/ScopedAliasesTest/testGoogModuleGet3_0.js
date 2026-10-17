goog.provide('foo.baz');
goog.provide('other.thing');

goog.scope(function() {

const a = goog.module.get('other.thing').b;
foo.baz = a

}); // goog.scope

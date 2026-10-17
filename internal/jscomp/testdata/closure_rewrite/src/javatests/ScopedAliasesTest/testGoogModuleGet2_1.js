goog.provide('foo.baz');
goog.provide('other.thing');
foo.baz = goog.module.get('other.thing').b;

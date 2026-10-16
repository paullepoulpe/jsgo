goog.provide('provided');
goog.provide('other.name.Foo');
/** @type {!other.name.Foo} */
provided.f = new (goog.module.get('other.name.Foo'));

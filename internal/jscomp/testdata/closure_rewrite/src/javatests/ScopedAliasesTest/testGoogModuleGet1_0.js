goog.provide('provided');
goog.provide('other.name.Foo');
goog.scope(function() {
  var Foo = goog.module.get('other.name.Foo')
  /** @type {!Foo} */
  provided.f = new Foo;
});

goog.provide('provided');
goog.scope(function() {
  var Foo = goog.module.get('other.name.Foo')
  /** @type {!Foo} */
  provided.f = new Foo;
});

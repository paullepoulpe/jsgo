goog.provide('provided');

goog.scope(function() {
  var Foo = goog.module.get('missing.name.Foo')
  /** @type {!Foo} */
  provided.f = new Foo;
});

goog.provide('provided');

goog.scope(function() {
  var Foo = goog.module.get('missing.name.Foo');
  var Bar = Foo.Bar
  /** @type {!Bar} */
  provided.f = new Bar;
});

goog.provide('provided');

goog.scope(function() {
  var Bar = goog.module.get('missing.name.Foo').Bar
  /** @type {!Bar} */
  provided.f = new Bar;
});

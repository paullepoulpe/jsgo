goog.provide('provided');

goog.scope(function() {
  var Foo = goog.module.get('missing.name.Foo')
  provided.Bar = class extends Foo {};
});

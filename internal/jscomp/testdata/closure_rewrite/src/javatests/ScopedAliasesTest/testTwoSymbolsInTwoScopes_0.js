var goog = {};
goog.scope(function() { var g = goog; g.Foo = function() {}; });
goog.scope(function() {
  var Foo = goog.Foo; goog.bar = function() { return new Foo(); };
});

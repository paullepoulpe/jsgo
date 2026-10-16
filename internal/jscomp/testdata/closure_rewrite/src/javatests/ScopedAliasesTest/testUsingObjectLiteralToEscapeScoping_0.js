var goog = {};
goog.bar = {};
goog.scope(function() {
  var bar = goog.bar;
  var baz = goog.bar.baz;
  goog.foo = function() {
    goog.bar = {baz: 3};
    return baz;
  };
});

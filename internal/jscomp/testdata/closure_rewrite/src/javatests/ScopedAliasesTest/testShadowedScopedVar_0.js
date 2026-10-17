var goog = {};
goog.bar = {};
goog.scope(function() {
  var bar = goog.bar;
  // This is bogus, because when the aliases are expanded, goog will shadow goog.bar.
  bar.newMethod = function(goog) { return goog + bar; };
});

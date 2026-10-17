var goog = {};
goog.bar = {};
goog.scope(function() {
  var bar = goog.bar;
  bar.newMethod = function(goog, a) { return bar + a; };
});
goog.scope(function() {
  var bar = goog.bar;
  bar.newMethod2 = function(goog, b) { return bar + b; };
});

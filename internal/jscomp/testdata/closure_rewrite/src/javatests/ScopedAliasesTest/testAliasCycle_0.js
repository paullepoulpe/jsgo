var x = {y: {}};
goog.scope(function() {
  var y = z.x;
  var z = y.x;
  y.ClassA = function() {};
  z.ClassB = function() {};
});

var x = {y: {}};
goog.scope(function() {
  var y = x.y;
  y.ClassA = function() { this.b = new ClassB(); };
  y.ClassB = function() {};
  var ClassB = y.ClassB;
});

goog.provide('a.b');
goog.scope(function() {
  var x = a.b;
  x.c = 1;
  debugger;
  this.y = 1;
});

goog.scope(function() {
  var x = goog.dom;
  x = goog.events;
  this.y = 1;
  if (x) { throw new Error(); }
  return;
});
goog.scope(function(a) {});
goog.scope(function named() {});
var y = goog.scope(function() {});
function f() { goog.scope(function() {}); }
goog.scope(function() {
  var a = b.c;
  var b = a.d;
});
goog.scope(function() {
  try {} catch (e) {}
  var goog2 = goog;
  function shadow() { var goog = 1; return goog; }
});

var goog = {};
goog.bar = {};
goog.foo = function(){
  goog.bar = {baz:3};
  return goog.bar.baz;
};

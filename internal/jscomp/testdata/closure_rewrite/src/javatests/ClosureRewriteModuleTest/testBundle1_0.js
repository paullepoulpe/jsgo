goog.loadModule(function(exports) {
  goog.module('ns.a');
  var b = goog.require('ns.b');
  exports.b = b;
  return exports;
});

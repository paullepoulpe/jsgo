goog.provide('a.b.c');
goog.provide('a.d');
goog.provide('e');
goog.require('x.y');
var e = {};
a.b.c = function() {};
goog.loadModule(function(exports) {
  goog.module('lm.mod');
  goog.module.declareLegacyNamespace();
  const local = 1;
  exports = {local};
  return exports;
});
goog.loadModule(function(exports) {
  goog.module('lm.other');
  exports.v = 1;
  return exports;
});
goog.provide();
goog.provide(1);

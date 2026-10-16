goog.loadModule(function(exports) {
  goog.module('ns.b');
  return exports;
});
goog.loadModule(function(exports) {
  'use strict';
  goog.module('ns.a');
  goog.module.declareLegacyNamespace();
  var b = goog.require('ns.b');
  return exports;
});

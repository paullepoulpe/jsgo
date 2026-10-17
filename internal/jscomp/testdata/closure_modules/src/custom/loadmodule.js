goog.loadModule(function(exports) {
  'use strict';
  goog.module('custom.loaded');
  goog.module.declareLegacyNamespace();
  const a = goog.require('custom.a');
  exports.x = a;
  return exports;
});
goog.loadModule(function(exports) {
  'use strict';
  goog.module('custom.loaded2');
  exports = 1;
  return exports;
});

goog.loadModule(function(exports) {
  goog.module('goog.asserts');
  return exports;
});
goog.loadModule(function(exports) {
  'use strict';
  goog.module('xid');
  goog.module.declareLegacyNamespace();
  var asserts = goog.require('goog.asserts');
  exports = function(id) {
    return xid.internal_(id);
  };
  var xid = exports;
  exports.internal_ = function(id) {};
  return exports;
});

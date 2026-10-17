goog.loadModule(function(exports) {
  goog.module('goog.asserts');
  return exports;
});
goog.loadModule(function(exports) {
  'use strict';
  goog.module('xid');
  goog.module.declareLegacyNamespace();
  var asserts = goog.require('goog.asserts');
  var xid = function(id) {
    return xid.internal_(id);
  };
  xid.internal_ = function(id) {};
  exports = xid;
  return exports;
});

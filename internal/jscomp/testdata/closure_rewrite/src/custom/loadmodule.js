goog.loadModule(function(exports) {
  'use strict';
  goog.module('lm.a');
  /** @param {number} x */
  exports.double = function(x) { return 2 * x; };
  exports.Klass = class { m() { return 1; } };
  return exports;
});
goog.loadModule(function(exports) {
  goog.module('lm.b');
  const a = goog.require('lm.a');
  const {Klass} = goog.require('lm.a');
  /** @type {!Klass} */
  const k = new Klass();
  exports = a.double(k.m());
  return exports;
});
goog.loadModule(function(exports) {
  goog.module('lm.c');
  goog.module.declareLegacyNamespace();
  const b = goog.require('lm.b');
  exports.value = b;
  return exports;
});
goog.loadModule(function(exports) {
  goog.module('lm.noreturn');
  exports.x = 1;
});
goog.require('lm.c');
use(lm.c.value, goog.module.get('lm.a').double(2));

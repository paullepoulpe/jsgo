goog.loadModule(function(exports) { 'use strict';
  goog.module('mod_B');

  /** @interface */ function B(){}

  exports.B = B;
  return exports;
});
goog.loadModule(function(exports) { 'use strict';
  goog.module('mod_A');

  var {B} = goog.require('mod_B');

  /** @constructor @implements {B} */
  function A() {}
  return exports;
});

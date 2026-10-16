goog.module('ns.a');

function f(exports, a) {
// test the various syntactic froms of doing goog.module exports
  exports.prop = 0;
  exports = {a};
  exports = function() {};
  if (true) {
    const exports = {};
  }
  return exports;
}

goog.module('my.module');
goog.module.declareLegacyNamespace();
const {assert} = goog.require('goog.asserts');
const dom = goog.require('goog.dom');
$jscomp.scope.someVar = 3;
$jscomp.scope.someFn = function () {};
/** @const */
exports.x = 1;
exports = {a, b};
goog.scope(function () {
  var D = goog.dom;
});

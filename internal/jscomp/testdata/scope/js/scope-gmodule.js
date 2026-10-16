goog.module('my.mod');
goog.module.declareLegacyNamespace();
const dep = goog.require('other.dep');
const {a, b: bb} = goog.require('other.destructured');
function f() { return exports; }
exports.f = f;
exports = {f};
let later = goog.weakUsage(dep);

goog.module('app.props');
goog.module.declareLegacyNamespace();

exports.a = 1;
/** @const */ exports.b = function() { return exports.a; };
exports.C = class {};
/** @typedef {number} */
exports.T;
var local = 5;
exports.local = local;
function f() { return exports.b(); }
exports.f = f;

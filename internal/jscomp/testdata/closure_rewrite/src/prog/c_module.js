goog.module('p.C');
goog.module.declareLegacyNamespace();
const B = goog.require('p.B');
const late = goog.require('p.late');
/** @type {!B} */
exports.b = new B();
exports.late = late;

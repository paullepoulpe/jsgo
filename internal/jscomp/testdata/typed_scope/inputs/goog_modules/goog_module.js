goog.module('my.mod');
goog.module.declareLegacyNamespace();
const dep = goog.require('other.dep');
const {Thing} = goog.requireType('other.things');
/** @constructor */
function Impl() {}
/** @enum {number} */
const Kind = {A: 1};
exports = {Impl, Kind};

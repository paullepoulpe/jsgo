goog.module('custom.gmodule');
goog.module.declareLegacyNamespace();
const {assert, fail} = goog.require('goog.asserts');
const dom = goog.require('goog.dom');
let letReq = goog.require('goog.array');
const Unused = goog.require('custom.Unused');
const {Unused2} = goog.require('custom.Unused2');
goog.requireType('custom.Type');
const T = goog.requireType('custom.T2');
const F = goog.forwardDeclare('custom.F');
var a = goog.require('x.y'), b = 1;
const {c = 1} = goog.require('x.z');
goog.require('goog.array');
goog.array.forEach([], () => {});
custom.gmodule.x = 1;
/** @type {goog.dom.TagName} */ let tag;
/** @type {!T} */ let t;
function f() { goog.require('late'); }
exports = class {};
exports.x = 1;
exports.x = 2;
exports.y = this;
if (true) { exports.z = 1; }
throw new Error();
goog.module('custom.another');

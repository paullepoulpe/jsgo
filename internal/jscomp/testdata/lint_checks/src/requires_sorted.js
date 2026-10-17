// leading comment
goog.module('x.y');
/** @suppress {extraRequire} */
const {Zed, Alpha: Al, Beta: Beta} = goog.require('ns.b');
const b = goog.requireType('ns.a');
goog.require('ns.c');
goog.forwardDeclare('ns.c');
const {Gamma} = goog.require('ns.b');
const a2 = goog.require('ns.a');
const {} = goog.require('ns.e');
goog.requireType('ns.d');
goog.require('ns.d');
let other = 1;
goog.require('ns.after');
/** @param {!Al} p @return {!Gamma} */
function use(p) { return a2 && goog.module.get('ns.f') && Zed; }
exports = {use};

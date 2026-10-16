goog.module('custom.extra');
const Used = goog.require('custom.Used');
const Unused = goog.require('custom.Unused');
const {UsedD, UnusedD} = goog.require('custom.D');
const {AllUnused} = goog.require('custom.E');
/** @suppress {extraRequire} */
const Suppressed = goog.require('custom.S');
const TypeOnly = goog.requireType('custom.TypeOnly');
const Fwd = goog.forwardDeclare('custom.Fwd');
goog.forwardDeclare('custom.Fwd2');
goog.require('goog.testing.asserts');
/** @type {!TypeOnly} */ let t;
new Used(UsedD);
goog.module.get('custom.Got');

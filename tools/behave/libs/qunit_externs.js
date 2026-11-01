/**
 * @fileoverview Externs of the behavioural oracle's library suites
 * (tools/behave/libs.py): the QUnit subset that tools/behave/runner.js
 * provides with --env qunit, and the host names the suites' own environment
 * checks read (node's require/process/module, AMD's define, PhantomJS).
 * @externs
 */

/** @const */
var QUnit = {};
/** @type {!Object} */ QUnit.config;
QUnit.config.asyncRetries;
QUnit.config.hidepassed;
QUnit.config.noglobals;
QUnit.config.testTimeout;
QUnit.config.reorder;
/** @param {string} name @param {*=} hooks */ QUnit.module = function(name, hooks) {};
/** @param {string} name @param {function(!QUnitAssert)} fn */ QUnit.test = function(name, fn) {};
/** @param {string} name @param {function(!QUnitAssert)} fn */ QUnit.only = function(name, fn) {};
/** @param {string} name @param {function(!QUnitAssert)=} fn */ QUnit.skip = function(name, fn) {};
/** @param {string} name @param {function(!QUnitAssert)=} fn */ QUnit.todo = function(name, fn) {};
QUnit.load = function() {};
QUnit.start = function() {};
/** @param {function(?)} cb */ QUnit.done = function(cb) {};
/** @param {function(?)} cb */ QUnit.begin = function(cb) {};
/** @param {function(?)} cb */ QUnit.testStart = function(cb) {};
/** @param {function(?)} cb */ QUnit.testDone = function(cb) {};
/** @param {function(?)} cb */ QUnit.moduleStart = function(cb) {};
/** @param {function(?)} cb */ QUnit.moduleDone = function(cb) {};
/** @param {function(?)} cb */ QUnit.log = function(cb) {};
/** @return {boolean} */ QUnit.equiv = function(a, b) {};

/** @constructor */
function QUnitAssert() {}
QUnitAssert.prototype.ok = function(v, m) {};
QUnitAssert.prototype.notOk = function(v, m) {};
QUnitAssert.prototype.equal = function(x, y, m) {};
QUnitAssert.prototype.equals = function(x, y, m) {};
QUnitAssert.prototype.notEqual = function(x, y, m) {};
QUnitAssert.prototype.strictEqual = function(x, y, m) {};
QUnitAssert.prototype.notStrictEqual = function(x, y, m) {};
QUnitAssert.prototype.deepEqual = function(x, y, m) {};
QUnitAssert.prototype.same = function(x, y, m) {};
QUnitAssert.prototype.notDeepEqual = function(x, y, m) {};
QUnitAssert.prototype.propEqual = function(x, y, m) {};
QUnitAssert.prototype.notPropEqual = function(x, y, m) {};
QUnitAssert.prototype.throws = function(f, e, m) {};
QUnitAssert.prototype.raises = function(f, e, m) {};
QUnitAssert.prototype.expect = function(n) {};
/** @return {function()} */ QUnitAssert.prototype.async = function(n) {};
QUnitAssert.prototype.pushResult = function(r) {};
QUnitAssert.prototype.step = function(m) {};
QUnitAssert.prototype.verifySteps = function(s, m) {};
QUnitAssert.prototype.timeout = function(ms) {};

/** @type {?} */ var lodashStable;
/** @type {?} */ var lodashBizarro;
/** @type {?} */ var ui;
/** @type {?} */ var require;
/** @type {?} */ var module;
/** @type {?} */ var exports;
/** @type {?} */ var define;
/** @type {?} */ var global;
/** @type {?} */ var process;
/** @type {?} */ var phantom;
/** @type {?} */ var callPhantom;
/** @type {?} */ var Buffer;
/** @type {?} */ var WinRTError;
/** @type {?} */ var __coverage__;
/** @type {?} */ var toString;
/** @type {?} */ var lodashModule;
/** @type {?} */ var shimmedModule;
/** @type {?} */ var underscoreModule;
/** @param {...*} var_args */ function print(var_args) {}

/**
 * @fileoverview A test file.
 * @suppress {missingRequire}
 */

/** @license Copyright 2026 Foo
 *    Bar baz.
 */

/*! important comment
 *  kept */

/**
 * Says hello.
 * @param {string} name The name to greet.
 * @param {number=} opt_count How many times.
 * @return {!Array<string>} The greetings.
 * @template T
 * @deprecated Use something else.
 */
function hello(name, opt_count) {}

/** @type {?function(this:Foo, ...number): (string|undefined)} */
var x = /** number */ (y);

/** @enum {string} */
var E = {A: 'a'};

/** @typedef {{a: number, b: (string|null), c}} */
var T;

/** @constructor @extends {Base} @implements {I1} @implements {I1} */
function C() {}

/** @param {Array.<string> */
function bad1(a) {}

/** @unknownTag @private {number} desc */
var z;

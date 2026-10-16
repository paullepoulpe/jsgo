/**
 * @fileoverview Test file.
 * @suppress {checkTypes}
 */

/** @const {number} */
var x = 1;

/**
 * @param {string} a
 * @return {string}
 */
function f(a) { return /** @type {string} */ (a); }

/** @constructor @extends {Base} */
function C() {}

/** @enum {string} */
var E = {A: 'a', B: 'b'};

/** @typedef {{a: number}} */
var T;

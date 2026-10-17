goog.module('my.mod');
const Other = goog.require('other.Mod');
/** @const */
exports.VALUE = 1;
/** @param {!Array<number>} arr @return {number} */
exports.sum = function(arr) { return arr.reduce((a, b) => a + b, 0); };
/** @suppress {checkTypes} */
const /** !Object */ mismatch = 1;

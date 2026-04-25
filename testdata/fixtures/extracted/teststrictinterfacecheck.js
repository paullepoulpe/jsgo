/** @interface */
function EventTarget() {}
/** @constructor
 * @implements {EventTarget} */
function Node() {}
/** @type {number} */ Node.prototype.index;
/** @param {EventTarget} x
 * @return {string} */
function foo(x) { return x.index; }

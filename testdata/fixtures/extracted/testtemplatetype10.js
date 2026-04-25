/** @constructor
 * @template T
 */
function Bar() {}

/** @type {!Bar<?>} */ var x;
/** @type {!Bar<number>} */ var y;
y = x;

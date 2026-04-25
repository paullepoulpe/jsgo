/** @constructor
 * @template S,T
 */
function Bar() {}
/**
 * @param {Bar<T>} bar
 * @return {T}
 * @template T
 */
function fn(bar) {}
/** @param {Bar<number>} bar */ function g(bar) {
  /** @type {!Object} */ var x = fn(bar);
}

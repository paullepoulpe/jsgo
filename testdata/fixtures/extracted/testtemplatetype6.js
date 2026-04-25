/**
 * @param {Array<T>} arr
 * @param {?function(T)} f
 * @return {T}
 * @template T
 */
function fn(arr, f) { return arr[0]; }
/** @param {Array<number>} arr */ function g(arr) {
  /** @type {!Object} */ var x = fn.call(null, arr, null);
}

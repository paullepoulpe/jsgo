/**
 * @param {!Object} x
 * @return {!Object}
 */
 function useNonNull(x) {
   return x;
 };
/** @type {?Object} */ var x;
x ?? useNonNull(x);

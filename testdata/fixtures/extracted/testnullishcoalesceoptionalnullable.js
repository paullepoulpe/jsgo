/**
 * @param {?Object=} optionalNullableObject
 * @return {!Object}
 */
function f(optionalNullableObject = undefined) {
 return optionalNullableObject ?? {}; // always returns !Object type
};

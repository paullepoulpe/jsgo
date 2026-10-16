var internalName;

/**
 * @param {number} a
 * @return {!Iterator<number>}
 */
internalName = function *(a) { yield a; };
goog.exportSymbol('externalName', internalName)

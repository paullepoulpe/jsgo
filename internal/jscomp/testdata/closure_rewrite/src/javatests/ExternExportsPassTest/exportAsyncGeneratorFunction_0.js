var internalName;

/**
 * @param {number} a
 * @return {!AsyncGenerator<number>}
 */
internalName = async function *(a) { yield a; };
goog.exportSymbol('externalName', internalName)

var internalName;

/**
 * @param {number} a
 * @return {!Promise<number>}
 */
internalName = async function(a) { return a; };
goog.exportSymbol('externalName', internalName)

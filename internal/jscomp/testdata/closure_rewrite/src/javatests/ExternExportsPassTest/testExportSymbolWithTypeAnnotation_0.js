var internalName;
/**
 * @param {string} param1
 * @param {number} param2
 * @return {string}
 */
internalName = function(param1, param2) {
return param1 + param2;
};
goog.exportSymbol('externalName', internalName)

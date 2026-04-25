/** @constructor */ var T = function() {};
/** @param {T|string} x
@return {T} */
var f = function(x) {
if (x instanceof T) { return x; } else { return new T(); }
};

/** @constructor */ var T = function() {};
/** @param {!T|string} x
@return {string} */
var f = function(x) {
if (x instanceof T) { return ''; } else { return x; }
};

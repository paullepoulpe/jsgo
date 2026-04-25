var a = {};
/** @constructor */ a.N = function() {};
a.N.prototype.p = 1;
/** @constructor */ a.S = function() {};
a.S.prototype.p = 'a';
/** @param {!a.N|!a.S} x
@return {string} */
var f = function(x) { return x.p; };

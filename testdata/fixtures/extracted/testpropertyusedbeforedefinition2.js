var n = {};
/** @constructor */ n.T = function() {};
/** @return {string} */
n.T.prototype.f = function() { return this.g(); };
/** @return {number} */ n.T.prototype.g = function() { return 1; };

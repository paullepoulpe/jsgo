/** @constructor */ var T = function() {};
/** @return {string} */
T.prototype.f = function() { return this.g(); };
/** @return {number} */ T.prototype.g = function() { return 1; };

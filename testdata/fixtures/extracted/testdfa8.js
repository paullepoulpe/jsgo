/** @constructor */ var T = function() {};
/** @type {number|string} */ T.prototype.x = '';
function u() {}
/** @param {!T} t
@return {number} */ function f(t) {
if (u()) { t.x = 1; } else { t.x = 2; } return t.x;
}

/** @param {Object} o
@param {string} x */
function s1(o,x) { o.x = x; }
/** @param {Object} o
@return {string} */
function g1(o) { return typeof o.x == 'undefined' ? '' : o.x; }
/** @param {Object} o
@param {number} x */
function s2(o,x) { o.x = x; }
/** @param {Object} o
@return {number} */
function g2(o) { return typeof o.x == 'undefined' ? 0 : o.x; }

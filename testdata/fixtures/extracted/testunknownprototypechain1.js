/**
* @param {Object} co
 * @return {Object}
 */
function inst(co) {
 /** @constructor */
 var c = function() {};
 c.prototype = co.prototype;
 return new c;
}

/** @param {string} x
 * @constructor */
var Foo = function(x) { this.bar.call(null, x); };
/** @type {function(number)} */ Foo.prototype.bar;

/**
 * @param {?} d
 * @param {?} e
 * @param {?} f
 * @constructor
 */
var foobar = function(d, e, f) {
};
/**
 * @return {undefined}
 * @this {(typeof a.b.c)}
 */
foobar.exportedStaticMethod = function() {
};
/**
 * @param {string} a
 * @return {number}
 * @this {!a.b.c}
 */
foobar.prototype.exportedMethod = function(a) {
};
foobar.prototype.thisProp;

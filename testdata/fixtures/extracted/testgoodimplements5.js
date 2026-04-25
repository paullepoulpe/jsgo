/** @interface */
goog.Disposable = function() {};
/** @type {Function} */
goog.Disposable.prototype.dispose = function() {};
/** @implements {goog.Disposable}
 * @constructor */
goog.SubDisposable = function() {};
/** @param {number} key
 * @override */
goog.SubDisposable.prototype.dispose = function(key) {};

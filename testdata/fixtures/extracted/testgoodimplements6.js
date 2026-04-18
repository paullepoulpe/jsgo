var myNullFunction = function() {};
/** @interface */
goog.Disposable = function() {};
/** @return {number} */
goog.Disposable.prototype.dispose = myNullFunction;
/** @implements {goog.Disposable}
 * @constructor */
goog.SubDisposable = function() {};
/** @return {number}
 * @override */
goog.SubDisposable.prototype.dispose = function() { return 0; };

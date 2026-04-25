var goog = {};
/** @type {!Function} */
goog.abstractMethod = function() {};
/** @interface */
goog.Disposable = function() {};
goog.Disposable.prototype.dispose = goog.abstractMethod;
/** @implements {goog.Disposable}
 * @constructor */
goog.SubDisposable = function() {};
/** @inheritDoc */
goog.SubDisposable.prototype.dispose = function() {};

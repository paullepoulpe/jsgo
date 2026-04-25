var goog = {};
/** @type {!Function} */ goog.abstractMethod = function(){};
/** @interface */ var Disposable = function() {};
Disposable.prototype.method = goog.abstractMethod;
/** @implements {Disposable}
 * @constructor */function f() {}

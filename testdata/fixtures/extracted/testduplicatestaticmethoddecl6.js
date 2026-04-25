var goog = goog || {};
goog.foo = function(x) {};
/**
 * @suppress {duplicate}
 * @return {undefined}
 */
goog.foo =
   /** @type {!Function} */ (function(x) {});

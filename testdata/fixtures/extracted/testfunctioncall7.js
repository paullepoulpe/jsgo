/** @param {Function} handler
 * @param {Object} opt_context */
var Foo = function(handler, opt_context) {
  handler.call(opt_context, x);
};

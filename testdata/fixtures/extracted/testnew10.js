var goog = {};
/** @param {Function} opt_constructor */
goog.Foo = function(opt_constructor) {
new (opt_constructor || Array);
}

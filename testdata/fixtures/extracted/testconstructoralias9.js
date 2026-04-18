var goog = {};
/**
 * @param {number} x
 * @constructor */
goog.Foo = function(x) {};
/** @constructor */ goog.FooAlias = goog.Foo;
/** @return {number} */ function foo() {
  return new goog.FooAlias(1); }

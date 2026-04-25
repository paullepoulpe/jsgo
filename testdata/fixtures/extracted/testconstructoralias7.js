var goog = {};
/** @constructor */ goog.Foo = function() {};
/** @constructor */ goog.FooAlias = goog.Foo;
/** @return {number} */ function foo() {
  return new goog.FooAlias(); }

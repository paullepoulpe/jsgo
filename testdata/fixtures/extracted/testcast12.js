var goog = {}; goog.foo = function() {};
function f() { return /** @type {goog.foo} */ (new Object()); }

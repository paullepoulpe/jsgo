/** @interface */var I = function() {};
// Note that we didn't declare the type but we still expect JsCompiler to recognize the
// property.
I.prototype.foo;
var v = /** @type {I} */ (null);
v.foo = 5;

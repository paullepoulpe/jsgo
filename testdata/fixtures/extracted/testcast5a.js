/** @constructor */function foo() {}
/** @constructor */function bar() {}
var barInstance = new bar;
var baz = /** @type {!foo} */(barInstance);

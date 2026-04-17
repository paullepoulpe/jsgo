/** @constructor */function foo() {}
/** @constructor\s
 @extends foo */function bar() {}
var baz = /** @type {!bar} */(new bar);
var baz = /** @type {!foo} */(new foo);
var baz = /** @type {bar} */(new bar);
var baz = /** @type {foo} */(new foo);
var baz = /** @type {!foo} */(new bar);
var baz = /** @type {!bar} */(new foo);
var baz = /** @type {foo} */(new bar);
var baz = /** @type {bar} */(new foo);

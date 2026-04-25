/** @constructor @abstract */
function Bar() {};
/** @return {function(new:Bar)} */
function foo() {}
var Foo = foo();
var f = new Foo;

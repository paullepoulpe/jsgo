function f() { (new Foo(true)); }
/**
 * @constructor
 * @param {number} x */
var YourFoo = function(x) {};
/**
 * @constructor
 * @param {number} x */
var Foo = YourFoo;

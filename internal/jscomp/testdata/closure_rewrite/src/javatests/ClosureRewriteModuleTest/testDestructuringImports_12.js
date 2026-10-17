goog.module('ns.a');

var {Foo: Bar} = goog.require('ns.b');

/** @type {Bar} */
var f = new Bar;

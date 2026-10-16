goog.module('m');
const Foo = goog.requireType('a.Foo');
const {Bar} = goog.requireType('a.Bar');
const Baz = goog.require('a.Baz');
/** @type {Foo} */ let x;
new Foo();
Bar.x;
Baz.y;
function f(Foo) { return Foo; }

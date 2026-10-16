goog.provide('as.Main');
goog.require('goog.asserts');
/** @constructor */
as.Main = function() {};
as.Main.prototype.abs = goog.abstractMethod;
as.Main.prototype.a2 = as.Main.prototype.a3 = goog.abstractMethod;
/** @abstract */
as.Main.prototype.a4 = function() {};
var x = goog.asserts.assert(foo(), 'm');
goog.asserts.assert(x);
var y = goog.asserts.assertString(bar());
goog.asserts.fail('no');
var z = goog.asserts.assert();
class Abs { /** @abstract */ m() {} n() {} }

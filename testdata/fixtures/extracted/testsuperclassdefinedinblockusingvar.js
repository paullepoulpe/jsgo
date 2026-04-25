{
  /** @constructor */
  var Base = function() {};
  /** @param {number} x */
  Base.prototype.baz = function(x) {};
}
/** @constructor @extends {Base} */
var Foo = function() {};
/** @override */
Foo.prototype.baz = function(x) {
  var /** string */ y = x;
};

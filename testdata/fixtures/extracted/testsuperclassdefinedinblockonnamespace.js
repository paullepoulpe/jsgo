/** @const */
var ns = {};
{
  /** @constructor */
  ns.Base = function() {};
  /** @param {number} x */
  ns.Base.prototype.baz = function(x) {};
}
/** @constructor @extends {ns.Base} */
ns.Foo = function() {};
/** @override */
ns.Foo.prototype.baz = function(x) {
  var /** string */ y = x;
};

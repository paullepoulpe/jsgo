goog.module('my.mod');
const Foo = goog.require('my.Foo');
const {bar} = goog.require('my.bar');
/** @const */ exports.x = 1;
/** @param {number} n */
function f(n) {
  var y = n, z = Math.floor(n);
  y = z + 1;
  var w = y.length;
  if (y.indexOf(z) != -1) { return w; }
  var q = Foo.create(), s = bar(q);
  return startMarker() && s;
}
exports.f = f;

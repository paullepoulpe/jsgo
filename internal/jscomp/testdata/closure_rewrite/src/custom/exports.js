/** @export */
function exportedFn(a, b) { return a + b; }
/** @constructor @export */
function Ctor(x) { /** @export {number} */ this.prop = x; }
/** @export */
Ctor.prototype.method = function(y) { return y; };
/** @export */
Ctor.staticFn = function() {};
/** @const */ var ns = {};
/** @export */
ns.sub = {a: 1, b: function() {}};
/** @export */
class Klass {
  /** @export */ method(a, ...rest) {}
  /** @export */ static sm({x}, y = 2) {}
  constructor(p) { /** @export */ this.field = p; }
}
var E = {/** @export */ A: 1, B: 2};
goog.exportSymbol('manual', exportedFn);
goog.exportProperty(Ctor.prototype, 'other', Ctor.prototype.method);
goog.exportSymbol('ns.renamed', ns.sub);
goog.exportProperty(ns.renamed, 'b', ns.sub.b);
goog.inherits(Ctor, Object);

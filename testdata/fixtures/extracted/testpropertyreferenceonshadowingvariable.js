/** @const */
var ns = {};
ns.fn = function(/** number */ a){};
class Namespace {
  fn(/** string */ a){}
}
function test() {
// The local 'ns' shadows the global 'ns'.
  const /** !Namespace */ ns = /** @type {!Namespace} */ (unknown);
// Verify that ns.fn resolves to Namespace.prototype.fn, not the global ns.fn.
  ns.fn(0);
}

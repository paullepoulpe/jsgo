/** @const */
var ns = {};
ns.fn = function(/** number */ a){};
class Namespace {
  fn(/** string */ a){}
}
function test(/** !Namespace */ ns) { // The parameter 'ns' shadows the global 'ns'.
// Verify that ns.fn resolves to Namespace.prototype.fn, not the global ns.fn.
  ns.fn(0);
}

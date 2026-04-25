/** @const */
var ns = {};
ns.fn = function(/** number */ a){};
/** @unrestricted */
class Namespace {}
function test(/** !Namespace */ ns) { // The parameter 'ns' shadows the global 'ns'.
  ns.fn = function(/** string */ a) {};
// Verify that ns.fn resolves to the locally assigned ns.fn, not the global ns.fn.
  ns.fn(0);
}

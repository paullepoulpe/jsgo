/** @unrestricted */
class Namespace {
  fn(/** string */ a) {}
}
function test(/** !Namespace */ ns) {
// Assign ns.fn to a narrow type than Namespace.prototype.fn.
  ns.fn = function(/** string|number */ a) {};
// Verify that ns.fn resolves to the narrower local type, which can accept `number`.
  ns.fn(0);
}

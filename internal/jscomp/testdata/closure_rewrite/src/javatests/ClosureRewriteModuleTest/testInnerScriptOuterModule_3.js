goog.provide('A.B');
/** @constructor */
A.B = function() {};
function main() {
// Note A.B was NOT written to module$exports$A.B
  /** @type {A.B} */
  var l = new A.B();
}

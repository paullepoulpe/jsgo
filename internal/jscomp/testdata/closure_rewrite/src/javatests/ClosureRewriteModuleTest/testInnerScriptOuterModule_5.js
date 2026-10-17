goog.provide('A.b.c.D');
/** @constructor */
A.b.c.L = function() {}; // Note L not D
function main() {
// Note A.b.c.L was NOT written to module$exports$A.b.c.L.
  /** @type {A.b.c.L} */
  var l = new A.b.c.L();
}

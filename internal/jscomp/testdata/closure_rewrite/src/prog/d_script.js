goog.require('p.C');
goog.require('p.A');
goog.scope(function() {
  var C = p.C;
  var Sub = p.A.Sub;
  /** @type {!Sub} */
  C.sub = new Sub();
});
var x = goog.module.get('p.B');
/** @type {p.B} */
var y = new x();
/** @export */
p.C.exported = function() {};

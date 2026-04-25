/** @constructor */
var C = function() { }; // (1)

var fn = function() {
  /** @type {?C} */ var innerC; // (2)

  /** @constructor */
  var C = function() { }; // (3)

  innerC = new C(); // (4)
}

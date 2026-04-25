/** @constructor */
var C = function() { }; // (1)

var fn = function() {
  /** @constructor */
  var C = function() { }; // (2)

// This function will be hoisted above, and therefore type-analysed before, (2).
  /** @return {!C} */ // (3)
  function hoisted() { return new C(); }; // (4)
}

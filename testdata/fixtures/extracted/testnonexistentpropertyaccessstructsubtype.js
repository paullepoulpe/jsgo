/**
 * @constructor
 * @struct
 */
var A = function() {};
/**
 * @constructor
 * @struct
 * @extends {A}
 */
var B = function() { this.bar = function(){}; };
/** @param {A} a */
function foo(a) {
  if (a.bar) { a.bar(); }
}

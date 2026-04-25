/**
 * @interface
 * @struct
 */
var A = function() {};

/**
 * @interface
 * @struct
 * @extends {A}
 */
var B = function() {};
/** @return {void} */ B.prototype.bar = function(){};

/** @param {A} a */
function foo(a) {
  if (a.bar) { a.bar(); }
}

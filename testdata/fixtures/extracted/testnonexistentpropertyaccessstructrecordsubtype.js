/**
 * @record
 * @struct
 */
var A = function() {};

/**
 * @record
 * @struct
 * @extends {A}
 */
var B = function() {};
/** @return {void} */ B.prototype.bar = function(){};

/** @param {A} a */
function foo(a) {
  if (a.bar) {
    a.bar();
  }
}

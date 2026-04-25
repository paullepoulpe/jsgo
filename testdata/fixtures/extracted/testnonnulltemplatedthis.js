/** @constructor */
function C() {}

/**
  @return {THIS}
  @this {THIS}
  @template THIS
*/
C.prototype.method = function() {};

/** @return {C|null} */
function f() {
  return x;
};

/** @type {string} */ var s = f().method();

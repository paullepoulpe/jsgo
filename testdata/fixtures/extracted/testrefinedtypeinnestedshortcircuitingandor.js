/** @constructor */ function B() {}
/** @constructor @extends {B} */ function C() {}
/** @return {string} */ C.prototype.foo = function() {};
/** @return {boolean} */ C.prototype.bar = function() {};
/** @constructor @extends {B} */ function D() {}
/** @return {number} */ D.prototype.foo = function() {};
/** @param {!B} arg
    @return {null} */
function f(arg) {
  if ((arg instanceof C && arg.bar()) || arg instanceof D) {
    return arg.foo();
  }
  return null;
}

/** @constructor */ var C = function() {
  /** @type {string} */ this.a = 'str';
};
/** @type {C} */ var x = new C();
/** @param {?} a
 @return {number} */
C.prototype.f = function(a) {
  x.a = a;
  return x.a;
}

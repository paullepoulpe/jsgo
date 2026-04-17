/** @constructor */ var C = function() {
  /** @type {string} */ this.a = 'str'};
/** @param {?} a
 @return {number} */
C.prototype.f = function(a) {
  this.a = a;
  return this.a;
}

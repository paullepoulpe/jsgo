/** @type {!Function} */ var f = function() {
  /** @type {number} */ this.prop = 3;
};
/**
 * @param {Object} x
 * @return {string}
 */ var g = function(x) { return x.prop; };

/** @constructor */ function F() {}
/** jsdoc
 * @return {number} */ F.prototype.bar = function() {
  return 3;
};
/** jsdoc
 * @return {string}
 * @suppress {duplicate} */
F.prototype.bar = function() { return ''; };

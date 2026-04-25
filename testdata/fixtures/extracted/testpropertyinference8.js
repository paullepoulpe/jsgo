/** @constructor */ function F() {
  /** @type {string} */ this.x_ = 'x';
}
(new F).x_ = 3;
/** @return {string} */
F.prototype.bar = function() { return this.x_; };

/** @param {number} x */ function f(x) {}
/** @constructor */ function Ctor() {
  /** @this {Date} */
  this.method = function() {
    f(this);
  };
}

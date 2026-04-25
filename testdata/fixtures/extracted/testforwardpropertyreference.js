/** @constructor */ var Foo = function() { this.init(); };
/** @return {string} */
Foo.prototype.getString = function() {
  return this.number_;
};
Foo.prototype.init = function() {
  /** @type {number} */
  this.number_ = 3;
};

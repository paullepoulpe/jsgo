/** @constructor */ function F() { }
F.prototype.baz = function() { this.x_ = 3; };
/** @return {string} */
F.prototype.bar = function() { if (this.x_) return this.x_; };

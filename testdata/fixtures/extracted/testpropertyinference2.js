/** @constructor */ function F() { this.x_ = true; }
F.prototype.baz = function() { this.x_ = null; };
/** @return {string} */
F.prototype.bar = function() { if (this.x_) return this.x_; };

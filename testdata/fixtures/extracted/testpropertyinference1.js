/** @constructor */ function F() { this.x_ = true; }
/** @return {string} */
F.prototype.bar = function() { if (this.x_) return this.x_; };

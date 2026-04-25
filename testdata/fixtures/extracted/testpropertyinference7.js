/** @constructor */ function F() { this.x_ = true; }
(new F).x_ = 3;
/** @return {string} */
F.prototype.bar = function() { return this.x_; };

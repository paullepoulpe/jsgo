/** @constructor */ function Point(x, y) { /** @type {number} */ this.x = x; /** @type {number} */ this.y = y; }
Point.prototype.getX = function() { return this.x; };
window['res'] = new Point(10, 20).getX();

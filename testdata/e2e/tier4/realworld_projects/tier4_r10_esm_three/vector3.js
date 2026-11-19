/** @constructor */ function Vector3(x, y, z) { this.x = x; this.y = y; this.z = z; }
Vector3.prototype.dot = function(v) { return this.x * v.x + this.y * v.y + this.z * v.z; };
window['dotRes'] = new Vector3(1, 2, 3).dot(new Vector3(4, 5, 6));

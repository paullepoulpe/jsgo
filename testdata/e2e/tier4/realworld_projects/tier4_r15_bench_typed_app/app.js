/** @constructor */ function User(name, age) { /** @type {string} */ this.name = name; /** @type {number} */ this.age = age; }
/** @return {string} */ User.prototype.info = function() { return this.name + ' (' + this.age + ')'; };
var u = new User('Alice', 30);

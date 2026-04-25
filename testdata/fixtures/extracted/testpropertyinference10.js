/** @constructor */ function A() {}
/** @return {function(): ?} */ function f() {
  return function() {};
}
var g = f();
/** @type {number} */ g.prototype.bar_ = 1;
var h = f();
/** @type {string} */ h.prototype.bar_ = 1;

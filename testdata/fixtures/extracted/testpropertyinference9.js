/** @constructor */ function A() {}
/** @return {function(): ?} */ function f() {
  return function() {};
}
var g = f();
/** @type {number} */ g.prototype.bar_ = null;

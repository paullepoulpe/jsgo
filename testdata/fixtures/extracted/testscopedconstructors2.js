/** @param {Function} f */
function foo1(f) {
  /** @param {Function} g */
  f.prototype.bar = function(g) {};
}

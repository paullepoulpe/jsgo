/** @constructor */ function C() {
  /** @type {void} */ this.x;
}
/**
* @param {T} x
* @param {function(this:T, ...)} y
* @template T
*/
function f(x, y) {}
f(new C, /** @param {number} a */ function(a) {this.x = a;});

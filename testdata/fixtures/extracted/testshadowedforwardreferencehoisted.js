/** @constructor */ var C = function() { /** @type {string} */ this.prop = 's'; };
var fn = function() {
  /** @type {C} */ var x = f();
  /** @type {number} */ var n1 = x.prop;
  /** @type {number} */ var n2 = new C().prop;
  /** @constructor */ function C() { /** @type {number} */ this.prop = 1; };
  /** @return {C} */ function fn() { return new C(); };

}

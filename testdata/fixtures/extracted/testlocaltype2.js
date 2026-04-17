/** @constructor */ function C(){ /** @const */ this.a = true;}
function f() {
  // shadow
  /** @constructor */ function C(){ /** @const */ this.a = 1;}
  /** @type {null} */ var x = new C().a;
}

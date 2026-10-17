/** @constructor */
function C() {
  /** @private */
  this.foo_ = 1;
  /** @private */
  this.bar_ = 2;
  /** @private @const */
  this.baz_ = 3;
  /** @private */
  this['qux_'] = 4;
  /** @private */
  this.fn_ = function() {};
}
C.prototype.inc = function() { this.bar_++; };
/** @private */
C.static_ = 5;
delete C.static_;
class D {
  constructor() {
    /** @private */
    this.p_ = [];
    /** @private */
    this.q_ = 0;
  }
  m() { this.q_ = 1; return this.p_; }
}
/** @private */
D.s_ = 1;
/** @private @typedef {number} */
D.T_;
/** @interface */
function I() {}
/** @private */
I.prototype.x_;
/** @private */
I.y_ = 2;
var obj = {a: 1, get b() { return 1; }, set c(v) {}, d() {}};
JSCompiler_renameProperty('zz_', obj);
goog.reflect.objectProperty('yy_', obj);
/** @private */
D.prototype.zz_ = 1;
/** @private */
D.prototype.yy_ = 2;
class E { /** @private */ m_() {} /** @private */ n_() {} k() { this.n_(); } }
x = (this.a_ += 1);

//// file: w5/access/base.js
/** @constructor */
function Base() {
  /** @private {number} */ this.priv_ = 1;
  /** @protected {number} */ this.prot = 2;
  /** @package {number} */ this.pkg = 3;
  /** @const {number} */ this.k = 4;
  /** @final {number} */ this.fin = 5;
}
/** @private */ Base.prototype.privMethod_ = function() {};
/** @protected */ Base.prototype.protMethod = function() {};
/** @package */ Base.prototype.pkgMethod = function() {};
/** @deprecated Use something else. */ Base.prototype.old = function() {};
/** @deprecated */ Base.prototype.oldProp;
/** @deprecated */ function oldFn() {}
/** @constructor @deprecated Gone. */ function OldCtor() {}
/** @private */ var privGlobal_ = 1;
/** @const */ var CONST_GLOBAL = {};
/** @const */ CONST_GLOBAL.x = 1;
/** @protected */ Base.staticProt = function() {};

class Klass {
  constructor() {
    /** @private {number} */ this.p_ = 1;
    /** @const {number} */ this.c = 2;
  }
  /** @protected */ prot() {}
  /** @private */ static sp_() {}
  /** @package */ pk() {}
  /** @deprecated Old. */ get dep() { return 1; }
}
/** @final */ class FinalK {}
/** @constructor @final */ function FinalCtor() {}

var b0 = new Base();
b0.priv_;
b0.k = 11;

//// file: w5/access/user.js
var b = new Base();
b.priv_;
b.prot;
b.pkg;
b.k = 10;
b.fin = 6;
b.privMethod_();
b.protMethod();
b.pkgMethod();
b.old();
b.oldProp;
oldFn();
new OldCtor();
privGlobal_;
CONST_GLOBAL.x = 2;
delete b.k;
Base.staticProt();

/** @constructor @extends {Base} */
function Sub() {
  Base.call(this);
  this.prot;
  this.priv_ = 3;
  this.fin = 7;
}
/** @override */ Sub.prototype.protMethod = function() {};
/** @override */ Sub.prototype.privMethod_ = function() {};
/** @public @override */ Sub.prototype.pkgMethod = function() {};
/** @private @override */ Sub.prototype.old = function() {};

class K2 extends Klass {
  constructor() {
    super();
    this.c = 3;
    this.p_;
  }
  /** @override */ prot() { super.prot(); }
}
class F2 extends FinalK {}
/** @constructor @extends {FinalCtor} */ function F3() {}
new Klass().prot();
Klass.sp_();
new Klass().pk();
new Klass().dep;

/** @deprecated Do not use. */
function deprecatedCaller() {
  oldFn();
  b.old();
}

//// file: w5/access/other/pkg_user.js
new Base().pkg;
new Base().pkgMethod();
new Klass().pk();
/** @suppress {visibility} */
function suppressed() {
  new Base().priv_;
  new Base().pkg;
}

//// file: w5/access/fileoverview.js
/** @fileoverview @package */
/** @constructor */ function PkgByDefault() {}
PkgByDefault.prototype.method = function() {};
/** @public */ PkgByDefault.prototype.pub = function() {};
var pkgGlobal = 1;

//// file: w5/access/other/fileoverview_user.js
new PkgByDefault().method();
new PkgByDefault().pub();
pkgGlobal;
/** @override */ PkgByDefault.prototype.method = function() {};

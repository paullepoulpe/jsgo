goog.provide('custom.exports');
/** @export */ function exportedFn() {}
/** @export */ var exportedVar = 1;
/** @export */ custom.exports.prop = function() {};
/** @export */ custom.exports.Klass = class { /** @export */ method() {} /** @export */ static sm() {} };
/** @export */ class TopClass { /** @export */ m() {} /** @export */ static x = 1; /** @export */ y = 2; }
var o = { /** @export */ key: 1, /** @export */ get g() { return 1; }, /** @export */ set s(v) {} };
function local() { /** @export */ var notGlobal = 1; /** @export */ this.p = 1; }
/** @export */ Foo.prototype.myprop;
/** @export */ x = y = 5;
window.a = 1; window['b'] = 2; self.c = 3; goog.global.d = 4; top['e'] = 5; this.f = 6;
goog$global.g = 7; globalThis.h = 8; $jscomp.global.i = 9; window?.j; window?.['k'];
goog.exportSymbol('exp.sym', 1);
goog.exportSymbol('top', 2);

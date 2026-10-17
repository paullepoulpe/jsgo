/**
 * @fileoverview Hand-written input of the w5-checks oracle: CheckJSDoc and
 * CheckJSDocStyle diagnostics.
 */

/** @type {number} */ function misplacedType() {}
/* @type {number} */ var inBlockComment;
/** @suppress {visibility} */ var suppressOnVar = 1;
/** @suppress {visibility} */ suppressOnVar = 2;
/** @constructor */ var Arrow = () => {};
/** @param {number} a */ function def(a = 1) {}
/** @param {number=} a */ function defOk(a = 1) {}
/** @nosideeffects */ function nse() {}
/** @modifies {this} */ function mod() {}
/** @define {boolean} */ let DEF = true;
function ret() { /** @type {number} */ return 1; }
/** @desc not a message */ var notMsg = 1;
/** @param {number} r */ function rest(...r) {}
/** @param {...number} r */ function restOk(...r) {}
/** @param {number} a @param {number} b */ function wrongNum(a) {}
/** @param {number} b */ function wrongName(a) {}
function missingJsdoc(a) { return a; }
/** @param {number} x */ function optNotMarked(opt_x) {}
/** Use @code{x} here. */ function backticks() {}
/** @param {number} a */ function mixed(/** number */ a, b) {}
/** @return {number} */ function missingReturn() {}
/** @param {number} a */ function missingReturnJsdoc(a) { return a; }
/** @constructor */ class DisallowedOnClass {}
/** @abstract */ function abstractFn() {}
/** @enum {string} */ var E = {A: 'a'};
var obj = {/** @type {number} */ a: 1, /** @export */ b: 2};
var ns = {};
/** @private */ ns.priv_ = 1;
/** @deprecated */ var dep = 1;
/** @expose */ ns.exposed = 1;
/** @return {number} */ var notAFunction = 1;

class C {
  /** @param {number} x */ get g() { return 1; }
  /** @return {number} */ set s(v) {}
  /** @type {number} */ method() {}
  /** @type {number} */ field = 1;
  constructor() {
    /** @type {number} */ this.x = 1;
  }
}

/**
 * Has a param
 * @param {string} name The name.
 * @return {string}
 */
function documented(name) {
  return name;
}

/** @interface */ function I() {}
/** @param {number} x */ I.prototype.m = function(x) {};

/** @const */ var v = /** @type {number} */ (1);
var cast = /** @type {string} */ ('x');
var notCast = /** @type {string} */ 'x';

/** @param {number} a */
var arrowFn = (a) => a;
/** @this {Object} */ var arrowThis = () => {};

if (true) {
  /** @type {number} */ var inBlock = 1;
}
/** @type {number} */ (function() {})();

/**
 * @fileoverview A file with every kind of JSDoc the printer knows.
 * @suppress {missingRequire,checkTypes} Because.
 */

/** @const */
var ns = {};

/** @const */
ns.sub = ns.sub || {};

/**
 * Block description of the class.
 * More description.
 * @constructor
 * @struct
 * @extends {Base}
 * @implements {Iface}
 * @implements {Other.Iface}
 * @param {number} a The a param
 *     continued.
 * @param {string=} b Optional.
 * @param {...*} var_args Rest.
 * @template T, U
 * @template {!Array<number>} V
 */
function Klass(a, b, var_args) {}

/**
 * @interface
 * @extends {A}
 * @extends {B}
 */
function I() {}

/** @record */
function R() {}

/**
 * @return {?Object|undefined} The thing.
 * @throws {Error} When bad.
 * @this {Klass}
 * @deprecated Use other.
 * @nosideeffects
 */
Klass.prototype.method = function () {};

/** @private @const {number} */
Klass.prototype.priv_ = 1;

/** @protected */ Klass.prototype.prot;
/** @package */ Klass.prototype.pack;
/** @public */ Klass.prototype.pub;
/** @export */ Klass.prototype.exp;
/** @override */ Klass.prototype.over = function () {};
/** @final */ Klass.prototype.fin;
/** @abstract */ Klass.prototype.abs = function () {};

/** @typedef {{a: number, b: (string|undefined), c}} */
var Rec;

/** @enum {string} */
var E = {A: 'a', B: 'b'};

/** @define {boolean} */
var DEBUG = true;

/** @type {function(new:Klass, number, ...string): void} */
var ctor;

/** @type {function(this:Klass, ?): ?} */
var fn;

/** @type {function(): number} */
var fn2;

/** @type {function()} */
var fn3;

/** @type {!Map<string, !Array<?number>>} */
var map;

/** @type {typeof Klass} */
var tof;

/** @type {*} */ var star;
/** @type {?} */ var q;
/** @type {void} */ var vd;
/** @type {undefined} */ var und;
/** @dict */ var dict;
/** @unrestricted */ var unr;
/** @nocollapse */ var noc;
/** @noinline */ var noi;
/** @implicitCast */ var ic;
/** @ngInject */ function inj() {}
/** @closurePrimitive {asserts.truthy} */ function prim() {}
/** @tsType ():string */ var ts;
/** @polymer @polymerBehavior @mixinFunction @mixinClass @customElement */ var poly;
/** @lends {Klass.prototype} */ var lends = {};
/** @desc A message. */ var MSG_X = goog.getMsg('x');
/**
 * @see http://example.com
 * @author someone@example.com (Some One)
 */
var seen;
/** @idGenerator {consistent} */ var idg;
/** @idGenerator {stable} */ var idg2;
/** @idGenerator {xid} */ var idg3;
/** @idGenerator {mapped} */ var idg4;
/** @idGenerator */ var idg5;
/** @wizaction */ function wiz() {}
/** @externs */
/** @template T := cond(isUnknown(T), 'a', 'b') =: */ function ttl() {}

var cast = /** @type {number} */ (x);
var cast2 = /** @type {!Array} */ ([]);
function inline(/** number */ a, /** string= */ b, /** ...number */ c) {}
function inlineRet(a) /** number */ {}
var arrowDoc = /** @return {number} */ () => 1;
class C2 {
  /** @param {number} x */
  constructor(x) {
    /** @const {number} */
    this.x = x;
  }
  /** @return {string} */
  get g() { return ''; }
  /** @type {number} */
  field = 1;
}

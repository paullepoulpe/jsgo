// Port of com/google/javascript/jscomp/CommandLineRunner's externs
// bundle + com/google/javascript/jscomp/DefaultExterns.js
// @ closure-compiler v20260415 (narrow; hand-curated subset).
//
// Closure ships a collection of externs .js files (es6.js, html5.js,
// w3c_dom1.js, …) that declare the shapes of the runtime environment
// for the type checker. jsgo's pipeline loads this one synthetic
// source as the "externs scope" parent of every user source —
// matches how Compiler#compile accepts both externs and sources.
// The content below is deliberately narrow: it declares the names
// the TypeCheckTest corpus references at checker level (Disposable,
// IArrayLike, IThenable, common DOM globals) without committing to a
// full port of closure-compiler's externs .jar.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// parseExternsSource parses the externs bundle into a SCRIPT node.
// Split out so callers can substitute alternate extern bundles once
// the real externs loader lands. Returns nil on parse errors — the
// caller is expected to fall back to the builtins scope.
func parseExternsSource(src string) (*rhino.Node, error) {
	return parser.Parse(src)
}

// defaultExternsSrc is the JS source that every compile's externs
// scope is built from. Declarations here become the parent scope of
// the user source, so VarCheck sees them as already-declared.
//
// Format rules:
//   - Every top-level binding must carry a @const or @constructor or
//     function declaration so TypedScopeCreator registers it.
//   - Prefer @interface / @record for duck-typed shapes (IArrayLike,
//     IThenable) so structural subtyping applies.
//   - Use valid JS syntax — the source goes through the same parser
//     as user code.
const defaultExternsSrc = `
// ---- ES5 / primitives ------------------------------------------ //

/** @const */ var undefined;
/** @const {number} */ var Infinity;
/** @const {number} */ var NaN;
/** @const */ var globalThis;
/** @const */ var arguments;

/** @return {number} */ function parseInt(s, opt_radix) {}
/** @return {number} */ function parseFloat(s) {}
/** @return {boolean} */ function isNaN(n) {}
/** @return {boolean} */ function isFinite(n) {}
/** @return {string} */ function encodeURI(s) {}
/** @return {string} */ function encodeURIComponent(s) {}
/** @return {string} */ function decodeURI(s) {}
/** @return {string} */ function decodeURIComponent(s) {}
/** @param {string} s @return {*} */ function eval(s) {}

// ---- Core constructors ----------------------------------------- //

/** @constructor */ function Object(opt_v) {}
Object.prototype.hasOwnProperty = function(key) {};
Object.prototype.toString = function() {};
Object.prototype.valueOf = function() {};
Object.keys = function(obj) {};
Object.values = function(obj) {};
Object.entries = function(obj) {};
Object.assign = function(target, var_args) {};
Object.create = function(proto, opt_props) {};
Object.freeze = function(obj) {};
Object.defineProperty = function(obj, prop, descriptor) {};
Object.getPrototypeOf = function(obj) {};
Object.setPrototypeOf = function(obj, proto) {};
Object.getOwnPropertyNames = function(obj) {};
Object.getOwnPropertyDescriptor = function(obj, prop) {};
Object.getOwnPropertyDescriptors = function(obj) {};
Object.getOwnPropertySymbols = function(obj) {};
Object.defineProperties = function(obj, props) {};
Object.is = function(a, b) {};
Object.fromEntries = function(iterable) {};
Object.isFrozen = function(obj) {};
Object.isSealed = function(obj) {};
Object.seal = function(obj) {};

/** @constructor */ function Function(var_args) {}
Function.prototype.call = function(thisArg, var_args) {};
Function.prototype.apply = function(thisArg, opt_args) {};
Function.prototype.bind = function(thisArg, var_args) {};

/**
 * @constructor
 * @template T
 * @param {...*} var_args
 * @return {!Array<T>}
 * @implements {IArrayLike<T>}
 */
function Array(var_args) {}
/** @type {number} */ Array.prototype.length;
/** @param {...T} var_args @return {number} */ Array.prototype.push = function(var_args) {};
/** @return {T|undefined} */ Array.prototype.pop = function() {};
/** @return {T|undefined} */ Array.prototype.shift = function() {};
/** @param {...T} var_args @return {number} */ Array.prototype.unshift = function(var_args) {};
/** @param {T} searchElement @param {number=} fromIndex @return {number} */ Array.prototype.indexOf = function(searchElement, fromIndex) {};
/** @param {T} searchElement @param {number=} fromIndex @return {number} */ Array.prototype.lastIndexOf = function(searchElement, fromIndex) {};
/** @param {number=} begin @param {number=} end @return {!Array<T>} */ Array.prototype.slice = function(begin, end) {};
/** @param {number} start @param {number=} deleteCount @param {...T} var_args @return {!Array<T>} */ Array.prototype.splice = function(start, deleteCount, var_args) {};
/** @param {...(!Array<T>|T)} var_args @return {!Array<T>} */ Array.prototype.concat = function(var_args) {};
/** @param {string=} separator @return {string} */ Array.prototype.join = function(separator) {};
/** @return {!Array<T>} */ Array.prototype.reverse = function() {};
/** @param {function(T,T):number=} compareFn @return {!Array<T>} */ Array.prototype.sort = function(compareFn) {};
/** @param {function(*): ?} callback @param {*=} thisArg @return {void} */ Array.prototype.forEach = function(callback, thisArg) {};
/** @template U @param {function(*): U} callback @param {*=} thisArg @return {!Array<U>} */ Array.prototype.map = function(callback, thisArg) {};
/** @param {function(*): boolean} callback @param {*=} thisArg @return {!Array<T>} */ Array.prototype.filter = function(callback, thisArg) {};
/** @template U @param {function(U, T): U} callback @param {U=} initialValue @return {U} */ Array.prototype.reduce = function(callback, initialValue) {};
/** @template U @param {function(U, T): U} callback @param {U=} initialValue @return {U} */ Array.prototype.reduceRight = function(callback, initialValue) {};
/** @param {function(*): boolean} callback @param {*=} thisArg @return {boolean} */ Array.prototype.some = function(callback, thisArg) {};
/** @param {function(*): boolean} callback @param {*=} thisArg @return {boolean} */ Array.prototype.every = function(callback, thisArg) {};
/** @param {function(*): boolean} callback @param {*=} thisArg @return {T|undefined} */ Array.prototype.find = function(callback, thisArg) {};
/** @param {function(*): boolean} callback @param {*=} thisArg @return {number} */ Array.prototype.findIndex = function(callback, thisArg) {};
/** @param {T} value @param {number=} start @param {number=} end @return {!Array<T>} */ Array.prototype.fill = function(value, start, end) {};
/** @param {T} searchElement @param {number=} fromIndex @return {boolean} */ Array.prototype.includes = function(searchElement, fromIndex) {};
/** @param {number=} depth @return {!Array<*>} */ Array.prototype.flat = function(depth) {};
/** @template U @param {function(T,number,!Array<T>):(U|!Array<U>)} callback @return {!Array<U>} */ Array.prototype.flatMap = function(callback) {};
/** @return {string} */ Array.prototype.toString = function() {};
/** @param {!IArrayLike<*>} arrayLike @return {!Array<*>} */ Array.from = function(arrayLike) {};
/** @param {...*} var_args @return {!Array<*>} */ Array.of = function(var_args) {};
/** @param {*} value @return {boolean} */ Array.isArray = function(value) {};

/** @constructor @param {*=} arg @return {string} */ function String(arg) {}
String.prototype.charAt = function(i) {};
String.prototype.charCodeAt = function(i) {};
String.prototype.indexOf = function(s, opt_from) {};
String.prototype.lastIndexOf = function(s, opt_from) {};
String.prototype.slice = function(start, opt_end) {};
String.prototype.substring = function(start, opt_end) {};
String.prototype.substr = function(start, opt_len) {};
String.prototype.toLowerCase = function() {};
String.prototype.toUpperCase = function() {};
String.prototype.trim = function() {};
String.prototype.trimStart = function() {};
String.prototype.trimEnd = function() {};
String.prototype.split = function(sep, opt_limit) {};
String.prototype.replace = function(search, replace) {};
String.prototype.match = function(regexp) {};
String.prototype.startsWith = function(s, opt_pos) {};
String.prototype.endsWith = function(s, opt_pos) {};
String.prototype.includes = function(s, opt_pos) {};
String.prototype.padStart = function(len, opt_fill) {};
String.prototype.padEnd = function(len, opt_fill) {};
String.prototype.repeat = function(n) {};
String.prototype.concat = function(var_args) {};
String.prototype.toString = function() {};
String.prototype.valueOf = function() {};
String.prototype.at = function(i) {};
String.prototype.replaceAll = function(search, replace) {};
String.prototype.matchAll = function(regexp) {};
String.prototype.codePointAt = function(i) {};
String.prototype.normalize = function(opt_form) {};
String.fromCharCode = function(var_args) {};
String.fromCodePoint = function(var_args) {};
String.raw = function(template, var_args) {};

/** @constructor @param {*=} arg @return {number} */ function Number(arg) {}
Number.prototype.toFixed = function(digits) {};
Number.prototype.toString = function(opt_radix) {};
Number.prototype.valueOf = function() {};
Number.isFinite = function(n) {};
Number.isInteger = function(n) {};
Number.isNaN = function(n) {};
Number.parseFloat = function(s) {};
Number.parseInt = function(s, opt_radix) {};
/** @const {number} */ Number.MAX_VALUE;
/** @const {number} */ Number.MIN_VALUE;
/** @const {number} */ Number.MAX_SAFE_INTEGER;
/** @const {number} */ Number.MIN_SAFE_INTEGER;
/** @const {number} */ Number.POSITIVE_INFINITY;
/** @const {number} */ Number.NEGATIVE_INFINITY;
/** @const {number} */ Number.NaN;
/** @const {number} */ Number.EPSILON;
Number.isSafeInteger = function(n) {};
Number.prototype.toPrecision = function(digits) {};
Number.prototype.toExponential = function(digits) {};

/** @constructor @param {*=} arg @return {boolean} */ function Boolean(arg) {}
Boolean.prototype.toString = function() {};
Boolean.prototype.valueOf = function() {};

/** @constructor */ function Date(opt_year, opt_month, opt_day) {}
/** @return {number} */ Date.prototype.getTime = function() {};
/** @return {number} */ Date.prototype.getFullYear = function() {};
/** @return {number} */ Date.prototype.getMonth = function() {};
/** @return {number} */ Date.prototype.getDate = function() {};
/** @return {number} */ Date.prototype.getDay = function() {};
/** @return {number} */ Date.prototype.getHours = function() {};
/** @return {number} */ Date.prototype.getMinutes = function() {};
/** @return {number} */ Date.prototype.getSeconds = function() {};
/** @return {number} */ Date.prototype.getMilliseconds = function() {};
/** @return {number} */ Date.prototype.getTimezoneOffset = function() {};
/** @return {number} */ Date.prototype.getUTCFullYear = function() {};
/** @return {number} */ Date.prototype.getUTCMonth = function() {};
/** @return {number} */ Date.prototype.getUTCDate = function() {};
/** @return {number} */ Date.prototype.getUTCDay = function() {};
/** @return {number} */ Date.prototype.getUTCHours = function() {};
/** @return {number} */ Date.prototype.getUTCMinutes = function() {};
/** @return {number} */ Date.prototype.getUTCSeconds = function() {};
/** @return {number} */ Date.prototype.getUTCMilliseconds = function() {};
/** @param {number} value @return {number} */ Date.prototype.setTime = function(value) {};
/** @param {number} year @param {number=} month @return {number} */ Date.prototype.setFullYear = function(year, month) {};
/** @param {number} month @return {number} */ Date.prototype.setMonth = function(month) {};
/** @param {number} date @return {number} */ Date.prototype.setDate = function(date) {};
/** @param {number} hours @return {number} */ Date.prototype.setHours = function(hours) {};
/** @param {number} minutes @return {number} */ Date.prototype.setMinutes = function(minutes) {};
/** @param {number} seconds @return {number} */ Date.prototype.setSeconds = function(seconds) {};
/** @param {number} ms @return {number} */ Date.prototype.setMilliseconds = function(ms) {};
/** @return {string} */ Date.prototype.toString = function() {};
/** @return {string} */ Date.prototype.toISOString = function() {};
/** @return {string} */ Date.prototype.toDateString = function() {};
/** @return {string} */ Date.prototype.toTimeString = function() {};
/** @return {string} */ Date.prototype.toLocaleDateString = function() {};
/** @return {string} */ Date.prototype.toLocaleTimeString = function() {};
/** @return {string} */ Date.prototype.toLocaleString = function() {};
/** @return {string} */ Date.prototype.toUTCString = function() {};
/** @return {number} */ Date.prototype.valueOf = function() {};
/** @return {number} */ Date.now = function() {};
/** @param {string} dateString @return {number} */ Date.parse = function(dateString) {};

/** @constructor */ function RegExp(pattern, opt_flags) {}
/** @type {string} */ RegExp.prototype.source;
/** @type {string} */ RegExp.prototype.flags;
/** @type {boolean} */ RegExp.prototype.global;
/** @type {boolean} */ RegExp.prototype.ignoreCase;
/** @type {boolean} */ RegExp.prototype.multiline;
/** @type {boolean} */ RegExp.prototype.sticky;
/** @type {boolean} */ RegExp.prototype.unicode;
/** @type {number} */ RegExp.prototype.lastIndex;
/** @param {string} s @return {boolean} */ RegExp.prototype.test = function(s) {};
/** @param {string} s @return {Array<string>} */ RegExp.prototype.exec = function(s) {};
/** @return {string} */ RegExp.prototype.toString = function() {};

/** @constructor */ function Error(opt_msg) {}
/** @type {string} */ Error.prototype.message;
/** @type {string} */ Error.prototype.stack;
/** @type {string} */ Error.prototype.name;
/** @return {string} */ Error.prototype.toString = function() {};
/** @constructor @extends {Error} */ function TypeError(opt_msg) {}
/** @constructor @extends {Error} */ function RangeError(opt_msg) {}
/** @constructor @extends {Error} */ function SyntaxError(opt_msg) {}
/** @constructor @extends {Error} */ function ReferenceError(opt_msg) {}
/** @constructor @extends {Error} */ function URIError(opt_msg) {}
/** @constructor @extends {Error} */ function EvalError(opt_msg) {}

// ---- ES6 iterable / iterator / promise ------------------------- //

/** @record */ function IArrayLike() {}
/** @type {number} */ IArrayLike.prototype.length;

/** @interface @template T */ function Iterable() {}
/** @interface @template VALUE,RETURN,NEXT */ function Iterator() {}
/** @interface @template VALUE,RETURN,NEXT */ function IteratorResult() {}
/** @interface @template VALUE,RETURN,NEXT */ function Generator() {}
/** @interface @template VALUE,RETURN,NEXT */ function IteratorIterable() {}
/** @interface @template T */ function AsyncIterator() {}
/** @interface @template T */ function AsyncIterable() {}
/** @interface @template VALUE,RETURN,NEXT */ function AsyncGenerator() {}
/** @interface @template VALUE,RETURN,NEXT */ function AsyncIteratorIterable() {}

/**
 * @record
 * @template VALUE
 */
function IThenable() {}
/** @param {VALUE} value @return {IThenable<VALUE>} */
IThenable.prototype.then = function(opt_onFulfilled, opt_onRejected) {};

/**
 * @constructor
 * @template TYPE
 * @implements {IThenable<TYPE>}
 * @param {function(function(TYPE),function(*))} executor
 */
function Promise(executor) {}
/**
 * @template RESULT
 * @param {function(TYPE):RESULT=} opt_onFulfilled
 * @param {function(*):*=} opt_onRejected
 * @return {!Promise<RESULT>}
 */
Promise.prototype.then = function(opt_onFulfilled, opt_onRejected) {};
/** @param {function(*):*} onRejected @return {!Promise} */
Promise.prototype.catch = function(onRejected) {};
/** @param {function():*} onFinally @return {!Promise} */
Promise.prototype.finally = function(onFinally) {};
/** @param {TYPE} value @return {!Promise<TYPE>} @template TYPE */
Promise.resolve = function(value) {};
/** @param {*=} reason @return {!Promise} */
Promise.reject = function(reason) {};
/** @param {!Iterable} iterable @return {!Promise} */
Promise.all = function(iterable) {};
/** @param {!Iterable} iterable @return {!Promise} */
Promise.race = function(iterable) {};
/** @param {!Iterable} iterable @return {!Promise} */
Promise.allSettled = function(iterable) {};
/** @param {!Iterable} iterable @return {!Promise} */
Promise.any = function(iterable) {};

/**
 * @constructor
 * @template KEY, VALUE
 * @param {!Iterable=} opt_iter
 */
function Map(opt_iter) {}
/** @param {KEY} key @return {VALUE|undefined} */ Map.prototype.get = function(key) {};
/** @param {KEY} key @param {VALUE} value @return {!Map<KEY,VALUE>} */ Map.prototype.set = function(key, value) {};
/** @param {KEY} key @return {boolean} */ Map.prototype.has = function(key) {};
/** @param {KEY} key @return {boolean} */ Map.prototype.delete = function(key) {};
/** @return {void} */ Map.prototype.clear = function() {};
/** @param {function(VALUE,KEY,!Map<KEY,VALUE>)} cb */ Map.prototype.forEach = function(cb) {};
/** @return {!Iterator} */ Map.prototype.keys = function() {};
/** @return {!Iterator} */ Map.prototype.values = function() {};
/** @return {!Iterator} */ Map.prototype.entries = function() {};
/** @type {number} */ Map.prototype.size;

/**
 * @constructor
 * @template VALUE
 * @param {!Iterable=} opt_iter
 */
function Set(opt_iter) {}
/** @param {VALUE} value @return {!Set<VALUE>} */ Set.prototype.add = function(value) {};
/** @param {VALUE} value @return {boolean} */ Set.prototype.has = function(value) {};
/** @param {VALUE} value @return {boolean} */ Set.prototype.delete = function(value) {};
/** @return {void} */ Set.prototype.clear = function() {};
/** @param {function(VALUE,VALUE,!Set<VALUE>)} cb */ Set.prototype.forEach = function(cb) {};
/** @return {!Iterator} */ Set.prototype.values = function() {};
/** @return {!Iterator} */ Set.prototype.keys = function() {};
/** @return {!Iterator} */ Set.prototype.entries = function() {};
/** @type {number} */ Set.prototype.size;

/**
 * @constructor
 * @template KEY, VALUE
 * @param {!Iterable=} opt_iter
 */
function WeakMap(opt_iter) {}
/** @param {KEY} key @return {VALUE|undefined} */ WeakMap.prototype.get = function(key) {};
/** @param {KEY} key @param {VALUE} value @return {!WeakMap<KEY,VALUE>} */ WeakMap.prototype.set = function(key, value) {};
/** @param {KEY} key @return {boolean} */ WeakMap.prototype.has = function(key) {};
/** @param {KEY} key @return {boolean} */ WeakMap.prototype.delete = function(key) {};

/**
 * @constructor
 * @template VALUE
 * @param {!Iterable=} opt_iter
 */
function WeakSet(opt_iter) {}
/** @param {VALUE} value @return {!WeakSet<VALUE>} */ WeakSet.prototype.add = function(value) {};
/** @param {VALUE} value @return {boolean} */ WeakSet.prototype.has = function(value) {};
/** @param {VALUE} value @return {boolean} */ WeakSet.prototype.delete = function(value) {};

/**
 * @constructor
 * @template T
 * @param {T} target
 */
function WeakRef(target) {}
/** @return {T|undefined} */ WeakRef.prototype.deref = function() {};

/**
 * @constructor
 * @template T
 * @param {function(T): void} callback
 */
function FinalizationRegistry(callback) {}
/** @param {T} target @param {*} heldValue @param {*=} unregisterToken */ FinalizationRegistry.prototype.register = function(target, heldValue, unregisterToken) {};
/** @param {*} unregisterToken @return {boolean} */ FinalizationRegistry.prototype.unregister = function(unregisterToken) {};

/** @constructor */ function Proxy(target, handler) {}
/** @const */ var Reflect;

// Symbol matches Closure's es3.js externs: @constructor + @return {symbol}.
// new Symbol('x') throws at runtime, but the type system treats Symbol as
// a constructor so !Symbol resolves to the wrapper instance type used by
// SymbolType.testForEquality / canTestForShallowEqualityWith.
/** @constructor @return {symbol} */ function Symbol(opt_desc) {}
/** @type {symbol} */ Symbol.iterator;
/** @type {symbol} */ Symbol.hasInstance;
/** @type {symbol} */ Symbol.isConcatSpreadable;
/** @type {symbol} */ Symbol.toPrimitive;
/** @type {symbol} */ Symbol.toStringTag;
/** @type {symbol} */ Symbol.asyncIterator;
/** @type {symbol} */ Symbol.match;
/** @type {symbol} */ Symbol.replace;
/** @type {symbol} */ Symbol.search;
/** @type {symbol} */ Symbol.species;
/** @type {symbol} */ Symbol.split;
/** @type {symbol} */ Symbol.unscopables;
Symbol.for = function(key) {};
Symbol.keyFor = function(sym) {};
/** @constructor */ function ArrayBuffer(opt_len) {}
/** @type {number} */ ArrayBuffer.prototype.byteLength;
/** @param {number=} begin @param {number=} end @return {!ArrayBuffer} */ ArrayBuffer.prototype.slice = function(begin, end) {};
/** @param {*} arg @return {boolean} */ ArrayBuffer.isView = function(arg) {};

/** @constructor */ function DataView(buf, opt_off, opt_len) {}
/** @type {!ArrayBuffer} */ DataView.prototype.buffer;
/** @type {number} */ DataView.prototype.byteLength;
/** @type {number} */ DataView.prototype.byteOffset;

/** @constructor */ function Int8Array(v) {}
/** @type {number} */ Int8Array.prototype.length;
/** @type {!ArrayBuffer} */ Int8Array.prototype.buffer;
/** @type {number} */ Int8Array.prototype.byteOffset;
/** @type {number} */ Int8Array.prototype.byteLength;
/** @const {number} */ Int8Array.BYTES_PER_ELEMENT;
/** @param {number=} begin @param {number=} end @return {!Int8Array} */ Int8Array.prototype.slice = function(begin, end) {};
/** @param {function(number,number,!Int8Array):void} cb @return {void} */ Int8Array.prototype.forEach = function(cb) {};

/** @constructor */ function Uint8Array(v) {}
/** @type {number} */ Uint8Array.prototype.length;
/** @type {!ArrayBuffer} */ Uint8Array.prototype.buffer;
/** @type {number} */ Uint8Array.prototype.byteOffset;
/** @type {number} */ Uint8Array.prototype.byteLength;
/** @const {number} */ Uint8Array.BYTES_PER_ELEMENT;
/** @param {number=} begin @param {number=} end @return {!Uint8Array} */ Uint8Array.prototype.slice = function(begin, end) {};
/** @param {function(number,number,!Uint8Array):void} cb @return {void} */ Uint8Array.prototype.forEach = function(cb) {};

/** @constructor */ function Uint8ClampedArray(v) {}
/** @type {number} */ Uint8ClampedArray.prototype.length;
/** @type {!ArrayBuffer} */ Uint8ClampedArray.prototype.buffer;

/** @constructor */ function Int16Array(v) {}
/** @type {number} */ Int16Array.prototype.length;
/** @type {!ArrayBuffer} */ Int16Array.prototype.buffer;
/** @const {number} */ Int16Array.BYTES_PER_ELEMENT;

/** @constructor */ function Uint16Array(v) {}
/** @type {number} */ Uint16Array.prototype.length;
/** @type {!ArrayBuffer} */ Uint16Array.prototype.buffer;
/** @const {number} */ Uint16Array.BYTES_PER_ELEMENT;

/** @constructor */ function Int32Array(v) {}
/** @type {number} */ Int32Array.prototype.length;
/** @type {!ArrayBuffer} */ Int32Array.prototype.buffer;
/** @const {number} */ Int32Array.BYTES_PER_ELEMENT;

/** @constructor */ function Uint32Array(v) {}
/** @type {number} */ Uint32Array.prototype.length;
/** @type {!ArrayBuffer} */ Uint32Array.prototype.buffer;
/** @const {number} */ Uint32Array.BYTES_PER_ELEMENT;

/** @constructor */ function Float32Array(v) {}
/** @type {number} */ Float32Array.prototype.length;
/** @type {!ArrayBuffer} */ Float32Array.prototype.buffer;
/** @const {number} */ Float32Array.BYTES_PER_ELEMENT;
/** @param {number=} begin @param {number=} end @return {!Float32Array} */ Float32Array.prototype.slice = function(begin, end) {};
/** @param {function(number,number,!Float32Array):void} cb @return {void} */ Float32Array.prototype.forEach = function(cb) {};

/** @constructor */ function Float64Array(v) {}
/** @type {number} */ Float64Array.prototype.length;
/** @type {!ArrayBuffer} */ Float64Array.prototype.buffer;
/** @const {number} */ Float64Array.BYTES_PER_ELEMENT;
/** @param {number=} begin @param {number=} end @return {!Float64Array} */ Float64Array.prototype.slice = function(begin, end) {};
/** @param {function(number,number,!Float64Array):void} cb @return {void} */ Float64Array.prototype.forEach = function(cb) {};

/** @constructor */ function BigInt(opt_v) {}
/** @constructor */ function BigInt64Array(v) {}
/** @type {number} */ BigInt64Array.prototype.length;
/** @type {!ArrayBuffer} */ BigInt64Array.prototype.buffer;

/** @constructor */ function BigUint64Array(v) {}
/** @type {number} */ BigUint64Array.prototype.length;
/** @type {!ArrayBuffer} */ BigUint64Array.prototype.buffer;

// ---- Closure-library-shaped names used by the corpus ---------- //

/** @interface */ function Disposable() {}
/** @interface */ function IDisposable() {}

// ---- Thenable — structural duck-type for promise-like values --- //
// Java's externs ship this as a @typedef so it resolves as a
// non-nullable structural type (no implied-null from nominal lookup).
/** @typedef {{then: ?}} */
var Thenable;

// ---- Math / JSON stubs ---------------------------------------- //

/** @const @type {{
  E: number,
  LN10: number,
  LN2: number,
  LOG2E: number,
  LOG10E: number,
  PI: number,
  SQRT1_2: number,
  SQRT2: number,
  abs: function(number): number,
  ceil: function(number): number,
  floor: function(number): number,
  round: function(number): number,
  trunc: function(number): number,
  sign: function(number): number,
  sqrt: function(number): number,
  cbrt: function(number): number,
  pow: function(number, number): number,
  exp: function(number): number,
  expm1: function(number): number,
  log: function(number): number,
  log2: function(number): number,
  log10: function(number): number,
  log1p: function(number): number,
  sin: function(number): number,
  cos: function(number): number,
  tan: function(number): number,
  asin: function(number): number,
  acos: function(number): number,
  atan: function(number): number,
  atan2: function(number, number): number,
  sinh: function(number): number,
  cosh: function(number): number,
  tanh: function(number): number,
  asinh: function(number): number,
  acosh: function(number): number,
  atanh: function(number): number,
  hypot: function(...number): number,
  max: function(...number): number,
  min: function(...number): number,
  random: function(): number,
  clz32: function(number): number,
  imul: function(number, number): number,
  fround: function(number): number
}} */ var Math;

/** @const */ var JSON;
JSON.parse = function(text, opt_reviver) {};
JSON.stringify = function(value, opt_replacer, opt_space) {};

// ---- Minimal DOM / BOM --------------------------------------- //

/** @const */ var window;
/** @const */ var document;
/** @const */ var console;
/** @const */ var navigator;
/** @const */ var location;
/** @const */ var history;
/** @const */ var localStorage;
/** @const */ var sessionStorage;
/** @const */ var self;
/** @const */ var top;

/** @return {void} */ function alert(msg) {}
/** @return {boolean} */ function confirm(msg) {}
/** @return {?string} */ function prompt(msg, opt_default) {}
/** @return {number} */ function setTimeout(cb, opt_ms) {}
/** @return {number} */ function setInterval(cb, opt_ms) {}
/** @return {void} */ function clearTimeout(id) {}
/** @return {void} */ function clearInterval(id) {}
/** @return {number} */ function requestAnimationFrame(cb) {}
/** @return {void} */ function cancelAnimationFrame(id) {}
/** @return {*} */ function fetch(url, opt_init) {}

/** @constructor */ function XMLHttpRequest() {}

/** @interface */
function EventTarget() {}
EventTarget.prototype.addEventListener = function(type, listener, opt_capture) {};
EventTarget.prototype.removeEventListener = function(type, listener, opt_capture) {};
EventTarget.prototype.dispatchEvent = function(evt) {};

/** @constructor */ function Event(type, opt_init) {}
/** @constructor @extends {Event} */ function MouseEvent(type, opt_init) {}
/** @constructor @extends {Event} */ function KeyboardEvent(type, opt_init) {}
/** @constructor @extends {Event} */ function UIEvent(type, opt_init) {}
/** @constructor @extends {Event} */ function FocusEvent(type, opt_init) {}
/** @constructor @extends {Event} */ function InputEvent(type, opt_init) {}

/** @constructor @implements {EventTarget} */ function Node() {}
/** @constructor @extends {Node} */ function Element() {}
/** @type {string} */ Element.prototype.id;
/** @type {string} */ Element.prototype.className;
/** @type {string} */ Element.prototype.tagName;
/** @type {string} */ Element.prototype.innerHTML;
/** @type {string} */ Element.prototype.outerHTML;
/** @type {string} */ Element.prototype.textContent;
Element.prototype.getAttribute = function(name) {};
Element.prototype.setAttribute = function(name, value) {};
Element.prototype.removeAttribute = function(name) {};
Element.prototype.querySelector = function(selector) {};
Element.prototype.querySelectorAll = function(selector) {};
Element.prototype.appendChild = function(node) {};
Element.prototype.removeChild = function(node) {};
Element.prototype.addEventListener = function(type, listener, opt_capture) {};
Element.prototype.removeEventListener = function(type, listener, opt_capture) {};

/** @constructor @extends {Element} */ function HTMLElement() {}
/** @type {string} */ HTMLElement.prototype.style;
/** @type {string} */ HTMLElement.prototype.innerText;
/** @type {string} */ HTMLElement.prototype.title;
/** @type {boolean} */ HTMLElement.prototype.hidden;
/** @type {boolean} */ HTMLElement.prototype.draggable;
/** @type {string} */ HTMLElement.prototype.lang;
/** @type {string} */ HTMLElement.prototype.dir;
/** @type {!DOMStringMap} */ HTMLElement.prototype.dataset;
/** @type {number} */ HTMLElement.prototype.tabIndex;
/** @type {number} */ HTMLElement.prototype.offsetWidth;
/** @type {number} */ HTMLElement.prototype.offsetHeight;
/** @type {number} */ HTMLElement.prototype.offsetTop;
/** @type {number} */ HTMLElement.prototype.offsetLeft;
/** @type {number} */ HTMLElement.prototype.scrollTop;
/** @type {number} */ HTMLElement.prototype.scrollLeft;
/** @type {number} */ HTMLElement.prototype.scrollWidth;
/** @type {number} */ HTMLElement.prototype.scrollHeight;
/** @type {number} */ HTMLElement.prototype.clientWidth;
/** @type {number} */ HTMLElement.prototype.clientHeight;
HTMLElement.prototype.click = function() {};
/** @param {FocusOptions=} options */ HTMLElement.prototype.focus = function(options) {};
HTMLElement.prototype.blur = function() {};
/** @return {!DOMRect} */ HTMLElement.prototype.getBoundingClientRect = function() {};
/** @param {*=} options */ HTMLElement.prototype.scrollIntoView = function(options) {};
/** @return {boolean} */ HTMLElement.prototype.contains = function(other) {};
HTMLElement.prototype.remove = function() {};
/** @param {string} cls */ HTMLElement.prototype.dispatchEvent = function(evt) {};

/** @constructor @extends {HTMLElement} */ function HTMLDivElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLSpanElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLInputElement() {}
/** @type {string} */ HTMLInputElement.prototype.value;
/** @type {string} */ HTMLInputElement.prototype.type;
/** @type {boolean} */ HTMLInputElement.prototype.checked;
/** @constructor @extends {HTMLElement} */ function HTMLButtonElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLAnchorElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLFormElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLSelectElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLTextAreaElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLImageElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLScriptElement() {}

/** @constructor @extends {Node} @implements {EventTarget} */
function Document() {}
/** @type {!HTMLElement} */ Document.prototype.body;
/** @type {!HTMLElement} */ Document.prototype.head;
/** @type {string} */ Document.prototype.title;
Document.prototype.getElementById = function(id) {};
Document.prototype.getElementsByClassName = function(cls) {};
Document.prototype.getElementsByTagName = function(tag) {};
Document.prototype.querySelector = function(selector) {};
Document.prototype.querySelectorAll = function(selector) {};
Document.prototype.createElement = function(tagName) {};
Document.prototype.createTextNode = function(data) {};
Document.prototype.createDocumentFragment = function() {};

/** @constructor @extends {Document} */
function Window() {}
/** @type {!Document} */ Window.prototype.document;
/** @type {!Window} */ Window.prototype.window;
/** @type {string} */ Window.prototype.location;
/** @type {number} */ Window.prototype.innerWidth;
/** @type {number} */ Window.prototype.innerHeight;
Window.prototype.alert = function(msg) {};
Window.prototype.confirm = function(msg) {};
Window.prototype.setTimeout = function(cb, opt_ms) {};
Window.prototype.setInterval = function(cb, opt_ms) {};
Window.prototype.clearTimeout = function(id) {};
Window.prototype.clearInterval = function(id) {};
Window.prototype.addEventListener = function(type, listener, opt_capture) {};
Window.prototype.removeEventListener = function(type, listener, opt_capture) {};
/** @return {*} */ Window.prototype.fetch = function(url, opt_init) {};
/** @return {number} */ Window.prototype.requestAnimationFrame = function(cb) {};
/** @return {void} */ Window.prototype.cancelAnimationFrame = function(id) {};
/** @return {void} */ Window.prototype.cancelTimeout = function(id) {};
/** @return {void} */ Window.prototype.postMessage = function(msg, targetOrigin, opt_transfer) {};
Window.prototype.dispatchEvent = function(evt) {};
Window.prototype.prompt = function(msg, opt_default) {};

// ---- ES2022 / Web Platform additions ----------------------------- //

/** @constructor */ function DOMRect() {}
/** @type {number} */ DOMRect.prototype.x;
/** @type {number} */ DOMRect.prototype.y;
/** @type {number} */ DOMRect.prototype.width;
/** @type {number} */ DOMRect.prototype.height;
/** @type {number} */ DOMRect.prototype.top;
/** @type {number} */ DOMRect.prototype.right;
/** @type {number} */ DOMRect.prototype.bottom;
/** @type {number} */ DOMRect.prototype.left;

/** @constructor */ function DOMStringMap() {}

/** @constructor */ function Response() {}
/** @type {boolean} */ Response.prototype.ok;
/** @type {number} */ Response.prototype.status;
/** @type {string} */ Response.prototype.statusText;
/** @type {string} */ Response.prototype.url;
/** @type {boolean} */ Response.prototype.redirected;
/** @type {string} */ Response.prototype.type;
/** @return {!Promise} */ Response.prototype.json = function() {};
/** @return {!Promise} */ Response.prototype.text = function() {};
/** @return {!Promise} */ Response.prototype.arrayBuffer = function() {};
/** @return {!Promise} */ Response.prototype.blob = function() {};
/** @return {!Promise} */ Response.prototype.formData = function() {};

/** @constructor */ function Request(input, opt_init) {}
/** @type {string} */ Request.prototype.url;
/** @type {string} */ Request.prototype.method;

/** @constructor */ function Headers(opt_init) {}
/** @param {string} name @return {string|null} */ Headers.prototype.get = function(name) {};
/** @param {string} name @param {string} value */ Headers.prototype.set = function(name, value) {};
/** @param {string} name @param {string} value */ Headers.prototype.append = function(name, value) {};
/** @param {string} name */ Headers.prototype.delete = function(name) {};
/** @param {string} name @return {boolean} */ Headers.prototype.has = function(name) {};

/** @constructor */ function URL(url, opt_base) {}
/** @type {string} */ URL.prototype.href;
/** @type {string} */ URL.prototype.origin;
/** @type {string} */ URL.prototype.protocol;
/** @type {string} */ URL.prototype.host;
/** @type {string} */ URL.prototype.hostname;
/** @type {string} */ URL.prototype.port;
/** @type {string} */ URL.prototype.pathname;
/** @type {string} */ URL.prototype.search;
/** @type {string} */ URL.prototype.hash;
/** @type {string} */ URL.prototype.username;
/** @type {string} */ URL.prototype.password;
/** @return {string} */ URL.prototype.toString = function() {};
URL.createObjectURL = function(obj) {};
URL.revokeObjectURL = function(url) {};

/** @constructor */ function URLSearchParams(opt_init) {}
/** @param {string} name @return {string|null} */ URLSearchParams.prototype.get = function(name) {};
/** @param {string} name @return {!Array<string>} */ URLSearchParams.prototype.getAll = function(name) {};
/** @param {string} name @param {string} value */ URLSearchParams.prototype.set = function(name, value) {};
/** @param {string} name @param {string} value */ URLSearchParams.prototype.append = function(name, value) {};
/** @param {string} name */ URLSearchParams.prototype.delete = function(name) {};
/** @param {string} name @return {boolean} */ URLSearchParams.prototype.has = function(name) {};
/** @return {string} */ URLSearchParams.prototype.toString = function() {};

/** @constructor */ function AbortController() {}
/** @type {!AbortSignal} */ AbortController.prototype.signal;
AbortController.prototype.abort = function() {};

/** @constructor */ function AbortSignal() {}
/** @type {boolean} */ AbortSignal.prototype.aborted;

/** @constructor */ function MutationObserver(callback) {}
MutationObserver.prototype.observe = function(target, options) {};
MutationObserver.prototype.disconnect = function() {};

/** @constructor */ function IntersectionObserver(callback, opt_options) {}
IntersectionObserver.prototype.observe = function(target) {};
IntersectionObserver.prototype.unobserve = function(target) {};
IntersectionObserver.prototype.disconnect = function() {};

/** @constructor */ function ResizeObserver(callback) {}
ResizeObserver.prototype.observe = function(target, opt_options) {};
ResizeObserver.prototype.unobserve = function(target) {};
ResizeObserver.prototype.disconnect = function() {};

/** @constructor */ function CustomEvent(type, opt_init) {}
/** @type {*} */ CustomEvent.prototype.detail;

/** @constructor @extends {HTMLElement} */ function HTMLCanvasElement() {}
/** @type {number} */ HTMLCanvasElement.prototype.width;
/** @type {number} */ HTMLCanvasElement.prototype.height;
HTMLCanvasElement.prototype.getContext = function(type, opt_attrs) {};
HTMLCanvasElement.prototype.toDataURL = function(opt_type, opt_quality) {};

/** @constructor @extends {HTMLElement} */ function HTMLVideoElement() {}
/** @type {string} */ HTMLVideoElement.prototype.src;
/** @type {number} */ HTMLVideoElement.prototype.currentTime;
/** @type {number} */ HTMLVideoElement.prototype.duration;
/** @type {boolean} */ HTMLVideoElement.prototype.paused;
/** @type {number} */ HTMLVideoElement.prototype.volume;
HTMLVideoElement.prototype.play = function() {};
HTMLVideoElement.prototype.pause = function() {};

/** @constructor @extends {HTMLElement} */ function HTMLAudioElement() {}
/** @type {string} */ HTMLAudioElement.prototype.src;
/** @type {boolean} */ HTMLAudioElement.prototype.paused;
/** @type {number} */ HTMLAudioElement.prototype.volume;
HTMLAudioElement.prototype.play = function() {};
HTMLAudioElement.prototype.pause = function() {};

`

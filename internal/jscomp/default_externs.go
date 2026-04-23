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

/** @constructor */ function Function(var_args) {}
Function.prototype.call = function(thisArg, var_args) {};
Function.prototype.apply = function(thisArg, opt_args) {};
Function.prototype.bind = function(thisArg, var_args) {};

// Array intentionally omitted here — jsgo's builtins() scope
// ships an Array stub with length/push/pop/indexOf/slice/concat
// populated so property-access checks find them.
// Additional Array.prototype methods:
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
String.fromCharCode = function(var_args) {};

/** @constructor @param {*=} arg @return {number} */ function Number(arg) {}
Number.prototype.toFixed = function(digits) {};
Number.prototype.toString = function(opt_radix) {};
Number.prototype.valueOf = function() {};
Number.isFinite = function(n) {};
Number.isInteger = function(n) {};
Number.isNaN = function(n) {};
Number.parseFloat = function(s) {};
Number.parseInt = function(s, opt_radix) {};

/** @constructor @param {*=} arg @return {boolean} */ function Boolean(arg) {}
Boolean.prototype.toString = function() {};
Boolean.prototype.valueOf = function() {};

/** @constructor */ function Date(opt_year, opt_month, opt_day) {}
Date.prototype.getTime = function() {};
Date.prototype.getFullYear = function() {};
Date.prototype.getMonth = function() {};
Date.prototype.getDate = function() {};
Date.prototype.toString = function() {};
Date.now = function() {};

/** @constructor */ function RegExp(pattern, opt_flags) {}
RegExp.prototype.test = function(s) {};
RegExp.prototype.exec = function(s) {};

/** @constructor */ function Error(opt_msg) {}
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

/** @constructor */ function Proxy(target, handler) {}
/** @const */ var Reflect;

/**
 * @constructor
 * @param {string=} opt_desc
 */
function Symbol(opt_desc) {}
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
/** @constructor */ function DataView(buf, opt_off, opt_len) {}
/** @constructor */ function Int8Array(v) {}
/** @constructor */ function Uint8Array(v) {}
/** @constructor */ function Uint8ClampedArray(v) {}
/** @constructor */ function Int16Array(v) {}
/** @constructor */ function Uint16Array(v) {}
/** @constructor */ function Int32Array(v) {}
/** @constructor */ function Uint32Array(v) {}
/** @constructor */ function Float32Array(v) {}
/** @constructor */ function Float64Array(v) {}
/** @constructor */ function BigInt(opt_v) {}
/** @constructor */ function BigInt64Array(v) {}
/** @constructor */ function BigUint64Array(v) {}

// ---- Closure-library-shaped names used by the corpus ---------- //

/** @interface */ function Disposable() {}
/** @interface */ function IDisposable() {}

// ---- Thenable — structural duck-type for promise-like values --- //
// Java's externs ship this as a @typedef so it resolves as a
// non-nullable structural type (no implied-null from nominal lookup).
/** @typedef {{then: ?}} */
var Thenable;

// ---- Math / JSON stubs ---------------------------------------- //

/** @const */ var Math;

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
HTMLElement.prototype.click = function() {};
HTMLElement.prototype.focus = function() {};
HTMLElement.prototype.blur = function() {};

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
`

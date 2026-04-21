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
/** @constructor */ function Function(var_args) {}
// Array intentionally omitted here — jsgo's builtins() scope
// ships an Array stub with length/push/pop/indexOf/slice/concat
// populated so property-access checks find them.
/** @constructor @param {*=} arg @return {string} */ function String(arg) {}
/** @constructor @param {*=} arg @return {number} */ function Number(arg) {}
/** @constructor @param {*=} arg @return {boolean} */ function Boolean(arg) {}
/** @constructor */ function Date(opt_year, opt_month, opt_day) {}
/** @constructor */ function RegExp(pattern, opt_flags) {}
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

/** @record */ function Iterable() {}
/** @record */ function Iterator() {}
/** @record */ function IteratorResult() {}
/** @record */ function Generator() {}
/** @record */ function AsyncIterator() {}
/** @record */ function AsyncIterable() {}
/** @record */ function AsyncGenerator() {}

/** @record */ function IThenable() {}
/** @constructor @implements {IThenable} */ function Promise(executor) {}

/** @constructor */ function Map(opt_iter) {}
/** @constructor */ function Set(opt_iter) {}
/** @constructor */ function WeakMap(opt_iter) {}
/** @constructor */ function WeakSet(opt_iter) {}
/** @constructor */ function Proxy(target, handler) {}
/** @const */ var Reflect;

/** @constructor */ function Symbol(opt_desc) {}
/** @constructor */ function ArrayBuffer(opt_len) {}
/** @type {number} */ ArrayBuffer.prototype.byteLength;
/** @return {!ArrayBuffer} */ ArrayBuffer.prototype.slice = function(begin, opt_end) {};
/** @constructor */ function DataView(buf, opt_off, opt_len) {}
/** @type {!ArrayBuffer} */ DataView.prototype.buffer;
/** @type {number} */ DataView.prototype.byteLength;
/** @type {number} */ DataView.prototype.byteOffset;
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

// ---- Typed-array prototype surface ----------------------------- //
//
// DefaultExterns.js declares .length / .buffer / .byteLength /
// .byteOffset / .set / .subarray on every typed-array constructor.
// The corpus only exercises a handful of the numeric typed arrays
// but stubbing each one keeps the checker from falling back to the
// unknown type when fixtures index into them or read .length.

/** @type {number} */ Int8Array.prototype.length;
/** @type {!ArrayBuffer} */ Int8Array.prototype.buffer;
/** @type {number} */ Int8Array.prototype.byteLength;
/** @type {number} */ Int8Array.prototype.byteOffset;
/** @return {void} */ Int8Array.prototype.set = function(array, opt_offset) {};
/** @return {!Int8Array} */ Int8Array.prototype.subarray = function(opt_begin, opt_end) {};

/** @type {number} */ Uint8Array.prototype.length;
/** @type {!ArrayBuffer} */ Uint8Array.prototype.buffer;
/** @type {number} */ Uint8Array.prototype.byteLength;
/** @type {number} */ Uint8Array.prototype.byteOffset;
/** @return {void} */ Uint8Array.prototype.set = function(array, opt_offset) {};
/** @return {!Uint8Array} */ Uint8Array.prototype.subarray = function(opt_begin, opt_end) {};

/** @type {number} */ Uint8ClampedArray.prototype.length;
/** @type {!ArrayBuffer} */ Uint8ClampedArray.prototype.buffer;
/** @type {number} */ Uint8ClampedArray.prototype.byteLength;
/** @type {number} */ Uint8ClampedArray.prototype.byteOffset;
/** @return {void} */ Uint8ClampedArray.prototype.set = function(array, opt_offset) {};
/** @return {!Uint8ClampedArray} */ Uint8ClampedArray.prototype.subarray = function(opt_begin, opt_end) {};

/** @type {number} */ Int16Array.prototype.length;
/** @type {!ArrayBuffer} */ Int16Array.prototype.buffer;
/** @type {number} */ Int16Array.prototype.byteLength;
/** @type {number} */ Int16Array.prototype.byteOffset;
/** @return {void} */ Int16Array.prototype.set = function(array, opt_offset) {};
/** @return {!Int16Array} */ Int16Array.prototype.subarray = function(opt_begin, opt_end) {};

/** @type {number} */ Uint16Array.prototype.length;
/** @type {!ArrayBuffer} */ Uint16Array.prototype.buffer;
/** @type {number} */ Uint16Array.prototype.byteLength;
/** @type {number} */ Uint16Array.prototype.byteOffset;
/** @return {void} */ Uint16Array.prototype.set = function(array, opt_offset) {};
/** @return {!Uint16Array} */ Uint16Array.prototype.subarray = function(opt_begin, opt_end) {};

/** @type {number} */ Int32Array.prototype.length;
/** @type {!ArrayBuffer} */ Int32Array.prototype.buffer;
/** @type {number} */ Int32Array.prototype.byteLength;
/** @type {number} */ Int32Array.prototype.byteOffset;
/** @return {void} */ Int32Array.prototype.set = function(array, opt_offset) {};
/** @return {!Int32Array} */ Int32Array.prototype.subarray = function(opt_begin, opt_end) {};

/** @type {number} */ Uint32Array.prototype.length;
/** @type {!ArrayBuffer} */ Uint32Array.prototype.buffer;
/** @type {number} */ Uint32Array.prototype.byteLength;
/** @type {number} */ Uint32Array.prototype.byteOffset;
/** @return {void} */ Uint32Array.prototype.set = function(array, opt_offset) {};
/** @return {!Uint32Array} */ Uint32Array.prototype.subarray = function(opt_begin, opt_end) {};

/** @type {number} */ Float32Array.prototype.length;
/** @type {!ArrayBuffer} */ Float32Array.prototype.buffer;
/** @type {number} */ Float32Array.prototype.byteLength;
/** @type {number} */ Float32Array.prototype.byteOffset;
/** @return {void} */ Float32Array.prototype.set = function(array, opt_offset) {};
/** @return {!Float32Array} */ Float32Array.prototype.subarray = function(opt_begin, opt_end) {};

/** @type {number} */ Float64Array.prototype.length;
/** @type {!ArrayBuffer} */ Float64Array.prototype.buffer;
/** @type {number} */ Float64Array.prototype.byteLength;
/** @type {number} */ Float64Array.prototype.byteOffset;
/** @return {void} */ Float64Array.prototype.set = function(array, opt_offset) {};
/** @return {!Float64Array} */ Float64Array.prototype.subarray = function(opt_begin, opt_end) {};

// ---- String.prototype completeness ----------------------------- //
//
// jsgo's builtins() scope ships a handful of String-prototype names
// but the ES2015 surface is broader. Declaring here makes the type
// checker accept chain calls like s.trim().startsWith('x') without
// falling to unknown on the second hop.

/** @type {number} */ String.prototype.length;
/** @return {string} */ String.prototype.charAt = function(index) {};
/** @return {number} */ String.prototype.charCodeAt = function(index) {};
/** @return {number} */ String.prototype.codePointAt = function(index) {};
/** @return {string} */ String.prototype.concat = function(var_args) {};
/** @return {boolean} */ String.prototype.endsWith = function(search, opt_end) {};
/** @return {boolean} */ String.prototype.includes = function(search, opt_start) {};
/** @return {number} */ String.prototype.indexOf = function(search, opt_start) {};
/** @return {number} */ String.prototype.lastIndexOf = function(search, opt_start) {};
/** @return {number} */ String.prototype.localeCompare = function(other) {};
/** @return {?Array<string>} */ String.prototype.match = function(regexp) {};
/** @return {string} */ String.prototype.normalize = function(opt_form) {};
/** @return {string} */ String.prototype.padEnd = function(target_len, opt_pad) {};
/** @return {string} */ String.prototype.padStart = function(target_len, opt_pad) {};
/** @return {string} */ String.prototype.repeat = function(count) {};
/** @return {string} */ String.prototype.replace = function(pattern, replacement) {};
/** @return {string} */ String.prototype.replaceAll = function(pattern, replacement) {};
/** @return {number} */ String.prototype.search = function(regexp) {};
/** @return {string} */ String.prototype.slice = function(opt_begin, opt_end) {};
/** @return {!Array<string>} */ String.prototype.split = function(opt_sep, opt_limit) {};
/** @return {boolean} */ String.prototype.startsWith = function(search, opt_start) {};
/** @return {string} */ String.prototype.substring = function(begin, opt_end) {};
/** @return {string} */ String.prototype.substr = function(begin, opt_len) {};
/** @return {string} */ String.prototype.toLowerCase = function() {};
/** @return {string} */ String.prototype.toLocaleLowerCase = function(opt_locale) {};
/** @return {string} */ String.prototype.toUpperCase = function() {};
/** @return {string} */ String.prototype.toLocaleUpperCase = function(opt_locale) {};
/** @return {string} */ String.prototype.toString = function() {};
/** @return {string} */ String.prototype.trim = function() {};
/** @return {string} */ String.prototype.trimEnd = function() {};
/** @return {string} */ String.prototype.trimStart = function() {};
/** @return {string} */ String.prototype.valueOf = function() {};

// ---- Closure-library-shaped names used by the corpus ---------- //

/** @interface */ function Disposable() {}
/** @interface */ function IDisposable() {}

// ---- Math / JSON stubs ---------------------------------------- //

/** @const */ var Math;
/** @const */ var JSON;

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
/** @constructor */ function Event(type, opt_init) {}
/** @constructor @extends {Event} */ function MouseEvent(type, opt_init) {}
/** @constructor @extends {Event} */ function KeyboardEvent(type, opt_init) {}
/** @constructor */ function Node() {}
/** @constructor @extends {Node} */ function Element() {}
/** @constructor @extends {Element} */ function HTMLElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLDivElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLInputElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLAnchorElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLButtonElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLCanvasElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLFormElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLImageElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLSpanElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLTableElement() {}
/** @constructor @extends {HTMLElement} */ function HTMLTextAreaElement() {}

// ---- Node / Element / HTMLElement prototype surface ------------ //
//
// Only the properties the corpus actually reaches for. Strict ports
// of w3c_dom1.js / html5.js would run thousands of lines.

/** @type {?Node} */ Node.prototype.parentNode;
/** @type {?Node} */ Node.prototype.firstChild;
/** @type {?Node} */ Node.prototype.lastChild;
/** @type {?Node} */ Node.prototype.nextSibling;
/** @type {?Node} */ Node.prototype.previousSibling;
/** @type {string} */ Node.prototype.nodeName;
/** @type {number} */ Node.prototype.nodeType;
/** @type {?string} */ Node.prototype.nodeValue;
/** @type {?string} */ Node.prototype.textContent;
/** @return {!Node} */ Node.prototype.appendChild = function(child) {};
/** @return {!Node} */ Node.prototype.removeChild = function(child) {};
/** @return {!Node} */ Node.prototype.insertBefore = function(newNode, ref) {};
/** @return {!Node} */ Node.prototype.replaceChild = function(newChild, oldChild) {};
/** @return {boolean} */ Node.prototype.hasChildNodes = function() {};
/** @return {!Node} */ Node.prototype.cloneNode = function(opt_deep) {};
/** @return {boolean} */ Node.prototype.contains = function(other) {};

/** @type {string} */ Element.prototype.id;
/** @type {string} */ Element.prototype.className;
/** @type {string} */ Element.prototype.tagName;
/** @type {string} */ Element.prototype.innerHTML;
/** @type {string} */ Element.prototype.outerHTML;
/** @type {number} */ Element.prototype.clientWidth;
/** @type {number} */ Element.prototype.clientHeight;
/** @type {number} */ Element.prototype.scrollTop;
/** @type {number} */ Element.prototype.scrollLeft;
/** @return {?string} */ Element.prototype.getAttribute = function(name) {};
/** @return {void} */ Element.prototype.setAttribute = function(name, value) {};
/** @return {void} */ Element.prototype.removeAttribute = function(name) {};
/** @return {boolean} */ Element.prototype.hasAttribute = function(name) {};
/** @return {?Element} */ Element.prototype.querySelector = function(selector) {};
/** @return {!Array<!Element>} */ Element.prototype.querySelectorAll = function(selector) {};
/** @return {void} */ Element.prototype.addEventListener = function(type, listener, opt_opts) {};
/** @return {void} */ Element.prototype.removeEventListener = function(type, listener, opt_opts) {};

/** @type {string} */ HTMLElement.prototype.title;
/** @type {string} */ HTMLElement.prototype.lang;
/** @type {string} */ HTMLElement.prototype.dir;
/** @type {boolean} */ HTMLElement.prototype.hidden;
/** @type {number} */ HTMLElement.prototype.tabIndex;
/** @type {string} */ HTMLElement.prototype.innerText;
/** @type {number} */ HTMLElement.prototype.offsetWidth;
/** @type {number} */ HTMLElement.prototype.offsetHeight;
/** @type {number} */ HTMLElement.prototype.offsetTop;
/** @type {number} */ HTMLElement.prototype.offsetLeft;
/** @type {?HTMLElement} */ HTMLElement.prototype.offsetParent;
/** @return {void} */ HTMLElement.prototype.click = function() {};
/** @return {void} */ HTMLElement.prototype.focus = function() {};
/** @return {void} */ HTMLElement.prototype.blur = function() {};

// ---- Window / document surface --------------------------------- //
//
// The corpus references window.document / window.location and a few
// timer globals. Declaring these as constructors lets user code do
// /** @type {!Window} */ without tripping UNRECOGNIZED_TYPE.

/** @constructor */ function Window() {}
/** @type {!Document} */ Window.prototype.document;
/** @type {!Location} */ Window.prototype.location;
/** @type {!Console} */ Window.prototype.console;
/** @type {!Navigator} */ Window.prototype.navigator;
/** @type {!Window} */ Window.prototype.self;
/** @type {!Window} */ Window.prototype.top;
/** @type {!Window} */ Window.prototype.parent;
/** @type {!History} */ Window.prototype.history;
/** @type {!Storage} */ Window.prototype.localStorage;
/** @type {!Storage} */ Window.prototype.sessionStorage;
/** @type {number} */ Window.prototype.innerWidth;
/** @type {number} */ Window.prototype.innerHeight;
/** @type {number} */ Window.prototype.outerWidth;
/** @type {number} */ Window.prototype.outerHeight;
/** @return {number} */ Window.prototype.setTimeout = function(cb, opt_ms) {};
/** @return {number} */ Window.prototype.setInterval = function(cb, opt_ms) {};
/** @return {void} */ Window.prototype.clearTimeout = function(id) {};
/** @return {void} */ Window.prototype.clearInterval = function(id) {};
/** @return {number} */ Window.prototype.requestAnimationFrame = function(cb) {};
/** @return {void} */ Window.prototype.cancelAnimationFrame = function(id) {};
/** @return {void} */ Window.prototype.addEventListener = function(type, listener, opt_opts) {};
/** @return {void} */ Window.prototype.removeEventListener = function(type, listener, opt_opts) {};

/** @constructor @extends {Node} */ function Document() {}
/** @type {?HTMLElement} */ Document.prototype.body;
/** @type {?HTMLElement} */ Document.prototype.head;
/** @type {?HTMLElement} */ Document.prototype.documentElement;
/** @type {string} */ Document.prototype.title;
/** @type {string} */ Document.prototype.URL;
/** @type {string} */ Document.prototype.cookie;
/** @return {?Element} */ Document.prototype.getElementById = function(id) {};
/** @return {!Array<!Element>} */ Document.prototype.getElementsByTagName = function(name) {};
/** @return {!Array<!Element>} */ Document.prototype.getElementsByClassName = function(name) {};
/** @return {?Element} */ Document.prototype.querySelector = function(selector) {};
/** @return {!Array<!Element>} */ Document.prototype.querySelectorAll = function(selector) {};
/** @return {!HTMLElement} */ Document.prototype.createElement = function(tag) {};
/** @return {!Text} */ Document.prototype.createTextNode = function(text) {};

/** @constructor @extends {Node} */ function Text() {}

/** @constructor */ function Location() {}
/** @type {string} */ Location.prototype.href;
/** @type {string} */ Location.prototype.protocol;
/** @type {string} */ Location.prototype.host;
/** @type {string} */ Location.prototype.hostname;
/** @type {string} */ Location.prototype.port;
/** @type {string} */ Location.prototype.pathname;
/** @type {string} */ Location.prototype.search;
/** @type {string} */ Location.prototype.hash;
/** @return {void} */ Location.prototype.assign = function(url) {};
/** @return {void} */ Location.prototype.replace = function(url) {};
/** @return {void} */ Location.prototype.reload = function() {};

/** @constructor */ function Console() {}
/** @return {void} */ Console.prototype.log = function(var_args) {};
/** @return {void} */ Console.prototype.info = function(var_args) {};
/** @return {void} */ Console.prototype.warn = function(var_args) {};
/** @return {void} */ Console.prototype.error = function(var_args) {};
/** @return {void} */ Console.prototype.debug = function(var_args) {};
/** @return {void} */ Console.prototype.trace = function(var_args) {};

/** @constructor */ function Navigator() {}
/** @type {string} */ Navigator.prototype.userAgent;
/** @type {string} */ Navigator.prototype.language;
/** @type {string} */ Navigator.prototype.platform;
/** @type {boolean} */ Navigator.prototype.onLine;

/** @constructor */ function History() {}
/** @type {number} */ History.prototype.length;
/** @return {void} */ History.prototype.back = function() {};
/** @return {void} */ History.prototype.forward = function() {};
/** @return {void} */ History.prototype.go = function(opt_delta) {};
/** @return {void} */ History.prototype.pushState = function(state, title, opt_url) {};
/** @return {void} */ History.prototype.replaceState = function(state, title, opt_url) {};

/** @constructor */ function Storage() {}
/** @type {number} */ Storage.prototype.length;
/** @return {?string} */ Storage.prototype.getItem = function(key) {};
/** @return {void} */ Storage.prototype.setItem = function(key, value) {};
/** @return {void} */ Storage.prototype.removeItem = function(key) {};
/** @return {void} */ Storage.prototype.clear = function() {};
/** @return {?string} */ Storage.prototype.key = function(index) {};
`

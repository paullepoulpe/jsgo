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
/** @constructor */ function Array(var_args) {}
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
`

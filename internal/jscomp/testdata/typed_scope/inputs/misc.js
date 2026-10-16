/** @constructor */ function Window() {}
var w = new Window();
var sym = Symbol('foo');
var big = 1n;
var notb = !w;
var numr = +w;
var strr = '' + w;
var x = x || {};
try { f(); } catch ({message, errno}) {}
try { g(); } catch (/** @type {string} */ e) {}
for (const {p, q} of []) {}
for (var i = 0; i < 3; i++) {}
label: { let inner = 1; }
/** @const */ var castType = /** @type {number} */ (w);
/** @const */ var assigned = w = 3;
var Foo_1;
let Foo = Foo_1 = class Foo {};
/** @const */ var newFoo = new Foo();
const {a, b: {c}} = {a: 1, b: {c: 2}};
/** @const */ const {d} = {d: 1};
var reflected = goog.reflect.object(Window, {a: 1});
var notACtor = goog.reflect.object(sym, {a: 1});
this.thisProp = 3;
/** @type {number} */ this.thisTyped = 3;

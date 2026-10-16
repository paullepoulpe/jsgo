var m = new Map();
var s = new Set([1, 2]);
var p = Promise.resolve(1).then(function(v) { return v; });
var a = Array.from([1, 2, 3]).includes(2);
var str = 'abc'.startsWith('a') && 'abc'.includes('b');
if (typeof Symbol !== 'undefined') { var sym = Symbol('x'); }
if (Promise.allSettled) { Promise.allSettled([]); }
var w = window.Map;
var g = globalThis.Promise;
var o = Object.assign({}, {a: 1});
var e = Object.entries(o);
function local() { var Map = 1; return Map; }
obj?.includes(1);
getStr().includes('x');
[].includes.call(arr, 0);
var n = Number.isFinite(1) && Math.trunc(1.5);
var r = String.raw`a\nb`;
var wm = new WeakMap();

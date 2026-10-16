/** @constructor @param {number} x */
function Foo(x) { this.x = x; }
/** @return {string} */
Foo.prototype.bar = function() { return 'a'; };
var ns = {};
/** @typedef {{a: number}} */
ns.T;
/** @enum {string} */
var E = {A: 'a', B: 'b'};
class Bar extends Foo { constructor() { super(1); this.y = 2; } /** @param {T} t @template T */ baz(t) {} }
let z = new Bar();
const k = 3;
if (k) { let q = 1; }

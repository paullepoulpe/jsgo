/** @constructor @param {number} x */
function Foo(x) { /** @type {number} */ this.x = x; }
/** @return {string} */
Foo.prototype.name = function() { return 'foo' + this.x; };
class Bar extends Foo {
  constructor() { super(1); this.y = `a${this.x}bé`; }
  get z() { return 2; }
  set z(v) {}
  static s() { return new Bar(); }
  ['comp' + 1]() {}
  field = 3;
  static #priv = 4;
}
/** @enum {string} */
const E = {A: 'a', 'B': 'b'};
/** @typedef {{a: number, b: (string|undefined)}} */
let Rec;
const bar = /** @type {!Foo} */ (new Bar());
var big = 123n, re = /a+b/gi, n = -0.5e3;

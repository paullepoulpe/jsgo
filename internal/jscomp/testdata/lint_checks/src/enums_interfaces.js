/** @enum {number} */
var N = {A: 1, B: 1, C: 2.5, D: 2.5e3, ['E']: 3, F, g: 4};
/** @enum {string} */
const S = {A: 'a', B: `b`, C: `c${x}`, D: 'a', E: 1 + 2};
/** @enum {?} */
let Q = {A: 1};
/** @enum {Foo|Bar} */
let R = {A: 1};
/** @enum */
var T = {A: 0, B: 0};
/** @record */
function Rec() {
  /** @type {number} */
  this.a;
  this.b;
  foo();
}
/** @interface */
function If(x, y) { this.z = 1; }
/** @record */
let RecCls = class extends Base { constructor() { /** @type {string} */ this.s; } m() {} n() { return 1; } static q() {} };
/** @interface @extends {Base} */
class If2 { constructor() {} }

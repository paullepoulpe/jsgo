/** @template T */
class Base {
  /** @param {T} t */
  constructor(t) { /** @const {T} */ this.t = t; }
  /** @return {T} */ get() { return this.t; }
  static create() { return new Base(1); }
  get prop() { return 1; }
  set prop(v) {}
  /** @param {number} n */ set other(n) {}
}
/** @extends {Base<string>} */
class Derived extends Base {
  constructor() { super('x'); this.z = 3; }
  /** @override */ get() { return 'y'; }
  method(a = 1, ...rest) { return a; }
  [Symbol.iterator]() {}
}
class WithFields {
  x = 1;
  /** @type {string} */ y;
  static s = 'q';
  z;
  [Symbol.iterator] = 2;
}
const Expr = class Named { m() { return Named; } };
let anon = class {};
var obj = { cls: class {}, fn: function() {} };
class NoCtorExtends extends Derived {}
function makeMixin(superClass) { return class extends superClass {}; }
class Dyn extends makeMixin(Base) {}
/** @extends {Base<number>} */
class DynOk extends makeMixin(Base) {}

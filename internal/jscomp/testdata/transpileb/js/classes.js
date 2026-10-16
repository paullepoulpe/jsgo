/** @constructor */ function Base(x) { this.x = x; }
Base.prototype.m = function() { return this.x; };
class A extends Base {
  constructor(x, y) { super(x); this.y = y; }
  m() { return super.m() + this.y; }
  static s() { return 1; }
  get g() { return this.y; }
  set g(v) { this.y = v; }
  static get sg() { return 2; }
  ['comp' + 1]() { return 3; }
}
var B = class extends A { n() { return 1; } };
const ns = {};
ns.C = class Inner { static make() { return new Inner(); } };
class E extends Error { constructor(msg) { super(msg); this.name = 'E'; } }
class Arr extends Array {}
class Obj extends Object { constructor() { super(); } }
function makeClass() { return class extends (Math.random() ? A : Base) {}; }
foo(class { bar() {} });
class Fields { a = 1; static b = 2; ['c'] = 3; static { this.d = this.b; } }
class Statics extends A { static x = super.s(); static y = this.x; }
class Spread extends Base { constructor(...args) { super(...args); } }
class Ret extends Base { constructor() { if (!Ret) return; super(1); } }

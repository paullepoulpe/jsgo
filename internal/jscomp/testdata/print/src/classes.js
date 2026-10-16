class A {}
class B extends A { constructor() { super(); super.x(); new.target; } }
class C extends (a, b) {}
class D extends (a = b) {}
class E extends f() {}
class F extends (class {}) {}
var G = class {};
var H = class Named extends A {};
class I {
  field = 1;
  static sfield;
  #priv = 2;
  static #spriv;
  [computed] = 3;
  static [scomputed];
  'quoted' = 4;
  42 = 5;
  get g() { return this.#priv; }
  set s(v) { this.#priv = v; }
  static get sg() {}
  static set ss(v) {}
  async am() {}
  *gm() {}
  async *agm() {}
  static async sam() {}
  [cm]() {}
  get [cg]() {}
  set [cs](v) {}
  async [cam]() {}
  *[cgm]() {}
  'quoted method'() {}
  123() {}
  if() {}
  static { init(); }
  static {}
  #pm() { return #priv in this; }
}
var o = {
  a, b: c, 'd': e, [f]: g, ...h, i() {}, get j() {}, set j(v) {}, async k() {}, *l() {},
  async *m() {}, [n]() {}, get [o]() {}, set [p](v) {}, if: 1, default: 2, class: 3,
  "x-y": 4, x: function () {}, y: () => {}, z: class {},
};
var o2 = {a: 1,};
var arr = [1, , 3, , ];
var arr2 = [, , ];
var arr3 = [...a, ...b];
var kw = a.default + a.if + a.class + a.new + a.delete;

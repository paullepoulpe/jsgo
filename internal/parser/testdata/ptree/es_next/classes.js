class A {}
class B extends A {
  constructor(x, y,) { super(x); super.y = y; }
  static s = 1;
  #priv = 2;
  #m() { return this.#priv; }
  static #sm() {}
  get g() { return 1; }
  set g(v) {}
  static get sg() { return 2; }
  static set sg(v) {}
  get [Symbol.iterator]() {}
  set ['x' + 1](v) {}
  *gen() { yield 1; yield* other(); yield; }
  async am() { await x; }
  async *agm() { for await (const x of y) {} }
  static async *sagm() {}
  'string'() {}
  42() {}
  [computed]() {}
  async 'as'() {}
  async [c2]() {}
  if() {}
  static() {}
  static static() {}
  get() {}
  set() {}
  async() {}
  get;
  set;
  static;
  a;
  b = 1
  c = 2;
  ['d'] = 3;
  static ['e'];
  static { let z = 1; }
  ;;
}
let C = class {};
let D = class Named extends (foo()) {};
var E = class extends B { m() { return super.m(); } };

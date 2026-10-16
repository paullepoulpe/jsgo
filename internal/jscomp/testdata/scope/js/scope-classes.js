class C extends Object {
  static #p = 1;
  static { var s = this; let t = super.toString; }
  x = this;
  [C.y] = 2;
  [k];
  static z = () => this;
  constructor(a) { super(); this.a = a; }
  m(b) { return super.m(b) + this.x; }
  get g() { return 1; }
  set g(v) {}
  static sm() { return this; }
  [Symbol.iterator]() {}
}
const D = class Named { n() { return Named; } };
let E = class { [this] = 1; };
var k;
function outer() {
  class Inner extends outer {
    f = arguments;
    static { arguments; }
  }
  return new Inner();
}

class A {}
class B extends A {
  constructor() { this.x = 1; super(); }
}
class C extends A {
  constructor() { super.foo(); super(); }
}
class D extends A {
  constructor() { }
}
class E extends A {
  constructor() { return {}; }
}
class F {
  constructor() { super(); }
  method() { super(); }
  get g() { super(); return 1; }
  [computed]() { super(); }
}
class G extends A {
  constructor() { const f = () => super(); f(); }
}
class H extends A {
  constructor() { const f = () => super.x; super(); }
}
var o = { m() { super.x; super(); } };
function plain() { super.x; super(); }
super.y;
class I extends A { field = super.z; static { super.w; super(); } }
class J extends A { constructor() { super(); } method() { return () => super.q(); } }
function k() { var s = super; }

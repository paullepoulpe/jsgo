function decl(a, b = a.c, ...rest) { return a + b; }
var fe = function named() { named(); };
var arrow = (x) => x.y;
var arrow2 = (x) => { if (x) return 1; return x.z; };
function* gen() { yield 1; yield* other(); }
async function as() { await p; try { await q; } catch (e) { await r; } }
class A extends B {
  static s = 1;
  f = this.g();
  [computed()] = 2;
  static { if (A.x) init(); }
  static { try { a(); } finally { b(); } }
  constructor() { super(); this.x = 1; }
  method() { return () => this.x; }
  get g() { return 1; }
  set g(v) {}
  static [Symbol.iterator]() {}
}
var C = class { m() {} static { x = 1; } };
function outer() {
  function inner() { return 1; }
  var a = inner();
  function inner2() {}
  return a;
}
new Foo().bar[baz]++;
a = b instanceof C;
c = 'k' in d;
tag`template ${x}`;
x?.y?.(z);
delete a.b;
void 0;
throw new Error('done');

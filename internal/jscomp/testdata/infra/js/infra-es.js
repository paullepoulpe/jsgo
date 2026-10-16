/** @fileoverview Modern syntax for the l11-infra oracle. */
async function* gen(a, ...rest) {
  for await (const x of a) { yield x; }
  const [p, , q = 2, ...others] = rest;
  const {u, v: w = 3, ...more} = others;
  return await Promise.all([p, q, u, w, more]);
}
class Base {
  static count = 0;
  #secret = 1;
  field = this.#secret;
  static { Base.count++; }
  get value() { return this.field; }
  set value(v) { this.field = v; }
  ['computed' + 1]() { return new.target; }
  static create() { return new this(); }
}
class Derived extends Base {
  constructor() { super(); this.x = super.value ?? 0; }
  m() { return super.m?.(); }
}
const o = {a, b: 1, [c]: 2, ...d, get g() { return 1; }, set s(x) {}, method() {}};
let t = `x${o.a}y${o.b}z`;
let r = /ab+c/gi, big = 123n, e = 2 ** 10;
a ||= b; a &&= c; a ??= d;
x?.y?.[z]?.(w);
label: for (const k in o) { if (!k) continue label; else break; }
switch (t) { case 'a': t = 1; break; default: t = 2; case 'b': t = 3; }
try { throw new Error('e'); } catch { t = 4; } finally { t = 5; }
try { f(); } catch ({message}) { g(message); }
do { t--; } while (t > 0 && !(t % 2 || t === 3));
while (!(a && b) || (c ? !d : e)) { a = !a; }
for (var i = 0, j = 10; i < j; i++, j--) { if (i == j) { var same = true; } }
var hook = a ? b : c, comma = (a, b, !c);
if (!(x === 1 || y !== 2)) { z(); }
import('mod').then(m => m.default);
var fn = function named() { return arguments.length; };
var arrow = (p = 1, {q} = {}) => p + q;

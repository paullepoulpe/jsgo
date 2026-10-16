async function f(a, b) {
  const x = await a;
  const y = await b(this, arguments);
  return x + y;
}
var arrow = async (p) => await p;
class C {
  async m() { return super.m2() + await this.n; }
  m2() { return 1; }
}
async function nested() {
  return async () => { await this; return arguments; };
}
async function withTry() {
  try { await f(1, 2); } catch (e) { throw await e; } finally { await 0; }
}

async function* gen(a, b = 1, ...rest) {
  for await (const x of a) { yield* rest; }
  label: for (let i = 0; i < 3; i++) { if (i) continue label; else break; }
  switch (b) { case 1: return; default: throw new Error('x'); }
}
const {p, q: {r = 2} = {}, ...others} = {p: 1};
const [s, , t = 3, ...u] = [1, 2];
let o = {p, [s]: t, m() { return this?.p ?? o?.['q']?.(); }, get g() { return 1; }};
try { delete o.p; o.p ||= 1; o.q &&= 2; o.r ??= 3; } catch (e) { void e; } finally { debugger; }
do { s **= 2; } while (typeof s !== 'number' && !(s in o) || s instanceof Object);
with (o) {}
function nt() { return new.target; }

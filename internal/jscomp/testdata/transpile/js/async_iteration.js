async function* agen() {
  yield 1;
  const v = await Promise.resolve(2);
  yield* other();
  return v;
}
async function* other() { yield this; yield arguments.length; }
async function consume(it) {
  for await (const x of it) { console.log(x); }
  label: for await (let [a, b] of it) { if (a) continue label; }
  let z;
  for await (z of it) {}
}
class K { async *m() { yield super.toString(); } }

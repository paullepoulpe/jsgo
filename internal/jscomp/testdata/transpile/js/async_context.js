async function af(x) {
  try {
    await x;
  } catch (e) {
    console.log(e);
  }
  for await (const y of x) { await y; }
  return await x;
}
function* gf() { yield 1; yield this; }
class M { *gm() { yield super.toString(); yield arguments[0]; } }
async function* agf() { yield await 1; }
const arrow = async () => await 1;

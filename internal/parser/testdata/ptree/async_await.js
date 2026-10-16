async function f() {
  await x;
  await (async () => { await y; })();
  const v = await z ?? 1;
  for await (const q of r) {}
}
var await_ident = await;
async
function notAsync() {}
var a = async;
var b = async(1, 2);
var c = async (x)
=> x;
var d = async () => await 1;
var e = { async f() {}, async: 1, async() {} };
class K { async; async m() {} }

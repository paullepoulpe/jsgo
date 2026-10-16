var x = 1;
let y = 2, z;
const k = 3;
function f(a, b = x, ...rest) {
  var a;
  var arguments;
  let q = a + b;
  for (const z of [q]) { x += z; }
  for (let i = 0; i < 1; i++) { var j = i; }
  for (var p in {}) {}
  switch (q) { case 1: let s = 1; break; default: function sw() {} }
  label: { break label; }
  return function g() { return arguments.length + this.q + g; };
}
var f;
if (x) { function h() { return h; } } else { var h2 = () => this; }
try { f(); } catch (e) { var e2 = e; let e3; } finally { var fin; }
try { f(); } catch ({message, stack: [s0]}) { message; }
try { f(); } catch { }
async function* ag(...[u, v]) { for await (const w of ag()) { yield w; } }
(function iife(a, {b, c: [d]} = {}, [e] = []) { a = b; d++; e--; })();
with (x) { var inWith = 1; }
{ let blockLet = 1; class BlockClass {} function blockFn() {} }
var dup = 1; var dup = 2;
function shadow(param) { let param2; { let param = 1; } }

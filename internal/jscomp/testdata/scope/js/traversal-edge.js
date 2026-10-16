function a({x = y, y} = {}, [z = x] = []) {}
const {m = n, n = 1} = {};
[a1 = b1, b1] = [];
function b(x = function() { return this; }, y = () => this) { 'use strict'; }
class S { static [Symbol.x] = 1; #priv() {} static async *gen() {} }
var obj = { get g() {}, set g(v) {}, m() {}, [k]: 1, ...spread, async am() {}, *gm() {} };
tag`a${x}b${y}`;
function nt() { return new.target; }
import.meta;
for (const [k2, v2] of Object.entries(obj)) {}
(async () => { await 1; for await (x of y) {} })();
x ??= y; x ||= y; x &&= y;

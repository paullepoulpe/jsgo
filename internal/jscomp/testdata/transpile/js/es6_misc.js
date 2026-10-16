for (const x of [1, 2, 3]) { console.log(x); }
label: for (let y of iter()) { if (y) break label; }
var t = `a${1}b${2}c`;
var tagged = tag`x${y}z`;
var raw = String.raw`\n${1}`;
var computed = {[k]: 1, a: 2, ['b' + c]: 3, m() { return 4; }, get g() { return 5; }};
var shorthand = {a, b, m() {}};
class P { static s() {} m() {} get g() { return 1; } set g(v) {} }
class Q extends P { constructor() { super(); } m() { return super.m(); } }
let bs = 1;
{ let bs = 2; const cs = 3; function inner() { return bs + cs; } }
for (let i = 0; i < 3; i++) { setTimeout(() => i); }

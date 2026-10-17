var a = 0, b = foo(), c;
let d = 1, e = 2;
const F = 3, G = 4;
while (a < 10) { a++; }
for (var i = 0, j = 1; i < 10; i++) {}
for (i = 0; i < 1; i++) {}
for (var k in obj) {}
for (var [x1, y1 = 3] of arr) {}
for (var {p1, q1} in obj) {}
async function fa() { for await (var aw of gen()) {} }
lbl: var lv = 1;
lbl2: for (var z = 0;;) { break lbl2; }
lbl3: lbl4: for (var w in o) { continue lbl4; }
var [m, n = 2, ...rest] = arr;
var {o1, o2: {o3}, ...orest} = obj;
let [lm] = arr;
a += 1; b -= 2; c *= 3; d **= 2; a.x += 1; a[0] |= 4;
a ||= 1; b &&= 2; c ??= 3; a.b ||= 4; a[b] ??= 5;
f(1, 2,);
var arr2 = [1, 2,];
var o4 = {a: 1,};
new Foo(1,);
x?.(1,);
function hoisted() {}
if (a) { function inBlock() {} }

if (x); act();
if (x) {} else;
while (x);
for (;;);
for (var k in o);
for (var v of o);
async function af() { for await (var w of o); }
if (x == NaN) {}
if (NaN !== y) {}
if (x < Number.NaN) {}
if ('a' in 3) {}
if (1 instanceof Object) {}
if (('a' + 'b') in o) {}
if (!x in o) {}
var r1 = true && f();
var r2 = 0 || f();
var r3 = (a, void 0) && b;
var r4 = (x ? 1 : 2) || y;
var r5 = function() {} || y;
var r6 = !{} && y;
var r7 = (a?.b).c;
var r8 = (a?.b)();
var r9 = (a?.[b])[c];

function dup(a, a, [b, {c: a}], ...b) {}
lbl: lbl: x;
lbl2: { break lbl2; }
lbl3: { continue lbl3; }
break;
continue;
for (;;) { continue undefinedLabel; }
function inner() { return; }
new.target;
function nt() { return new.target; }
class C { static { return; } }
class D { static { yield; } }
class E { static { await x; } x = await y; }
function notAsync() { await 5; }
function* gen() { yield; function nested() { yield; } }
async function af() { for await (const x of y) {} }
for await (const x of y) {}
if (x) function blockFn() {}
{ function blockFn2() {} }

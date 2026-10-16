if (a) b(); else c();
if (a) { if (b) c(); } else d();
if (a) if (b) c(); else d();
if (a) ; else ;
if (a) {} else {}
if (a) { let x = 1; } else { const y = 2; }
if (a) function f() {}
if (a) do b(); while (c);
label: for (;;) { break label; continue label; }
outer: inner: while (true) break outer;
lbl: { break lbl; }
lbl2: g();
do x++; while (x < 10)
while (a) b();
for (;;) ;
for (var i = 0, j = 1; i < j; i++, j--) ;
for (let [a, b] of c) {}
for (const {a, b} in c) {}
async function af() { for await (const x of y) {} await z; }
switch (a) { case 1: b(); break; case 2: case 3: { c(); } default: d(); }
switch (a) {}
try { a(); } catch (e) { b(); } finally { c(); }
try { a(); } catch { b(); }
try { a(); } finally { b(); }
throw new Error("x");
debugger;
;;;
{ let block = 1; }
{}
function ret() { return; }
function ret2() { return a, b; }
function ret3() { return (/** @type {number} */ (x)); }
var v1, v2 = 1, [v3] = [], {v4} = {};
let l1 = 1; const c1 = 2;

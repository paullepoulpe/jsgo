{ let x = 1; }
{ const y = 2; }
{ class K {} }
{ function f() {} }
{ var v = 1; }
{ { inner(); } }
if (x) { a(); } { b(); }
outer: for (;;) { inner: for (;;) { break outer; } }
l1: { l2: { break l1; } }
l3: while (true) { continue l3; }
l4: x = () => { a: 2 };
function g() { ; }
for (;;);
switch (x) { case f(): case f(): case 'a': case "a": default: }
do { x++ } while (x)
export2 = 1
class Cls { m() {} }

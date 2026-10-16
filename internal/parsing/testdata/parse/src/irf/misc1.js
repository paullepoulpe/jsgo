'use strict';
"use strict";
'not a directive';
function f() { 'use strict'; "use strict"; return 1; }
(0, f)(); (0, a.b)(); (1, 2, g)(); eval('x'); (0, eval)('x'); a.b(); a[b](); a?.b(); a?.[b](); (a.b)(); (/** @type {?} */ (f))();
(a, b, c, d, e);
x = a ? b : c ? d : e;
x ||= 1; x &&= 2; x ??= 3; x **= 2; y = a ?? b; z = a ** b;
with (o) { p; }
switch (x) { case 1: a; break; case 2: default: b; }
try { a(); } catch { b(); } finally { c(); }
try { a(); } catch (e) {} 
try { a(); } finally {}
do x++; while (x < 10)
while (x) x--;
for (var i in o) ; for (const [k, v] of m) ; for (let {a} of b) ;
if (a) b(); else if (c) d(); else e();
debugger;
new Foo; new Foo(1, 2,); new new Bar()();
f(a, ...b, c,); [1, , 3, ...d,]; ({...o, p, q: 1, [r]: 2, s() {}, get t() {}, set t(v) {}, async u() {}, *v() {}, async *w() {}});
(a, ...b) => {}; async x => x; () => ({}); (...[a]) => a;
x = function named() {}; x = class Named extends Base {}; x = class {};
label1: label2: x;
lbl: function labelledFn() {}
lbl: class LabelledClass {}
lbl: let letIdent;

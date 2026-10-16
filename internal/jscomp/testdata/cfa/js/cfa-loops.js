var x = 0;
if (x) { f(); } else if (y) { g(); } else { h(); }
if (a) b();
while (x < 10) { x++; if (x == 5) continue; if (x == 7) break; }
while (true) { if (q()) break; }
do { x--; } while (x > 0);
do { if (x) break; } while (true);
for (var i = 0; i < 10; i++) { if (i == 3) continue; g(i); }
for (;;) { if (z()) break; }
for (let j = 0; ; j++) { if (j > 3) break; }
for (var k in obj) { if (k) continue; use(k); }
for (const v of list) { if (!v) break; use(v); }
async function af() { for await (const w of gen()) { await w; } }
outer: for (var a1 = 0; a1 < 3; a1++) {
  inner: for (var b1 = 0; b1 < 3; b1++) {
    if (b1 == 1) continue outer;
    if (b1 == 2) break outer;
    if (a1 == b1) continue inner;
  }
}
blk: { if (x) break blk; f(); }
lbl: if (x) { break lbl; }
switch (x) { case 1: f(); case 2: g(); break; default: h(); case 3: i(); }
switch (x) { default: h(); }
switch (x) {}
switch (x) { case 1: case 2: break; }

for (var i = 0; i < 10; i++) {}
for (let j = 0, k = 1; ; ) {}
for (;;) break;
for (const x in obj) {}
for (var y in obj) {}
for (z in obj) {}
for (a.b in obj) {}
for (const x of arr) {}
for (let [a, b] of arr) {}
for (x of arr) ;
async function f() {
  for await (const x of gen()) {}
  for await (x of gen()) {}
}
for (var q = 1 in obj) {}
for (let r = 1 of arr) {}
for (var s, t of arr) {}
for (var u, w in obj) {}
for (const c; ; ) {}
for (let [d]; ; ) {}
for (x = 1 in obj) {}
for ((x) in obj) {}
for (f() in obj) {}
for (a, b in obj) {}

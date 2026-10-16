try { f(); } catch (e) { g(e); }
try { f(); } finally { cleanup(); }
try { f(); } catch (e) { g(); } finally { h(); }
try {} catch (e) {}
try { throw new Error('x'); } catch ({message}) { log(message); }
try { a(); } catch { b(); }
while (x) {
  try {
    try {
      break;
    } catch (a) {
    } finally {
      foo();
    }
    fooFollow();
  } catch (b) {
  } finally {
    bar();
  }
  barFollow();
}
END();
for (var i = 0; i < 3; i++) {
  try {
    if (i) break; else continue;
  } finally {
    fin();
  }
}
function r() {
  try {
    try {
      return a();
    } finally {
      b();
    }
  } finally {
    c();
  }
}
function r2(x) {
  try {
    if (x) return 1;
    if (x > 2) throw x;
    return 2;
  } catch (e) {
    return 3;
  } finally {
    d();
  }
}
function r3() {
  for (;;) {
    try {
      try { break; } finally { if (q) continue; }
    } finally { w(); }
  }
  try { return; } catch (e) { x = e.y; } finally { z = 1; }
}
lab: try { break lab; } finally { f(); }

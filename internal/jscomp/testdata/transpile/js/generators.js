function* simple() { yield 1; yield 2; return 3; }
function* withLoops(n) {
  for (var i = 0; i < n; i++) {
    if (i % 2) { continue; }
    yield i;
  }
  var o = {a: 1, b: 2};
  for (var k in o) { yield k; }
  do { yield n--; } while (n > 0);
  outer: for (;;) { inner: for (;;) { if (yield 'x') break outer; continue inner; } }
}
function* withTry() {
  try {
    yield 1;
  } catch (e) {
    yield e;
  } finally {
    yield 'finally';
  }
  try { yield 2; } finally { return 5; }
}
function* withSwitch(x) {
  switch (x) {
    case 1: yield 'one'; break;
    case 2: yield 'two';
    default: yield 'other';
  }
}
function* delegating() { var r = yield* simple(); return r + (yield 7); }
var gen = function*() { var a = this.x + arguments[0]; yield a; };
var obj = { *m() { yield this; } };

outer: for (var i = 0; i < 1; i++) {
  inner: while (true) {
    if (i) continue outer;
    break inner;
  }
}
block: { break block; }
a: b: c: ;
yield: 1;

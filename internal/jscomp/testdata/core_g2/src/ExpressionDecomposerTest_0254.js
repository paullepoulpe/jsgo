function *f() {
  var result$jscomp$0 = yield 1;
  return { a: result$jscomp$0, c: foo(yield 2, yield 3) };
}

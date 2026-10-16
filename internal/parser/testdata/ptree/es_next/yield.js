function* g() {
  yield;
  yield 1;
  yield* g();
  var x = yield;
  var y = yield 2, z;
  f(yield, yield 3);
  yield
  4;
  yield*
  g();
}
function* h() { yield* ; }
function notGen() { yield; }
var gen = function* () { yield (yield 1); };

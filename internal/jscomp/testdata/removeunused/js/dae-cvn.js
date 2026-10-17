function deadAssigns(a) {
  var x = 1;
  x = 2;
  var y = a;
  y = y + 1;
  a = 5;
  return y;
}
function coalesce() {
  var first = 1;
  alert(first);
  var second = 2;
  alert(second);
  for (var i = 0; i < 10; i++) { var t = i * 2; alert(t); }
  let blockScoped = 3;
  alert(blockScoped);
}
function loops(n) {
  var acc = 0;
  while (n > 0) { acc += n; n--; }
  do { acc = acc * 2; } while (acc < 100);
  return acc;
}
function withInner() {
  var captured = 1;
  captured = 2;
  return function() { return captured; };
}
function destructure(o) {
  var {a, b} = o;
  var [c, d] = [a, b];
  c = 1;
  return d;
}
async function asyncFn(p) { var r = await p; r = 3; return 1; }
function* gen() { var g = yield 1; g = 2; }
deadAssigns(1); coalesce(); loops(3); withInner(); destructure({}); asyncFn(); gen();

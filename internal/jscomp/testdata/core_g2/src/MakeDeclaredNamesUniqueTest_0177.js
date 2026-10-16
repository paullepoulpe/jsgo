const x = 'outer';
function foo(a = x) {
  const x = 'inner';
  return a;
}

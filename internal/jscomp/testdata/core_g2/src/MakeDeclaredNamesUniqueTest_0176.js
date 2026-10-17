let x = 'outer';
function foo(bar = baz => x) {
  let x$jscomp$1 = 'inner';
  console.log(bar());
}

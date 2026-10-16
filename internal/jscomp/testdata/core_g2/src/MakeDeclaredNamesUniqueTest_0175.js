let x = 'outer';
function foo(bar = baz => x) {
  let x = 'inner';
  console.log(bar());
}

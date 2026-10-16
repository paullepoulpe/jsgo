const x = 'outerouter';
{
  const x = 'outer';
  function foo(a = x) {
    return a;
  }
foo();
}

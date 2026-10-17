function foo() {
  function bar() {
    return 1;
  }
  if (true) {
    function bar() {
      return 2;
    }
  }
}

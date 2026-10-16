function foo() {
  function bar() {
    return 1;
  }
  if (true) {
    function bar$jscomp$1() {
      return 2;
    }
  }
}

function f() {
  var x = 1;
  function g(/** typeof x */ a) {}
  x = 'str';
  g(null);
}

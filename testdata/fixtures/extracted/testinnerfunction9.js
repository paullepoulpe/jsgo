function f() {
  var x = 3;
  function g() { x = null; };
  function h() { return x == null; }
  return h();
}

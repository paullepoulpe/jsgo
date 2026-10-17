function outer(p1, p2,) {
  foo();
  function inner1() { return 1; }
  bar();
  function inner2() {}
  var v1 = 1, v2;
  return (x) => x + v1;
}
var fe = function named() { return named; };
var arrow = (a, b) => a + b;
var arrow2 = async (a) => await a;
var arrow3 = () => ({a: 1});
class C {
  static {
    foo();
    function sb1() {}
    var sbv = 1, sbw = 2;
  }
  m() { var mv = 1, mw = 2; }
}
function dup() { var d1 = 1; var d1 = 2; var d1; }
function dupParam(a) { var a = 3; var a; }
function dupFn() { var g = 1; function g() {} }
var gdup = 1;
var gdup = 2;
var gdup;
var window = 1;
var Arguments;
function shadow() { var window = 2; { let window = 3; } }
function args() { var arguments; return arguments; }

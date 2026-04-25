/** @constructor */ function Foo() {}
function f() {
  var x = new Foo();
  x.onload = function() {
    x.onload = null;
  };
}

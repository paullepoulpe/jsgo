/** @constructor */ function Foo() {}
function f(z) {
  var x = new Foo();
  {
    x.onload = function() {};
  }
  {
    x.onload = null;
  };
}

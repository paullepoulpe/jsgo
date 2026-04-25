/** @constructor */ function Foo() {}
function f(z) {
  var x = new Foo();
  if (z) {
    x.onload = function() {};
  } else {
    x.onload = null;
  };
}

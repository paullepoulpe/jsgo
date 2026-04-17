/** @constructor */
function Foo() {}
function f(/** ? */ x) {
  return (new Foo) < x;
}

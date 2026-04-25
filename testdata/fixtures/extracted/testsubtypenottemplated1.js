/** @interface @template T */ function A() {}
/** @constructor @implements {A<U>} @template U */ function Foo() {}
function f(/** (!Object|!Foo<string>) */ x) {
  var /** null */ n = x;
}

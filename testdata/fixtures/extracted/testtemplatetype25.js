/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Foo(x) {
  /** @type {T} */
  this.p = x;
}
/** @return {!Foo<!Foo<T>>} */
Foo.prototype.m = function() {
  return new Foo(new Foo(new Object));
};
/** @return {T} */
Foo.prototype.get = function() {
  return this.p;
};
var /** null */ n = new Foo(new Object).m().get();

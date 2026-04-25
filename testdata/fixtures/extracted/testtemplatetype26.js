/**
 * @param {T} bar
 * @constructor
 * @template T
 */
function Bar(bar) {
  /** @type {T} */
  this.bar = bar;
}
/** @return {T} */
Bar.prototype.getBar = function() {
  return this.bar;
};
/**
 * @param {T} foo
 * @constructor
 * @template T
 * @extends {Bar<!Array<T>>}
 */
function Foo(foo) {
  /** @type {T} */
  this.foo = foo;
}
var /** null */ n = new Foo(new Object).getBar();

/**
 * @template {number|string} T
 */
class Foo {
  constructor() {
    /** @type {T} */
    this.x;
  }
  m() {
  this.x = 0;
  }
}
function f(t) { t = new C(); }
